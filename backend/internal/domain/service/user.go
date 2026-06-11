package service

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepo repository.UserRepository
	BaseRepo repository.BaseRepository
}

func NewUserService(userRepo repository.UserRepository, baseRepo repository.BaseRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		BaseRepo: baseRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.UserRepo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.UserRepo.DeleteUser(ctx, id)
}

func (s *UserService) GetEconomy(ctx context.Context, id uuid.UUID) (model.UserEconomy, error) {
	return s.UserRepo.GetEconomy(ctx, id)
}

func (s *UserService) EconomyCollect(ctx context.Context, id uuid.UUID) (model.UserEconomy, error) {
	return model.UserEconomy{}, nil
}
