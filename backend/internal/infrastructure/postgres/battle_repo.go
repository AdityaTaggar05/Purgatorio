package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BattleRepository struct {
	DB *pgxpool.Pool
}

func NewBattleRepository(db *pgxpool.Pool) *BattleRepository {
	return &BattleRepository{DB: db}
}

func (r *BattleRepository) GetMatchList(ctx context.Context, terraceLevel int, excludeUserID uuid.UUID) ([]model.MatchListEntry, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, username, terrace_level FROM users
		 WHERE terrace_level = $1 AND id != $2
		 ORDER BY username`,
		terraceLevel, excludeUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.MatchListEntry
	for rows.Next() {
		var e model.MatchListEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.TerraceLevel); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *BattleRepository) CreateBattle(ctx context.Context, battle model.Battle) (uuid.UUID, error) {
	err := r.DB.QueryRow(ctx,
		`INSERT INTO battles (attacker_id, defender_id, outcome, base_snapshot_id, started_at)
		 VALUES ($1, $2, 'pending', $3, $4)
		 RETURNING id`,
		battle.AttackerID, battle.DefenderID, battle.BaseSnapshotID, battle.StartedAt,
	).Scan(&battle.ID)
	return battle.ID, err
}

func (r *BattleRepository) GetBattle(ctx context.Context, battleID uuid.UUID) (model.Battle, error) {
	var b model.Battle
	err := r.DB.QueryRow(ctx,
		`SELECT id, attacker_id, defender_id, outcome, destruction, loot, duration,
		        base_snapshot_id, started_at, finished_at
		 FROM battles WHERE id = $1`,
		battleID,
	).Scan(&b.ID, &b.AttackerID, &b.DefenderID, &b.Outcome, &b.Destruction,
		&b.Loot, &b.Duration, &b.BaseSnapshotID, &b.StartedAt, &b.FinishedAt)
	return b, err
}

func (r *BattleRepository) UpdateBattleOutcome(ctx context.Context, battleID uuid.UUID, outcome string, destruction float64, loot, duration int, snapshotID uuid.UUID) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE battles SET outcome = $2, destruction = $3, loot = $4, duration = $5,
		        base_snapshot_id = $6, finished_at = now()
		 WHERE id = $1`,
		battleID, outcome, destruction, loot, duration, snapshotID,
	)
	return err
}

func (r *BattleRepository) CreateBaseSnapshot(ctx context.Context, userID uuid.UUID, buildings []engine.BuildingSnapshot) (uuid.UUID, error) {
	data, err := json.Marshal(buildings)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = r.DB.QueryRow(ctx,
		`INSERT INTO base_snapshots (user_id, buildings) VALUES ($1, $2) RETURNING id`,
		userID, data,
	).Scan(&id)
	return id, err
}

func (r *BattleRepository) GetBaseSnapshot(ctx context.Context, snapshotID uuid.UUID) ([]engine.BuildingSnapshot, error) {
	var data []byte
	err := r.DB.QueryRow(ctx,
		`SELECT buildings FROM base_snapshots WHERE id = $1`,
		snapshotID,
	).Scan(&data)
	if err != nil {
		return nil, err
	}

	var buildings []engine.BuildingSnapshot
	if err := json.Unmarshal(data, &buildings); err != nil {
		return nil, err
	}
	return buildings, nil
}

func (r *BattleRepository) StoreReplay(ctx context.Context, battleID uuid.UUID, data model.ReplayData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(ctx,
		`INSERT INTO battle_replays (battle_id, data) VALUES ($1, $2)
		 ON CONFLICT (battle_id) DO UPDATE SET data = EXCLUDED.data`,
		battleID, jsonData,
	)
	return err
}

