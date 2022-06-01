package sql

import (
	"context"
	"database/sql"

	"github.com/ychengcloud/cre"
)

var _ cre.Driver = (*Driver)(nil)

type Driver struct {
	dialect string
	db      *sql.DB
}

func Open(dialect, dsn string) (*Driver, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}
	return &Driver{dialect, db}, nil
}

func OpenDB(dialect string, db *sql.DB) *Driver {
	return &Driver{dialect, db}
}

func (d *Driver) Dialect() string { return d.dialect }

func (d *Driver) Close() error {
	return d.db.Close()
}

func (d *Driver) Exec(ctx context.Context, query string, args ...any) (any, error) {

	return d.db.ExecContext(ctx, query, args...)
}

// Query implements the cre.Query method.
func (d *Driver) Query(ctx context.Context, query string, args ...any) (any, error) {
	return d.db.QueryContext(ctx, query, args...)

}

// QueryRow implements the cre.QueryRow method.
func (d *Driver) QueryRow(ctx context.Context, query string, args ...any) (any, error) {
	return d.db.QueryRowContext(ctx, query, args...), nil

}
