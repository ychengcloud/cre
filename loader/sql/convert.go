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

	fb := spec.Builder(c.Name).
		Type(c.Type).
		Nullable(c.Nullable).
		Comment(c.Comment).
		Unique(c.Unique).
		PrimaryKey(c.Primary).
		AutoIncrement(c.AutoIncrement).
		OnUpdate(c.OnUpdate)

	convertIndexes(c, fb)
	convertForeignKeys(c, fb)

	return fb.Build(), nil
}

func contains(indexColumns []*IndexColumn, c *Column) bool {
	for _, ics := range indexColumns {
		if ics.Column == c.Name {
			return true
		}
	}
	return false
}
func convertIndexes(c *Column, b *spec.FieldBuilder) {
	for _, index := range c.Table.Indexes {
		if contains(index.IndexColumns, c) {
			b = b.Index(true)
			if index.Primary {
				b = b.PrimaryKey(true)
			}

			if index.Unique && len(index.IndexColumns) == 1 {
				b = b.Unique(true)
			}
		}

	}
}

func convertForeignKeys(c *Column, b *spec.FieldBuilder) {
	for _, foreighKey := range c.Table.ForeignKeys {
		for _, fc := range foreighKey.Columns {
			if fc.Name == c.Name {
				b = b.ForeignKey(true)
			}
		}

	}
}
