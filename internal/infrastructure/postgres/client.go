package postgres

import (
	"context"
	"log"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg config.PostgresConfig) *pgxpool.Pool {
	pgCfg, err := pgxpool.ParseConfig(cfg.URL)

	if err != nil {
		log.Fatal("[ERR] Unable to parse config for database: ", err)
	}

	pgCfg.MaxConns = int32(cfg.MaxOpenConns)
	pgCfg.MinConns = 5
	pgCfg.MaxConnLifetime = time.Hour

	DB, err := pgxpool.NewWithConfig(ctx, pgCfg)

	if err != nil {
		log.Fatal("[ERR] Unable to connect to database: ", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatal("[ERR] Could not ping the database: ", err)
	}
	log.Println("[DEBUG] Connected to the database!")

	return DB
}
