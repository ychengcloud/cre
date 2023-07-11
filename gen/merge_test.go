package gen

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ychengcloud/cre/spec"
)

var (

	// Post { id, user_id, reluser }
	postIDField = spec.Builder("id").
			Type(&spec.IntegerType{Name: "int", Size: 32}).
			PrimaryKey(true).
			Unique(true).Build()

	postNameField = spec.Builder("name").
			Type(&spec.StringType{Name: "char", Size: 64}).Build()

	postSkipField   = spec.Builder("skip").Build()
	postUserIDField = spec.Builder("user_id").Type(&spec.IntegerType{Name: "int", Size: 32}).Build()

	postRelUserField       = spec.Builder("reluser").Type(&spec.ObjectType{Name: "reluser"}).Build()
	postRelCategoriesField = spec.Builder("relcategories").Type(&spec.ObjectType{Name: "relcategories"}).Build()
	postRelRemoteIDField   = spec.Builder("remote_id").Type(&spec.IntegerType{Name: "int", Size: 32}).Build()
	postRelRemoteField     = spec.Builder("relremote").Type(&spec.ObjectType{Name: "relremote"}).Remote(true).Build()

	postTable = &spec.Table{Name: "post"}

	// Category { id }
	categoryIDField = spec.Builder("id").
			Type(&spec.IntegerType{Name: "int", Size: 32}).
			PrimaryKey(true).
			Unique(true).Build()

	categoryTable = &spec.Table{Name: "category"}

	// PostCategory { id, post_id, category_id }
	postCategoryIDField = spec.Builder("id").
				Type(&spec.IntegerType{Name: "int", Size: 32}).
				PrimaryKey(true).
				Unique(true).Build()

	postCategoryPostIDField     = spec.Builder("post_id").Type(&spec.IntegerType{Name: "int", Size: 32}).Build()
	postCategoryCategoryIDField = spec.Builder("category_id").Type(&spec.IntegerType{Name: "int", Size: 32}).Build()

	postCategoryTable = &spec.Table{Name: "post_category"}

	// User { id }
	userIDField          = spec.Builder("id").Type(&spec.IntegerType{Name: "int", Size: 32}).PrimaryKey(true).Unique(true).Build()
	userRelOnePostField  = spec.Builder("relonepost").Type(&spec.ObjectType{Name: "relonepost"}).Build()
	UserRelManyPostField = spec.Builder("relmanyposts").Type(&spec.ObjectType{Name: "relmanyposts"}).Build()

	userTable = &spec.Table{Name: "user"}

	skipTable = &spec.Table{Name: "skip"}

	// Remote { id }
	remoteIDField = spec.Builder("id").Build()

	remoteTable = &spec.Table{Name: "remote"}
)

func setup() {
	postTable.AddFields(postIDField, postNameField, postSkipField, postUserIDField, postRelUserField, postRelCategoriesField, postRelRemoteIDField, postRelRemoteField)
	categoryTable.AddFields(categoryIDField)
	categoryTable.ID = categoryTable.GetField("id")
	postCategoryTable.AddFields(postCategoryIDField, postCategoryPostIDField, postCategoryCategoryIDField)
	userTable.AddFields(userIDField, userRelOnePostField)

	remoteTable.AddFields(remoteIDField)
	remoteTable.ID = remoteIDField

}
func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestMergeSchema(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
		Attrs: []spec.Attribute{
			&Attr{
				name:  "test",
				value: "value",
			},
		},
	}
	// remote table don't add to schema
	schema.AddTables(categoryTable, postTable, userTable, skipTable)

	trueValue := true

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
					Filterable: &trueValue,
					Operations: []string{"EQ", "In"},
				},
				{Name: "skip", Skip: true},
			},
		},
		{
			Name: "skip",
			Skip: true,
		},
		{Name: "user"},
	}

	cfg := &Config{
		Tables: tablesInCfg,
		Attrs: map[string]any{
			"test": "value",
		},
	}

	s, err := mergeSchema(schema, cfg)
	r := require.New(t)
	r.NoError(err)
	r.NotNil(s)
	fmt.Printf("%#v\n", s.Tables())
	// Check skip condition
	r.Equal(3, len(s.Tables()))
	r.Equal(0, len(s.JoinTables()))
	r.Equal("category", s.Tables()[0].Name)
	r.Equal("post", s.Tables()[1].Name)

	category := s.Tables()[0]
	matchField(t, category.ID, categoryTable.ID)

	post := s.Table("post")
	actualPostIDField := post.GetField("id")
	actualPostNameField := post.GetField("name")
	// Check skip condition
	skipField := post.GetField("skip")
	r.Nil(skipField)

	expectedIDField := spec.Builder("id").
		Type(&spec.IntegerType{Name: "int", Size: 32}).
		PrimaryKey(true).
		Unique(true).
		Ops(spec.NumericOps).Build()

	expectedNameField := spec.Builder("name").
		Type(&spec.StringType{Name: "char", Size: 64}).
		Nullable(true).
		Optional(true).
		Comment("post name").
		Alias("postAlias").
		Sortable(true).
		Filterable(true).
		Ops([]spec.Op{spec.Eq, spec.In}).Build()

	matchField(t, expectedIDField, actualPostIDField)
	matchField(t, expectedNameField, actualPostNameField)
}

