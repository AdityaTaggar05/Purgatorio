package repository

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	// Auth Functions
	CreateUser(ctx context.Context, email string, hash string, username string) (model.User, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Time) error
	GetAuthAndUserByEmail(ctx context.Context, email string) (model.AuthAndUser, error)
	GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error

	// User Functions
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Economy Functions
	GetEconomy(ctx context.Context, id uuid.UUID) (model.UserEconomy, error)
	EconomyCollect(ctx context.Context, id uuid.UUID) (model.UserEconomy, error)
}
