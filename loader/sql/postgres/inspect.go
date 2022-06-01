package postgres

import (
	"context"
	"database/sql"
	"fmt"

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
	ctype   string
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
	rows, err := i.querySqlRows(ctx, VariablesQuery)
	if err != nil {
		return fmt.Errorf("postgres/dbInfo: query fail: %w", err)
	}

	defer rows.Close()

	settings := make([]string, 0, 3)
	for rows.Next() {
		var setting string
		if err := rows.Scan(&setting); err != nil {
			return fmt.Errorf("postgres/dbInfo: scan fail: %w", err)
		}
		settings = append(settings, setting)
	}

	i.collate, i.ctype, i.version = settings[0], settings[1], settings[2]
	if len(i.version) != 6 {
		return fmt.Errorf("postgres/dbInfo: malformed version: %s", i.version)
	}

	i.version = fmt.Sprintf("%s.%s.%s", i.version[:2], i.version[2:4], i.version[4:])
	if semver.Compare("v"+i.version, "v10.0.0") != -1 {
		return fmt.Errorf("postgres/dbInfo: unsupported version: %s", i.version)
	}

	return nil
}

func (i *inspect) tables(ctx context.Context) ([]*schema.Table, error) {
	var (
		name    string
		comment sql.NullString
	)

	rows, err := i.querySqlRows(ctx, TablesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*schema.Table
	for rows.Next() {
		if err := rows.Scan(&name, &comment); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("postgres/tables: no tables found for schema [%q]", i.schema.Name)
			}
			return nil, fmt.Errorf("postgres/tables: Scan [%w]", err)
		}

		table := &schema.Table{
			Name:    name,
			Schema:  i.schema,
			Comment: comment.String,
		}
		tables = append(tables, table)

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres/tables: rows [%w]", err)
	}

	return tables, nil
}

func (i *inspect) inspectColumns(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		columns, err := i.columns(ctx, table)
		if err != nil {
			return fmt.Errorf("postgres/columns: columns [%w]", err)
		}
		table.Columns = columns
	}
	return nil
}

func (i *inspect) inspectIndexes(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		indexes, err := i.indexes(ctx, table.Name)
		if err != nil {
			return fmt.Errorf("postgres/indexes: indexes [%w]", err)
		}
		table.Indexes = indexes
	}
	return nil
}

func (i *inspect) inspectForeignKeys(ctx context.Context) error {
	for _, table := range i.schema.Tables {
		foreignKeys, err := i.foreignKeys(ctx, table)
		if err != nil {
			return fmt.Errorf("postgres/fks: fks [%w]", err)
		}
		table.ForeignKeys = foreignKeys
	}
	return nil
}

func (i *inspect) columns(ctx context.Context, t *schema.Table) ([]*schema.Column, error) {
	args := []any{t.Name}

	rows, err := i.querySqlRows(ctx, ColumnsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres/columns: query rows [%w]", err)
	}
	defer rows.Close()

	var columns []*schema.Column
	for rows.Next() {
		var (
			typid, maxlen, precision, scale                                               sql.NullInt64
			name, dataType, nullable, defaults, udt, charset, collation, comment, typtype sql.NullString
		)

		if err := rows.Scan(&name, &dataType, &comment, &nullable, &defaults, &charset, &collation, &precision, &scale, &maxlen, &udt, &typtype, &typid); err != nil {
			return nil, fmt.Errorf("postgres/columns: scan rows [%w]", err)
		}

		ci := &columnInfo{
			dataType:  dataType.String,
			nullable:  nullable.String == "YES",
			precision: precision.Int64,
			scale:     scale.Int64,
			size:      maxlen.Int64,
			charset:   charset.String,
			collation: collation.String,
			udt:       udt.String,
			typtype:   typtype.String,
		}
		ct, err := parseType(ci)
		if err != nil {
			return nil, fmt.Errorf("postgres/columns: parse type [%s, %w]", t.Name, err)
		}
		column := &schema.Column{
			Name:      name.String,
			Type:      ct,
			Comment:   comment.String,
			Nullable:  nullable.String == "YES",
			Charset:   charset.String,
			Collation: collation.String,

			Table: t,
		}
		switch column.Type.(type) {
		case *spec.FloatType:
			column.Precision = int(precision.Int64)
			column.Scale = int(scale.Int64)

		}
		i.enumValues(ctx, typid.Int64, column)
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres/columns: rows [%w]", err)
	}

	return columns, nil
}

func (i *inspect) enumValues(ctx context.Context, id int64, column *schema.Column) error {
	if _, ok := column.Type.(*spec.EnumType); !ok {
		return fmt.Errorf("postgres/enumValues: column [%s] is not an enum", column.Name)
	}

	rows, err := i.querySqlRows(ctx, EnumQuery, id)
	if err != nil {
		return fmt.Errorf("postgres/enumValues: query rows [%w]", err)
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("postgres/enumValues: scan rows [%w]", err)
		}
		values = append(values, v)
	}
	column.Type.(*spec.EnumType).Values = values

	return nil
}

func (i *inspect) indexes(ctx context.Context, name string) ([]*schema.Index, error) {
	args := []any{name}

	rows, err := i.querySqlRows(ctx, IndexesQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexs []*schema.Index
	indexMap := make(map[string]*schema.Index)
	for rows.Next() {
		var (
			name, idxType                                          string
			column, constraintType, predicate, expression, comment sql.NullString
			primary, unique                                        bool
			asc, desc, nullsFirst, nullsLast                       sql.NullBool
		)

		if err := rows.Scan(&name, &idxType, &column, &primary, &unique, &constraintType, &predicate, &expression, &asc, &desc, &nullsFirst, &nullsLast, &comment); err != nil {
			return nil, fmt.Errorf("postgres/indexes: scanning index: %w", err)
		}

		index, ok := indexMap[name]
		if !ok {
			index = &schema.Index{
				Name:    name,
				Unique:  unique,
				Type:    idxType,
				Comment: comment.String,
			}
			if primary {
				index.Primary = true
			}

			indexMap[name] = index
			indexs = append(indexs, index)
		}

		idxCol := &schema.IndexColumn{
			SeqNo:  len(index.IndexColumns) + 1,
			Column: column.String,
		}

		index.IndexColumns = append(index.IndexColumns, idxCol)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return indexs, nil
}

func (i *inspect) foreignKeys(ctx context.Context, t *schema.Table) ([]*schema.ForeignKey, error) {
	args := []any{t.Name}

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
			return nil, fmt.Errorf("postgres/fks: scanning fk: %w", err)
		}

		foreignKey, ok := foreignKeysMap[name]
		rt := t.Schema.Table(refTable)
		if rt == nil {

			return nil, fmt.Errorf("postgres/fks: ref table %q not found for fk %q", refTable, name)
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
			return nil, fmt.Errorf("postgres/fks: column %q not found for fk %q", column, foreignKey.Name)

		}
		foreignKey.Columns = append(foreignKey.Columns, c)

		rc := foreignKey.RefTable.Column(refColumn)
		if rc == nil {
			return nil, fmt.Errorf("postgres/fks: ref column %q not found for fk %q", refColumn, foreignKey.Name)
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
		return nil, fmt.Errorf("postgres: invalid type %T. expect *sql.Rows for result", result)
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
		return nil, fmt.Errorf("postgres: invalid type %T. expect *sqlld.Row for result", result)
	}
	return row, nil
}
