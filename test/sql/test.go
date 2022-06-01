package testsql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mysql"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/ychengcloud/cre"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type option struct {
	dialect      string
	mysqlOptions mysqlOptions
}
type mysqlOptions struct {
	user     string `yaml:"user"`
	password string `yaml:"password"`
	host     string `yaml:"host"`
	port     int    `yaml:"port"`
	name     string `yaml:"name"`
	charset  string `yaml:"charset"`
}
type postgresOptions struct {
	user       string `yaml:"user"`
	password   string `yaml:"password"`
	host       string `yaml:"host"`
	port       int    `yaml:"port"`
	name       string `yaml:"name"`
	searchPath string `yaml:"search_path"`
}

var (
	defaultMysqlOptions = mysqlOptions{
		user:     "test",
		password: "test",
		host:     "localhost",
		port:     3306,
		name:     "test",
		charset:  "utf8mb4",
	}
	defaultPostgresOptions = postgresOptions{
		user:     "test",
		password: "test",
		host:     "localhost",
		port:     3306,
		name:     "test",
	}
)

func RunContainerForTest(ctx context.Context, dialect string, queriesFile string, version string) (*gnomock.Container, error) {
	var p gnomock.Preset
	switch dialect {
	case cre.MySQL:
		p = mysql.Preset(
			mysql.WithUser(defaultMysqlOptions.user, defaultMysqlOptions.password),
			mysql.WithDatabase(defaultMysqlOptions.name),
			mysql.WithQueriesFile(queriesFile),
			mysql.WithVersion(version),
		)
	case cre.Postgres:
		p = postgres.Preset(
			postgres.WithUser(defaultPostgresOptions.user, defaultPostgresOptions.password),
			postgres.WithDatabase(defaultPostgresOptions.name),
			postgres.WithQueriesFile(queriesFile),
			postgres.WithVersion(version),
		)
	default:
		return nil, fmt.Errorf("unsupported dialect: %v", dialect)
	}

	// path, _ := os.Getwd()
	// fmt.Println("Path:", path)
	return gnomock.Start(p)
	// return gnomock.Start(p, gnomock.WithDebugMode())

}

func ConnectForTest(ctx context.Context, dialect string, addr string) (db *sql.DB, err error) {
	address := strings.Split(addr, ":")
	port, _ := strconv.Atoi(address[1])

	var dsn string
	switch dialect {
	case cre.MySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", defaultMysqlOptions.user, defaultMysqlOptions.password, address[0], port, defaultMysqlOptions.name, defaultMysqlOptions.charset)
	case cre.Postgres:
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", defaultPostgresOptions.user, defaultPostgresOptions.password, address[0], port, defaultPostgresOptions.name)
		if defaultPostgresOptions.searchPath != "" {
			dsn += "&search_path=" + defaultPostgresOptions.searchPath
		}
	default:
		return nil, fmt.Errorf("unsupported dialect: %v", dialect)
	}
	db, err = sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}
	return
}
