package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Config     config.JWTConfig
	UserRepo   repository.UserRepository
	SigningKey *model.SigningKey
}

func NewAuthService(cfg config.JWTConfig, key *model.SigningKey, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		Config:     cfg,
		UserRepo:   userRepo,
		SigningKey: key,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (model.User, model.TokenPair, error) {
	var user model.User
	var tokens model.TokenPair

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return user, tokens, err
	}

	user, err = s.UserRepo.CreateUser(ctx, email, string(hash), username)
	if err != nil {
		fmt.Println(err)
		return user, tokens, ErrUserAlreadyExists
	}

	tokens.AccessToken, err = model.GenerateJWT(user, s.SigningKey, s.Config.AccessTTL)
	if err != nil {
		return user, tokens, err
	}

	refreshToken, err := model.GenerateRefreshToken(user.ID, s.Config.RefreshTTL)
	if err != nil {
		return user, tokens, err
	}
	tokens.RefreshToken = refreshToken.Token

	err = s.UserRepo.CreateRefreshToken(ctx, user.ID, tokens.RefreshToken, time.Now().Add(s.Config.RefreshTTL))
	if err != nil {
		return user, tokens, err
	}

	return user, tokens, nil
}
