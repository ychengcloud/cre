package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ychengcloud/cre/spec"
)

type col struct {
	typ   string
	ext   []string
	attrs []string
}

func TestParseTypeAttrs(t *testing.T) {
	type testcase struct {
		def      string
		expected col
	}
	tests := []testcase{
		{
			def: "bit(1)",
			expected: col{
				typ:   "bit",
				ext:   []string{"1"},
				attrs: []string{},
			},
		},
		{
			def: "tinyint(1)",
			expected: col{
				typ:   "tinyint",
				ext:   []string{"1"},
				attrs: []string{},
			},
		},
		{
			def: "tinyint(-128) unsigned",
			expected: col{
				typ:   "tinyint",
				ext:   []string{"-128"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "tinyint(127) unsigned zerofill",
			expected: col{
				typ:   "tinyint",
				ext:   []string{"127"},
				attrs: []string{"unsigned", "zerofill"},
			},
		},
		{
			def: "bool",
			expected: col{
				typ:   "bool",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "boolean",
			expected: col{
				typ:   "boolean",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "smallint",
			expected: col{
				typ:   "smallint",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "smallint unsigned",
			expected: col{
				typ:   "smallint",
				ext:   []string{},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "smallint(5) unsigned",
			expected: col{
				typ:   "smallint",
				ext:   []string{"5"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "mediumint",
			expected: col{
				typ:   "mediumint",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "mediumint(9) unsigned",
			expected: col{
				typ:   "mediumint",
				ext:   []string{"9"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "int",
			expected: col{
				typ:   "int",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "int(10) unsigned",
			expected: col{
				typ:   "int",
				ext:   []string{"10"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "integer",
			expected: col{
				typ:   "integer",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "integer(10) unsigned",
			expected: col{
				typ:   "integer",
				ext:   []string{"10"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "decimal(10,2)",
			expected: col{
				typ:   "decimal",
				ext:   []string{"10", "2"},
				attrs: []string{},
			},
		},
		{
			def: "decimal(10,2) unsigned",
			expected: col{
				typ:   "decimal",
				ext:   []string{"10", "2"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "float",
			expected: col{
				typ:   "float",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "float(10,2)",
			expected: col{
				typ:   "float",
				ext:   []string{"10", "2"},
				attrs: []string{},
			},
		},
		{
			def: "float(10,2) unsigned",
			expected: col{
				typ:   "float",
				ext:   []string{"10", "2"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "double",
			expected: col{
				typ:   "double",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "double(10,2)",
			expected: col{
				typ:   "double",
				ext:   []string{"10", "2"},
				attrs: []string{},
			},
		},
		{
			def: "double(10,2) unsigned",
			expected: col{
				typ:   "double",
				ext:   []string{"10", "2"},
				attrs: []string{"unsigned"},
			},
		},
		{
			def: "datetime",
			expected: col{
				typ:   "datetime",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "timestamp",
			expected: col{
				typ:   "timestamp",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "timestamp(6)",
			expected: col{
				typ:   "timestamp",
				ext:   []string{"6"},
				attrs: []string{},
			},
		},
		{
			def: "char",
			expected: col{
				typ:   "char",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "char(10)",
			expected: col{
				typ:   "char",
				ext:   []string{"10"},
				attrs: []string{},
			},
		},
		{
			def: "varchar",
			expected: col{
				typ:   "varchar",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "varchar(10)",
			expected: col{
				typ:   "varchar",
				ext:   []string{"10"},
				attrs: []string{},
			},
		},
		{
			def: "binary",
			expected: col{
				typ:   "binary",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "binary(10)",
			expected: col{
				typ:   "binary",
				ext:   []string{"10"},
				attrs: []string{},
			},
		},
		{
			def: "varbinary",
			expected: col{
				typ:   "varbinary",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "varbinary(10)",
			expected: col{
				typ:   "varbinary",
				ext:   []string{"10"},
				attrs: []string{},
			},
		},
		{
			def: "varbinary(10) character set utf8mb4_bin collate utf8mb4_bin",
			expected: col{
				typ: "varbinary",
				ext: []string{"10"},
				attrs: []string{
					"character",
					"set",
					"utf8mb4_bin",
					"collate",
					"utf8mb4_bin",
				},
			},
		},
		{
			def: "tinyblob",
			expected: col{
				typ:   "tinyblob",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "tinytext",
			expected: col{
				typ:   "tinytext",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "blob",
			expected: col{
				typ:   "blob",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "text",
			expected: col{
				typ:   "text",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "mediumblob",
			expected: col{
				typ:   "mediumblob",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "mediumtext",
			expected: col{
				typ:   "mediumtext",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "longblob",
			expected: col{
				typ:   "longblob",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "longtext",
			expected: col{
				typ:   "longtext",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "enum('a','b','c')",
			expected: col{
				typ:   "enum",
				ext:   []string{"a", "b", "c"},
				attrs: []string{},
			},
		},
		{
			def: "enum('a', 'b', 'c' , 'd')",
			expected: col{
				typ:   "enum",
				ext:   []string{"a", "b", "c", "d"},
				attrs: []string{},
			},
		},
		{
			def: "set('a','b','c')",
			expected: col{
				typ:   "set",
				ext:   []string{"a", "b", "c"},
				attrs: []string{},
			},
		},
		{
			def: "json",
			expected: col{
				typ:   "json",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "jsonb",
			expected: col{
				typ:   "jsonb",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "date",
			expected: col{
				typ:   "date",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "time",
			expected: col{
				typ:   "time",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "geometry",
			expected: col{
				typ:   "geometry",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "point",
			expected: col{
				typ:   "point",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "linestring",
			expected: col{
				typ:   "linestring",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "polygon",
			expected: col{
				typ:   "polygon",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "multipoint",
			expected: col{
				typ:   "multipoint",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "multilinestring",
			expected: col{
				typ:   "multilinestring",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "multipolygon",
			expected: col{
				typ:   "multipolygon",
				ext:   []string{},
				attrs: []string{},
			},
		},
		{
			def: "geometrycollection",
			expected: col{
				typ:   "geometrycollection",
				ext:   []string{},
				attrs: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.def, func(t *testing.T) {
			typ, ext, attrs, err := parseTypeAttrs(test.def)
			require.NoError(t, err)
			require.Equal(t, test.expected.typ, typ)
			require.EqualValues(t, test.expected.ext, ext)
			require.EqualValues(t, test.expected.attrs, attrs)
		})

	}

}

func TestParseType(t *testing.T) {
	tests := []struct {
		def      string
		expected spec.Type
	}{
		{
			def: "bit(1)",
			expected: &spec.BitType{
				Name: TypeBit,
				Len:  1,
			},
		},
		{
			def: "bigint",
			expected: &spec.IntegerType{
				Name: TypeBigInt,
				Size: 64,
			},
		},
		{
			def: "float(10,2)",
			expected: &spec.FloatType{
				Name:      TypeFloat,
				Precision: 10,
				Scale:     2,
			},
		},
		{
			def: "decimal(10,2)",
			expected: &spec.FloatType{
				Name:      TypeDecimal,
				Precision: 10,
				Scale:     2,
			},
		},
		{
			def: "decimal(10,2) unsigned",
			expected: &spec.FloatType{
				Name:      TypeDecimal,
				Precision: 10,
				Scale:     2,
			},
		},
		{
			def: "char(255)",
			expected: &spec.StringType{
				Name: TypeChar,
				Size: int(255),
			},
		},
		{
			def: "varchar(255)",
			expected: &spec.StringType{
				Name: TypeVarchar,
				Size: int(255),
			},
		},
		{
			def: "varchar(255) character set utf8mb4_bin collate utf8mb4_bin",
			expected: &spec.StringType{
				Name:      TypeVarchar,
				Size:      int(255),
				Charset:   "utf8mb4_bin",
				Collation: "utf8mb4_bin",
			},
		},
		{
			def: "enum('a','b')",
			expected: &spec.EnumType{
				Name: TypeEnum,
				Values: []string{
					"a", "b",
				},
			},
		},
		{
			def: "date",
			expected: &spec.TimeType{
				Name: TypeDate,
			},
		},
		{
			def: "time",
			expected: &spec.TimeType{
				Name: TypeTime,
			},
		},
		{
			def: "datetime",
			expected: &spec.TimeType{
				Name: TypeDateTime,
			},
		},
		{
			def: "datetime(5)",
			expected: &spec.TimeType{
				Name: TypeDateTime,
			},
		},
		{
			def: "timestamp",
			expected: &spec.TimeType{
				Name: TypeTimestamp,
			},
		},
		{
			def: "timestamp(5)",
			expected: &spec.TimeType{
				Name: TypeTimestamp,
			},
		},
		{
			def: "point",
			expected: &spec.SpatialType{
				Name: TypePoint,
			},
		},
	}

	for _, test := range tests {
		typ, err := ParseType(test.def)
		require.NoError(t, err)
		require.Equal(t, test.expected, typ)
	}
}
