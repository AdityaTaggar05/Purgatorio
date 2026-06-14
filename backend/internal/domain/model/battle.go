package model

import (
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/google/uuid"
)

type MatchListEntry struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	TerraceLevel int       `json:"terrace_level"`
}

type InitiateRequest struct {
	DefenderID string `json:"defender_id" validate:"required,uuid"`
}

type InitiateResponse struct {
	BattleID     uuid.UUID `json:"battle_id"`
	DefenderName string    `json:"defender_name"`
}

type BattleResultResponse struct {
	BattleID    uuid.UUID            `json:"battle_id"`
	Outcome     engine.BattleOutcome `json:"outcome"`
	Destruction float64              `json:"destruction"`
	Loot        int                  `json:"loot"`
	Duration    int                  `json:"duration"`
	SinMeter    int                  `json:"sin_meter"`
}

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
}

type BattleReplay struct {
	BattleID   uuid.UUID               `json:"battle_id"`
	AttackerID uuid.UUID               `json:"attacker_id"`
	DefenderID uuid.UUID               `json:"defender_id"`
	Outcome    string                  `json:"outcome"`
	Data       ReplayData              `json:"data"`
}

type ReplayData struct {
	Deployment     []engine.TroopDeployment `json:"deployment"`
	Seed           int64                    `json:"seed"`
	BaseSnapshotID uuid.UUID                `json:"base_snapshot_id"`
}

type ReplayResponse struct {
	Ticks      []engine.TickResult       `json:"ticks"`
	Result     engine.BattleResult       `json:"result"`
	Deployment []engine.TroopDeployment  `json:"deployment"`
}
