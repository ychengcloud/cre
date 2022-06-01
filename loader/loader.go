// Package loader is the interface for loading an external resource into cre
package loader

import (
	"context"
	"fmt"

	"github.com/ychengcloud/cre"
	"github.com/ychengcloud/cre/loader/sql"
	"github.com/ychengcloud/cre/loader/sql/mysql"
	"github.com/ychengcloud/cre/loader/sql/postgres"
	"github.com/ychengcloud/cre/spec"
)

type SQLLoader struct {
	driver cre.Driver
}

func NewSQLLoader(driver cre.Driver) *SQLLoader {
	return &SQLLoader{driver: driver}
}
func (l *SQLLoader) Load(ctx context.Context, name string) (*spec.Schema, error) {
	var i sql.Inspector
	switch l.driver.Dialect() {
	case cre.MySQL:
		i = mysql.NewInspector(l.driver)
	case cre.Postgres:
		i = postgres.NewInspector(l.driver)
	default:
		return nil, fmt.Errorf("load: unsupported dialect: %v", l.driver.Dialect())
	}
	schema, err := i.Inspect(ctx, name)
	if err != nil {
		return nil, err
	}

	specSchema, err := schema.Convert()
	if err != nil {
		return nil, err
	}

	return specSchema, nil
}

func (l *SQLLoader) Dialect() string {
	return l.driver.Dialect()
}

type LoaderOption func(*loaderOptions)

type loaderOptions struct {
	resources []string
}

func WithResources(resources ...string) LoaderOption {
	return func(l *loaderOptions) {
		l.resources = append(l.resources, resources...)
	}
}

func NewLoader(driver cre.Driver, opts ...LoaderOption) (cre.Loader, error) {
	lo := loaderOptions{}

	for _, apply := range opts {
		apply(&lo)
	}
	switch driver.Dialect() {
	case cre.MySQL, cre.Postgres:
		return NewSQLLoader(driver), nil
	default:
		return nil, fmt.Errorf("load/create: unsupported dialect: %v", driver.Dialect())
	}
}
