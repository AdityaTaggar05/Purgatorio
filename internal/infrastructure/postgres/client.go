package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(logger *slog.Logger, ctx context.Context, cfg config.PostgresConfig) *pgxpool.Pool {
	pgCfg, err := pgxpool.ParseConfig(cfg.URL)

	if err != nil {
		logger.Error(fmt.Sprintf("unable to parse config for database: %v", err))
	}

	pgCfg.MaxConns = int32(cfg.MaxOpenConns)
	pgCfg.MinConns = 5
	pgCfg.MaxConnLifetime = time.Hour

	DB, err := pgxpool.NewWithConfig(ctx, pgCfg)

	if err != nil {
		logger.Error(fmt.Sprintf("unable to connect to database: %v", err))
	}

	err = DB.Ping(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("could not ping the database: %v", err))
	}
	logger.Debug("connected to the database!")

	return DB
}
