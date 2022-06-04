package mysql

import (
	"fmt"
	"strconv"
	"strings"

	schema "github.com/ychengcloud/cre/loader/sql"
	"github.com/ychengcloud/cre/spec"
)

// parseTypeAttrs parses the column type attributes from the column definition.
// 返回值全部为小写格式
// If ext or attrs does not appear in the column definition, an empty slice is returned.
// eg: TINYINT[(M)] [UNSIGNED] [ZEROFILL]
//     colType ext   attr        attr
func parseTypeAttrs(colDef string) (colType string, ext []string, attrs []string, err error) {
	colDef = strings.ToLower(colDef)
	colDef = strings.TrimSpace(colDef)

	ext = make([]string, 0)
	attrs = make([]string, 0)
	idx := strings.IndexAny(colDef, "( ")
	if idx == -1 {
		colType = colDef
		return
	}
	colType = colDef[0:idx]

	before, after, found := strings.Cut(colDef[idx+1:], ")")
	if !found {
		attrs = strings.Fields(before)
		return
	}

	ext = strings.FieldsFunc(before, func(c rune) bool {
		return c == '\'' || c == ',' || c == ' '
	})
	if len(ext) == 0 {
		ext = append(ext, before)
	}
	attrs = strings.Fields(after)
	return
}

// parseInteger parses the integer column type from the column definition.
func parseInteger(colType string, ext []string, attrs []string) (spec.Type, error) {
	t := &spec.IntegerType{Name: colType}

	for _, attr := range attrs {
		if strings.Contains(attr, "unsigned") {
			t.Unsigned = true
		}

		// if strings.Contains(attr, "zerofill") {
		// 	t.ZeroFill = true
		// }
	}

	var size int64
	if len(ext) > 0 {
		var err error
		size, err = strconv.ParseInt(ext[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size: %v", ext[0])
		}
	}

	switch colType {
	case TypeTinyInt:
		if size == 1 {
			return &spec.BoolType{Name: colType}, nil
		}

		t.Size = 8
	case TypeSmallInt:
		t.Size = 16
	case TypeMediumInt, TypeInt:
		t.Size = 32
	case TypeBigInt:
		t.Size = 64
	default:
		return nil, fmt.Errorf("invalid col type: %v", colType)
	}

	return t, nil

}

// parseDecimal parses the decimal column type from the column definition.
func parseDecimal(colType string, ext []string) (spec.Type, error) {
	t := &spec.FloatType{
		Name: colType,
	}

	if len(ext) > 0 {
		precision, err := strconv.ParseInt(ext[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("decimal: invalid precision: %v", ext[0])
		}
		t.Precision = int(precision)
	}
	if len(ext) > 1 {
		scale, err := strconv.ParseInt(ext[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("decimal: invalid scale: %v", ext[1])
		}
		t.Scale = int(scale)
	}
	return t, nil

}

// parseFloat parses the float column type from the column definition.
func parseFloat(colType string, ext []string) (spec.Type, error) {
	t := &spec.FloatType{
		Name: colType,
	}

	if len(ext) > 0 {
		precision, err := strconv.ParseInt(ext[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("float: invalid precision: %v", ext[0])
		}
		t.Precision = int(precision)
	}
	if len(ext) > 1 {
		scale, err := strconv.ParseInt(ext[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("float: invalid scale: %v", ext[1])
		}
		t.Scale = int(scale)
	}
	return t, nil

}

// parseString parses the string column type from the column definition.
func parseString(colType string, ext []string, attrs []string) (spec.Type, error) {

	t := &spec.StringType{
		Name: colType,
	}

	for i, attr := range attrs {
		if strings.Contains(attr, "character") {
			t.Charset = attrs[i+2]
		}
		if strings.Contains(attr, "collate") {
			t.Collation = attrs[i+1]
		}
	}

	if len(ext) == 0 {
		return t, nil
	}

	size, err := strconv.ParseInt(ext[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size: %v", ext[0])
	}
	t.Size = int(size)

	return t, nil
}

// parseEnum parses the enum and set column type from the column definition.
func parseEnum(colType string, ext []string) (spec.Type, error) {

	t := &spec.EnumType{
		Name: colType,
	}
	if len(ext) == 0 {
		return t, nil
	}
	t.Values = append(t.Values, ext...)

	return t, nil
}

// parseBit parses the bit and set column type from the column definition.
func parseBit(colType string, ext []string) (spec.Type, error) {

	t := &spec.BitType{
		Name: colType,
	}
	if len(ext) == 0 {
		return t, nil
	}
	len, err := strconv.ParseInt(ext[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid len: %v", ext[0])
	}
	t.Len = int(len)

	return t, nil
}

// parseTime parses the time column type from the column definition.
func parseTime(colType string, ext []string) (spec.Type, error) {
	t := &spec.TimeType{
		Name: colType,
	}

	return t, nil
}

// parseSpatial parses the time column type from the column definition.
func parseSpatial(colType string, ext []string) (spec.Type, error) {
	t := &spec.SpatialType{
		Name: colType,
	}

	return t, nil
}

// parseJSON parses the time column type from the column definition.
func parseJSON(colType string, ext []string) (spec.Type, error) {
	t := &spec.JSONType{
		Name: colType,
	}

	return t, nil
}

// ParseColumn parses the column from the column definition.
func ParseType(colDef string) (spec.Type, error) {
	colType, ext, attrs, err := parseTypeAttrs(colDef)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(colType) {
	case TypeInt, TypeTinyInt, TypeSmallInt, TypeMediumInt, TypeBigInt:
		return parseInteger(colType, ext, attrs)
	case TypeBit:
		return parseBit(colType, ext)
	case TypeDecimal, TypeDouble:
		return parseDecimal(colType, ext)
	case TypeFloat:
		return parseFloat(colType, ext)
	case TypeChar, TypeVarchar, TypeBinary, TypeVarBinary, TypeTinyBlob, TypeTinyText, TypeBlob, TypeText, TypeMediumBlob, TypeMediumText, TypeLongBlob, TypeLongText:
		return parseString(colType, ext, attrs)
	case TypeEnum, TypeSet:
		return parseEnum(colType, ext)
	case TypeDate, TypeTime, TypeDateTime, TypeTimestamp, TypeYear:
		return parseTime(colType, ext)
	case TypeGeometry, TypePoint, TypeLineString, TypePolygon, TypeMultiPoint, TypeMultiLineString, TypeMultiPolygon, TypeGeometryCollection:
		return parseSpatial(colType, ext)
	case TypeJSON:
		return parseJSON(colType, ext)
	default:
		return nil, fmt.Errorf("invalid column type: %v", colDef)
	}

}

func parseExtra(c *schema.Column, extra string) {
	extra = strings.ToLower(extra)
	switch extra {
	case "auto_increment":
		c.AutoIncrement = true
	case "default_generated on update current_timestamp", "on update current_timestamp", "on update current_timestamp()":
		c.OnUpdate = true
	}
}
