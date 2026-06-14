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

func (s *BattleService) InitiateBattle(ctx context.Context, attackerID, defenderID uuid.UUID) (*model.InitiateResponse, error) {
	if attackerID == defenderID {
		return nil, purgerr.Wrap(ErrCannotAttackSelf, ErrCannotAttackSelf)
	}

	attacker, err := s.UserRepo.GetUserByID(ctx, attackerID)
	if err != nil {
		return nil, purgerr.Wrap(ErrUserNotFound, err)
	}

	defender, err := s.UserRepo.GetUserByID(ctx, defenderID)
	if err != nil {
		return nil, purgerr.Wrap(ErrDefenderNotFound, err)
	}

	if attacker.TerraceLevel != defender.TerraceLevel {
		return nil, purgerr.Wrap(ErrTerraceLevelMismatch, ErrTerraceLevelMismatch)
	}

	defenderCombat, err := s.BattleRepo.GetUserCombat(ctx, defenderID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get defender combat state"), err)
	}

	if defenderCombat.ShieldExpiresAt != nil && time.Now().Before(*defenderCombat.ShieldExpiresAt) {
		return nil, purgerr.Wrap(ErrDefenderShieldActive, ErrDefenderShieldActive)
	}

	buildings, err := s.snapshotDefenderBase(ctx, defenderID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to snapshot defender base"), err)
	}

	snapshotID, err := s.BattleRepo.CreateBaseSnapshot(ctx, defenderID, buildings)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to create base snapshot"), err)
	}

	battle := model.Battle{
		AttackerID:     attackerID,
		DefenderID:     defenderID,
		Outcome:        "pending",
		BaseSnapshotID: &snapshotID,
		StartedAt:      time.Now(),
	}

	battleID, err := s.BattleRepo.CreateBattle(ctx, battle)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to create battle"), err)
	}

	return &model.InitiateResponse{
		BattleID:     battleID,
		DefenderName: defender.Username,
	}, nil
}

func (s *BattleService) snapshotDefenderBase(ctx context.Context, userID uuid.UUID) ([]engine.BuildingSnapshot, error) {
	layout, err := s.BaseLayoutRepo.GetLayout(ctx, userID)
	if err != nil {
		return nil, err
	}

	var snapshots []engine.BuildingSnapshot
	for _, pb := range layout {
		stats, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, pb.BuildingID, pb.Level)
		if err != nil {
			return nil, err
		}

		building, err := s.ShopRepo.GetBuildingByID(ctx, pb.BuildingID)
		if err != nil {
			return nil, err
		}

		dps := 0
		if stats.DamagePerSec != nil {
			dps = *stats.DamagePerSec
		}
		atkRange := 0.0
		if stats.AttackRange != nil {
			atkRange = *stats.AttackRange
		}

		hp := 0
		if stats.HP != nil {
			hp = *stats.HP
		}

		snapshots = append(snapshots, engine.BuildingSnapshot{
			ID:           fmt.Sprintf("%s_%d_%d", pb.BuildingID, pb.X, pb.Y),
			BuildingType: pb.BuildingID,
			Category:     building.Category.String(),
			Level:        pb.Level,
			HP:           hp,
			MaxHP:        hp,
			DPS:          dps,
			Range:        atkRange,
			Position:     engine.Point{X: float64(pb.X), Y: float64(pb.Y)},
			Size:         building.Size,
		})
	}
	return snapshots, nil
}

