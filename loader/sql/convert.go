package sql

import (
	"github.com/ychengcloud/cre/spec"
)

func (s *Schema) Convert() (*spec.Schema, error) {

	schema := &spec.Schema{
		Name: s.Name,
	}

	for _, t := range s.Tables {
		table, err := t.convert()
		if err != nil {
			return nil, err
		}
		schema.AddTables(table)
	}
	return schema, nil
}

func (t *Table) convert() (*spec.Table, error) {
	table := &spec.Table{
		Name: t.Name,
	}
	for _, c := range t.Columns {
		field, err := c.convert()
		if err != nil {
			return nil, err
		}

		table.AddFields(field)
	}
	return table, nil
}

func (c *Column) convert() (*spec.Field, error) {

	field := &spec.Field{
		Name:          c.Name,
		Type:          c.Type,
		Nullable:      c.Nullable,
		Comment:       c.Comment,
		Unique:        c.Unique,
		PrimaryKey:    c.Primary,
		AutoIncrement: c.AutoIncrement,
		OnUpdate:      c.OnUpdate,
	}

	convertIndexes(c, field)
	convertForeignKeys(c, field)

	if field.PrimaryKey {
		field.Filterable = true
		field.Sortable = true
	}
	return field, nil
}

func contains(indexColumns []*IndexColumn, c *Column) bool {
	for _, ics := range indexColumns {
		if ics.Column == c.Name {
			return true
		}
	}
	return false
}
func convertIndexes(c *Column, field *spec.Field) {
	for _, index := range c.Table.Indexes {
		if contains(index.IndexColumns, c) {

			if index.Primary {
				field.PrimaryKey = true
			}

			if index.Unique && len(index.IndexColumns) == 1 {
				field.Unique = true
			}
		}

	}
}

func convertForeignKeys(c *Column, field *spec.Field) {
	for _, foreighKey := range c.Table.ForeignKeys {
		for _, fc := range foreighKey.Columns {
			if fc.Name == c.Name {
				field.ForeignKey = true
			}
		}

	}
}
