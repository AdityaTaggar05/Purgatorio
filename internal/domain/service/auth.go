package service

import (
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
)

type AuthService struct {
	Config config.JWTConfig
	Repo repository.UserRepository
}

func NewAuthService(cfg config.JWTConfig, repo repository.UserRepository) *AuthService {
	return &AuthService{
		Config: cfg,
		Repo: repo,
	}
}
