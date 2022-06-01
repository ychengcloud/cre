package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/ychengcloud/cre"
	schema "github.com/ychengcloud/cre/loader/sql"
	"github.com/ychengcloud/cre/spec"
)

type inspect struct {
	Driver cre.Driver

	schema *schema.Schema

	// Database information
	version string
	collate string
	charset string
}

var _ schema.Inspector = (*inspect)(nil)

func NewInspector(drv cre.Driver) schema.Inspector {
	return &inspect{Driver: drv}
}
func (i *inspect) Inspect(ctx context.Context, name string) (*schema.Schema, error) {
	i.schema = &schema.Schema{Name: name}

	err := i.dbInfo(ctx)
	if err != nil {
		return nil, err
	}

	tables, err := i.tables(ctx)
	if err != nil {
		return nil, err
	}
	i.schema.Tables = tables

	err = i.inspectColumns(ctx)
	if err != nil {
		return nil, err
	}

	err = i.inspectIndexes(ctx)
	if err != nil {
		return nil, err
	}

	err = i.inspectForeignKeys(ctx)
	if err != nil {
		return nil, err
	}
	return i.schema, nil
}

func (i *inspect) dbInfo(ctx context.Context) error {
	row, err := i.querySqlRow(ctx, VariablesQuery)
	if err != nil {
		return err
	}

	if err := row.Scan(&i.version, &i.collate, &i.charset); err != nil {
		return err
	}
	return nil
}

// compareVersion returns an integer comparing two versions according to
// semantic version precedence.
func (i *inspect) compareVersion(version string) int {
	v := i.version
	if i.mariadb() {
		v = v[:strings.Index(v, "MariaDB")-1]
	}
	return semver.Compare("v"+v, "v"+version)
}

// indexExpr check if the connected database supports
// index expressions (functional key part).
func (i *inspect) indexExpr() bool {
	return !i.mariadb() && i.compareVersion("8.0.13") >= 0
}

// mariadb check if the Driver is connected to a MariaDB database.
func (i *inspect) mariadb() bool {
	return strings.Index(i.version, "MariaDB") > 0
}

func (i *inspect) tables(ctx context.Context) ([]*schema.Table, error) {
	var (
		name                                 string
		autoIncrement                        sql.NullInt64
		charset, collation, comment, options sql.NullString
	)

	rows, err := i.querySqlRows(ctx, TablesQuery, i.schema.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*schema.Table
	for rows.Next() {
		if err := rows.Scan(&name, &charset, &collation, &autoIncrement, &comment, &options); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("mysql/tables: no tables found for schema [%q]", i.schema.Name)
			}
			return nil, err
		}

		table := &schema.Table{
			Name:          name,
			Schema:        i.schema,
			Charset:       charset.String,
			Collation:     collation.String,
			AutoIncrement: int(autoIncrement.Int64),
			Comment:       comment.String,
			Options:       options.String,
		}
		tables = append(tables, table)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (i *inspect) inspectColumns(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		columns, err := i.columns(ctx, table)
		if err != nil {
			return err
		}
		table.Columns = columns
	}
	return nil
}

func (i *inspect) inspectIndexes(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		indexes, err := i.indexes(ctx, table.Name)
		if err != nil {
			return err
		}
		table.Indexes = indexes
	}
	return nil
}

func (i *inspect) inspectForeignKeys(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		foreignKeys, err := i.foreignKeys(ctx, table)
		if err != nil {
			return err
		}
		table.ForeignKeys = foreignKeys
	}
	return nil
}

