package repository

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
)

type UserRepository interface {
	CreateUser(context.Context, string, string, string) (model.User, error)
	CreateRefreshToken(context.Context, string, string, time.Time) error
}
