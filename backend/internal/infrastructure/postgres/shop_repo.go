package postgres

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopRepository struct {
	DB *pgxpool.Pool
}

func NewShopRepository(db *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{DB: db}
}

func (r *ShopRepository) GetAllBuildings(ctx context.Context) ([]model.Building, error) {
	query := `SELECT id, name, size, price, currency, category FROM buildings ORDER BY name`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.Building
	for rows.Next() {
		var b model.Building
		if err := rows.Scan(&b.ID, &b.Name, &b.Size, &b.Price, &b.Currency, &b.Category); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, rows.Err()
}

func (r *ShopRepository) GetUserBuildingCounts(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT building_id, quantity FROM user_buildings WHERE user_id = $1`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var id string
		var qty int
		if err := rows.Scan(&id, &qty); err != nil {
			return nil, err
		}
		counts[id] = qty
	}
	return counts, rows.Err()
}

func (r *ShopRepository) GetLimitsByTerrace(ctx context.Context, terraceLevel int) (map[string]int, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT building_id, max_allowed FROM building_limits WHERE terrace_level = $1`, terraceLevel,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	limits := make(map[string]int)
	for rows.Next() {
		var id string
		var max int
		if err := rows.Scan(&id, &max); err != nil {
			return nil, err
		}
		limits[id] = max
	}
	return limits, rows.Err()
}

func (r *ShopRepository) GetBuildingLevels(ctx context.Context, buildingID string) ([]model.BuildingLevel, error) {
	query := `
		SELECT building_id, level, hp, damage_per_second, production_rate,
		       storage_capacity, attack_range, upgrade_cost, upgrade_time
		FROM building_levels
		WHERE building_id = $1
		ORDER BY level`

	rows, err := r.DB.Query(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []model.BuildingLevel
	for rows.Next() {
		var l model.BuildingLevel
		if err := rows.Scan(&l.BuildingID, &l.Level, &l.HP, &l.DamagePerSec,
			&l.ProductionRate, &l.StorageCapacity, &l.AttackRange,
			&l.UpgradeCost, &l.UpgradeTime); err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}
	return levels, rows.Err()
}
