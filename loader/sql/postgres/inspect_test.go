package postgres

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/ychengcloud/cre"
	schema "github.com/ychengcloud/cre/loader/sql"
	"github.com/ychengcloud/cre/spec"
)

type postgresMock struct {
	sqlmock.Sqlmock
}

func (m postgresMock) info() {
	m.ExpectQuery(Escape(VariablesQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"setting"}).
			AddRow("en_US.utf8").
			AddRow("en_US.utf8").
			AddRow("100000"))
}

func (m postgresMock) noColumns() {
	m.ExpectQuery(Escape(ColumnsQuery)).
		WillReturnRows(sqlmock.NewRows(ColumnsQueryFields))

}

func (m postgresMock) noIndexes() {
	m.ExpectQuery(Escape(IndexesQuery)).
		WillReturnRows(sqlmock.NewRows(IndexesQueryFields))

}

func (m postgresMock) noForeignKeys() {
	m.ExpectQuery(Escape(ForeignKeysQuery)).
		WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields))
}
func TestInspectTable(t *testing.T) {
	schemaName, tableName, fkTableName := "test", "table", "fk"
	tests := []struct {
		name     string
		before   func(postgresMock)
		expected func() *schema.Schema
		wantErr  bool
	}{
		{
			name: "no table",
			before: func(mock postgresMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME", "COMMENT"}))
			},
			expected: func() *schema.Schema {
				return &schema.Schema{
					Name: "test",
				}
			},
		},
		{
			name: "custom schema",
			before: func(mock postgresMock) {
				mock.info()

				mock.ExpectQuery(Escape(TablesQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME", "COMMENT"}).
						AddRow(tableName, "Comment"))

					// http://www.postgres.cn/docs/14/datatype.html
					// column_name, data_type, comment, is_nullable, column_default, character_set_name, collation_name, numeric_precision, numeric_scale, character_maximum_length, udt_name, typtype, oid
				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs(tableName).
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("bigint", "bigint", "bigint comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1).
						AddRow("bigserial", "bigserial", "bigserial comment", "NO", nil, nil, nil, 64, 0, nil, "serial8", "b", 2).
						AddRow("bit", "bit", "bit comment", "YES", nil, nil, nil, nil, nil, 1, "bit", "b", 3).
						AddRow("bit varying", "bit varying", "bit varying comment", "YES", nil, nil, nil, nil, nil, 255, "varbit", "b", 4).
						AddRow("boolean", "boolean", "boolean comment", "YES", nil, nil, nil, nil, nil, nil, "bool", "b", 5).
						AddRow("box", "box", "box comment", "YES", nil, nil, nil, nil, nil, nil, nil, "b", 6).
						AddRow("bytea", "bytea", "bytea comment", "YES", nil, nil, nil, nil, nil, 1, nil, "b", 7).
						AddRow("character", "character", "character comment", "YES", nil, nil, nil, nil, nil, 1, "bpchar", "b", 8).
						AddRow("character varying", "character varying", "character varying comment", "YES", nil, nil, nil, nil, nil, 255, "varchar", "b", 9).
						AddRow("cidr", "cidr", "cidr comment", "YES", nil, nil, nil, nil, nil, nil, "cidr", "b", 10).
						AddRow("circle", "circle", "circle comment", "YES", nil, nil, nil, nil, nil, nil, "circle", "b", 11).
						AddRow("date", "date", "date comment", "YES", nil, nil, nil, nil, nil, nil, "date", "b", 12).
						AddRow("double precision", "double precision", "double precision comment", "YES", nil, nil, nil, nil, nil, nil, "float8", "b", 13).
						AddRow("inet", "inet", "inet comment", "YES", nil, nil, nil, nil, nil, nil, "inet", "b", 14).
						AddRow("integer", "integer", "integer comment", "YES", nil, nil, nil, 16, 0, nil, "int4", "b", 15).
						AddRow("interval", "interval", "interval comment", "YES", nil, nil, nil, nil, nil, nil, "interval", "b", 16).
						AddRow("json", "json", "json comment", "YES", nil, nil, nil, nil, nil, nil, "json", "b", 17).
						AddRow("jsonb", "jsonb", "jsonb comment", "YES", nil, nil, nil, nil, nil, nil, "jsonb", "b", 18).
						AddRow("line", "line", "line comment", "YES", nil, nil, nil, nil, nil, nil, "line", "b", 19).
						AddRow("lseg", "lseg", "lseg comment", "YES", nil, nil, nil, nil, nil, nil, "lseg", "b", 20).
						AddRow("macaddr", "macaddr", "macaddr comment", "YES", nil, nil, nil, nil, nil, nil, "macaddr", "b", 21).
						AddRow("macaddr8", "macaddr8", "macaddr8 comment", "YES", nil, nil, nil, nil, nil, nil, "macaddr8", "b", 22).
						AddRow("money", "money", "money comment", "YES", nil, nil, nil, nil, nil, nil, "money", "b", 23).
						AddRow("numeric", "numeric", "numeric comment", "YES", nil, nil, nil, nil, nil, nil, "numeric", "b", 24).
						AddRow("path", "path", "path comment", "YES", nil, nil, nil, nil, nil, nil, "path", "b", 25).
						AddRow("pg_lsn", "pg_lsn", "pg_lsn comment", "YES", nil, nil, nil, nil, nil, nil, "pg_lsn", "b", 26).
						AddRow("pg_snapshot", "pg_snapshot", "pg_snapshot comment", "YES", nil, nil, nil, nil, nil, nil, "pg_snapshot", "b", 27).
						AddRow("point", "point", "point comment", "YES", nil, nil, nil, nil, nil, nil, "point", "b", 28).
						AddRow("polygon", "polygon", "polygon comment", "YES", nil, nil, nil, nil, nil, nil, "polygon", "b", 29).
						AddRow("real", "real", "real comment", "YES", nil, nil, nil, nil, nil, nil, "float4", "b", 30).
						AddRow("smallint", "smallint", "smallint comment", "YES", nil, nil, nil, 16, 0, nil, "int2", "b", 31).
						AddRow("smallserial", "smallserial", "smallserial comment", "YES", nil, nil, nil, 16, 0, nil, "serial2", "b", 32).
						AddRow("serial", "serial", "serial comment", "YES", nil, nil, nil, 32, 0, nil, "serial4", "b", 33).
						AddRow("text", "text", "text comment", "YES", nil, nil, nil, nil, nil, nil, "text", "b", 34).
						AddRow("time", "time without time zone", "time comment", "YES", nil, nil, nil, nil, nil, nil, "time", "b", 35).
						AddRow("time without time zone", "time without time zone", "time without time zone comment", "YES", nil, nil, nil, nil, nil, nil, "time", "b", 36).
						AddRow("time with time zone", "time with time zone", "time with time zone comment", "YES", nil, nil, nil, nil, nil, nil, "time", "b", 37).
						AddRow("timestamp", "timestamp without time zone", "timestamp comment", "YES", nil, nil, nil, nil, nil, nil, "timestamp", "b", 38).
						AddRow("timestamp without time zone", "timestamp without time zone", "timestamp without time zone comment", "YES", nil, nil, nil, nil, nil, nil, "timestamp", "b", 39).
						AddRow("timestamp with time zone", "timestamp with time zone", "timestamp with time zone comment", "YES", nil, nil, nil, nil, nil, nil, "timestamptz", "b", 40).
						AddRow("tsquery", "tsquery", "tsquery comment", "YES", nil, nil, nil, nil, nil, nil, "tsquery", "b", 41).
						AddRow("tsvector", "tsvector", "tsvector comment", "YES", nil, nil, nil, nil, nil, nil, "tsvector", "b", 42).
						AddRow("uuid", "uuid", "uuid comment", "YES", nil, nil, nil, nil, nil, nil, "uuid", "b", 43).
						AddRow("xml", "xml", "xml comment", "YES", nil, nil, nil, nil, nil, nil, "xml", "b", 44).
						AddRow("user-defined", "user-defined", "user-defined comment", "YES", nil, nil, nil, nil, nil, nil, "ltree", "b", 45).
						AddRow("enum", "user-defined", "enum comment", "YES", nil, nil, nil, nil, nil, nil, "enum", "e", 46))

				mock.ExpectQuery(Escape(EnumQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"enumlabel"}).
						AddRow("a").AddRow("b").AddRow("c"))

				mock.noIndexes()
				mock.noForeignKeys()

			},
			expected: func() *schema.Schema {
				s := &schema.Schema{
					Name: schemaName,
				}
				tables := []*schema.Table{
					{
						Name:    tableName,
						Comment: "Comment",
						Schema:  s,
					},
				}
				columns := []*schema.Column{
					{Name: "bigint", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Comment: "bigint comment", Nullable: false, Table: tables[0]},
					{Name: "bigserial", Type: &spec.IntegerType{Name: "bigserial", Size: 64}, Comment: "bigserial comment", Nullable: false, Table: tables[0]},
					{Name: "bit", Type: &spec.BitType{Name: "bit", Len: 1}, Comment: "bit comment", Nullable: true, Table: tables[0]},
					{Name: "bit varying", Type: &spec.BitType{Name: "bit varying", Len: 255}, Comment: "bit varying comment", Nullable: true, Table: tables[0]},
					{Name: "boolean", Type: &spec.BoolType{Name: "boolean"}, Comment: "boolean comment", Nullable: true, Table: tables[0]},
					{Name: "box", Type: &spec.SpatialType{Name: "box"}, Comment: "box comment", Nullable: true, Table: tables[0]},
					{Name: "bytea", Type: &spec.BinaryType{Name: "bytea", Size: 1}, Comment: "bytea comment", Nullable: true, Table: tables[0]},
					{Name: "character", Type: &spec.StringType{Name: "character", Size: 1}, Comment: "character comment", Nullable: true, Table: tables[0]},
					{Name: "character varying", Type: &spec.StringType{Name: "character varying", Size: 255}, Comment: "character varying comment", Nullable: true, Table: tables[0]},
					{Name: "cidr", Type: &spec.SpatialType{Name: "cidr"}, Comment: "cidr comment", Nullable: true, Table: tables[0]},
					{Name: "circle", Type: &spec.SpatialType{Name: "circle"}, Comment: "circle comment", Nullable: true, Table: tables[0]},
					{Name: "date", Type: &spec.TimeType{Name: "date"}, Comment: "date comment", Nullable: true, Table: tables[0]},
					{Name: "double precision", Type: &spec.FloatType{Name: "double precision"}, Comment: "double precision comment", Nullable: true, Table: tables[0]},
					{Name: "inet", Type: &spec.SpatialType{Name: "inet"}, Comment: "inet comment", Nullable: true, Table: tables[0]},
					{Name: "integer", Type: &spec.IntegerType{Name: "integer", Size: 32}, Comment: "integer comment", Nullable: true, Table: tables[0]},
					{Name: "interval", Type: &spec.TimeType{Name: "interval"}, Comment: "interval comment", Nullable: true, Table: tables[0]},
					{Name: "json", Type: &spec.JSONType{Name: "json"}, Comment: "json comment", Nullable: true, Table: tables[0]},
					{Name: "jsonb", Type: &spec.JSONType{Name: "jsonb"}, Comment: "jsonb comment", Nullable: true, Table: tables[0]},
					{Name: "line", Type: &spec.SpatialType{Name: "line"}, Comment: "line comment", Nullable: true, Table: tables[0]},
					{Name: "lseg", Type: &spec.SpatialType{Name: "lseg"}, Comment: "lseg comment", Nullable: true, Table: tables[0]},
					{Name: "macaddr", Type: &spec.SpatialType{Name: "macaddr"}, Comment: "macaddr comment", Nullable: true, Table: tables[0]},
					{Name: "macaddr8", Type: &spec.SpatialType{Name: "macaddr8"}, Comment: "macaddr8 comment", Nullable: true, Table: tables[0]},
					{Name: "money", Type: &spec.FloatType{Name: "money"}, Comment: "money comment", Nullable: true, Table: tables[0]},
					{Name: "numeric", Type: &spec.FloatType{Name: "numeric"}, Comment: "numeric comment", Nullable: true, Table: tables[0]},
					{Name: "path", Type: &spec.SpatialType{Name: "path"}, Comment: "path comment", Nullable: true, Table: tables[0]},
					{Name: "pg_lsn", Type: &spec.SpatialType{Name: "pg_lsn"}, Comment: "pg_lsn comment", Nullable: true, Table: tables[0]},
					{Name: "pg_snapshot", Type: &spec.SpatialType{Name: "pg_snapshot"}, Comment: "pg_snapshot comment", Nullable: true, Table: tables[0]},
					{Name: "point", Type: &spec.SpatialType{Name: "point"}, Comment: "point comment", Nullable: true, Table: tables[0]},
					{Name: "polygon", Type: &spec.SpatialType{Name: "polygon"}, Comment: "polygon comment", Nullable: true, Table: tables[0]},
					{Name: "real", Type: &spec.FloatType{Name: "real"}, Comment: "real comment", Nullable: true, Table: tables[0]},
					{Name: "smallint", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Comment: "smallint comment", Nullable: true, Table: tables[0]},
					{Name: "smallserial", Type: &spec.IntegerType{Name: "smallserial", Size: 16}, Comment: "smallserial comment", Nullable: true, Table: tables[0]},
					{Name: "serial", Type: &spec.IntegerType{Name: "serial", Size: 32}, Comment: "serial comment", Nullable: true, Table: tables[0]},
					{Name: "text", Type: &spec.StringType{Name: "text"}, Comment: "text comment", Nullable: true, Table: tables[0]},
					{Name: "time", Type: &spec.TimeType{Name: "time without time zone"}, Comment: "time comment", Nullable: true, Table: tables[0]},
					{Name: "time without time zone", Type: &spec.TimeType{Name: "time without time zone"}, Comment: "time without time zone comment", Nullable: true, Table: tables[0]},
					{Name: "time with time zone", Type: &spec.TimeType{Name: "time with time zone"}, Comment: "time with time zone comment", Nullable: true, Table: tables[0]},
					{Name: "timestamp", Type: &spec.TimeType{Name: "timestamp without time zone"}, Comment: "timestamp comment", Nullable: true, Table: tables[0]},
					{Name: "timestamp without time zone", Type: &spec.TimeType{Name: "timestamp without time zone"}, Comment: "timestamp without time zone comment", Nullable: true, Table: tables[0]},
					{Name: "timestamp with time zone", Type: &spec.TimeType{Name: "timestamp with time zone"}, Comment: "timestamp with time zone comment", Nullable: true, Table: tables[0]},
					{Name: "tsquery", Type: &spec.StringType{Name: "tsquery"}, Comment: "tsquery comment", Nullable: true, Table: tables[0]},
					{Name: "tsvector", Type: &spec.StringType{Name: "tsvector"}, Comment: "tsvector comment", Nullable: true, Table: tables[0]},
					{Name: "uuid", Type: &spec.StringType{Name: "uuid"}, Comment: "uuid comment", Nullable: true, Table: tables[0]},
					{Name: "xml", Type: &spec.StringType{Name: "xml"}, Comment: "xml comment", Nullable: true, Table: tables[0]},
					{Name: "user-defined", Type: &spec.SpatialType{Name: "user-defined"}, Comment: "user-defined comment", Nullable: true, Table: tables[0]},
					{Name: "enum", Type: &spec.EnumType{Name: "enum", Values: []string{"a", "b", "c"}}, Comment: "enum comment", Nullable: true, Table: tables[0]},
				}

				tables[0].Schema = s
				tables[0].Columns = columns
				s.Tables = tables
				return s

			},
		},
		{
			name: "indexes",
			before: func(mock postgresMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME", "COMMENT"}).
						AddRow(tableName, "Comment"))

				mock.noColumns()

				mock.ExpectQuery(Escape(IndexesQuery)).
					WithArgs(tableName).
					WillReturnRows(sqlmock.NewRows(IndexesQueryFields).
						AddRow("bigint", "btree", "bigint", 1, 0, "p", nil, nil, 0, 0, 0, 0, "comment").
						AddRow("character_idx", "btree", "character", 0, 0, "u", nil, nil, 0, 0, 0, 0, nil).
						AddRow("subpart", "btree", "character1", 0, 0, nil, nil, nil, 0, 0, 0, 0, nil).
						AddRow("non_unique", "btree", "char", 0, 0, nil, nil, nil, 0, 0, 0, 0, nil).
						AddRow("unique", "btree", "char", 0, 0, nil, nil, nil, 0, 0, 0, 0, nil).
						AddRow("unique_union", "btree", "char", 0, 0, nil, nil, nil, 0, 0, 0, 0, nil).
						AddRow("unique_union", "btree", "character1", 0, 0, nil, nil, nil, 0, 0, 0, 0, nil))

				mock.noForeignKeys()

			},
			expected: func() *schema.Schema {
				s := &schema.Schema{
					Name: schemaName,
				}

				table := &schema.Table{
					Name:    tableName,
					Comment: "Comment",
					Schema:  s,
				}

				tables := []*schema.Table{
					table,
				}

				indexes := []*schema.Index{
					{
						Name: "bigint",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "bigint",
							},
						},
						Type:    "btree",
						Unique:  false,
						Primary: true,
						Comment: "comment",
					},
					{
						Name: "character_idx",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "character",
								Sub:    0,
								Expr:   nil,
							},
						},
						Type: "btree",
					},
					{
						Name: "subpart",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "character1",
								Sub:    0,
								Expr:   nil,
							},
						},
						Type: "btree",
					},
					{
						Name: "non_unique",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "char",
								Sub:    0,
								Expr:   nil,
							},
						},
						Type:   "btree",
						Unique: false,
					},
					{
						Name: "unique",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "char",
								Sub:    0,
								Expr:   nil,
							},
						},
						Type: "btree",
					},
					{
						Name: "unique_union",
						IndexColumns: []*schema.IndexColumn{
							{
								SeqNo:  1,
								Column: "char",
								Sub:    0,
								Expr:   nil,
							},
							{
								SeqNo:  2,
								Column: "character1",
								Sub:    0,
								Expr:   nil,
							},
						},
						Type: "btree",
					},
				}

				table.Indexes = indexes
				s.Tables = tables

				return s
			},
		},
		{
			name: "foreign keys",
			before: func(mock postgresMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME", "COMMENT"}).
						AddRow(tableName, "Comment").
						AddRow(fkTableName, "Comment"))

				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs(tableName).
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("id", "bigint", "id comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1).
						AddRow("gid", "bigint", "gid comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1).
						AddRow("cid", "bigint", "cid comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1).
						AddRow("uid", "bigint", "uid comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1))

				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs(fkTableName).
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("id", "bigint", "id comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1).
						AddRow("cid", "bigint", "cid comment", "NO", nil, nil, nil, 64, 0, nil, "int8", "b", 1))

				mock.noIndexes()
				mock.noIndexes()

				mock.ExpectQuery(Escape(ForeignKeysQuery)).
					WithArgs(tableName).
					WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields).
						AddRow("multi_column", tableName, "gid", schemaName, fkTableName, "id", schemaName, "NO ACTION", "CASCADE").
						AddRow("multi_column", tableName, "cid", schemaName, fkTableName, "cid", schemaName, "NO ACTION", "CASCADE").
						AddRow("self_reference", tableName, "uid", schemaName, tableName, "id", schemaName, "NO ACTION", "CASCADE"))

				mock.ExpectQuery(Escape(ForeignKeysQuery)).
					WithArgs(fkTableName).
					WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields))

			},
			expected: func() *schema.Schema {
				s := &schema.Schema{
					Name: schemaName,
				}

				table := &schema.Table{
					Name:    tableName,
					Comment: "Comment",
					Schema:  s,
				}

				tableFK := &schema.Table{
					Name:    fkTableName,
					Comment: "Comment",
					Schema:  s,
				}

				columns := []*schema.Column{
					{
						Name:    "id",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "id comment",
						Table:   table,
					},
					{
						Name:    "gid",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "gid comment",
						Charset: "",
						Table:   table,
					},
					{
						Name:    "cid",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "cid comment",
						Table:   table,
					},
					{
						Name:    "uid",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "uid comment",
						Table:   table,
					},
				}

				refColumns := []*schema.Column{
					{
						Name:    "id",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "id comment",
						Table:   tableFK,
					},
					{
						Name:    "cid",
						Type:    &spec.IntegerType{Name: "bigint", Size: 64},
						Comment: "cid comment",
						Table:   tableFK,
					},
				}
				tables := []*schema.Table{
					table,
					tableFK,
				}

				fks := []*schema.ForeignKey{
					{
						Name:       "multi_column",
						Table:      tables[0],
						Columns:    columns[1:3],
						RefTable:   tables[1],
						RefColumns: refColumns[0:2],
						OnUpdate:   schema.ReferenceOption(schema.NoAction),
						OnDelete:   schema.ReferenceOption(schema.Cascade),
					},
					{
						Name:       "self_reference",
						Table:      tables[0],
						Columns:    columns[3:4],
						RefTable:   tables[0],
						RefColumns: columns[0:1],
						OnUpdate:   schema.ReferenceOption(schema.NoAction),
						OnDelete:   schema.ReferenceOption(schema.Cascade),
					},
				}

				table.Columns = columns
				table.ForeignKeys = fks
				tableFK.Columns = refColumns
				s.Tables = tables

				return s
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			test.before(postgresMock{mock})
			l := &inspect{
				Driver: schema.OpenDB(cre.MySQL, db),
			}

			require.NoError(t, err)
			s, err := l.Inspect(context.Background(), schemaName)
			require.Equal(t, test.wantErr, err != nil, err)
			require.EqualValues(t, test.expected(), s)
		})

	}
}

func Escape(query string) string {
	rows := strings.Split(query, "\n")
	for i := range rows {
		rows[i] = strings.TrimPrefix(rows[i], " ")
	}
	query = strings.Join(rows, " ")
	return strings.TrimSpace(regexp.QuoteMeta(query)) + "$"
}
