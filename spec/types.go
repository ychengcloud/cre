package spec

import "reflect"

// Type represents a field type definition.
type Type interface {
	GetName() string
	Kind() string
	ProtobufKind() string
}

var _ Type = (*BinaryType)(nil)
var _ Type = (*BitType)(nil)
var _ Type = (*IntegerType)(nil)
var _ Type = (*FloatType)(nil)
var _ Type = (*StringType)(nil)
var _ Type = (*EnumType)(nil)
var _ Type = (*TimeType)(nil)
var _ Type = (*SpatialType)(nil)
var _ Type = (*JSONType)(nil)
var _ Type = (*ObjectType)(nil)
var _ Type = (*UUIDType)(nil)

type BinaryType struct {
	Name string
	Size int
}

type BitType struct {
	Name string

	// Len indicates the Length of bits
	Len int
}

type BoolType struct {
	Name string
}

type IntegerType struct {
	Name     string
	Size     int
	Unsigned bool
}

type FloatType struct {
	Name      string
	Precision int
	Scale     int
}

type StringType struct {
	Name      string
	Size      int
	Charset   string
	Collation string
}

type EnumType struct {
	Name   string
	Values []string
}

type UUIDType struct {
	Name    string
	Version string
}

type TimeType struct {
	Name string
	Size int
}

type SpatialType struct {
	Name string
}

type JSONType struct {
	Name string
}

type ObjectType struct {
	Name     string
	Exported bool
}

func (b *BinaryType) GetName() string {
	return b.Name
}
func (b *BitType) GetName() string {
	return b.Name
}
func (b *BoolType) GetName() string {
	return b.Name
}

func (i *IntegerType) GetName() string {
	return i.Name
}

func (f *FloatType) GetName() string {
	return f.Name
}
func (s *StringType) GetName() string {
	return s.Name
}
func (e *EnumType) GetName() string {
	return e.Name
}
func (t *TimeType) GetName() string {
	return t.Name
}
func (s *SpatialType) GetName() string {
	return s.Name
}
func (j *JSONType) GetName() string {
	return j.Name
}

func (o *ObjectType) GetName() string {
	return o.Name
}

func (o *UUIDType) GetName() string {
	return o.Name
}

func (*BinaryType) Kind() string {
	return reflect.Slice.String()
}
func (*BitType) Kind() string {
	return reflect.Uint8.String()
}
func (*BoolType) Kind() string {
	return reflect.Bool.String()
}

func (i *IntegerType) Kind() string {
	switch i.Size {
	case 8:
		if i.Unsigned {
			return reflect.Uint8.String()
		}
		return reflect.Int8.String()
	case 16:
		if i.Unsigned {
			return reflect.Uint16.String()
		}
		return reflect.Int16.String()
	case 32:
		if i.Unsigned {
			return reflect.Uint32.String()
		}
		return reflect.Int32.String()
	case 64:
		if i.Unsigned {
			return reflect.Uint64.String()
		}
		return reflect.Int64.String()
	default:
		if i.Unsigned {
			return reflect.Uint.String()
		}
		return reflect.Int.String()
	}

}

func (f *FloatType) Kind() string {
	if f.Precision > 24 {
		return reflect.Float64.String()
	}
	return reflect.Float32.String()
}
func (*StringType) Kind() string {
	return reflect.String.String()
}
func (*EnumType) Kind() string {
	return reflect.String.String()
}
func (*TimeType) Kind() string {
	return "time.Time"
}

func (*UUIDType) Kind() string {
	return "uuid"
}

func (s *SpatialType) Kind() string {
	return s.Name
}
func (*JSONType) Kind() string {
	return "[]byte"
}

func (o *ObjectType) Kind() string {
	return o.Name
}

func (*BinaryType) ProtobufKind() string {
	return "bytes"
}
func (*BitType) ProtobufKind() string {
	return "bool"
}
func (*BoolType) ProtobufKind() string {
	return "bool"
}

func (i *IntegerType) ProtobufKind() string {
	switch i.Size {
	case 8, 16, 32:
		if i.Unsigned {
			return "uint32"
		}
		return "int32"
	case 64:
		if i.Unsigned {
			return "uint64"
		}
		return "int64"
	default:
		if i.Unsigned {
			return "uint32"
		}
		return "int32"
	}

}

func (f *FloatType) ProtobufKind() string {
	if f.Precision > 24 {
		return "double"
	}
	return "float"
}
func (*StringType) ProtobufKind() string {
	return "string"
}
func (*EnumType) ProtobufKind() string {
	return "string"
}
func (*TimeType) ProtobufKind() string {
	return "google.protobuf.Timestamp"
}
func (s *SpatialType) ProtobufKind() string {
	return s.Name
}
func (*JSONType) ProtobufKind() string {
	return "bytes"
}

func (o *ObjectType) ProtobufKind() string {
	return o.Name
}

func (*UUIDType) ProtobufKind() string {
	return "uuid"
}
