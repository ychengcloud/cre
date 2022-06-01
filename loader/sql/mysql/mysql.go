package mysql

import (
	"strings"
)

// MySQL data types
// https://dev.mysql.com/doc/refman/8.0/en/data-types.html
const (
	// MySQL numeric data types
	TypeBit       = "bit"
	TypeInt       = "int"
	TypeTinyInt   = "tinyint"
	TypeSmallInt  = "smallint"
	TypeMediumInt = "mediumint"
	TypeBigInt    = "bigint"
	TypeDecimal   = "decimal"
	TypeNumeric   = "numeric"
	TypeFloat     = "float"
	TypeDouble    = "double"
	TypeReal      = "real"

	// MySQL time data types
	TypeTimestamp = "timestamp"
	TypeDate      = "date"
	TypeTime      = "time"
	TypeDateTime  = "datetime"
	TypeYear      = "year"

	// MySQL string data types
	TypeChar       = "char"
	TypeVarchar    = "varchar"
	TypeBinary     = "binary"
	TypeVarBinary  = "varbinary"
	TypeTinyBlob   = "tinyblob"
	TypeTinyText   = "tinytext"
	TypeBlob       = "blob"
	TypeText       = "text"
	TypeMediumBlob = "mediumblob"
	TypeMediumText = "mediumtext"
	TypeLongBlob   = "longblob"
	TypeLongText   = "longtext"
	TypeEnum       = "enum"
	TypeSet        = "set"
	TypeJSON       = "json"

	// MySQL spatial data types
	TypeGeometry           = "geometry"
	TypePoint              = "point"
	TypeMultiPoint         = "multipoint"
	TypeLineString         = "linestring"
	TypeMultiLineString    = "multilinestring"
	TypePolygon            = "polygon"
	TypeMultiPolygon       = "multipolygon"
	TypeGeoCollection      = "geomcollection"
	TypeGeometryCollection = "geometrycollection"
)

// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// From: https://github.com/ariga/atlas/tree/v0.4.1/sql/mysql/inspect.go#646
var (

	// VariablesQuery 系统变量查询语句
	VariablesQueryFields = []string{"@@version", "@@collation_server", "@@character_set_server"}
	VariablesQuery       = "SELECT " + strings.Join(VariablesQueryFields, ",")

	// TableQuery 表查询语句
	TablesQueryFields = []string{"t1.table_name", "t2.character_set_name", "t1.table_collation", "t1.auto_increment", "t1.table_comment", "t1.create_options"}
	TablesQuery       = "SELECT " + strings.Join(TablesQueryFields, ",") + " FROM information_schema.tables AS t1 JOIN information_schema.collations AS t2 ON t1.table_collation = t2.collation_name WHERE table_schema = ?"

	// ColumnsQuery 列查询语句
	ColumnsQueryFields = []string{"column_name", "column_type", "column_comment", "is_nullable", "column_key", "column_default", "extra", "character_set_name", "collation_name", "numeric_precision", "numeric_scale"}
	ColumnsQuery       = "SELECT " + strings.Join(ColumnsQueryFields, ",") + " FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position"

	// IndexesQuery 索引查询语句
	IndexesQueryFields = []string{"index_name", "column_name", "non_unique", "seq_in_index", "index_type", "collation", "index_comment", "sub_part", "expression"}
	IndexesQuery       = `
SELECT 
	index_name,
	column_name,
	non_unique,
	seq_in_index,
	index_type,
	collation,
	index_comment,
	sub_part,
 	null AS expression 
FROM 
	information_schema.statistics 
WHERE 
	table_schema = ? 
	AND table_name = ? 
ORDER BY 
	index_name, seq_in_index
`
	IndexesExprQuery = `
SELECT 
	index_name,
	column_name,
	non_unique,
	seq_in_index,
	index_type,
	collation,
	index_comment,
	sub_part,
	expression 
FROM 
	information_schema.statistics 
WHERE 
	table_schema = ? 
	AND table_name = ? 
ORDER BY 
	index_name, seq_in_index
`

	// ForeignKeysQuery 外键查询语句
	ForeignKeysQueryFields = []string{"t1.constraint_name", "t1.table_name", "t1.column_name", "t1.table_schema", "t1.referenced_table_name", "t1.referenced_column_name", "t1.referenced_table_schema", "t3.update_rule", "t3.delete_rule"}
	ForeignKeysQuery       = `
SELECT ` + strings.Join(ForeignKeysQueryFields, ",") +
		` FROM
	information_schema.key_column_usage AS t1
	JOIN information_schema.table_constraints AS t2
	JOIN information_schema.referential_constraints AS t3
	ON t1.constraint_name = t2.constraint_name
	AND t1.constraint_name = t3.constraint_name
	AND t1.table_schema = t2.table_schema
	AND t1.table_schema = t3.constraint_schema
WHERE
	t2.constraint_type = 'FOREIGN KEY'
	AND t1.table_schema = ?
	AND t1.table_name = ?
ORDER BY
	t1.constraint_name,
	t1.ordinal_position`
)
