package service

import (
	"context"
	"errors"
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
	var (
		user   model.User
		tokens model.TokenPair
	)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return user, tokens, err
	}

	user, err = s.UserRepo.CreateUser(ctx, email, string(hash), username)
	if err != nil {
		return user, tokens, ErrUserAlreadyExists
	}

	tokens.AccessToken, err = model.GenerateJWT(user.ID, s.SigningKey, s.Config.AccessTTL)
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

func (s *AuthService) Login(ctx context.Context, email, password string) (model.User, model.TokenPair, error) {
	var (
		user   model.User
		tokens model.TokenPair
	)

	if password == "" {
		return user, tokens, ErrIncorrectPassword
	}

	user, err := s.UserRepo.GetAuthAndUserByEmail(ctx, email)
	if err != nil {
		return user, tokens, errors.Join(err, ErrUserNotFound)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return user, tokens, errors.Join(err, ErrIncorrectPassword)
	}

	tokens.AccessToken, err = model.GenerateJWT(user.ID, s.SigningKey, s.Config.AccessTTL)
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

func (s *AuthService) Logout(ctx context.Context, oldToken string) error {
	// TODO: Implement custom validation logic
	//if !IsValidRefreshToken(oldToken) {
	//return tokenservice.ErrInvalidRefreshTokenFormat
	//}

	return s.UserRepo.RevokeRefreshToken(ctx, oldToken)
}

func (s *AuthService) Refresh(ctx context.Context, oldToken string) (model.TokenPair, error) {
	// TODO: Implement custom validation logic
	//if !IsValidRefreshToken(oldToken) {
	//return model.TokenPair{}, ErrInvalidRefreshTokenFormat
	//}

	tokens := model.TokenPair{}

	rt, err := s.UserRepo.GetRefreshToken(ctx, oldToken)

	if err != nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return tokens, ErrInvalidRefreshToken
	}

	err = s.UserRepo.RevokeRefreshToken(ctx, oldToken)
	if err != nil {
		return tokens, err
	}

	refreshToken, err := model.GenerateRefreshToken(rt.UserID, s.Config.RefreshTTL)
	if err != nil {
		return tokens, err
	}
	tokens.RefreshToken = refreshToken.Token

	_ = s.UserRepo.CreateRefreshToken(ctx, rt.UserID, tokens.RefreshToken, time.Now().Add(s.Config.RefreshTTL))
	tokens.AccessToken, err = model.GenerateJWT(rt.UserID, s.SigningKey, s.Config.AccessTTL)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
