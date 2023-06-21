package gen

import (
	"fmt"

	"github.com/ychengcloud/cre/spec"
)

func mergeOps(f *spec.Field, ops []string) (*spec.Field, error) {
	if len(ops) == 0 {
		return f, nil
	}

	f.Ops = make([]spec.Op, len(ops))
	for i, opc := range ops {
		op := spec.GetOP(opc)
		if op == spec.Unknown {
			return nil, fmt.Errorf("unknown operation: %s", opc)
		}
		f.Ops[i] = op
	}
	return f, nil
}

func mergeJoinTable(f *spec.Field, joinTableInCfg *JoinTable) (*spec.Field, error) {
	if f.Rel == nil {
		return f, nil
	}
	if joinTableInCfg == nil {
		joinTableInCfg = &JoinTable{}
	}

	// 关联表如没有配置相关信息，则使用默认值
	// eg: 表 a 和 b 关联 ，主键名 id,  则默认值为 关联表名 a_b , 关联字段 a_id , 关联引用字段 b_id
	if joinTableInCfg.Name == "" {
		joinTableInCfg.Name = f.Table.Name + "_" + f.Rel.RefTable.Name
	}
	if joinTableInCfg.Field == "" {
		joinTableInCfg.Field = f.Table.Name + "_" + f.Table.ID.Name

	}
	if joinTableInCfg.RefField == "" {
		joinTableInCfg.RefField = f.Rel.RefTable.Name + "_" + f.Rel.RefField.Name
	}

	joinTable := f.Table.Schema.Table(joinTableInCfg.Name)
	if joinTable == nil {
		return nil, fmt.Errorf("join table %s not found", joinTableInCfg.Name)
	}

	joinTable.IsJoinTable = true

	joinField := joinTable.GetField(joinTableInCfg.Field)
	if joinField == nil {
		return nil, fmt.Errorf("join field %s not found", joinTableInCfg.Field)
	}

	joinRefField := joinTable.GetField(joinTableInCfg.RefField)
	if joinRefField == nil {
		return nil, fmt.Errorf("join ref field %s not found", joinTableInCfg.RefField)
	}

	jt := &spec.JoinTable{
		Name:         joinTableInCfg.Name,
		JoinField:    joinField,
		JoinRefField: joinRefField,
	}

	f.Rel.JoinTable = jt
	joinTable.JoinTable = jt

	return f, nil
}

func refTable(f *spec.Field, refTableName string) (*spec.Table, error) {
	if f.Remote {
		return &spec.Table{Name: refTableName}, nil
	}

	rt := f.Table.Schema.Table(refTableName)
	if rt == nil {
		return nil, fmt.Errorf("mergeRelation/table %s not found", refTableName)
	}
	return rt, nil
}

func refField(f *spec.Field, refTableName string, refFieldName string) (*spec.Field, error) {
	if f.Remote {
		rf := spec.Builder(refFieldName).Build()
		f.Rel.RefTable.AddFields(rf)
		f.Rel.RefTable.ID = rf
		return rf, nil
	}

	f.Rel.RefField = f.Rel.RefTable.GetField(refFieldName)
	if f.Rel.RefField == nil {
		return nil, fmt.Errorf("ref field  not found [  table: %s, field: %s, ref field: %s ]", f.Table.Name, f.Name, refFieldName)
	}
	return f.Rel.RefField, nil
}

func defaultRelField(f *spec.Field) string {
	var name string
	switch f.Rel.Type {
	case spec.RelTypeBelongsTo:
		name = f.Rel.RefTable.Name + "_" + DefaultIDName
	case spec.RelTypeHasOne, spec.RelTypeHasMany:
		name = DefaultIDName
	case spec.RelTypeManyToMany:
		name = DefaultIDName
	}
	return name
}

func defaultRelRefField(f *spec.Field) string {
	var name string
	switch f.Rel.Type {
	case spec.RelTypeBelongsTo:
		name = DefaultIDName
	case spec.RelTypeHasOne, spec.RelTypeHasMany:
		name = f.Table.Name + "_" + DefaultIDName
	case spec.RelTypeManyToMany:
		name = DefaultIDName
	}
	return name
}

