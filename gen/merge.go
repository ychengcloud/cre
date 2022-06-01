package gen

import (
	"fmt"

	"github.com/ychengcloud/cre/spec"
)

// ops returns all operations for given field.
func getOps(f *spec.Field) (o []spec.Op) {

	var ops []spec.Op
	switch f.Type.(type) {
	case *spec.BoolType:
		ops = spec.BoolOps
	case *spec.EnumType:
		ops = spec.EnumOps
	case *spec.IntegerType:
		ops = spec.NumericOps
	case *spec.StringType:
		ops = spec.StringOps
	case *spec.TimeType:
		ops = spec.NumericOps
	default:
		return
	}

	if f.Optional {
		ops = append(ops, spec.NullableOps...)
	}
	for _, op := range ops {
		o = append(o, op)
	}
	return
}

func mergeOps(f *spec.Field, ops []string) (*spec.Field, error) {
	if len(ops) > 0 {
		for _, opc := range ops {
			op := spec.GetOP(opc)
			if op == spec.Unknown {
				return nil, fmt.Errorf("unknown operation: %s", opc)
			}
			f.Ops = append(f.Ops, op)
		}
		return f, nil
	}

	f.Ops = getOps(f)

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

	joinTable.JoinTable = true

	joinField := joinTable.GetField(joinTableInCfg.Field)
	if joinField == nil {
		return nil, fmt.Errorf("join field %s not found", joinTableInCfg.Field)
	}

	joinRefField := joinTable.GetField(joinTableInCfg.RefField)
	if joinRefField == nil {
		return nil, fmt.Errorf("join ref field %s not found", joinTableInCfg.RefField)
	}

	f.Rel.JoinTable = &spec.JoinTable{
		Name:         joinTableInCfg.Name,
		JoinField:    joinField,
		JoinRefField: joinRefField,
	}

	return f, nil
}

func refTable(f *spec.Field, refTableName string) (*spec.Table, error) {
	if f.Remote {
		return &spec.Table{Name: refTableName}, nil
	}

	refTable := f.Table.Schema.Table(refTableName)
	if refTable == nil {
		return nil, fmt.Errorf("mergeRelation/table %s not found", refTableName)
	}
	return refTable, nil
}

func refField(f *spec.Field, refTableName string, refFieldName string) (*spec.Field, error) {
	if refFieldName == "" {
		refFieldName = "id"
	}
	if f.Remote {
		rf := &spec.Field{Name: refFieldName}
		f.Rel.RefTable.AddFields(rf)
		return rf, nil
	}

	f.Rel.RefField = f.Rel.RefTable.GetField(refFieldName)
	if f.Rel.RefField == nil {
		return nil, fmt.Errorf("ref field %s not found", refFieldName)
	}
	return f.Rel.RefField, nil
}

func mergeRelation(f *spec.Field, rel *Relation) (*spec.Field, error) {
	if rel == nil {
		return f, nil
	}

	var err error
	rt := spec.GetRelType(rel.Type)
	if rt == nil {
		return nil, fmt.Errorf("unknown relation type: %s", rel.Type)
	}
	f.Rel = &spec.Relation{Type: *rt}

	if rel.Field == "" {
		f.Rel.Field = f.Table.ID
		if f.Table.ID == nil {
			return nil, fmt.Errorf("relation field table id %#v not found", f.Table)
		}
	} else {
		f.Rel.Field = f.Table.GetField(rel.Field)
		if f.Rel.Field == nil {
			return nil, fmt.Errorf("relation field %s not found", rel.Field)
		}
	}

	f.Rel.RefTable, err = refTable(f, rel.RefTable)
	if err != nil {
		return nil, err
	}

	f.Rel.RefField, err = refField(f, rel.RefTable, rel.RefField)
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
	if len(fc.Alias) > 0 {
		f.Alias = fc.Alias
	}

	if fc.Sortable {
		f.Sortable = fc.Sortable
	}
	if fc.Filterable {
		f.Filterable = fc.Filterable
	}

	f.Remote = fc.Remote

	f, err := mergeOps(f, fc.Operations)
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
