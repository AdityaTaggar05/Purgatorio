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
	BattleRepo     repository.BattleRepository
	UserRepo       repository.UserRepository
	ArmyRepo       repository.ArmyRepository
	BaseLayoutRepo repository.BaseLayoutRepository
	ShopRepo       repository.ShopRepository
	Catalog        engine.TroopCatalog
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
		BattleRepo:     battleRepo,
		UserRepo:       userRepo,
		ArmyRepo:       armyRepo,
		BaseLayoutRepo: baseLayoutRepo,
		ShopRepo:       shopRepo,
		Catalog:        catalog,
	}
}

func (s *BattleService) GetMatchList(ctx context.Context, userID uuid.UUID) ([]model.MatchPlayer, error) {
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(ErrUserNotFound, err)
	}
	return s.BattleRepo.GetMatchList(ctx, user.TerraceLevel, userID)
}

type InitiateBattleResult struct {
	BattleID      uuid.UUID
	DefenderName  string
	TerraceLevel  int
	Buildings     []engine.BuildingSnapshot
}

func (s *BattleService) InitiateBattle(ctx context.Context, attackerID, defenderID uuid.UUID) (*InitiateBattleResult, error) {
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

	return &InitiateBattleResult{
		BattleID:     battleID,
		DefenderName: defender.Username,
		TerraceLevel: defender.TerraceLevel,
		Buildings:    buildings,
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

func (s *BattleService) ValidateFullDeployment(ctx context.Context, userID uuid.UUID, deployments []engine.TroopDeployment) error {
	army, err := s.BattleRepo.GetUserArmyForBattle(ctx, userID)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get user army"), err)
	}

	totals := make(map[string]int)
	for _, dep := range deployments {
		totals[dep.TroopType] += dep.Count
	}

	for troopType, count := range totals {
		if owned := army.Troops[troopType]; count > owned {
			return purgerr.Wrap(ErrInsufficientArmyTroops, ErrInsufficientArmyTroops)
		}
	}
	return nil
}

func (s *BattleService) PrepareSimulation(ctx context.Context, battleID, userID uuid.UUID, deployments []engine.TroopDeployment) (*engine.Simulation, error) {
	battle, err := s.BattleRepo.GetBattle(ctx, battleID)
	if err != nil {
		return nil, purgerr.Wrap(ErrBattleNotFound, err)
	}

	if battle.Outcome != "pending" {
		return nil, purgerr.Wrap(ErrBattleNotPending, ErrBattleNotPending)
	}

	if battle.AttackerID != userID {
		return nil, purgerr.Wrap(ErrBattleNotFound, ErrBattleNotFound)
	}

	if battle.BaseSnapshotID == nil {
		return nil, purgerr.Wrap(fmt.Errorf("battle has no snapshot"), fmt.Errorf("battle has no snapshot"))
	}

	buildings, err := s.BattleRepo.GetBaseSnapshot(ctx, *battle.BaseSnapshotID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to load base snapshot"), err)
	}

	input := engine.BattleInput{
		Seed:        time.Now().UnixNano(),
		Deployments: deployments,
		Buildings:   buildings,
		Catalog:     s.Catalog,
		TicksPerSec: 20,
		MaxDuration: 3600,
	}

	return engine.NewSimulation(input), nil
}

func (s *BattleService) ResolveAndStore(ctx context.Context, battleID uuid.UUID, sim *engine.Simulation, deployment []engine.TroopDeployment) (*model.BattleOutcome, error) {
	battle, err := s.BattleRepo.GetBattle(ctx, battleID)
	if err != nil {
		return nil, purgerr.Wrap(ErrBattleNotFound, err)
	}

	if battle.BaseSnapshotID == nil {
		return nil, purgerr.Wrap(fmt.Errorf("battle has no snapshot"), fmt.Errorf("battle has no snapshot"))
	}

	engineResult := sim.Result()

	attackerCombat, err := s.BattleRepo.GetUserCombat(ctx, battle.AttackerID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get attacker combat state"), err)
	}

	finalOutcome, newSin := resolveSinMeter(engineResult, attackerCombat.SinMeter)
	engineResult.Outcome = finalOutcome

	now := time.Now()

	loot := 0
	if finalOutcome == engine.Victory {
		defenderEco, err := s.BattleRepo.GetUserEconomyForBattle(ctx, battle.DefenderID)
		if err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to get defender economy"), err)
		}
		loot = computeLoot(defenderEco.Penitence, engineResult.Destruction)
	}

	if err := s.BattleRepo.UpdateBattleOutcome(ctx, battleID, string(finalOutcome), engineResult.Destruction, loot, engineResult.Duration, *battle.BaseSnapshotID); err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to update battle outcome"), err)
	}

	if err := s.BattleRepo.UpdateUserCombat(ctx, battle.AttackerID, newSin, &now); err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to update attacker combat"), err)
	}

	switch finalOutcome {
	case engine.Victory:
		if err := s.BattleRepo.DeductDefenderPenitence(ctx, battle.DefenderID, loot); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to deduct defender penitence"), err)
		}
		if err := s.BattleRepo.AddAttackerLoot(ctx, battle.AttackerID, loot); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to add attacker loot"), err)
		}
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.AttackerID, true, true); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment attacker stats"), err)
		}
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.DefenderID, false, false); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment defender stats"), err)
		}

		shieldDuration := time.Duration(12*3600+int(engineResult.Destruction*30*60)) * time.Second
		shieldExpires := time.Now().Add(shieldDuration)
		if err := s.BattleRepo.SetShield(ctx, battle.DefenderID, shieldExpires); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to set defender shield"), err)
		}
	case engine.Defeat:
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.AttackerID, true, false); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment attacker stats"), err)
		}
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.DefenderID, false, true); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment defender stats"), err)
		}
	default:
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.AttackerID, true, false); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment attacker stats"), err)
		}
		if err := s.BattleRepo.IncrementUserStats(ctx, battle.DefenderID, false, true); err != nil {
			return nil, purgerr.Wrap(fmt.Errorf("failed to increment defender stats"), err)
		}
	}

	deductions := make(map[string]int)
	for _, dep := range deployment {
		deductions[dep.TroopType] = dep.Count
	}
	if err := s.BattleRepo.DeductTroopsFromArmy(ctx, battle.AttackerID, deductions); err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to deduct troops"), err)
	}

	replayData := model.ReplayData{
		Deployment:     deployment,
		Seed:           sim.Seed(),
		BaseSnapshotID: *battle.BaseSnapshotID,
	}
	if err := s.BattleRepo.StoreReplay(ctx, battleID, replayData); err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to store replay"), err)
	}

	return &model.BattleOutcome{
		Outcome:     finalOutcome,
		Destruction: engineResult.Destruction,
		Loot:        loot,
		Duration:    engineResult.Duration,
		SinMeter:    newSin,
	}, nil
}

