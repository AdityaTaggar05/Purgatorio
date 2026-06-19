package model

import (
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/google/uuid"
)

type Battle struct {
	ID             uuid.UUID
	AttackerID     uuid.UUID
	DefenderID     uuid.UUID
	Outcome        string
	Destruction    float64
	Loot           int
	Duration       int
	BaseSnapshotID *uuid.UUID
	StartedAt      time.Time
	FinishedAt     *time.Time
}

type UserCombat struct {
	UserID            uuid.UUID
	SinMeter          int
	LastAttackAt      *time.Time
	ShieldExpiresAt   *time.Time
	ShieldMaxDuration int
	UpdatedAt         *time.Time
}

type BattleReplay struct {
	BattleID   uuid.UUID              `json:"battle_id"`
	AttackerID uuid.UUID              `json:"attacker_id"`
	DefenderID uuid.UUID              `json:"defender_id"`
	Outcome    string                 `json:"outcome"`
	Data       ReplayData             `json:"data"`
}

type ReplayData struct {
	Deployment     []engine.TroopDeployment `json:"deployment"`
	Seed           int64                    `json:"seed"`
	BaseSnapshotID uuid.UUID                `json:"base_snapshot_id"`
	EndTick        int                      `json:"end_tick,omitempty"` // 0 = full replay
}

type MatchPlayer struct {
	UserID       uuid.UUID
	Username     string
	TerraceLevel int
}

type BattleOutcome struct {
	Outcome     engine.BattleOutcome
	Destruction float64
	Loot        int
	Duration    int
	SinMeter    int
}

type ReplaySimResult struct {
	Ticks      []engine.TickResult
	Result     engine.BattleResult
	Deployment []engine.TroopDeployment
}
