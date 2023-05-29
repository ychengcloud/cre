package main

import (
	"context"
	"log"

	"github.com/ychengcloud/cre/gen"
)

//go:generate go run ./
func main() {
	trueValue := true

	cfg := &gen.Config{
		Project: "testproject",
		Package: "github.com/ychengcloud/cre",
		DSN:     "mysql://root:123456@tcp(localhost:3306)/test?charset=utf8",
		Root:    "templates",
		GenRoot: "actual",

		Templates: []*gen.Template{
			{
				Path:    "single.go.tmpl",
				GenPath: ".",
				Format:  "{{.Schema}}.go",
				Mode:    gen.TplModeSingle,
			},
			{
				Path:    "multi.tmpl",
				GenPath: "multi",
				Format:  "{{.Table}}.go",
				Mode:    gen.TplModeMulti,
			},
			{
				Path:    "path/path.tmpl",
				GenPath: "path",
				Format:  "{{.Schema}}.go",
				Mode:    gen.TplModeSingle,
			},
			{
				Path:    "format/go.tmpl",
				GenPath: "format",
				Format:  "{{.Schema}}.go",
				Mode:    gen.TplModeSingle,
			},
			{
				Path:    "format/proto.tmpl",
				GenPath: "format",
				Format:  "{{.Schema}}.proto",
				Mode:    gen.TplModeSingle,
			},
		},

		Tables: []*gen.Table{
			{
				Name: "user",
				Fields: []*gen.Field{
					{
						Name:       "name",
						Comment:    "user name",
						Operations: []string{"EQ", "In"},
					},
					{
						Name: "post",
						Relation: &gen.Relation{
							Name:     "Post",
							Type:     "HasMany",
							RefTable: "post",
							RefField: "id",
						},
					},
				},
			},
			{
				Name: "post",
				Fields: []*gen.Field{
					{Name: "id", Filterable: &trueValue, Operations: []string{"EQ", "In"}},
					{
						Name: "author",
						Relation: &gen.Relation{
							Name:     "User",
							Type:     "HasOne",
							RefTable: "user",
							RefField: "id",
						},
					},
					{
						Name: "category",
						Relation: &gen.Relation{
							Name:     "Category",
							Type:     "ManyToMany",
							RefTable: "category",
							RefField: "id",
							JoinTable: &gen.JoinTable{
								Name:     "post_category",
								Field:    "post_id",
								RefField: "category_id",
							},
						},
					},
				},
			},
		},
	}
	g, err := gen.NewGenerator(cfg, NewFakeLoader())
	if err != nil {
		log.Fatal(err)
	}

	if err := g.Generate(context.Background()); err != nil {
		log.Fatal(err)
	}

}
