package service

import "github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}
