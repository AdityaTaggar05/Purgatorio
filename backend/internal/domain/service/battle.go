package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/google/uuid"
)

type BattleService struct {
	BattleRepo    repository.BattleRepository
	UserRepo      repository.UserRepository
	ArmyRepo      repository.ArmyRepository
	BaseLayoutRepo repository.BaseLayoutRepository
	ShopRepo      repository.ShopRepository
	Catalog       engine.TroopCatalog
}

func NewBattleService(
	battleRepo repository.BattleRepository,
	userRepo repository.UserRepository,
	armyRepo repository.ArmyRepository,
	baseLayoutRepo repository.BaseLayoutRepository,
	shopRepo repository.ShopRepository,
) *BattleService {
	catalog := engine.TroopCatalog{}
	troopList, _ := armyRepo.GetAllTroops(context.Background())
	for _, t := range troopList {
		catalog[t.ID] = engine.TroopStats{
			ID:              t.ID,
			HP:              t.HP,
			DPS:             t.DPS,
			Speed:           t.Speed,
			Range:           t.AttackRange,
			PreferredTarget: t.PreferredTarget,
		}
	}

	return &BattleService{
		BattleRepo:    battleRepo,
		UserRepo:      userRepo,
		ArmyRepo:      armyRepo,
		BaseLayoutRepo: baseLayoutRepo,
		ShopRepo:      shopRepo,
		Catalog:       catalog,
	}
}

func (s *BattleService) GetMatchList(ctx context.Context, userID uuid.UUID) ([]model.MatchListEntry, error) {
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(ErrUserNotFound, err)
	}
	return s.BattleRepo.GetMatchList(ctx, user.TerraceLevel, userID)
}

