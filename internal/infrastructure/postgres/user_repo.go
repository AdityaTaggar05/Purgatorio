package postgres

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, hash string) (model.User, error) {
	var user model.User
	return user, nil
}

func (r *UserRepository) CreateRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	return nil
}
