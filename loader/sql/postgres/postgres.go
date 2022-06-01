package postgres

// Postgres data types
// https://www.postgresql.org/docs/current/datatype.html

const (
	TypeBit     = "bit"
	TypeBitVar  = "bit varying"
	TypeBoolean = "boolean"
	TypeBool    = "bool" // boolean.
	TypeBytea   = "bytea"

	TypeCharacter = "character"
	TypeChar      = "char" // character
	TypeCharVar   = "character varying"
	TypeVarChar   = "varchar" // character varying
	TypeText      = "text"

	TypeSmallInt = "smallint"
	TypeInteger  = "integer"
	TypeBigInt   = "bigint"
	TypeInt      = "int"  // integer.
	TypeInt2     = "int2" // smallint.
	TypeInt4     = "int4" // integer.
	TypeInt8     = "int8" // bigint.

	TypeCIDR     = "cidr"
	TypeInet     = "inet"
	TypeMACAddr  = "macaddr"
	TypeMACAddr8 = "macaddr8"

	TypeCircle     = "circle"
	TypeLine       = "line"
	TypeLseg       = "lseg"
	TypeBox        = "box"
	TypePath       = "path"
	TypePolygon    = "polygon"
	TypePoint      = "point"
	TypePgLSN      = "pg_lsn"
	TypePgSnapshot = "pg_snapshot"
	TypeTSQuery    = "tsquery"
	TypeTSVector   = "tsvector"

	TypeDate               = "date"
	TypeTime               = "time" // time without time zone
	TypeTimeWithTZ         = "time with time zone"
	TypeTimeWithoutTZ      = "time without time zone"
	TypeTimestamp          = "timestamp" // timestamp without time zone
	TypeTimestampTZ        = "timestamptz"
	TypeTimestampWithTZ    = "timestamp with time zone"
	TypeTimestampWithoutTZ = "timestamp without time zone"

	TypeDouble = "double precision"
	TypeReal   = "real"
	TypeFloat8 = "float8" // double precision
	TypeFloat4 = "float4" // real

	TypeNumeric = "numeric"
	TypeDecimal = "decimal" // numeric

	TypeSmallSerial = "smallserial" // smallint with auto_increment.
	TypeSerial      = "serial"      // integer with auto_increment.
	TypeBigSerial   = "bigserial"   // bigint with auto_increment.
	TypeSerial2     = "serial2"     // smallserial
	TypeSerial4     = "serial4"     // serial
	TypeSerial8     = "serial8"     // bigserial

	TypeArray       = "array"
	TypeXML         = "xml"
	TypeJSON        = "json"
	TypeJSONB       = "jsonb"
	TypeUUID        = "uuid"
	TypeMoney       = "money"
	TypeInterval    = "interval"
	TypeUserDefined = "user-defined"
)

// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// From: https://github.com/ariga/atlas/tree/v0.4.1/sql/postgres/inspect.go#823
var (
	// VariablesQuery 系统变量查询语句
	VariablesQuery = `SELECT setting FROM pg_settings WHERE name IN ('lc_collate', 'lc_ctype', 'server_version_num') ORDER BY name`

	// TableQuery 表查询语句
	TablesQuery = `
SELECT 
	t1.TABLE_NAME, 
	pg_catalog.obj_description(t2.oid, 'pg_class') AS COMMENT 
FROM 
	INFORMATION_SCHEMA.TABLES AS t1 
	JOIN pg_catalog.pg_class AS t2 
	ON t1.table_name = t2.relname 
WHERE 
	t1.TABLE_TYPE = 'BASE TABLE' 
	AND t1.table_schema = (CURRENT_SCHEMA())
`

	// ColumnsQuery 列查询语句
	ColumnsQueryFields = []string{
		"column_name",
		"data_type",
		"comment",
		"is_nullable",
		"column_default",
		"character_set_name",
		"collation_name",
		"numeric_precision",
		"numeric_scale",
		"character_maximum_length",
		"udt_name",
		"typtype",
		"oid",
	}

	ColumnsQuery = `
SELECT
	t1.column_name,
	t1.data_type,
	col_description(to_regclass("table_schema" || '.' || "table_name")::oid, "ordinal_position") AS comment,
	t1.is_nullable,
	t1.column_default,
	t1.character_set_name,
	t1.collation_name,
	t1.numeric_precision,
	t1.numeric_scale,
	t1.character_maximum_length,
	t1.udt_name,
	t2.typtype,
	t2.oid
FROM
	"information_schema"."columns" AS t1
	LEFT JOIN pg_catalog.pg_type AS t2
	ON t1.udt_name = t2.typname
WHERE
	t2.typtype != 'c'
	AND t1.TABLE_NAME = $1
	AND t1.table_schema = (CURRENT_SCHEMA())
`

	EnumQuery = "SELECT enumlabel FROM pg_enum WHERE enumtypid = $1"

	// IndexesQuery 索引查询语句
	IndexesQueryFields = []string{
		"index_name",
		"index_type",
		"column_name",
		"primary",
		"unique",
		"constraint_type",
		"predicate",
		"expression",
		"asc",
		"desc",
		"nulls_first",
		"nulls_last",
		"comment",
	}
	IndexesQuery = `
SELECT
	i.relname AS index_name,
	am.amname AS index_type,
	a.attname AS column_name,
	idx.indisprimary AS primary,
	idx.indisunique AS unique,
	c.contype AS constraint_type,
	pg_get_expr(idx.indpred, idx.indrelid) AS predicate,
	pg_get_expr(idx.indexprs, idx.indrelid) AS expression,
	pg_index_column_has_property(idx.indexrelid, a.attnum, 'asc') AS asc,
	pg_index_column_has_property(idx.indexrelid, a.attnum, 'desc') AS desc,
	pg_index_column_has_property(idx.indexrelid, a.attnum, 'nulls_first') AS nulls_first,
	pg_index_column_has_property(idx.indexrelid, a.attnum, 'nulls_last') AS nulls_last,
	obj_description(to_regclass(CURRENT_SCHEMA() || i.relname)::oid) AS comment
FROM
	pg_index idx
	JOIN pg_class i
	ON i.oid = idx.indexrelid
	LEFT JOIN pg_constraint c
	ON idx.indexrelid = c.conindid
	LEFT JOIN pg_attribute a
	ON a.attrelid = idx.indexrelid
	JOIN pg_am am
	ON am.oid = i.relam
WHERE
	idx.indrelid = to_regclass(CURRENT_SCHEMA() || '.' || $1)::oid
	AND COALESCE(c.contype, '') <> 'f'
ORDER BY
	index_name, a.attnum
`

	// ForeignKeysQuery 外键查询语句
	ForeignKeysQueryFields = []string{
		"constraint_name",
		"table_name",
		"column_name",
		"table_schema",
		"referenced_table_name",
		"referenced_column_name",
		"referenced_schema_name",
		"update_rule",
		"delete_rule",
	}
	ForeignKeysQuery = `
SELECT
	t1.constraint_name,
	t1.table_name,
	t2.column_name,
	t1.table_schema,
	t3.table_name AS referenced_table_name,
	t3.column_name AS referenced_column_name,
	t3.table_schema AS referenced_schema_name,
	t4.update_rule,
	t4.delete_rule
FROM
	information_schema.table_constraints t1
	JOIN information_schema.key_column_usage t2
	ON t1.constraint_name = t2.constraint_name
	AND t1.table_schema = t2.constraint_schema
	JOIN information_schema.constraint_column_usage t3
	ON t1.constraint_name = t3.constraint_name
	AND t1.table_schema = t3.constraint_schema
	JOIN information_schema.referential_constraints t4
	ON t1.constraint_name = t4.constraint_name
	AND t1.table_schema = t4.constraint_schema
WHERE
	t1.constraint_type = 'FOREIGN KEY'
	AND t1.table_name = $1
ORDER BY
	t1.constraint_name,
	t2.ordinal_position
`
)
