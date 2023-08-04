package mysql

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

type mysqlMock struct {
	sqlmock.Sqlmock
}

func (m mysqlMock) info() {
	m.ExpectQuery(Escape(VariablesQuery)).
		WillReturnRows(sqlmock.NewRows(VariablesQueryFields).
			AddRow("8.0.20", " utf8mb4_0900_ai_ci", "utf8mb4"))

}

func TestTables(t *testing.T) {

	tests := []struct {
		name     string
		before   func(mysqlMock)
		expected func() *schema.Schema
		wantErr  bool
	}{
		{
			name: "empty schema name",
			before: func(mock mysqlMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WithArgs("").
					WillReturnRows(sqlmock.NewRows(TablesQueryFields))
			},
			expected: func() *schema.Schema {
				var (
					s1 = &schema.Schema{
						Name: "",
					}
				)

				return s1
			},
		},
		{
			name: "custom schema",
			before: func(mock mysqlMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WithArgs("test").
					WillReturnRows(sqlmock.NewRows(TablesQueryFields).
						AddRow("table", "utf8mb4", "utf8mb4_0900_ai_ci", nil, "Comment", "COMPRESSION=ZLIB"))

				//column_name, column_type, column_comment, is_nullable, column_key, column_default, extra, character_set_name, collation_name, numeric_precision, numeric_scale
				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("bigint", "bigint(20)", "中文bigint comment", "NO", "PRI", nil, "auto_increment", "", "", nil, nil).
						AddRow("varchar", "varchar(255)", "varchar comment", "YES", "YES", nil, "", "", "", nil, nil).
						AddRow("varchar1", "varchar(255) character set utf8mb4_bin collate utf8mb4_bin", "varchar1 comment", "YES", "YES", nil, "", "", "", nil, nil).
						AddRow("longtext", "longtext", "longtext comment", "YES", "YES", nil, "", "", "", nil, nil).
						AddRow("char", "char(36)", "char comment", "YES", "YES", nil, "", "", "utf8mb4_bin", nil, nil).
						AddRow("decimal", "decimal(6, 4)", "decimal comment", "NO", "YES", nil, "", "", "", "6", "4").
						AddRow("datetime", "datetime(5)", "datetime comment", "NO", "YES", nil, "", "", "", nil, nil).
						AddRow("point", "point", "point comment", "NO", "YES", nil, "", "", "", nil, nil).
						AddRow("json", "json", "json comment", "NO", "YES", nil, "", "", "", nil, nil).
						AddRow("enum", "enum('a','b','c')", "enum comment", "NO", "YES", nil, "", "", "", nil, nil))

				//INDEX_NAME, COLUMN_NAME, NON_UNIQUE, SEQ_IN_INDEX, INDEX_TYPE, COLLATION, INDEX_COMMENT, SUB_PART, EXPRESSION
				f := append(IndexesQueryFields)
				mock.ExpectQuery(Escape(IndexesExprQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(f).
						AddRow("PRIMARY", "bigint", "0", "1", "BTREE", "", "", nil, nil).
						AddRow("varchar_idx", "varchar", "0", "1", "BTREE", "", "", nil, nil).
						AddRow("subpart", "varchar1", "0", "1", "BTREE", "", "", 64, nil).
						AddRow("non_unique", "char", "1", "1", "BTREE", "", "", nil, nil).
						AddRow("unique", "char", "0", "1", "BTREE", "", "", nil, nil).
						AddRow("unique_union", "char", "0", "1", "BTREE", "", "", nil, nil).
						AddRow("unique_union", "varchar1", "0", "2", "BTREE", "", "", nil, nil))

				//"CONSTRAINT_NAME", "TABLE_NAME", "COLUMN_NAME", "TABLE_SCHEMA", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME", "REFERENCED_TABLE_SCHEMA", "UPDATE_RULE", "DELETE_RULE"
				mock.ExpectQuery(Escape(ForeignKeysQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields))

			},
			expected: func() *schema.Schema {
				s := &schema.Schema{
					Name: "test",
				}

				table := &schema.Table{
					Name:      "table",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Comment:   "Comment",
					Options:   "COMPRESSION=ZLIB",
					Schema:    s,
					Indexes: []*schema.Index{
						{
							Name: "PRIMARY",
							IndexColumns: []*schema.IndexColumn{
								{
									SeqNo:  1,
									Column: "bigint",
								},
							},
							Type:    "BTREE",
							Unique:  true,
							Primary: true,
						},
						{
							Name: "varchar_idx",
							IndexColumns: []*schema.IndexColumn{
								{
									SeqNo:  1,
									Column: "varchar",
									Sub:    0,
									Expr:   nil,
								},
							},
							Type:   "BTREE",
							Unique: true,
						},
						{
							Name: "subpart",
							IndexColumns: []*schema.IndexColumn{
								{
									SeqNo:  1,
									Column: "varchar1",
									Sub:    64,
									Expr:   nil,
								},
							},
							Type:   "BTREE",
							Unique: true,
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
							Type:   "BTREE",
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
							Type:   "BTREE",
							Unique: true,
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
									Column: "varchar1",
									Sub:    0,
									Expr:   nil,
								},
							},
							Type:   "BTREE",
							Unique: true,
						},
					},
				}

				columns := []*schema.Column{
					{
						Name: "bigint",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment:       "中文bigint comment",
						Primary:       true,
						AutoIncrement: true,
						Table:         table,
					},
					{
						Name: "varchar",
						Type: &spec.StringType{
							Name: "varchar",
							Size: 255,
						},
						Nullable: true,
						Comment:  "varchar comment",
						Table:    table,
					},
					{
						Name: "varchar1",
						Type: &spec.StringType{
							Name:      "varchar",
							Size:      255,
							Charset:   "utf8mb4_bin",
							Collation: "utf8mb4_bin",
						},
						Nullable: true,
						Comment:  "varchar1 comment",
						Table:    table,
					},
					{
						Name: "longtext",
						Type: &spec.StringType{
							Name: "longtext",
						},
						Nullable: true,
						Comment:  "longtext comment",
						Table:    table,
					},
					{
						Name: "char",
						Type: &spec.StringType{
							Name: "char",
							Size: 36,
						},
						Nullable:  true,
						Comment:   "char comment",
						Collation: "utf8mb4_bin",
						Table:     table,
					},
					{
						Name: "decimal",
						Type: &spec.FloatType{
							Name:      "decimal",
							Precision: 6,
							Scale:     4,
						},
						Nullable:  false,
						Precision: 6,
						Scale:     4,
						Comment:   "decimal comment",
						Table:     table,
					},
					{
						Name: "datetime",
						Type: &spec.TimeType{
							Name: "datetime",
						},
						Nullable: false,
						Comment:  "datetime comment",
						Table:    table,
					},
					{
						Name: "point",
						Type: &spec.SpatialType{
							Name: "point",
						},
						Nullable: false,
						Comment:  "point comment",
						Table:    table,
					},
					{
						Name: "json",
						Type: &spec.JSONType{
							Name: "json",
						},
						Nullable: false,
						Comment:  "json comment",
						Table:    table,
					},
					{
						Name: "enum",
						Type: &spec.EnumType{
							Name:   "enum",
							Values: []string{"a", "b", "c"},
						},
						Comment: "enum comment",
						Table:   table,
					},
				}

				table.Columns = columns

				tables := []*schema.Table{
					table,
				}

				s.Tables = tables

				return s
			},
		},
		{
			name: "foreign keys",
			before: func(mock mysqlMock) {
				mock.info()
				mock.ExpectQuery(Escape(TablesQuery)).
					WithArgs("test").
					WillReturnRows(sqlmock.NewRows(TablesQueryFields).
						AddRow("table", "utf8mb4", "utf8mb4_0900_ai_ci", nil, "Comment", "COMPRESSION=ZLIB").
						AddRow("fk", "utf8mb4", "utf8mb4_0900_ai_ci", nil, "Comment", "COMPRESSION=ZLIB"))

				//COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT, IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA, CHARACTER_SET_NAME, COLLATION_NAME, NUMERIC_PRECISION, NUMERIC_SCALE
				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("id", "bigint(20)", "id comment", "NO", "PRI", nil, "auto_increment", "", "", nil, nil).
						AddRow("gid", "bigint(20)", "gid comment", "NO", "MUL", nil, "", "", "", nil, nil).
						AddRow("cid", "bigint(20)", "cid comment", "NO", "MUL", nil, "", "", "", nil, nil).
						AddRow("uid", "bigint(20)", "uid comment", "NO", "MUL", nil, "", "", "", nil, nil))
				mock.ExpectQuery(Escape(ColumnsQuery)).
					WithArgs("test", "fk").
					WillReturnRows(sqlmock.NewRows(ColumnsQueryFields).
						AddRow("id", "bigint(20)", "id comment", "NO", "PRI", nil, "auto_increment", "", "", nil, nil).
						AddRow("cid", "bigint(20)", "cid comment", "NO", "MUL", nil, "", "", "", nil, nil))

				//INDEX_NAME, COLUMN_NAME, NON_UNIQUE, SEQ_IN_INDEX, INDEX_TYPE, COLLATION, INDEX_COMMENT, SUB_PART, EXPRESSION
				f := append(IndexesQueryFields, "`EXPRESSION`")
				mock.ExpectQuery(Escape(IndexesExprQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(f))
				mock.ExpectQuery(Escape(IndexesExprQuery)).
					WithArgs("test", "fk").
					WillReturnRows(sqlmock.NewRows(f))

				//"CONSTRAINT_NAME", "TABLE_NAME", "COLUMN_NAME", "TABLE_SCHEMA", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME", "REFERENCED_TABLE_SCHEMA", "UPDATE_RULE", "DELETE_RULE"
				mock.ExpectQuery(Escape(ForeignKeysQuery)).
					WithArgs("test", "table").
					WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields).
						AddRow("multi_column", "table", "gid", "test", "fk", "id", "test", "NO ACTION", "CASCADE").
						AddRow("multi_column", "table", "cid", "test", "fk", "cid", "test", "NO ACTION", "CASCADE").
						AddRow("self_reference", "table", "uid", "test", "table", "id", "test", "NO ACTION", "CASCADE"))

				mock.ExpectQuery(Escape(ForeignKeysQuery)).
					WithArgs("test", "fk").
					WillReturnRows(sqlmock.NewRows(ForeignKeysQueryFields))

			},
			expected: func() *schema.Schema {
				s := &schema.Schema{
					Name: "test",
				}

				table := &schema.Table{
					Name:      "table",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Comment:   "Comment",
					Options:   "COMPRESSION=ZLIB",
					Schema:    s,
				}

				tableFK := &schema.Table{
					Name:      "fk",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Comment:   "Comment",
					Options:   "COMPRESSION=ZLIB",
					Schema:    s,
				}

				columns := []*schema.Column{
					{
						Name: "id",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment:       "id comment",
						Primary:       true,
						AutoIncrement: true,
						Table:         table,
					},
					{
						Name: "gid",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment: "gid comment",
						Charset: "",
						Table:   table,
					},
					{
						Name: "cid",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment: "cid comment",
						Table:   table,
					},
					{
						Name: "uid",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment: "uid comment",
						Table:   table,
					},
				}

				refColumns := []*schema.Column{
					{
						Name: "id",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
						Comment:       "id comment",
						Primary:       true,
						AutoIncrement: true,
						Table:         tableFK,
					},
					{
						Name: "cid",
						Type: &spec.IntegerType{
							Name: "bigint",
							Size: 64,
						},
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
			test.before(mysqlMock{mock})
			l := &inspect{
				Driver: schema.OpenDB(cre.MySQL, db),
			}

			require.NoError(t, err)
			resource, err := l.Inspect(context.Background(), test.expected().Name)
			require.Equal(t, test.wantErr, err != nil, err)
			require.EqualValues(t, test.expected(), resource)
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
