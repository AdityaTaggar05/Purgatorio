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

func (r *BattleRepository) GetMatchList(ctx context.Context, terraceLevel int, excludeUserID uuid.UUID) ([]model.MatchPlayer, error) {
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

	var entries []model.MatchPlayer
	for rows.Next() {
		var e model.MatchPlayer
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

func (r *BattleRepository) GetReplay(ctx context.Context, battleID uuid.UUID) (model.BattleReplay, error) {
	var replay model.BattleReplay
	var data []byte
	err := r.DB.QueryRow(ctx,
		`SELECT br.battle_id, b.attacker_id, b.defender_id, b.outcome, br.data
		 FROM battle_replays br
		 JOIN battles b ON b.id = br.battle_id
		 WHERE br.battle_id = $1`,
		battleID,
	).Scan(&replay.BattleID, &replay.AttackerID, &replay.DefenderID, &replay.Outcome, &data)
	if err != nil {
		return replay, err
	}

	if err := json.Unmarshal(data, &replay.Data); err != nil {
		return replay, err
	}
	return replay, nil
}

func (r *BattleRepository) GetRecentBattles(ctx context.Context, userID uuid.UUID, limit int) ([]model.Battle, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, attacker_id, defender_id, outcome, destruction, loot, duration,
		        base_snapshot_id, started_at, finished_at
		 FROM battles
		 WHERE (attacker_id = $1 OR defender_id = $1) AND outcome != 'pending'
		 ORDER BY started_at DESC
		 LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []model.Battle
	for rows.Next() {
		var b model.Battle
		if err := rows.Scan(&b.ID, &b.AttackerID, &b.DefenderID, &b.Outcome,
			&b.Destruction, &b.Loot, &b.Duration, &b.BaseSnapshotID,
			&b.StartedAt, &b.FinishedAt); err != nil {
			return nil, err
		}
		battles = append(battles, b)
	}
	return battles, rows.Err()
}

func (r *BattleRepository) GetUserCombat(ctx context.Context, userID uuid.UUID) (model.UserCombat, error) {
	var c model.UserCombat
	c.UserID = userID
	err := r.DB.QueryRow(ctx,
		`SELECT sin_meter, last_attack_at, shield_expires_at, shield_max_duration
		 FROM user_combat WHERE user_id = $1`,
		userID,
	).Scan(&c.SinMeter, &c.LastAttackAt, &c.ShieldExpiresAt, &c.ShieldMaxDuration)
	return c, err
}

func (r *BattleRepository) UpdateUserCombat(ctx context.Context, userID uuid.UUID, sinMeter int, lastAttackAt *time.Time) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE user_combat SET sin_meter = $2, last_attack_at = $3, updated_at = now()
		 WHERE user_id = $1`,
		userID, sinMeter, lastAttackAt,
	)
	return err
}

func (r *BattleRepository) SetShield(ctx context.Context, userID uuid.UUID, shieldExpiresAt time.Time) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE user_combat SET shield_expires_at = $2, updated_at = now()
		 WHERE user_id = $1`,
		userID, shieldExpiresAt,
	)
	return err
}

func (r *BattleRepository) IncrementUserStats(ctx context.Context, userID uuid.UUID, isAttacker, isSuccess bool) error {
	if isAttacker {
		if isSuccess {
			_, err := r.DB.Exec(ctx,
				`UPDATE user_stats SET attacks = attacks + 1, attacks_success = attacks_success + 1, updated_at = now()
				 WHERE user_id = $1`, userID)
			return err
		}
		_, err := r.DB.Exec(ctx,
			`UPDATE user_stats SET attacks = attacks + 1, updated_at = now()
			 WHERE user_id = $1`, userID)
		return err
	}

	if isSuccess {
		_, err := r.DB.Exec(ctx,
			`UPDATE user_stats SET defenses = defenses + 1, defenses_success = defenses_success + 1, updated_at = now()
			 WHERE user_id = $1`, userID)
		return err
	}
	_, err := r.DB.Exec(ctx,
		`UPDATE user_stats SET defenses = defenses + 1, updated_at = now()
		 WHERE user_id = $1`, userID)
	return err
}

func (r *BattleRepository) GetUserEconomyForBattle(ctx context.Context, userID uuid.UUID) (model.UserEconomy, error) {
	var eco model.UserEconomy
	eco.ID = userID
	err := r.DB.QueryRow(ctx,
		`SELECT penitence, grace, max_penitence, collector_pending_penitence, collector_reset_at
		 FROM user_economy WHERE user_id = $1`,
		userID,
	).Scan(&eco.Penitence, &eco.Grace, &eco.MaxPenitence, &eco.CollectorPendingPenitence, &eco.CollectorResetAt)
	return eco, err
}

func (r *BattleRepository) DeductDefenderPenitence(ctx context.Context, userID uuid.UUID, amount int) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE user_economy SET
			penitence = GREATEST(penitence - $2, 0),
			collector_pending_penitence = GREATEST(collector_pending_penitence - GREATEST($2 - penitence, 0), 0),
			updated_at = now()
		 WHERE user_id = $1`,
		userID, amount,
	)
	return err
}

func (r *BattleRepository) AddAttackerLoot(ctx context.Context, userID uuid.UUID, amount int) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE user_economy SET penitence = penitence + $2, updated_at = now()
		 WHERE user_id = $1`,
		userID, amount,
	)
	return err
}

func (r *BattleRepository) GetUserArmyForBattle(ctx context.Context, userID uuid.UUID) (model.UserArmy, error) {
	var army model.UserArmy
	army.UserID = userID
	var troopsJSON []byte
	err := r.DB.QueryRow(ctx,
		`SELECT troops, used_capacity FROM user_army WHERE user_id = $1`,
		userID,
	).Scan(&troopsJSON, &army.UsedCapacity)
	if err != nil {
		return army, err
	}

	if err := json.Unmarshal(troopsJSON, &army.Troops); err != nil {
		return army, err
	}
	return army, nil
}

func (r *BattleRepository) DeductTroopsFromArmy(ctx context.Context, userID uuid.UUID, deductions map[string]int) error {
	army, err := r.GetUserArmyForBattle(ctx, userID)
	if err != nil {
		return err
	}

	for troopID, count := range deductions {
		current := army.Troops[troopID]
		army.Troops[troopID] = current - count
		if army.Troops[troopID] < 0 {
			army.Troops[troopID] = 0
		}
	}

	troopsJSON, err := json.Marshal(army.Troops)
	if err != nil {
		return err
	}

	newCapacity := 0
	troopList := []struct {
		id    string
		space int
	}{}
	rows, err := r.DB.Query(ctx, `SELECT id, space FROM troops`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id string
			var space int
			if err := rows.Scan(&id, &space); err == nil {
				troopList = append(troopList, struct {
					id    string
					space int
				}{id, space})
			}
		}
	}
	for _, t := range troopList {
		newCapacity += army.Troops[t.id] * t.space
	}

	_, err = r.DB.Exec(ctx,
		`UPDATE user_army SET troops = $2::jsonb, used_capacity = $3, updated_at = now()
		 WHERE user_id = $1`,
		userID, troopsJSON, newCapacity,
	)
	return err
}
