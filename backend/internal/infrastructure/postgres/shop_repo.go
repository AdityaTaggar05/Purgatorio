package postgres

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopRepository struct {
	DB *pgxpool.Pool
}

func NewShopRepository(db *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{DB: db}
}