func (s *BattleService) GetReplay(ctx context.Context, battleID uuid.UUID) (*model.ReplaySimResult, error) {
	replay, err := s.BattleRepo.GetReplay(ctx, battleID)
	if err != nil {
		return nil, purgerr.Wrap(ErrReplayNotFound, err)
	}

	buildings, err := s.BattleRepo.GetBaseSnapshot(ctx, replay.Data.BaseSnapshotID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to load base snapshot for replay"), err)
	}

	input := engine.BattleInput{
		Seed:        replay.Data.Seed,
		Deployments: replay.Data.Deployment,
		Buildings:   buildings,
		Catalog:     s.Catalog,
		TicksPerSec: 10,
		MaxDuration: 6000,
	}

	sim := engine.NewSimulation(input)
	var allTicks []engine.TickResult
	for !sim.IsDone() {
		allTicks = append(allTicks, sim.NextTick())
	}

	return &model.ReplaySimResult{
		Ticks:      allTicks,
		Result:     sim.Result(),
		Deployment: replay.Data.Deployment,
	}, nil
}

func (s *BattleService) GetRecentBattles(ctx context.Context, userID uuid.UUID, limit int) ([]model.Battle, error) {
	return s.BattleRepo.GetRecentBattles(ctx, userID, limit)
}

func resolveSinMeter(result engine.BattleResult, currentSin int) (engine.BattleOutcome, int) {
	if result.Outcome == engine.Defeat {
		return engine.Defeat, currentSin
	}

	if result.Destruction > float64(currentSin) {
		newSin := min(currentSin+10, 100)
		return engine.Victory, newSin
	}

	return engine.ThresholdFailed, 0
}

func computeLoot(defenderPenitence int, destruction float64) int {
	loot := int(float64(defenderPenitence) * 0.20 * destruction / 100.0)
	cap := int(float64(defenderPenitence) * 0.50)
	if loot > cap {
		loot = cap
	}
	return int(math.Min(float64(loot), float64(cap)))
}
