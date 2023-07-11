package gen

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	bufparser "github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/tools/imports"

	"github.com/ychengcloud/cre"
	"github.com/ychengcloud/cre/spec"
)

type Generator struct {
	Cfg *Config

	Loader cre.Loader
	Binder *Binder
	schema *spec.Schema

	templates map[string]*template.Template
	root      fs.FS
	assets    *assets
}

type schemaData struct {
	*spec.Schema

	Generator *Generator
	ImportPkg []string
	Project   string
	Package   string
}

type tableData struct {
	*spec.Table

	M2MField  *spec.Field
	Generator *Generator
	ImportPkg []string
	Project   string
	Package   string
}

type file struct {
	path    string
	content []byte
}
type assets struct {
	dirs  []string
	files []file
}

type assetName struct {
	Package string
	Schema  string
	Table   string
	Path    string
}

func NewGenerator(cfg *Config, loader cre.Loader) (*Generator, error) {
	g := &Generator{
		Cfg:    cfg,
		Loader: loader,
		Binder: &Binder{Dialect: loader.Dialect()},
		assets: &assets{},
	}
	g.templates = make(map[string]*template.Template)

	return g, nil
}

func schemaName(dialect, dsn string) (string, error) {
	switch dialect {
	case cre.MySQL:
		cfg, err := mysql.ParseDSN(dsn)
		if err != nil {
			return "", err
		}
		return cfg.DBName, nil
	case cre.Postgres:
		return "public", nil
	default:
		return "", fmt.Errorf("unsupported dialect: %s", dialect)
	}
}

func (g *Generator) Generate(ctx context.Context) error {

	if err := g.loadTemplates(); err != nil {
		return err
	}

	sn, err := schemaName(g.Loader.Dialect(), g.Cfg.DSN)
	if err != nil {
		return err
	}
	g.schema, err = g.Loader.Load(ctx, sn)
	if err != nil {
		return err
	}
	g.schema, err = mergeSchema(g.schema, g.Cfg)
	if err != nil {
		return err
	}

	if err := g.checkTables(); err != nil {
		return err
	}

	for _, t := range g.Cfg.Templates {
		switch t.Mode {

		case TplModeMulti:
			if t.M2M {
				if err := g.generateM2M(t); err != nil {
					return err
				}
			} else {
				if err := g.generateMulti(t); err != nil {
					return err
				}
			}

		default:
			if err := g.generateSingle(t); err != nil {
				return err
			}
		}

	}

	if err := g.assets.write(); err != nil {
		return err
	}
	if err := g.assets.format(); err != nil {
		return err
	}

	return nil
}

func MustParse(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return t
}

func (g *Generator) loadTemplates() error {
	g.root = os.DirFS(g.Cfg.Root)

	walkFn := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		t := template.New(filepath.Base(path)).Funcs(sprig.GenericFuncMap()).Funcs(Funcs)

		if g.Cfg.Delim.Left != "" && g.Cfg.Delim.Right != "" {
			t.Delims(g.Cfg.Delim.Left, g.Cfg.Delim.Right)
		}

		t, err = t.ParseFS(g.root, path)
		if err != nil {
			return fmt.Errorf("walk: %w", err)
		}
		g.templates[path] = t
		return nil
	}

	err := fs.WalkDir(g.root, ".", walkFn)
	if err != nil {
		return fmt.Errorf("loadTemplates: %w", err)
	}

	return nil
}

