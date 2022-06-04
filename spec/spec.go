package spec

import (
	"strings"
)

// RelType represents a relation type.
type RelType uint8

const (
	RelTypeNone RelType = iota
	RelTypeBelongsTo
	RelTypeHasOne
	RelTypeHasMany
	RelTypeManyToMany
)

type Op int

const (
	Unknown Op = iota
	Eq
	Neq
	In
	NotIn
	Gt
	Gte
	Lt
	Lte
	IsNil
	NotNil
	Contains
	StartsWith
	EndsWith
	AND
	OR
	NOT
)

var (
	opNames = [...]string{
		Eq:         "Eq",
		Neq:        "Neq",
		In:         "In",
		NotIn:      "NotIn",
		Gt:         "Gt",
		Gte:        "Gte",
		Lt:         "Lt",
		Lte:        "Lte",
		IsNil:      "IsNil",
		NotNil:     "NotNil",
		Contains:   "Contains",
		StartsWith: "StartsWith",
		EndsWith:   "EndsWith",
		AND:        "AND",
		OR:         "OR",
		NOT:        "NOT",
	}

	// operations collection.
	BoolOps     = []Op{Eq, Neq}
	EnumOps     = append(BoolOps, In, NotIn)
	NumericOps  = append(EnumOps, Gt, Gte, Lt, Lte)
	StringOps   = append(NumericOps, Contains, StartsWith, EndsWith)
	NullableOps = []Op{IsNil, NotNil}

	// relation names.
	relNames = [...]string{
		"None",
		"BelongsTo",
		"HasOne",
		"HasMany",
		"ManyToMany",
	}
)

func (o Op) Name() string {
	if int(o) < len(opNames) {
		return opNames[o]
	}
	return "Unknown"
}

func GetOP(name string) Op {
	for i, n := range opNames {
		if strings.ToLower(n) == strings.TrimSpace(strings.ToLower(name)) {
			return Op(i)
		}
	}
	return Unknown
}

func (r RelType) Name() string {
	if int(r) < len(relNames) {
		return relNames[r]
	}
	return "Unknown"
}

func GetRelType(name string) RelType {
	for i, n := range relNames {
		if strings.ToLower(n) == strings.TrimSpace(strings.ToLower(name)) {
			relType := RelType(i)
			return relType
		}
	}
	return RelTypeNone
}

type (
	// Schema represents an schema definition.
	Schema struct {
		Name   string
		tables []*Table
		Attrs  []Attribute
	}

	// Table represents a table definition.
	Table struct {
		Name   string
		fields []*Field
		ID     *Field

		JoinTable bool

		Schema *Schema
	}

	// Field represents a Field definition.
	Field struct {
		Name      string `json:"name,omitempty"`
		Type      Type   `json:"type,omitempty"`
		Nullable  bool   `json:"nullable,omitempty"`
		Optional  bool   `json:"optional,omitempty"`
		Sensitive bool   `json:"sensitive,omitempty"`
		Tag       string `json:"tag,omitempty"`
		Comment   string `json:"comment,omitempty"`

		Alias      string `json:"alias,omitempty"`
		Sortable   bool   `json:"sortable,omitempty"`
		Filterable bool   `json:"filterable,omitempty"`

		ForeignKey    bool `json:"foreignKey,omitempty"`
		PrimaryKey    bool `json:"primaryKey,omitempty"`
		Index         bool `json:"index,omitempty"`
		Unique        bool `json:"unique,omitempty"`
		AutoIncrement bool `json:"autoIncrement,omitempty"`
		OnUpdate      bool `json:"onUpdate,omitempty"`
		Remote        bool `json:"remote,omitempty"`

		Rel   *Relation `json:"rel,omitempty"`
		Ops   []Op      `json:"ops,omitempty"`
		Table *Table    `json:"table,omitempty"`

		Attrs []*Attribute `json:"attrs,omitempty"`
	}

	// Relation represents	a Relation definition.
	Relation struct {
		Type      RelType      `json:"type,omitempty"`
		Field     *Field       `json:"field,omitempty"`
		RefTable  *Table       `json:"ref_table,omitempty"`
		RefField  *Field       `json:"ref_field,omitempty"`
		JoinTable *JoinTable   `json:"join_table,omitempty"`
		Inverse   bool         `json:"inverse,omitempty"`
		Attrs     []*Attribute `json:"attrs,omitempty"`
	}

	JoinTable struct {
		Name         string `json:"name,omitempty"`
		JoinField    *Field `json:"join_field,omitempty"`
		JoinRefField *Field `json:"join_ref_field,omitempty"`
	}
	// Attribute represents an attribute definition.
	Attribute interface {
		Name() string
		Value() any
	}
)

