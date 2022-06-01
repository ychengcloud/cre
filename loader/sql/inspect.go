package sql

import (
	"context"
)

type Inspector interface {
	Inspect(ctx context.Context, name string) (*Schema, error)
}
