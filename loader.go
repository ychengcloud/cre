package cre

import (
	"context"

	"github.com/ychengcloud/cre/spec"
)

// Dialect names.
const (
	Fake     = "fake"
	MySQL    = "mysql"
	SQLite   = "sqlite3"
	Postgres = "postgres"
)

// ExecQuerier wraps the database operations.
type ExecQuerier interface {
	Exec(ctx context.Context, query string, args ...any) (any, error)
	Query(ctx context.Context, query string, args ...any) (any, error)
	QueryRow(ctx context.Context, query string, args ...any) (any, error)
}

type Driver interface {
	ExecQuerier

	Close() error

	Dialect() string
}

type Loader interface {
	// Load loads the schema from the database.
	Load(ctx context.Context, name string) (*spec.Schema, error)
	Dialect() string
}
