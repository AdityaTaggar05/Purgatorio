package service

import (
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
)

type AuthService struct {
	Config config.JWTConfig
	UserRepo repository.UserRepository
	SigningKey *model.SigningKey
}

func NewAuthService(cfg config.JWTConfig, key *model.SigningKey, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		Config: cfg,
		UserRepo: userRepo,
		SigningKey: key,
	}
}
