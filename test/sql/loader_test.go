package testsql

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"

	"github.com/ychengcloud/cre"
	"github.com/ychengcloud/cre/loader"
	"github.com/ychengcloud/cre/loader/sql"
	"github.com/ychengcloud/cre/spec"
)

const (
	deadlinePerTest                 = time.Duration(5 * time.Second)
	deadlineOnStartContanerForTests = time.Duration(60 * time.Second)
)

type driver struct {
	*sql.Driver

	// gnomock 支持的版本
	dialect  string
	version  string
	expected *spec.Schema
}

var (
	drivers   []*driver
	queryFile map[string]string = map[string]string{
		cre.MySQL:    "./mysql.sql",
		cre.Postgres: "./postgres.sql",
		cre.SQLite:   "./sqlite.sql",
	}
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))

}

func testMainWrapper(m *testing.M) int {
	ctx, cancelFn := context.WithTimeout(context.Background(), deadlineOnStartContanerForTests)
	defer cancelFn()

	pgNumericTable := &spec.Table{Name: "numeric"}
	pgStringTable := &spec.Table{Name: "string"}
	pgTimeTable := &spec.Table{Name: "time"}
	pgBinaryTable := &spec.Table{Name: "binary"}
	pgSpatialable := &spec.Table{Name: "spatial"}
	pgEnumTable := &spec.Table{Name: "enum"}
	pgFK1Table := &spec.Table{Name: "fk1"}
	pgFK2Table := &spec.Table{Name: "fk2"}

	pgExpected := &spec.Schema{Name: "test"}
	pgExpected.AddTables(pgNumericTable, pgStringTable, pgTimeTable, pgBinaryTable, pgSpatialable, pgEnumTable, pgFK1Table, pgFK2Table)

	pgNumericTable.AddFields([]*spec.Field{
		{Name: "smallint", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Nullable: false, Ops: spec.NumericOps},
		{Name: "smallint1", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Nullable: true, Ops: spec.NumericOps},
		{Name: "smallint2", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Comment: "smallint comment", Nullable: true, Index: true, Unique: true, Filterable: true, Ops: spec.NumericOps},
		{Name: "smallint3", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Nullable: true, Ops: spec.NumericOps},
		{Name: "integer", Type: &spec.IntegerType{Name: "integer", Size: 32}, Nullable: true, Ops: spec.NumericOps},
		{Name: "bigint", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Nullable: true, Ops: spec.NumericOps},
		{Name: "boolean", Type: &spec.BoolType{Name: "boolean"}, Nullable: true, Ops: spec.BoolOps},
		{Name: "numeric", Type: &spec.FloatType{Name: "numeric"}, Nullable: true, Ops: spec.NumericOps},
		{Name: "real", Type: &spec.FloatType{Name: "real", Precision: 24}, Nullable: true, Ops: spec.NumericOps},
		{Name: "double precision", Type: &spec.FloatType{Name: "double precision", Precision: 53}, Nullable: true, Ops: spec.NumericOps},
		{Name: "money", Type: &spec.FloatType{Name: "money"}, Nullable: true, Ops: spec.NumericOps},
	}...)

	pgStringTable.AddFields([]*spec.Field{
		{Name: "bit", Type: &spec.BitType{Name: "bit", Len: 1}},
		{Name: "bit varying", Type: &spec.BitType{Name: "bit varying", Len: 255}},
		{Name: "character", Type: &spec.StringType{Name: "character", Size: 1}, Ops: spec.StringOps},
		{Name: "character varying", Type: &spec.StringType{Name: "character varying"}, Ops: spec.StringOps},
		{Name: "character varying255", Type: &spec.StringType{Name: "character varying", Size: 255}, Ops: spec.StringOps},
		{Name: "text", Type: &spec.StringType{Name: "text"}, Ops: spec.StringOps},

		{Name: "tsquery", Type: &spec.StringType{Name: "tsquery"}, Ops: spec.StringOps},
		{Name: "tsvector", Type: &spec.StringType{Name: "tsvector"}, Ops: spec.StringOps},
		{Name: "uuid", Type: &spec.StringType{Name: "uuid"}, Ops: spec.StringOps},
		{Name: "xml", Type: &spec.StringType{Name: "xml"}, Ops: spec.StringOps},
		{Name: "json", Type: &spec.JSONType{Name: "json"}},
		{Name: "jsonb", Type: &spec.JSONType{Name: "jsonb"}},
	}...)

	pgTimeTable.AddFields([]*spec.Field{
		{Name: "date", Type: &spec.TimeType{Name: "date"}, Ops: spec.NumericOps},
		{Name: "time", Type: &spec.TimeType{Name: "time without time zone"}, Ops: spec.NumericOps},
		{Name: "timestamp", Type: &spec.TimeType{Name: "timestamp without time zone"}, Ops: spec.NumericOps},
		{Name: "timestamptz", Type: &spec.TimeType{Name: "timestamp with time zone"}, Ops: spec.NumericOps},
		{Name: "interval", Type: &spec.TimeType{Name: "interval"}, Ops: spec.NumericOps},
	}...)

	pgBinaryTable.AddFields([]*spec.Field{
		{Name: "bytea", Type: &spec.BinaryType{Name: "bytea"}},
	}...)

	pgSpatialable.AddFields([]*spec.Field{
		{Name: "cidr", Type: &spec.SpatialType{Name: "cidr"}},
		{Name: "inet", Type: &spec.SpatialType{Name: "inet"}},
		{Name: "macaddr", Type: &spec.SpatialType{Name: "macaddr"}},
		{Name: "box", Type: &spec.SpatialType{Name: "box"}},
		{Name: "circle", Type: &spec.SpatialType{Name: "circle"}},
		{Name: "line", Type: &spec.SpatialType{Name: "line"}},
		{Name: "lseg", Type: &spec.SpatialType{Name: "lseg"}},
		{Name: "path", Type: &spec.SpatialType{Name: "path"}},
		{Name: "point", Type: &spec.SpatialType{Name: "point"}},
		{Name: "polygon", Type: &spec.SpatialType{Name: "polygon"}},
	}...)

	pgEnumTable.AddFields([]*spec.Field{
		{Name: "enum", Type: &spec.EnumType{Name: "mood", Values: []string{"sad", "ok", "happy"}}, Nullable: true, Ops: spec.EnumOps},
	}...)

	pgFK1Table.AddFields([]*spec.Field{
		{Name: "id", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Nullable: false, PrimaryKey: true, Index: true, Unique: true, Filterable: true, Sortable: true, Ops: spec.NumericOps},
		{Name: "fkid", Type: &spec.IntegerType{Name: "bigint", Size: 64}, ForeignKey: true, Ops: spec.NumericOps},
	}...)
	pgFK1Table.ID = pgFK1Table.GetField("id")

	pgFK2Table.AddFields([]*spec.Field{
		{Name: "id", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Nullable: false, PrimaryKey: true, Index: true, Unique: true, Filterable: true, Sortable: true, Ops: spec.NumericOps},
	}...)
	pgFK2Table.ID = pgFK2Table.GetField("id")

	mysqlNumericTable := &spec.Table{Name: "numeric"}
	mysqlStringTable := &spec.Table{Name: "string"}
	mysqlTimeTable := &spec.Table{Name: "time"}
	mysqlSpatialTable := &spec.Table{Name: "spatial"}
	mysqlFK1Table := &spec.Table{Name: "fk1"}
	mysqlFK2Table := &spec.Table{Name: "fk2"}

	mysqlExpected := &spec.Schema{Name: "test"}

	mysqlExpected.AddTables(mysqlNumericTable, mysqlStringTable, mysqlTimeTable, mysqlSpatialTable, mysqlFK1Table, mysqlFK2Table)

	mysqlNumericTable.AddFields([]*spec.Field{
		{Name: "bigint", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Index: true, Nullable: false, Comment: "中文bigint comment", PrimaryKey: true, Unique: true, Filterable: true, Sortable: true, Ops: spec.NumericOps},
		{Name: "bigint1", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Ops: spec.NumericOps},
		{Name: "bit", Type: &spec.BitType{Name: "bit", Len: 1}},
		{Name: "int", Type: &spec.IntegerType{Name: "int", Size: 32}, Ops: spec.NumericOps},
		{Name: "tinyint", Type: &spec.IntegerType{Name: "tinyint", Size: 8}, Ops: spec.NumericOps},
		{Name: "smallint", Type: &spec.IntegerType{Name: "smallint", Size: 16}, Ops: spec.NumericOps},
		{Name: "mediumint", Type: &spec.IntegerType{Name: "mediumint", Size: 32}, Ops: spec.NumericOps},
		{Name: "decimal", Type: &spec.FloatType{Name: "decimal", Precision: 10, Scale: 0}, Ops: spec.NumericOps},
		{Name: "decimal1", Type: &spec.FloatType{Name: "decimal", Precision: 10, Scale: 2}, Ops: spec.NumericOps},
		{Name: "numeric", Type: &spec.FloatType{Name: "decimal", Precision: 10, Scale: 0}, Ops: spec.NumericOps},
		{Name: "float", Type: &spec.FloatType{Name: "float"}, Ops: spec.NumericOps},
		{Name: "float1", Type: &spec.FloatType{Name: "float", Precision: 10, Scale: 2}, Ops: spec.NumericOps},
		{Name: "double", Type: &spec.FloatType{Name: "double"}, Ops: spec.NumericOps},
		{Name: "double1", Type: &spec.FloatType{Name: "double", Precision: 10, Scale: 2}, Ops: spec.NumericOps},
		{Name: "real", Type: &spec.FloatType{Name: "double"}, Ops: spec.NumericOps},
		{Name: "real1", Type: &spec.FloatType{Name: "double", Precision: 10, Scale: 2}, Ops: spec.NumericOps},
	}...)

	mysqlStringTable.AddFields([]*spec.Field{
		{Name: "char", Type: &spec.StringType{Name: "char", Size: 1}, Ops: spec.StringOps},
		{Name: "char1", Type: &spec.StringType{Name: "char", Size: 10}, Ops: spec.StringOps},
		{Name: "varchar", Type: &spec.StringType{Name: "varchar", Size: 255}, Ops: spec.StringOps},
		{Name: "binary", Type: &spec.StringType{Name: "binary", Size: 1}, Ops: spec.StringOps},
		{Name: "varbinary", Type: &spec.StringType{Name: "varbinary", Size: 255}, Ops: spec.StringOps},
		{Name: "tinyblob", Type: &spec.StringType{Name: "tinyblob"}, Ops: spec.StringOps},
		{Name: "tinytext", Type: &spec.StringType{Name: "tinytext"}, Ops: spec.StringOps},
		{Name: "blob", Type: &spec.StringType{Name: "blob"}, Ops: spec.StringOps},
		{Name: "text", Type: &spec.StringType{Name: "text"}, Ops: spec.StringOps},
		{Name: "mediumblob", Type: &spec.StringType{Name: "mediumblob"}, Ops: spec.StringOps},
		{Name: "mediumtext", Type: &spec.StringType{Name: "mediumtext"}, Ops: spec.StringOps},
		{Name: "longblob", Type: &spec.StringType{Name: "longblob"}, Ops: spec.StringOps},
		{Name: "longtext", Type: &spec.StringType{Name: "longtext"}, Ops: spec.StringOps},
		{Name: "enum", Type: &spec.EnumType{Name: "enum", Values: []string{"a", "b", "c"}}, Ops: spec.EnumOps},
		{Name: "set", Type: &spec.EnumType{Name: "set", Values: []string{"a", "b", "c"}}, Ops: spec.EnumOps},
		{Name: "json", Type: &spec.JSONType{Name: "json"}},
	}...)

	mysqlTimeTable.AddFields([]*spec.Field{
		{Name: "date", Type: &spec.TimeType{Name: "date"}, Ops: spec.NumericOps},
		{Name: "time", Type: &spec.TimeType{Name: "time"}, Ops: spec.NumericOps},
		{Name: "timestamp", Type: &spec.TimeType{Name: "timestamp"}, Ops: spec.NumericOps},
		{Name: "datetime", Type: &spec.TimeType{Name: "datetime"}, Ops: spec.NumericOps},
		{Name: "year", Type: &spec.TimeType{Name: "year"}, Ops: spec.NumericOps},
	}...)

	mysqlSpatialTable.AddFields([]*spec.Field{
		{Name: "geometry", Type: &spec.SpatialType{Name: "geometry"}},
		{Name: "point", Type: &spec.SpatialType{Name: "point"}},
		{Name: "multipoint", Type: &spec.SpatialType{Name: "multipoint"}},
		{Name: "linestring", Type: &spec.SpatialType{Name: "linestring"}},
		{Name: "multilinestring", Type: &spec.SpatialType{Name: "multilinestring"}},
		{Name: "polygon", Type: &spec.SpatialType{Name: "polygon"}},
		{Name: "multipolygon", Type: &spec.SpatialType{Name: "multipolygon"}},
		{Name: "geomcollection", Type: &spec.SpatialType{Name: "geomcollection"}},
		{Name: "geometrycollection", Type: &spec.SpatialType{Name: "geometrycollection"}},
	}...)

	mysqlFK1Table.AddFields([]*spec.Field{
		{Name: "id", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Index: true, Nullable: false, PrimaryKey: true, Unique: true, Filterable: true, Sortable: true, Ops: spec.NumericOps},
		{Name: "fkid", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Index: true, ForeignKey: true, Filterable: true, Ops: spec.NumericOps},
	}...)
	mysqlFK1Table.ID = mysqlFK1Table.GetField("id")

	mysqlFK2Table.AddFields([]*spec.Field{
		{Name: "id", Type: &spec.IntegerType{Name: "bigint", Size: 64}, Index: true, Nullable: false, PrimaryKey: true, Unique: true, Filterable: true, Sortable: true, Ops: spec.NumericOps},
	}...)
	mysqlFK2Table.ID = mysqlFK2Table.GetField("id")

	drivers = []*driver{
		// {
		// 	dialect:  cre.MySQL,
		// 	version:  "5.7.32",
		// 	expected: mysqlExpected,
		// },
		// {
		// 	dialect:  cre.MySQL,
		// 	version:  "8.0.22",
		// 	expected: mysqlExpected,
		// },
		{
			dialect:  cre.Postgres,
			version:  "13.1",
			expected: pgExpected,
		},
		// {
		// 	dialect:  cre.Postgres,
		// 	version:  "12.5",
		// 	expected: pgExpected,
		// },
		// {
		// 	dialect:  cre.Postgres,
		// 	version:  "11.10",
		// 	expected: pgExpected,
		// },
		// {
		// 	dialect:  cre.Postgres,
		// 	version:  "10.15",
		// 	expected: pgExpected,
		// },
	}

	for _, driver := range drivers {
		log.Printf("create test container %s %s", driver.dialect, driver.version)
		testContainer, err := RunContainerForTest(ctx, driver.dialect, queryFile[driver.dialect], driver.version)
		if err != nil {
			log.Printf("Failed to create test container: %s", err)
			return 1
		}

		defer func() {
			err = gnomock.Stop(testContainer)
			if err != nil {
				log.Printf("Failed to Stop container: %s", err)
			}
		}()

		db, err := ConnectForTest(ctx, driver.dialect, testContainer.DefaultAddress())
		if err != nil {
			log.Fatal(err)
		}

		driver.Driver = sql.OpenDB(driver.dialect, db)
		log.Printf("created test container %s %s", driver.dialect, driver.version)

	}

	return m.Run()
}