// Tables 返回非关联表的 Table 列表
func (s *Schema) Tables() []*Table {
	tables := make([]*Table, 0)
	for _, t := range s.tables {
		if !t.JoinTable {
			tables = append(tables, t)
		}
	}
	return tables
}

// Table return the table with the given name.
func (s *Schema) Table(name string) *Table {
	for _, t := range s.tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// JoinTables 返回所有关联表
func (s *Schema) JoinTables() []*Table {
	tables := make([]*Table, 0)
	for _, t := range s.tables {
		if t.JoinTable {
			tables = append(tables, t)
		}
	}
	return tables
}

// JoinTable returns the table with the given name.
func (s *Schema) JoinTable(name string) *Table {
	for _, t := range s.tables {
		if t.Name == name && t.JoinTable {
			return t
		}
	}
	return nil
}

// AddTable adds a new table to the schema.
func (s *Schema) AddTables(tables ...*Table) {
	for _, t := range tables {
		t.Schema = s
		s.tables = append(s.tables, t)

	}
}

// RemoveTable removes the table with the given name.
func (s *Schema) RemoveTable(name string) {
	for i, t := range s.tables {
		if t.Name == name {
			s.tables = append(s.tables[:i], s.tables[i+1:]...)
		}
	}
}

func (t *Table) Fields() []*Field {
	return t.fields
}

// FilterFields returns the fields that is filterable.
func (t *Table) FilterFields() []*Field {
	fields := make([]*Field, 0)
	for _, f := range t.fields {
		if f.Filterable {
			fields = append(fields, f)
		}
	}
	return fields
}

// AddField adds a new field to the table.
func (t *Table) AddFields(fields ...*Field) {
	for _, f := range fields {
		// 主键且唯一即为 ID 字段（复合主键按约定视为关联表主键）
		if f.PrimaryKey && f.Unique {
			t.ID = f
		}
		f.Table = t
		t.fields = append(t.fields, f)
	}
}

// GetField returns the field with the given name.
func (t *Table) GetField(name string) *Field {
	for _, f := range t.fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// RemoveField removes the field with the given name.
func (t *Table) RemoveField(name string) {
	for i, f := range t.fields {
		if f.Name == name {
			t.fields = append(t.fields[:i], t.fields[i+1:]...)
		}
	}
}

// HasFilterField returns true if any field is filterable
func (t *Table) HasFilterField() bool {
	for _, f := range t.fields {
		if f.Filterable {
			return true
		}
	}
	return false
}

// HasParent returns true if the table has a parent table.
func (t *Table) HasParent() bool {
	for _, f := range t.fields {
		if f.RelBelongsTo() || f.RelManyToMany() {
			return true
		}
	}
	return false
}

// RelNone returns true if the relation is none.
func (f *Field) RelNone() bool {
	return f.Rel == nil || f.Rel.Type == RelTypeNone
}

// RelBelongsTo returns true if the relation is belongsto
func (f *Field) RelBelongsTo() bool {
	return f.Rel != nil && f.Rel.Type == RelTypeBelongsTo
}

// RelOne returns true if the relation is one
func (f *Field) RelHasOne() bool {
	return f.Rel != nil && f.Rel.Type == RelTypeHasOne
}

// RelMany returns true if the relation is many
func (f *Field) RelHasMany() bool {
	return f.Rel != nil && f.Rel.Type == RelTypeHasMany
}

// RelManyToMany returns true if the relation is many to many
func (f *Field) RelManyToMany() bool {
	return f.Rel != nil && f.Rel.Type == RelTypeManyToMany
}
