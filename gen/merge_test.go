package gen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ychengcloud/cre/spec"
)

func TestMergeSchema(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
	}

	categoryTable := &spec.Table{Name: "category"}
	postTable := &spec.Table{Name: "post"}
	skipTable := &spec.Table{Name: "skip"}
	postCategoryTable := &spec.Table{Name: "post_category"}
	userTable := &spec.Table{Name: "user"}
	remoteTable := &spec.Table{Name: "remote"}
	categoryTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
			Unique:     true,
		},
		{
			Name: "name",
			Type: &spec.StringType{Name: "char", Size: 64},
		},
	}...)

	categoryTable.ID = categoryTable.GetField("id")

	postTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
			Unique:     true,
		},
		{
			Name: "name",
			Type: &spec.StringType{Name: "char", Size: 64},
		},
		{Name: "skip"},
	}...)

	postCategoryTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
			Unique:     true,
		},
		{
			Name: "post_id",
			Type: &spec.IntegerType{Name: "int", Size: 32},
		},
		{
			Name: "category_id",
			Type: &spec.IntegerType{Name: "int", Size: 32},
		},
	}...)

	userTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
			Unique:     true,
		},
		{
			Name: "name",
			Type: &spec.StringType{Name: "char", Size: 64},
		},
	}...)

	remoteTable.AddFields([]*spec.Field{
		{
			Name: "id",
		},
	}...)

	expectedPostRelField := &spec.Field{
		Name: "categories",
		Type: &spec.ObjectType{Name: "categories"},
		Rel: &spec.Relation{
			Type:     spec.RelTypeManyToMany,
			Field:    postTable.Fields()[0],
			RefTable: categoryTable,
			RefField: categoryTable.Fields()[0],
			JoinTable: &spec.JoinTable{
				Name:         "post_category",
				JoinField:    postCategoryTable.Fields()[1],
				JoinRefField: postCategoryTable.Fields()[2],
			},
			Inverse: true,
		},
		Table: postTable,
	}
	expectedPostBelongsToField := &spec.Field{
		Name: "belongsto",
		Type: &spec.ObjectType{Name: "belongsto"},
		Rel: &spec.Relation{
			Type:     spec.RelTypeBelongsTo,
			Field:    postTable.Fields()[0],
			RefTable: userTable,
			RefField: userTable.Fields()[0],
		},
		Table: postTable,
	}

	expectedRemoteField := &spec.Field{
		Name:   "remote",
		Remote: true,
		Type:   &spec.ObjectType{Name: "remote"},
		Rel: &spec.Relation{
			Type:     spec.RelTypeBelongsTo,
			Field:    postTable.Fields()[0],
			RefTable: remoteTable,
			RefField: remoteTable.Fields()[0],
		},
		Table: postTable,
	}
	schema.AddTables(categoryTable, postTable, skipTable, postCategoryTable, userTable)

	tablesInCfg := []*Table{
		{
			Name: "category",
			Fields: []*Field{
				{
					Name: "id",
				},
			},
		},
		{
			Name: "post",
			Fields: []*Field{
				{
					Name: "id",
				},
				{
					Name:       "name",
					Nullable:   true,
					Optional:   true,
					Comment:    "post name",
					Alias:      "postAlias",
					Sortable:   true,
					Filterable: true,
					Operations: []string{"EQ", "In"},
				},
				{Name: "skip", Skip: true},
				{Name: "pid", Relation: &Relation{Type: "hasOne", RefTable: "post", RefField: "id"}},
				{Name: "cid", Relation: &Relation{Type: "hasOne", RefTable: "category", RefField: "id", Inverse: true}},
				{
					Name: "categories",
					Relation: &Relation{
						Type:     "manyToMany",
						RefTable: "category",
						RefField: "id",
						JoinTable: &JoinTable{
							Name:     "post_category",
							Table:    "post",
							RefTable: "category",
							Field:    "post_id",
							RefField: "category_id",
						},
						Inverse: true,
					},
				},
				{
					Name: "belongsto",
					Relation: &Relation{
						Type:     "BelongsTo",
						RefTable: "user",
					},
				},
				{
					Name:   "remote",
					Remote: true,
					Relation: &Relation{
						Type:     "BelongsTo",
						RefTable: "remote",
						RefField: "id",
					},
				},
			},
		},
		{
			Name: "skip",
			Skip: true,
		},
		{Name: "user"},
	}

	s, err := mergeSchema(schema, tablesInCfg)
	r := require.New(t)
	r.NoError(err)
	r.NotNil(s)
	// Check skip condition
	r.Equal(3, len(s.Tables()))
	r.Equal(1, len(s.JoinTables()))
	r.Equal("category", s.Tables()[0].Name)
	r.Equal("post", s.Tables()[1].Name)

	category := s.Tables()[0]
	r.EqualValues(category.ID, category.Fields()[0])

	post := s.Tables()[1]
	postFields := post.Fields()
	// Check skip condition
	r.Equal(len(tablesInCfg[1].Fields)-1, len(postFields))

	r.Equal(true, postFields[1].Nullable)
	r.Equal(true, postFields[1].Optional)
	r.Equal("post name", postFields[1].Comment)
	r.Equal("postAlias", postFields[1].Alias)
	r.Equal(true, postFields[1].Sortable)
	r.Equal(true, postFields[1].Filterable)
	r.Equal(len(spec.NumericOps), len(postFields[0].Ops))
	r.Equal(spec.Eq, postFields[0].Ops[0])
	r.Equal(2, len(postFields[1].Ops))
	r.Equal(spec.Eq, postFields[1].Ops[0])
	r.Equal(spec.In, postFields[1].Ops[1])
	r.Equal("HasOne", postFields[3].Rel.Type.Name())
	r.EqualValues(post, postFields[2].Rel.RefTable)
	r.EqualValues(category, postFields[3].Rel.RefTable)
	r.EqualValues(post.Fields()[0], postFields[2].Rel.RefField)
	r.EqualValues(category.Fields()[0], postFields[3].Rel.RefField)

	// post category relation
	postCategoryField := postFields[4]
	joinTable := s.JoinTables()[0]
	r.NotNil(postCategoryField.Rel)
	r.EqualValues(expectedPostRelField, postCategoryField)
	r.NotNil(postCategoryField.Rel.JoinTable)
	r.EqualValues(category, postCategoryField.Rel.RefTable)
	r.Equal(joinTable.Name, postCategoryField.Rel.JoinTable.Name)
	r.Equal("post_id", postCategoryField.Rel.JoinTable.JoinField.Name)
	r.Equal("category_id", postCategoryField.Rel.JoinTable.JoinRefField.Name)
	r.EqualValues(true, postCategoryField.Rel.Inverse)

	// post belongsto relation
	belongsToField := postFields[5]
	r.NotNil(belongsToField.Rel)
	r.EqualValues(expectedPostBelongsToField, belongsToField)

	// post user relation
	remoteField := postFields[6]
	r.NotNil(remoteField.Rel)
	r.EqualValues(expectedRemoteField, remoteField)
}

func TestGetOps(t *testing.T) {
	id := &spec.Field{
		Name:       "id",
		Type:       &spec.IntegerType{Name: "int", Size: 32},
		PrimaryKey: true,
	}

	name := &spec.Field{
		Name: "name",
		Type: &spec.StringType{Name: "char", Size: 64},
	}

	ops := getOps(id)
	r := require.New(t)
	r.Equal(len(spec.NumericOps), len(ops))
	r.Equal(spec.Eq, ops[0])
	r.Equal(spec.Neq, ops[1])
	r.Equal(spec.In, ops[2])
	r.Equal(spec.NotIn, ops[3])
	r.Equal(spec.Gt, ops[4])
	r.Equal(spec.Gte, ops[5])
	r.Equal(spec.Lt, ops[6])
	r.Equal(spec.Lte, ops[7])

	ops = getOps(name)
	r.Equal(len(spec.StringOps), len(ops))
}
