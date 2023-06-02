package sql

import (
	"database/sql"

	"github.com/ychengcloud/cre/spec"
)

// ReferenceOption represents constraint actions.
type ReferenceOption string

// https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html#foreign-key-referential-actions
// the referential action specified by ON UPDATE and ON DELETE subclauses of the FOREIGN KEY clause
const (
	NoAction   ReferenceOption = "NO ACTION"
	Restrict   ReferenceOption = "RESTRICT"
	Cascade    ReferenceOption = "CASCADE"
	SetNull    ReferenceOption = "SET NULL"
	SetDefault ReferenceOption = "SET DEFAULT"
)

type (
	// Attribute represents an attribute definition.
	Attribute interface {
		Name() string
		Value() any
	}

	// Schema represents an schema definition.
	Schema struct {
		Name   string
		Tables []*Table
		Attrs  []Attribute
	}

	// Table represents a table definition.
	Table struct {
		Name          string
		Charset       string
		Collation     string
		AutoIncrement int
		Comment       string
		Options       string
		Columns       []*Column
		Indexes       []*Index
		ForeignKeys   []*ForeignKey
		Attrs         []Attribute

		Schema *Schema
	}

	// Column represents a column definition.
	Column struct {
		Name          string
		Type          spec.Type
		Unique        bool
		Nullable      bool
		Sensitive     bool
		Comment       string
		Default       sql.NullString
		Charset       string
		Collation     string
		Precision     int
		Scale         int
		Primary       bool
		AutoIncrement bool
		OnUpdate      bool
		Attrs         []Attribute

		Table *Table
	}

	Index struct {
		Name         string
		Unique       bool
		Type         string
		Comment      string
		IndexColumns []*IndexColumn
		Primary      bool
	}

	IndexColumn struct {
		SeqNo  int
		Column string
		Sub    int
		Expr   Expression
	}

	Expression interface{}

	ForeignKey struct {
		Name       string
		Table      *Table
		Columns    []*Column
		RefTable   *Table
		RefColumns []*Column
		Inverse    bool
		OnUpdate   ReferenceOption
		OnDelete   ReferenceOption
		Attrs      []Attribute
	}
)

func (s *Schema) Table(name string) *Table {
	for _, t := range s.Tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (t *Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}
	return nil
}
