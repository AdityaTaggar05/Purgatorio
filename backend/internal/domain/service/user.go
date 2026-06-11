package service

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	var user model.User

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return nil
}
