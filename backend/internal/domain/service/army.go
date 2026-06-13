package service

import (
	"context"
	"fmt"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/google/uuid"
)

type ArmyService struct {
	ArmyRepo       repository.ArmyRepository
	UserRepo       repository.UserRepository
	BaseLayoutRepo repository.BaseLayoutRepository
}

func NewArmyService(armyRepo repository.ArmyRepository, userRepo repository.UserRepository, baseLayoutRepo repository.BaseLayoutRepository) *ArmyService {
	return &ArmyService{
		ArmyRepo:       armyRepo,
		UserRepo:       userRepo,
		BaseLayoutRepo: baseLayoutRepo,
	}
}

func (s *ArmyService) GetTroops(ctx context.Context) ([]model.Troop, error) {
	troops, err := s.ArmyRepo.GetAllTroops(ctx)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get troops"), err)
	}
	return troops, nil
}

