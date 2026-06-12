package service

import (
	"context"
	"time"

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
	eco, err := s.UserRepo.GetEconomy(ctx, id)
	if err != nil {
		return eco, err
	}

	collectors, err := s.BaseRepo.GetResourceGenerationInfo(ctx, id)

	if err != nil {
		return eco, err
	}

	collectionTime := time.Now()
	lastCollectionDuration := int(collectionTime.Sub(eco.CollectorResetAt).Seconds())

	var collectedAmt int

	for _, collector := range collectors {
		var collection int

		if collector.Metadata.UpgradeEndsAt != nil {
			if collector.Metadata.UpgradeEndsAt.After(collectionTime) {
				continue
			}

			collection += int(collectionTime.Sub(*collector.Metadata.UpgradeEndsAt).Seconds()) * collector.CurrentRate
			collector.Metadata.UpgradeEndsAt = nil
		} else {
			collection += lastCollectionDuration * collector.CurrentRate
		}

		collectedAmt += min(collection, collector.StorageCapacity)
	}

	eco.Penitence += collectedAmt + eco.CollectorPendingPenitence

	if eco.Penitence > eco.MaxPenitence {
		eco.CollectorPendingPenitence = eco.Penitence - eco.MaxPenitence
		eco.Penitence = eco.MaxPenitence
	}

	s.UserRepo.UpdateEconomy(ctx, eco)
	s.BaseRepo.RemoveUpgradeInfo(ctx, id)

	return eco, nil
}
