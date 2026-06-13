package repository

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/google/uuid"
)

type BattleRepository interface {
	GetMatchList(ctx context.Context, terraceLevel int, excludeUserID uuid.UUID) ([]model.MatchListEntry, error)

	CreateBattle(ctx context.Context, battle model.Battle) (uuid.UUID, error)
	GetBattle(ctx context.Context, battleID uuid.UUID) (model.Battle, error)
	UpdateBattleOutcome(ctx context.Context, battleID uuid.UUID, outcome string, destruction float64, loot, duration int, snapshotID uuid.UUID) error

	CreateBaseSnapshot(ctx context.Context, userID uuid.UUID, buildings []engine.BuildingSnapshot) (uuid.UUID, error)
	GetBaseSnapshot(ctx context.Context, snapshotID uuid.UUID) ([]engine.BuildingSnapshot, error)

	StoreReplay(ctx context.Context, battleID uuid.UUID, data model.ReplayData) error
	GetReplay(ctx context.Context, battleID uuid.UUID) (model.BattleReplay, error)
	GetRecentBattles(ctx context.Context, userID uuid.UUID, limit int) ([]model.Battle, error)

	GetUserCombat(ctx context.Context, userID uuid.UUID) (model.UserCombat, error)
	UpdateUserCombat(ctx context.Context, userID uuid.UUID, sinMeter int, lastAttackAt *time.Time) error
	SetShield(ctx context.Context, userID uuid.UUID, shieldExpiresAt time.Time) error
	IncrementUserStats(ctx context.Context, userID uuid.UUID, isAttacker, isSuccess bool) error

	GetUserEconomyForBattle(ctx context.Context, userID uuid.UUID) (model.UserEconomy, error)
	DeductDefenderPenitence(ctx context.Context, userID uuid.UUID, amount int) error
	AddAttackerLoot(ctx context.Context, userID uuid.UUID, amount int) error

	GetUserArmyForBattle(ctx context.Context, userID uuid.UUID) (model.UserArmy, error)
	DeductTroopsFromArmy(ctx context.Context, userID uuid.UUID, deductions map[string]int) error
}