func TestRelation(t *testing.T) {

}

func matchField(t *testing.T, expected *spec.Field, actual *spec.Field) {
	r := require.New(t)
	r.NotNil(expected)
	r.NotNil(actual)
	r.Equal(expected.Name, actual.Name)
	r.EqualValues(expected.Type, actual.Type)
	r.Equal(expected.PrimaryKey, actual.PrimaryKey)
	r.Equal(expected.Unique, actual.Unique)
	r.Equal(expected.Nullable, actual.Nullable)
	r.Equal(expected.Optional, actual.Optional)
	r.Equal(expected.Comment, actual.Comment)
	r.Equal(expected.Alias, actual.Alias)
	r.Equal(expected.Sortable, actual.Sortable)
	r.Equal(expected.Filterable, actual.Filterable)
	r.Equal(len(expected.Ops), len(actual.Ops))
	for i := range expected.Ops {
		r.Equal(expected.Ops[i], actual.Ops[i])
	}

}

func matchRelation(t *testing.T, expected *spec.Relation, actual *spec.Relation) {
	r := require.New(t)
	r.NotNil(expected)
	r.NotNil(actual)
	r.EqualValues(expected.Field, actual.Field)
	r.EqualValues(expected.RefTable, actual.RefTable)
	r.EqualValues(expected.RefField, actual.RefField)
	r.EqualValues(expected.JoinTable, actual.JoinTable)
	r.EqualValues(expected.Attrs, actual.Attrs)

}

func TestRelationBelongsTo(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
	}
	// remote table don't add to schema
	schema.AddTables(postTable, userTable)

	r := require.New(t)

	expected := &spec.Relation{
		Type:     spec.RelTypeBelongsTo,
		Field:    postUserIDField,
		RefTable: userTable,
		RefField: userIDField,
	}

	actual, err := mergeRelation(postRelUserField, &Relation{Type: "BelongsTo", RefTable: "user"})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

}

func TestRelationHasOne(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
	}
	schema.AddTables(postTable, userTable)

	r := require.New(t)

	expected := &spec.Relation{
		Type:     spec.RelTypeHasOne,
		Field:    userIDField,
		RefTable: postTable,
		RefField: postUserIDField,
	}

	actual, err := mergeRelation(userRelOnePostField, &Relation{Type: "HasOne", RefTable: "post"})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

}

func TestRelationHasMany(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
	}
	schema.AddTables(postTable, userTable)

	r := require.New(t)

	expected := &spec.Relation{
		Type:     spec.RelTypeHasOne,
		Field:    userIDField,
		RefTable: postTable,
		RefField: postUserIDField,
	}

	actual, err := mergeRelation(userRelOnePostField, &Relation{Type: "HasMany", RefTable: "post"})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

}

func TestRelationManyToMany(t *testing.T) {
	schema := &spec.Schema{
		Name: "test",
	}
	schema.AddTables(categoryTable, postTable, userTable, skipTable, postCategoryTable)

	r := require.New(t)

	expected := &spec.Relation{
		Type:     spec.RelTypeManyToMany,
		Field:    postIDField,
		RefTable: categoryTable,
		RefField: categoryIDField,
		JoinTable: &spec.JoinTable{
			Name:         "post_category",
			JoinField:    postCategoryPostIDField,
			JoinRefField: postCategoryCategoryIDField,
		},
		Inverse: true,
	}

	actual, err := mergeRelation(postRelCategoriesField, &Relation{
		Type:     "manyToMany",
		RefTable: "category",
		RefField: "id",
		JoinTable: &JoinTable{
			Name:     "post_category",
			Table:    "post",
			RefTable: "category",
		},
		Inverse: true,
	})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

	actual, err = mergeRelation(postRelCategoriesField, &Relation{
		Type:     "manyToMany",
		RefTable: "category",
		RefField: "id",
		JoinTable: &JoinTable{
			Name:     "post_category",
			Table:    "post",
			Field:    "post_id",
			RefTable: "category",
			RefField: "category_id",
		},
		Inverse: true,
	})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

}

