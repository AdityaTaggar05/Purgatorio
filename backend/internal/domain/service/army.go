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

func (s *ArmyService) GetMyTroops(ctx context.Context, userID uuid.UUID) (*model.MyTroopsResponse, error) {
	userArmy, err := s.ArmyRepo.GetUserArmy(ctx, userID)
	if err != nil {
		return &model.MyTroopsResponse{
			Troops:       map[string]int{},
			UsedCapacity: 0,
			MaxCapacity:  s.getMaxCapacity(ctx, userID),
		}, nil
	}

	return &model.MyTroopsResponse{
		Troops:       userArmy.Troops,
		UsedCapacity: userArmy.UsedCapacity,
		MaxCapacity:  s.getMaxCapacity(ctx, userID),
	}, nil
}

func (s *ArmyService) DetrainTroops(ctx context.Context, userID uuid.UUID, troopID string, count int) error {
	if count <= 0 {
		return purgerr.Wrap(fmt.Errorf("count must be positive"), fmt.Errorf("count must be positive"))
	}

	troop, err := s.ArmyRepo.GetTroopByID(ctx, troopID)
	if err != nil {
		return purgerr.Wrap(ErrTroopNotFound, err)
	}

	userArmy, err := s.ArmyRepo.GetUserArmy(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrInsufficientTroops, err)
	}

	owned := userArmy.Troops[troopID]
	if owned < count {
		return purgerr.Wrap(ErrInsufficientTroops, fmt.Errorf("owned %d, tried to detrain %d", owned, count))
	}

	newUsed := userArmy.UsedCapacity - count*troop.Space
	if newUsed < 0 {
		newUsed = 0
	}

	// Refund 50% of training cost
	eco, err := s.UserRepo.GetEconomy(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrUserNotFound, err)
	}

	refund := (troop.TrainingCost * count) / 2
	eco.Penitence += refund
	if eco.Penitence > eco.MaxPenitence {
		eco.CollectorPendingPenitence += eco.Penitence - eco.MaxPenitence
		eco.Penitence = eco.MaxPenitence
	}

	if err := s.UserRepo.UpdateEconomy(ctx, eco); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to refund detrain cost"), err)
	}

	if err := s.ArmyRepo.RemoveTroops(ctx, userID, troopID, count, newUsed); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to remove troops"), err)
	}

	return nil
}

func (s *ArmyService) getMaxCapacity(ctx context.Context, userID uuid.UUID) int {
	maxCapacity := 0
	placed, err := s.BaseLayoutRepo.GetLayout(ctx, userID)
	if err != nil {
		return 0
	}
	for _, pb := range placed {
		if pb.BuildingID == "barracks" {
			stats, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, pb.BuildingID, pb.Level)
			if err == nil && stats.StorageCapacity != nil {
				maxCapacity += *stats.StorageCapacity
			}
		}
	}
	return maxCapacity
}

func (s *ArmyService) TrainTroops(ctx context.Context, userID uuid.UUID, troopID string, count int) error {
	if count <= 0 {
		return purgerr.Wrap(fmt.Errorf("count must be positive"), fmt.Errorf("count must be positive"))
	}

	troop, err := s.ArmyRepo.GetTroopByID(ctx, troopID)
	if err != nil {
		return purgerr.Wrap(ErrTroopNotFound, err)
	}

	userArmy, err := s.ArmyRepo.GetUserArmy(ctx, userID)
	if err != nil {
		userArmy = &model.UserArmy{UserID: userID, Troops: map[string]int{}, UsedCapacity: 0}
	}

	maxCapacity := s.getMaxCapacity(ctx, userID)

	newUsed := userArmy.UsedCapacity + count*troop.Space
	if newUsed > maxCapacity {
		return purgerr.Wrap(ErrInsufficientArmyCapacity, ErrInsufficientArmyCapacity)
	}

	eco, err := s.UserRepo.GetEconomy(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrUserNotFound, err)
	}

	totalCost := troop.TrainingCost * count
	if eco.Penitence < totalCost {
		return purgerr.Wrap(ErrInsufficientResources, ErrInsufficientResources)
	}

	eco.Penitence -= totalCost
	if err := s.UserRepo.UpdateEconomy(ctx, eco); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to deduct training cost"), err)
	}

	if err := s.ArmyRepo.AddTroops(ctx, userID, troopID, count, newUsed); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to add troops"), err)
	}

	return nil
}
