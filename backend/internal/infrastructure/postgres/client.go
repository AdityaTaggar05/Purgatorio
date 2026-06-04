package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.URL)

	if err != nil {
		return nil, fmt.Errorf("unable to parse config for database: %v", err)
	}

	pgCfg.MaxConns = int32(cfg.MaxOpenConns)
	pgCfg.MinConns = 5
	pgCfg.MaxConnLifetime = time.Hour

	DB, err := pgxpool.NewWithConfig(ctx, pgCfg)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %v", err)
	}

	return DB, nil
}
