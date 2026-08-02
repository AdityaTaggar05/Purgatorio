package service

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/internal/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const sinDrainPerMinute = 1.0 / 6.0 // 10% per hour; 100 → 0 in ~10 hours

type UserService struct {
	UserRepo repository.UserRepository
	BaseRepo repository.BaseRepository
	DB       *pgxpool.Pool
}

func NewUserService(userRepo repository.UserRepository, baseRepo repository.BaseRepository, db *pgxpool.Pool) *UserService {
	return &UserService{
		UserRepo: userRepo,
		BaseRepo: baseRepo,
		DB:       db,
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

func (s *UserService) GetCombat(ctx context.Context, id uuid.UUID) (model.UserCombat, error) {
	combat, err := s.UserRepo.GetCombat(ctx, id)
	if err != nil {
		return combat, err
	}

	if combat.UpdatedAt != nil && combat.SinMeter > 0 {
		elapsed := time.Since(*combat.UpdatedAt)
		drained := int(elapsed.Minutes() * sinDrainPerMinute)
		if drained > 0 {
			combat.SinMeter = max(0, combat.SinMeter-drained)
			_ = s.UserRepo.UpdateCombat(ctx, id, combat.SinMeter)
		}
	}

	return combat, nil
}

func (s *UserService) EconomyCollect(ctx context.Context, id uuid.UUID) (model.UserEconomy, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.UserEconomy{}, err
	}
	defer tx.Rollback(ctx)

	txCtx := postgres.CtxWithTx(ctx, tx)

	eco, err := s.UserRepo.GetEconomyForUpdate(txCtx, id)
	if err != nil {
		return eco, err
	}

	collectors, err := s.BaseRepo.GetResourceGenerationInfo(txCtx, id)
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
	eco.CollectorPendingPenitence = 0

	if eco.Penitence > eco.MaxPenitence {
		eco.CollectorPendingPenitence = eco.Penitence - eco.MaxPenitence
		eco.Penitence = eco.MaxPenitence
	}

	eco.CollectorResetAt = collectionTime
	if err := s.UserRepo.UpdateEconomy(txCtx, eco); err != nil {
		return eco, err
	}
	if err := s.BaseRepo.RemoveUpgradeInfo(txCtx, id, model.BuildingResource); err != nil {
		return eco, err
	}

	if err := tx.Commit(ctx); err != nil {
		return eco, err
	}

	return eco, nil
}