func TestLoader(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), deadlineOnStartContanerForTests)
	defer cancelFn()

	for _, drv := range drivers {
		t.Run(drv.dialect+" "+drv.version, func(t *testing.T) {
			l, err := loader.NewLoader(drv.Driver)
			require.NoError(t, err)

			s, err := l.Load(ctx, "test")
			require.NoError(t, err)

			require.NoError(t, err)
			require.NotNil(t, s)
			matchSchema(t, drv.expected, s)
			// require.EqualValues(t, drv.expected, s)
		})

	}
}

func matchSchema(t *testing.T, expected *spec.Schema, s *spec.Schema) {
	require.Equal(t, expected.Name, s.Name)

	for _, table := range s.Tables() {
		expectedTable := expected.Table(table.Name)
		require.NotNil(t, expectedTable)

		matchTable(t, expectedTable, table)
	}
}

func matchTable(t *testing.T, expected *spec.Table, table *spec.Table) {
	require.Equal(t, expected.Name, table.Name)
	require.Equal(t, expected.JoinTable, table.JoinTable)
	// require.EqualValues(t, expected.ID, table.ID)

	for _, field := range table.Fields() {
		expectedField := expected.GetField(field.Name)
		require.NotNil(t, expectedField)
		matchField(t, expectedField, field)
	}
}