func relationField(f *spec.Field, rel *Relation) (*spec.Field, error) {
	var err error
	relField := rel.Field
	relRefField := rel.RefField

	f.Rel.RefTable, err = refTable(f, rel.RefTable)
	if err != nil {
		return nil, err
	}

	if relField == "" {
		relField = defaultRelField(f)
	}

	if relRefField == "" {
		relRefField = defaultRelRefField(f)
	}

	f.Rel.Field = f.Table.GetField(relField)
	if f.Rel.Field == nil {
		return nil, fmt.Errorf("relation not found [  table: %s, field: %s, relation field: %s ]", f.Table.Name, f.Name, relField)
	}

	f.Rel.RefField, err = refField(f, rel.RefTable, relRefField)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func mergeRelation(f *spec.Field, rel *Relation) (*spec.Field, error) {
	if rel == nil {
		return f, nil
	}

	var err error
	rt := spec.GetRelType(rel.Type)

	if f.Remote && (rt != spec.RelTypeBelongsTo && rt != spec.RelTypeManyToMany) {
		return nil, fmt.Errorf("remote field %s can only be belongs to or many to many", f.Name)
	}
	f.Rel = &spec.Relation{Type: rt}
	f, err = relationField(f, rel)
	if err != nil {
		return nil, err
	}

	if f.Rel.Type == spec.RelTypeManyToMany {
		f, err = mergeJoinTable(f, rel.JoinTable)
		if err != nil {
			return nil, err
		}
	}

	if rel.Inverse {
		f.Rel.Inverse = rel.Inverse
	}
	return f, nil
}

func mergeType(t Type) spec.Type {
	if t == "" {
		return nil
	}

	switch t {
	case Binary:
		return &spec.BinaryType{Name: string(t)}
	case Bit:
		return &spec.BitType{Name: string(t)}
	case Bool:
		return &spec.BoolType{Name: string(t)}
	case String:
		return &spec.StringType{Name: string(t)}
	case Int8:
		return &spec.IntegerType{Name: string(t), Size: 8}
	case Int16:
		return &spec.IntegerType{Name: string(t), Size: 16}
	case Int32:
		return &spec.IntegerType{Name: string(t), Size: 32}
	case Int64:
		return &spec.IntegerType{Name: string(t), Size: 64}
	case Uint8:
		return &spec.IntegerType{Name: string(t), Size: 8, Unsigned: true}
	case Uint16:
		return &spec.IntegerType{Name: string(t), Size: 16, Unsigned: true}
	case Uint32:
		return &spec.IntegerType{Name: string(t), Size: 32, Unsigned: true}
	case Uint64:
		return &spec.IntegerType{Name: string(t), Size: 64, Unsigned: true}
	case Float32:
		return &spec.FloatType{Name: string(t), Precision: 24}
	case Float64:
		return &spec.FloatType{Name: string(t), Precision: 32}
	case Time:
		return &spec.TimeType{Name: string(t)}
	case JSON:
		return &spec.JSONType{Name: string(t)}
	case UUID:
		return &spec.UUIDType{Name: string(t), Version: "v4"}
	case Enum:
		return &spec.EnumType{Name: string(t)}
	default:
		return nil
	}
}

func mergeField(f *spec.Field, fc *Field) (*spec.Field, error) {
	if fc == nil {
		return nil, fmt.Errorf("field config is nil")
	}

	if f == nil {
		return nil, fmt.Errorf("field is nil")
	}

	if fc.Nullable {
		f.Nullable = fc.Nullable
	}
	if fc.Optional {
		f.Optional = fc.Optional
	}
	if len(fc.Comment) > 0 {
		f.Comment = fc.Comment
	}

	f.Order = fc.Order

	if len(fc.Alias) > 0 {
		f.Alias = fc.Alias
	}

	if fc.Sortable {
		f.Sortable = fc.Sortable
	}
	if fc.Filterable != nil {
		f.Filterable = *fc.Filterable
	}

	f.Remote = fc.Remote

	if fc.Type != "" {
		f.Type = mergeType(fc.Type)
		if f.Type == nil {
			return nil, fmt.Errorf("unknown field type in config: %s", fc.Type)
		}
	}

	for k, v := range fc.Attrs {
		attr := NewAttr(k, v)
		f.Attrs = append(f.Attrs, attr)
	}

	var err error

	f, err = mergeOps(f, fc.Operations)
	if err != nil {
		return nil, err
	}

	f, err = mergeRelation(f, fc.Relation)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func mergeTable(t *spec.Table, tc *Table) (*spec.Table, error) {
	if t == nil {
		return nil, fmt.Errorf("table is nil")
	}

	if tc == nil {
		return nil, fmt.Errorf("table config is nil")
	}

	for _, fcs := range tc.Fields {
		if fcs.Skip {
			t.RemoveField(fcs.Name)
			continue
		}
		field := t.GetField(fcs.Name)
		if field == nil {
			if !fcs.Remote && fcs.Relation == nil {
				return nil, fmt.Errorf("field %s not found", fcs.Name)
			}
			field = &spec.Field{
				Name: fcs.Name,
				Type: &spec.ObjectType{Name: fcs.Name},
			}
			t.AddFields(field)
		}
		var err error
		field, err = mergeField(field, fcs)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range tc.Attrs {
		attr := NewAttr(k, v)
		t.Attrs = append(t.Attrs, attr)
	}

	return t, nil
}
func mergeSchema(s *spec.Schema, tables []*Table) (*spec.Schema, error) {
	if s == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	for _, tc := range tables {
		if tc.Skip {
			s.RemoveTable(tc.Name)
			continue
		}
		table := s.Table(tc.Name)
		if table == nil {
			return nil, fmt.Errorf("table %s not found", tc.Name)
		}
		table, err := mergeTable(table, tc)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}
