package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ychengcloud/cre/spec"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		before   func() *Schema
		expected func() *spec.Schema
	}{
		{
			name: "convert",
			before: func() *Schema {
				tables := []*Table{
					{
						Name: "t1",
					},
				}
				columns := []*Column{
					{
						Name:     "c1",
						Type:     &spec.IntegerType{},
						Nullable: true,
						Comment:  "c1 comment",
						Primary:  true,
						Unique:   true,
						Table:    tables[0],
					},
					{
						Name:     "c2",
						Type:     &spec.StringType{},
						Nullable: true,
						Comment:  "c2 comment",
						Table:    tables[0],
					},
					{
						Name:     "c3",
						Type:     &spec.JSONType{},
						Nullable: true,
						Comment:  "c3 comment",
						Table:    tables[0],
					},
					{
						Name:    "c4",
						Type:    &spec.StringType{},
						Comment: "c4 comment",
						Table:   tables[0],
					},
				}

				tables[0].Columns = columns

				s := &Schema{
					Tables: tables,
				}

				return s
			},

			expected: func() *spec.Schema {

				tables := []*spec.Table{
					{
						Name: "t1",
					},
				}
				tables[0].AddFields(
					spec.Builder("c1").
						Type(&spec.IntegerType{}).
						Nullable(true).
						Comment("c1 comment").
						PrimaryKey(true).
						Unique(true).
						Filterable(true).
						Sortable(true).
						Build(),
					spec.Builder("c2").
						Type(&spec.StringType{}).
						Nullable(true).
						Comment("c2 comment").
						Build(),
					spec.Builder("c3").
						Type(&spec.JSONType{}).
						Nullable(true).
						Comment("c3 comment").
						Build(),
					spec.Builder("c4").
						Type(&spec.StringType{}).
						Comment("c4 comment").
						Build(),
				)

				tables[0].ID = tables[0].Fields()[0]

				schema := &spec.Schema{}
				schema.AddTables(tables...)

				return schema
			},
		},
		{
			name: "convert indexes",
			before: func() *Schema {
				tables := []*Table{
					{
						Name: "t1",
					},
				}
				columns := []*Column{
					{
						Name:     "c1",
						Type:     &spec.IntegerType{},
						Nullable: true,
						Comment:  "c1 comment",
						Table:    tables[0],
					},
					{
						Name:     "c2",
						Type:     &spec.IntegerType{},
						Primary:  true,
						Nullable: true,
						Comment:  "c2 comment",
						Table:    tables[0],
					},
				}

				indexes := []*Index{
					{
						Name:   "idx1",
						Unique: true,
						IndexColumns: []*IndexColumn{
							{
								Column: columns[0].Name,
							},
						},
					},
				}

				tables[0].Columns = columns
				tables[0].Indexes = indexes

				s := &Schema{
					Tables: tables,
				}

				return s
			},

			expected: func() *spec.Schema {
				tables := []*spec.Table{
					{
						Name: "t1",
					},
				}

				tables[0].AddFields(
					spec.Builder("c1").
						Type(&spec.IntegerType{}).
						Nullable(true).
						Comment("c1 comment").
						Unique(true).
						Index(true).
						Filterable(true).
						Build(),
					spec.Builder("c2").
						Type(&spec.IntegerType{}).
						Nullable(true).
						Comment("c2 comment").
						PrimaryKey(true).
						Filterable(true).
						Sortable(true).
						Build(),
				)

				schema := &spec.Schema{}
				schema.AddTables(tables...)
				return schema
			},
		},
		{
			name: "convert fks",
			before: func() *Schema {
				tables := []*Table{
					{
						Name: "t1",
					},
					{
						Name: "t2",
					},
				}
				columns1 := []*Column{
					{
						Name:  "c1",
						Type:  &spec.IntegerType{},
						Table: tables[0],
					},
				}
				columns2 := []*Column{
					{
						Name:  "c2",
						Type:  &spec.IntegerType{},
						Table: tables[1],
					},
				}
				tables[0].Columns = columns1
				tables[1].Columns = columns2

				tables[0].ForeignKeys = []*ForeignKey{
					{
						Name:  "fk1",
						Table: tables[0],
						Columns: []*Column{
							columns1[0],
						},
						RefTable: tables[1],
						RefColumns: []*Column{
							columns2[0],
						},
					},
				}

				s := &Schema{
					Tables: tables,
				}

				return s
			},

			expected: func() *spec.Schema {
				tables := []*spec.Table{
					{
						Name: "t1",
					},
					{
						Name: "t2",
					},
				}

				tables[0].AddFields(spec.Builder("c1").Type(&spec.IntegerType{}).ForeignKey(true).Build())
				tables[1].AddFields(spec.Builder("c2").Type(&spec.IntegerType{}).Build())

				schema := &spec.Schema{}
				schema.AddTables(tables...)
				return schema
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.before().Convert()
			require.NoError(t, err)
			require.Equal(t, test.expected(), actual)
		})
	}
}