func (i *inspect) columns(ctx context.Context, t *schema.Table) ([]*schema.Column, error) {
	args := []any{i.schema.Name, t.Name}

	rows, err := i.querySqlRows(ctx, ColumnsQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*schema.Column
	for rows.Next() {
		var name, colType, comment, nullable, key, defaults, extra, charset, collation sql.NullString
		var precision, scale sql.NullInt64
		if err := rows.Scan(&name, &colType, &comment, &nullable, &key, &defaults, &extra, &charset, &collation, &precision, &scale); err != nil {
			return nil, err
		}

		ct, err := ParseType(colType.String)
		if err != nil {
			return nil, err
		}
		column := &schema.Column{
			Name:      name.String,
			Type:      ct,
			Comment:   comment.String,
			Nullable:  nullable.String == "YES",
			Charset:   charset.String,
			Collation: collation.String,
			Primary:   key.String == "PRI",
			Table:     t,
		}

		parseExtra(column, extra.String)

		switch column.Type.(type) {
		case *spec.FloatType:
			column.Precision = int(precision.Int64)
			column.Scale = int(scale.Int64)

		}
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (i *inspect) indexes(ctx context.Context, name string) ([]*schema.Index, error) {
	args := []any{i.schema.Name, name}

	query := IndexesQuery
	if i.indexExpr() {
		query = IndexesExprQuery
	}

	rows, err := i.querySqlRows(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexs []*schema.Index
	indexMap := make(map[string]*schema.Index)
	for rows.Next() {
		var (
			nonunique                      bool
			seqno                          int
			name, indexType                string
			subPart                        sql.NullInt64
			column, expr, comment, collate sql.NullString
		)
		if err := rows.Scan(&name, &column, &nonunique, &seqno, &indexType, &collate, &comment, &subPart, &expr); err != nil {
			return nil, fmt.Errorf("mysql/indexes: scanning index: %w", err)
		}

		index, ok := indexMap[name]
		if !ok {
			index = &schema.Index{
				Name:    name,
				Unique:  !nonunique,
				Type:    indexType,
				Comment: comment.String,
			}
			if name == "PRIMARY" {
				index.Primary = true
			}

			indexMap[name] = index
			indexs = append(indexs, index)
		}

		part := &schema.IndexColumn{
			SeqNo:  seqno,
			Column: column.String,
		}

		if subPart.Int64 > 0 {
			part.Sub = int(subPart.Int64)
		}
		index.IndexColumns = append(index.IndexColumns, part)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return indexs, nil
}

func (i *inspect) foreignKeys(ctx context.Context, t *schema.Table) ([]*schema.ForeignKey, error) {
	args := []any{i.schema.Name, t.Name}

	query := ForeignKeysQuery

	rows, err := i.querySqlRows(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []*schema.ForeignKey
	foreignKeysMap := make(map[string]*schema.ForeignKey)
	for rows.Next() {
		var (
			name, table, column, tableSchema, refTable, refColumn, refSchema, updateRule, deleteRule string
		)
		if err := rows.Scan(&name, &table, &column, &tableSchema, &refTable, &refColumn, &refSchema, &updateRule, &deleteRule); err != nil {
			return nil, fmt.Errorf("mysql/fks: scanning fk: %w", err)
		}

		foreignKey, ok := foreignKeysMap[name]
		rt := t.Schema.Table(refTable)
		if rt == nil {

			return nil, fmt.Errorf("mysql/fks: ref table %q not found for fk %q", refTable, name)
		}

		if !ok {
			foreignKey = &schema.ForeignKey{
				Name:  name,
				Table: t,
				//目前只支持引用同一个数据库的外键
				RefTable: rt,
				OnUpdate: schema.ReferenceOption(updateRule),
				OnDelete: schema.ReferenceOption(deleteRule),
			}

			foreignKeysMap[name] = foreignKey
			foreignKeys = append(foreignKeys, foreignKey)
		}

		c := t.Column(column)
		if c == nil {
			return nil, fmt.Errorf("mysql/fks: column %q not found for fk %q", column, foreignKey.Name)

		}
		foreignKey.Columns = append(foreignKey.Columns, c)

		rc := foreignKey.RefTable.Column(refColumn)
		if rc == nil {
			return nil, fmt.Errorf("mysql/fks: ref column %q not found for fk %q", refColumn, foreignKey.Name)
		}
		foreignKey.RefColumns = append(foreignKey.RefColumns, rc)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return foreignKeys, nil
}

func (i *inspect) querySqlRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	result, err := i.Driver.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	rows, ok := result.(*sql.Rows)
	if !ok {
		return nil, fmt.Errorf("mysql: invalid type %T. expect *sql.Rows for result", result)
	}
	return rows, nil
}
func (i *inspect) querySqlRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {

	result, err := i.Driver.QueryRow(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	row, ok := result.(*sql.Row)
	if !ok {
		return nil, fmt.Errorf("mysql: invalid type %T. expect *sqlld.Row for result", result)
	}
	return row, nil
}