func TestRelationRemote(t *testing.T) {
	r := require.New(t)

	expected := &spec.Relation{
		Type:     spec.RelTypeBelongsTo,
		Field:    postRelRemoteIDField,
		RefTable: remoteTable,
		RefField: remoteIDField,
	}

	actual, err := mergeRelation(postRelRemoteField, &Relation{
		Type:     "BelongsTo",
		RefTable: "remote",
		RefField: "id",
	})
	r.NoError(err)
	matchRelation(t, expected, actual.Rel)

	actual, err = mergeRelation(postRelRemoteField, &Relation{
		Type:     "HasOne",
		RefTable: "remote",
		RefField: "id",
	})
	r.Error(err)

	actual, err = mergeRelation(postRelRemoteField, &Relation{
		Type:     "HasMany",
		RefTable: "remote",
		RefField: "id",
	})
	r.Error(err)
}

func TestMergeOps(t *testing.T) {
	id := spec.Builder("id").
		Type(&spec.IntegerType{Name: "int", Size: 32}).
		PrimaryKey(true).Build()

	name := spec.Builder("name").
		Type(&spec.StringType{Name: "char", Size: 64}).Build()

	f, err := mergeOps(id, []string{})

	r := require.New(t)
	r.NoError(err)
	r.Equal(len(spec.NumericOps), len(f.Ops))
	r.Equal(spec.Eq, f.Ops[spec.Eq-1])

	f, err = mergeOps(name, []string{})
	r.NoError(err)
	r.Equal(len(spec.StringOps), len(f.Ops))

	f, err = mergeOps(id, []string{"Eq"})
	r.NoError(err)
	r.Equal(1, len(f.Ops))
	r.Equal(spec.Eq, f.Ops[0])

}

func TestMergeType(t *testing.T) {
	tests := []struct {
		name     string
		expected spec.Type
		input    Type
	}{

		{
			name:  "bool",
			input: "bool",
			expected: &spec.BoolType{
				Name: "bool",
			},
		},
		{
			name:  "binary",
			input: "binary",
			expected: &spec.BinaryType{
				Name: "binary",
			},
		},
		{
			name:  "bit",
			input: "bit",
			expected: &spec.BitType{
				Name: "bit",
			},
		},
		{
			name:  "int8",
			input: "int8",
			expected: &spec.IntegerType{
				Name: "int8",
				Size: 8,
			},
		},
		{
			name:  "int16",
			input: "int16",
			expected: &spec.IntegerType{
				Name: "int16",
				Size: 16,
			},
		},
		{
			name:  "int32",
			input: "int32",
			expected: &spec.IntegerType{
				Name: "int32",
				Size: 32,
			},
		},
		{
			name:  "int64",
			input: "int64",
			expected: &spec.IntegerType{
				Name: "int64",
				Size: 64,
			},
		},
		{
			name:  "uint8",
			input: "uint8",
			expected: &spec.IntegerType{
				Name:     "uint8",
				Size:     8,
				Unsigned: true,
			},
		},
		{
			name:  "uint16",
			input: "uint16",
			expected: &spec.IntegerType{
				Name:     "uint16",
				Size:     16,
				Unsigned: true,
			},
		},
		{
			name:  "uint32",
			input: "uint32",
			expected: &spec.IntegerType{
				Name:     "uint32",
				Size:     32,
				Unsigned: true,
			},
		},
		{
			name:  "uint64",
			input: "uint64",
			expected: &spec.IntegerType{
				Name:     "uint64",
				Size:     64,
				Unsigned: true,
			},
		},
		{
			name:  "float32",
			input: "float32",
			expected: &spec.FloatType{
				Name:      "float32",
				Precision: 24,
			},
		},
		{
			name:  "float64",
			input: "float64",
			expected: &spec.FloatType{
				Name:      "float64",
				Precision: 32,
			},
		},
		{
			name:  "string",
			input: "string",
			expected: &spec.StringType{
				Name: "string",
			},
		},
		{
			name:  "time",
			input: "time",
			expected: &spec.TimeType{
				Name: "time",
			},
		},
		{
			name:  "enum",
			input: "enum",
			expected: &spec.EnumType{
				Name: "enum",
			},
		},
		{
			name:  "uuid",
			input: "uuid",
			expected: &spec.UUIDType{
				Name:    "uuid",
				Version: "v4",
			},
		},
		{
			name:  "json",
			input: "json",
			expected: &spec.JSONType{
				Name: "json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(test.expected, mergeType(test.input))
		})
	}
}