func fileName(format string, data any) (string, error) {
	b := bytes.NewBuffer(nil)

	tlp := MustParse(
		template.New("assetName").
			Funcs(sprig.GenericFuncMap()).
			Funcs(Funcs).
			Parse(format))
	if err := tlp.ExecuteTemplate(b, "assetName", data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (g *Generator) file(t *Template, data any, content []byte) error {
	name, err := fileName(t.Format, data)
	if err != nil {
		return err
	}

	d := path.Dir(name)

	if d != "." && d != "/" {
		g.assets.dirs = append(g.assets.dirs, path.Join(t.GenPath, d))
	}

	g.assets.files = append(g.assets.files, file{
		path:    filepath.Join(g.Cfg.GenRoot, t.GenPath, name),
		content: content,
	})
	return nil
}

func goImportPkgs(r io.Reader) ([]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", r, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("parse: %s", err)
	}

	var pkgs []string
	for _, spec := range f.Imports {
		var pkg string
		name := spec.Name
		if name != nil {
			pkg = name.String()
		} else {
			pkg, err = strconv.Unquote(spec.Path.Value)
			if err != nil {
				return nil, err
			}
		}

		pkgs = append(pkgs, filepath.Base(pkg))
	}
	return pkgs, nil
}

func (g *Generator) generateSingle(tplCfg *Template) error {
	var err error

	g.assets.dirs = append(g.assets.dirs, filepath.Join(g.Cfg.GenRoot, tplCfg.GenPath))

	s := schemaData{
		Schema:    g.schema,
		Project:   g.Cfg.Project,
		Package:   g.Cfg.Package,
		Generator: g,
	}

	b := bytes.NewBuffer(nil)

	t, ok := g.templates[tplCfg.Path]
	if !ok {
		return fmt.Errorf("generateSingle load template %s fail", tplCfg.Path)
	}

	if err := t.Execute(b, s); err != nil {
		return fmt.Errorf("generateSingle Execute : %s : %s", tplCfg.Path, err.Error())
	}

	ext := filepath.Ext(tplCfg.Format)
	if ext == ".go" {
		if s.ImportPkg, err = goImportPkgs(b); err != nil {
			return fmt.Errorf("parse import pkgs : %s : %s", tplCfg.Path, err.Error())
		}
	}

	b.Reset()
	if err := t.Execute(b, s); err != nil {
		return err
	}

	if err := g.file(tplCfg, &s, b.Bytes()); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateMulti(tplCfg *Template) error {

	g.assets.dirs = append(g.assets.dirs, filepath.Join(g.Cfg.GenRoot, tplCfg.GenPath))

	if tplCfg.M2M {
		return fmt.Errorf("Is m2m template ? %s : %s", tplCfg.Path, tplCfg.Format)
	}

	for _, table := range g.schema.Tables() {
		td := tableData{
			Table:     table,
			Project:   g.Cfg.Project,
			Package:   g.Cfg.Package,
			Generator: g,
		}

		if err := render(g, tplCfg, &td); err != nil {
			return err
		}
	}
	return nil
}

// 处理 Many To Many 字段
func (g *Generator) generateM2M(tplCfg *Template) error {

	g.assets.dirs = append(g.assets.dirs, filepath.Join(g.Cfg.GenRoot, tplCfg.GenPath))

	if !tplCfg.M2M {
		return fmt.Errorf("Is not m2m template ? %s : %s", tplCfg.Path, tplCfg.Format)
	}

	for _, table := range g.schema.Tables() {
		for _, field := range table.SortedFields() {
			if !field.RelManyToMany() {
				continue
			}
			td := tableData{
				Table:     table,
				M2MField:  field,
				Project:   g.Cfg.Project,
				Package:   g.Cfg.Package,
				Generator: g,
			}

			if err := render(g, tplCfg, &td); err != nil {
				return err
			}
		}

	}
	return nil
}

func render(g *Generator, tplCfg *Template, td *tableData) error {

	b := bytes.NewBuffer(nil)

	t, ok := g.templates[tplCfg.Path]
	if !ok {
		return fmt.Errorf("generateMulti load template %s fail", tplCfg.Path)
	}

	if err := t.Execute(b, td); err != nil {
		return fmt.Errorf("generateMulti Execute : %s : %s", tplCfg.Path, err.Error())
	}
	ext := filepath.Ext(tplCfg.Format)
	if ext == ".go" {
		var err error
		if td.ImportPkg, err = goImportPkgs(b); err != nil {
			return fmt.Errorf("generateMulti import pkgs : %s : %s", tplCfg.Path, err.Error())
		}
	}

	b.Reset()
	if err := t.Execute(b, td); err != nil {
		return err
	}

	if err := g.file(tplCfg, td, b.Bytes()); err != nil {
		return err
	}

	return nil
}

func (g *Generator) checkTables() error {
	for _, table := range g.schema.Tables() {
		if table.Name == "" {
			return fmt.Errorf("table: name is empty")
		}
		if table.ID == nil {
			return fmt.Errorf("table [%s]: id is empty, Is Join Table?", table.Name)
		}
	}
	return nil
}

func (g *Generator) Template(name string, v any) (string, error) {
	b := bytes.NewBuffer(nil)

	t, ok := g.templates[name]
	if !ok {
		return "", fmt.Errorf("exec template %s fail", name)
	}

	if err := t.Execute(b, v); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (g *Generator) Schema() *spec.Schema {
	return g.schema
}

func (a assets) write() error {
	for _, d := range a.dirs {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	for _, f := range a.files {
		if err := os.WriteFile(f.path, f.content, 0644); err != nil {
			return fmt.Errorf("write file %q: %w", f.path, err)
		}
	}
	return nil
}

func (a assets) formatProto(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	fileNode, err := bufparser.Parse("", f, reporter.NewHandler(nil))
	if err != nil {
		return err
	}

	f.Close()

	f, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := newFormatter(f, fileNode).Run(); err != nil {
		return err
	}

	return nil
}

func (a assets) format() error {
	for _, file := range a.files {
		path := file.path
		content := file.content
		ext := filepath.Ext(path)

		var err error

		switch ext {
		case ".go":
			content, err = imports.Process(path, file.content, nil)
			if err != nil {
				return fmt.Errorf("format file %s: %v", path, err)
			}
			if err := os.WriteFile(path, content, 0644); err != nil {
				return fmt.Errorf("write file %s: %v", path, err)
			}
		case ".proto":
			if err := a.formatProto(path); err != nil {
				return fmt.Errorf("write file %s: %v", path, err)
			}
		}

	}
	return nil
}
