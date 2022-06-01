package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ychengcloud/cre"
	"github.com/ychengcloud/cre/gen"
	"github.com/ychengcloud/cre/loader"
	ldsql "github.com/ychengcloud/cre/loader/sql"
)

func Generate(cfg *gen.Config) error {
	var loaderInstance cre.Loader

	i := strings.Index(cfg.DSN, "://")
	if i == -1 {
		return fmt.Errorf("invalid data source name")
	}

	dialect := strings.TrimSpace(cfg.DSN[:i])
	dsn := cfg.DSN[i+3:]

	switch dialect {
	case gen.LoaderMysql, gen.LoaderPostgres:
		db, err := sql.Open(dialect, dsn)
		if err != nil {
			return err
		}
		drv := ldsql.OpenDB(dialect, db)
		loaderInstance, err = loader.NewLoader(drv)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported loader: %s", dialect)
	}

	g, err := gen.NewGenerator(cfg, loaderInstance)
	if err != nil {
		return err
	}
	return g.Generate(context.Background())
}