func matchField(t *testing.T, expected *spec.Field, field *spec.Field) {
	require.Equal(t, expected.Name, field.Name, "Name", field)
	require.Equal(t, expected.Type, field.Type, "Type", field)
	require.Equal(t, expected.Nullable, field.Nullable, "Nullable", field)
	require.Equal(t, expected.Optional, field.Optional, "Optional", field)
	require.Equal(t, expected.Sensitive, field.Sensitive, "Sensitive", field)
	require.Equal(t, expected.Tag, field.Tag, "Tag", field)
	require.Equal(t, expected.Comment, field.Comment, "Comment", field)
	require.Equal(t, expected.Alias, field.Alias, "Alias", field)
	require.Equal(t, expected.Sortable, field.Sortable, "Sortable", field)
	require.Equal(t, expected.Filterable, field.Filterable, "Filterable", field)
	require.Equal(t, expected.ForeignKey, field.ForeignKey, "ForeignKey", field)
	require.Equal(t, expected.PrimaryKey, field.PrimaryKey, "PrimaryKey", field)
	require.Equal(t, expected.Index, field.Index, "Index", field)
	require.Equal(t, expected.Unique, field.Unique, "Unique", field)
	require.EqualValues(t, expected.Rel, field.Rel, "Rel", field)
	require.EqualValues(t, expected.Ops, field.Ops, "Ops", field)
	require.EqualValues(t, expected.Attrs, field.Attrs, "Attrs", field)
	require.Equal(t, expected.Table.Name, field.Table.Name, "Name", field)
}
