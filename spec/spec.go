package spec

import (
	"database/sql"
	"sort"
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
		Name    string
		Comment string
		fields  []*Field
		ID      *Field
		Attrs   []Attribute

		IsJoinTable bool
		JoinTable   *JoinTable

		Schema *Schema
	}

	// Field represents a Field definition.
	Field struct {
		Name      string         `json:"name,omitempty"`
		Type      Type           `json:"type,omitempty"`
		Nullable  bool           `json:"nullable,omitempty"`
		Optional  bool           `json:"optional,omitempty"`
		Sensitive bool           `json:"sensitive,omitempty"`
		Tag       string         `json:"tag,omitempty"`
		Comment   string         `json:"comment,omitempty"`
		Default   sql.NullString `json:"default,omitempty"`
		Order     int            `json:"order,omitempty"`

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

		Attrs []Attribute `json:"attrs,omitempty"`
	}

	// Relation represents	a Relation definition.
	Relation struct {
		Type      RelType     `json:"type,omitempty"`
		Field     *Field      `json:"field,omitempty"`
		RefTable  *Table      `json:"ref_table,omitempty"`
		RefField  *Field      `json:"ref_field,omitempty"`
		JoinTable *JoinTable  `json:"join_table,omitempty"`
		Inverse   bool        `json:"inverse,omitempty"`
		Attrs     []Attribute `json:"attrs,omitempty"`
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

type (
	// FieldBuilder represents a field builder.
	FieldBuilder struct {
		field *Field
	}
)

func Builder(name string) *FieldBuilder {
	f := &Field{Name: name}
	f.Ops = defaultOps(f.Type, f.Optional)
	return &FieldBuilder{field: f}
}

func (fb *FieldBuilder) Type(t Type) *FieldBuilder {
	fb.field.Type = t
	fb.field.Ops = defaultOps(t, fb.field.Optional)
	return fb
}

func (fb *FieldBuilder) Nullable(nullable bool) *FieldBuilder {
	fb.field.Nullable = nullable
	return fb
}

func (fb *FieldBuilder) Optional(optional bool) *FieldBuilder {
	fb.field.Optional = optional
	fb.field.Ops = defaultOps(fb.field.Type, fb.field.Optional)
	return fb
}

func (fb *FieldBuilder) Sensitive(sensitive bool) *FieldBuilder {
	fb.field.Sensitive = sensitive
	return fb
}

func (fb *FieldBuilder) Tag(tag string) *FieldBuilder {
	fb.field.Tag = tag
	return fb
}

func (fb *FieldBuilder) Comment(comment string) *FieldBuilder {
	fb.field.Comment = comment
	return fb
}

func (fb *FieldBuilder) Default(d sql.NullString) *FieldBuilder {
	fb.field.Default = d
	return fb
}

func (fb *FieldBuilder) Alias(alias string) *FieldBuilder {
	fb.field.Alias = alias
	return fb
}

func (fb *FieldBuilder) Sortable(sortable bool) *FieldBuilder {
	fb.field.Sortable = sortable
	return fb
}

func (fb *FieldBuilder) Filterable(filterable bool) *FieldBuilder {
	fb.field.Filterable = filterable
	return fb
}

func (fb *FieldBuilder) ForeignKey(foreignKey bool) *FieldBuilder {
	fb.field.ForeignKey = foreignKey
	return fb
}

func (fb *FieldBuilder) PrimaryKey(primaryKey bool) *FieldBuilder {
	fb.field.PrimaryKey = primaryKey
	if primaryKey {
		fb.field.Filterable = true
		fb.field.Sortable = true
	}
	return fb
}

func (fb *FieldBuilder) Index(index bool) *FieldBuilder {
	fb.field.Index = index
	if index {
		fb.field.Filterable = true
	}
	return fb
}

func (fb *FieldBuilder) Unique(unique bool) *FieldBuilder {
	fb.field.Unique = unique
	return fb
}

func (fb *FieldBuilder) AutoIncrement(autoIncrement bool) *FieldBuilder {
	fb.field.AutoIncrement = autoIncrement
	return fb
}

func (fb *FieldBuilder) OnUpdate(onUpdate bool) *FieldBuilder {
	fb.field.OnUpdate = onUpdate
	return fb
}

func (fb *FieldBuilder) Remote(remote bool) *FieldBuilder {
	fb.field.Remote = remote
	return fb
}

func (fb *FieldBuilder) Rel(rel *Relation) *FieldBuilder {
	fb.field.Rel = rel
	return fb
}

func (fb *FieldBuilder) Ops(ops []Op) *FieldBuilder {
	fb.field.Ops = ops
	return fb
}

func (fb *FieldBuilder) Table(table *Table) *FieldBuilder {
	fb.field.Table = table
	return fb
}

func (fb *FieldBuilder) Attrs(attrs ...Attribute) *FieldBuilder {
	fb.field.Attrs = append(fb.field.Attrs, attrs...)
	return fb
}

func (fb *FieldBuilder) Build() *Field {
	return fb.field
}

// Tables 返回非关联表的 Table 列表
func (s *Schema) Tables() []*Table {
	return s.tables
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

// Tables 返回非关联表的 Table 列表
func (s *Schema) NoJoinTables() []*Table {
	tables := make([]*Table, 0)
	for _, t := range s.tables {
		if !t.IsJoinTable {
			tables = append(tables, t)
		}
	}
	return tables
}

// JoinTables 返回所有关联表
func (s *Schema) JoinTables() []*Table {
	tables := make([]*Table, 0)
	for _, t := range s.tables {
		if t.IsJoinTable {
			tables = append(tables, t)
		}
	}
	return tables
}

// JoinTable returns the table with the given name.
func (s *Schema) JoinTable(name string) *Table {
	for _, t := range s.tables {
		if t.Name == name && t.IsJoinTable {
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

// defaultOps returns default operations for given field.
func defaultOps(t Type, nullable bool) (ops []Op) {

	switch t.(type) {
	case *BoolType:
		ops = BoolOps
	case *EnumType:
		ops = EnumOps
	case *IntegerType:
		ops = NumericOps
	case *FloatType:
		ops = NumericOps
	case *StringType:
		ops = StringOps
	case *TimeType:
		ops = NumericOps
	default:
		return
	}

	if nullable {
		ops = append(ops, NullableOps...)
	}

	return
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

// HasRelations returns true if the table has a relation table.
func (t *Table) HasRelations() bool {
	for _, f := range t.fields {
		if !f.RelNone() {
			return true
		}
	}
	return false
}

// AutoIncrement returns true if the table has a auto increment primaryKey.
func (t *Table) AutoIncrement() bool {
	if t.ID == nil {
		return false
	}

	return t.ID.AutoIncrement
}

func (t *Table) SortedFields() []*Field {
	fields := make([]*Field, len(t.fields))
	copy(fields, t.fields)
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Order > fields[j].Order
	})
	return fields
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

// NameOrAlias returns the alias if it is not empty, otherwise returns the name.
func (f *Field) NameOrAlias() string {
	if f.Alias != "" {
		return f.Alias
	}
	return f.Name
}
