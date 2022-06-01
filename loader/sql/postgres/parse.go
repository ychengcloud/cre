package postgres

import (
	"fmt"
	"strings"

	"github.com/ychengcloud/cre/spec"
)

type columnInfo struct {
	dataType  string
	nullable  bool
	size      int64
	udt       string
	precision int64
	scale     int64
	collation string
	charset   string
	typtype   string
}

// parseInteger parses the integer column type from the column definition.
func parseInteger(ci *columnInfo) (spec.Type, error) {
	dt := strings.ToLower(ci.dataType)
	t := &spec.IntegerType{Name: dt}

	switch dt {
	case TypeSmallInt, TypeSmallSerial, TypeInt2:
		t.Size = 16
	case TypeInteger, TypeSerial, TypeInt4:
		t.Size = 32
	case TypeBigInt, TypeBigSerial, TypeInt8:
		t.Size = 64
	default:
		return nil, fmt.Errorf("invalid data type: %v", dt)
	}

	return t, nil

}

// ParseColumn parses the column from the column definition.
func parseType(ci *columnInfo) (spec.Type, error) {
	dt := strings.ToLower(ci.dataType)
	switch dt {
	case TypeBit, TypeBitVar:
		return &spec.BitType{Name: dt, Len: int(ci.size)}, nil
	case TypeBytea:
		return &spec.BinaryType{Name: dt, Size: int(ci.size)}, nil
	case TypeBoolean, TypeBool:
		return &spec.BoolType{Name: dt}, nil
	case TypeSmallInt, TypeInteger, TypeBigInt, TypeSmallSerial, TypeSerial, TypeBigSerial, TypeInt, TypeInt2, TypeInt4, TypeInt8:
		return parseInteger(ci)
	case TypeDecimal, TypeNumeric:
		return &spec.FloatType{Name: dt, Precision: int(ci.precision), Scale: int(ci.scale)}, nil
	case TypeReal, TypeDouble, TypeFloat4, TypeFloat8, TypeMoney:
		return &spec.FloatType{Name: dt, Precision: int(ci.precision), Scale: int(ci.scale)}, nil
	case TypeCharacter, TypeChar, TypeCharVar, TypeVarChar, TypeText, TypeUUID, TypeTSQuery, TypeTSVector, TypeXML:
		return &spec.StringType{Name: dt, Collation: ci.collation, Charset: ci.charset, Size: int(ci.size)}, nil
	case TypeCircle, TypeLine, TypeLseg, TypeBox, TypePath, TypePolygon, TypePoint, TypeCIDR, TypeInet, TypeMACAddr, TypeMACAddr8, TypeArray, TypePgLSN, TypePgSnapshot:
		return &spec.SpatialType{Name: dt}, nil
	case TypeDate, TypeTime, TypeTimeWithTZ, TypeTimeWithoutTZ,
		TypeTimestamp, TypeTimestampTZ, TypeTimestampWithTZ, TypeTimestampWithoutTZ, TypeInterval:
		return &spec.TimeType{Name: dt}, nil
	case TypeJSON, TypeJSONB:
		return &spec.JSONType{Name: dt}, nil
	case TypeUserDefined:
		if ci.typtype == "e" {
			return &spec.EnumType{Name: ci.udt}, nil
		}
		return &spec.SpatialType{Name: dt}, nil
	default:
		return nil, fmt.Errorf("invalid column type: %v", dt)
	}

}
