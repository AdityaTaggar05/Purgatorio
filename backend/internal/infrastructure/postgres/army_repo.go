package postgres

import (
	"context"
	"encoding/json"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArmyRepository struct {
	DB *pgxpool.Pool
}

func NewArmyRepository(db *pgxpool.Pool) *ArmyRepository {
	return &ArmyRepository{DB: db}
}

func (r *ArmyRepository) GetAllTroops(ctx context.Context) ([]model.Troop, error) {
	query := `SELECT id, name, training_cost, space, hp, dps, speed, attack_range, preferred_target FROM troops ORDER BY id`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var troops []model.Troop
	for rows.Next() {
		var t model.Troop
		if err := rows.Scan(&t.ID, &t.Name, &t.TrainingCost, &t.Space, &t.HP, &t.DPS,
			&t.Speed, &t.AttackRange, &t.PreferredTarget); err != nil {
			return nil, err
		}
		troops = append(troops, t)
	}
	return troops, rows.Err()
}

func (r *ArmyRepository) GetTroopByID(ctx context.Context, troopID string) (*model.Troop, error) {
	var t model.Troop
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, training_cost, space, hp, dps, speed, attack_range, preferred_target
		 FROM troops WHERE id = $1`, troopID,
	).Scan(&t.ID, &t.Name, &t.TrainingCost, &t.Space, &t.HP, &t.DPS,
		&t.Speed, &t.AttackRange, &t.PreferredTarget)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ArmyRepository) GetUserArmy(ctx context.Context, userID uuid.UUID) (*model.UserArmy, error) {
	var troopsJSON []byte
	var usedCapacity int

	err := r.DB.QueryRow(ctx,
		`SELECT troops, used_capacity FROM user_army WHERE user_id = $1`, userID,
	).Scan(&troopsJSON, &usedCapacity)
	if err != nil {
		return nil, err
	}

	var troopsMap map[string]int
	if err := json.Unmarshal(troopsJSON, &troopsMap); err != nil {
		return nil, err
	}

	return &model.UserArmy{
		UserID:       userID,
		Troops:       troopsMap,
		UsedCapacity: usedCapacity,
	}, nil
}

