package main

import (
	"context"

	"github.com/ychengcloud/cre/spec"
)

type FakeLoader struct {
	schema *spec.Schema
}

func NewFakeLoader() *FakeLoader {
	schema := &spec.Schema{Name: "test"}

	categoryTable := &spec.Table{Name: "category"}
	postTable := &spec.Table{Name: "post"}
	userTable := &spec.Table{Name: "user"}
	postCategoryTable := &spec.Table{Name: "post_category"}

	schema.AddTables(categoryTable, postTable, userTable, postCategoryTable)

	categoryTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
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
		},
		{
			Name: "name",
			Type: &spec.StringType{Name: "char", Size: 64},
		},
	}...)
	postTable.ID = postTable.GetField("id")

	userTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
		},
		{
			Name: "name",
			Type: &spec.StringType{Name: "char", Size: 64},
		},
	}...)
	userTable.ID = userTable.GetField("id")

	postCategoryTable.AddFields([]*spec.Field{
		{
			Name:       "id",
			Type:       &spec.IntegerType{Name: "int", Size: 32},
			PrimaryKey: true,
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
	postCategoryTable.ID = postCategoryTable.GetField("id")

	return &FakeLoader{
		schema: schema,
	}
}
func (f *FakeLoader) Load(ctx context.Context, name string) (*spec.Schema, error) {
	return f.schema, nil
}

func (f *FakeLoader) Dialect() string {
	return "mysql"
}
