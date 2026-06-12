package postgres

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseLayoutRepository struct {
	DB *pgxpool.Pool
}

func NewBaseLayoutRepository(db *pgxpool.Pool) *BaseLayoutRepository {
	return &BaseLayoutRepository{DB: db}
}

func (r *BaseLayoutRepository) GetLayout(ctx context.Context, userID uuid.UUID) ([]model.PlacedBuilding, error) {
	query := `
		SELECT building_id, x, y, level, metadata
		FROM base_layouts
		WHERE user_id = $1`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.PlacedBuilding
	for rows.Next() {
		var pb model.PlacedBuilding
		pb.UserID = userID
		if err := rows.Scan(&pb.BuildingID, &pb.X, &pb.Y, &pb.Level, &pb.Metadata); err != nil {
			return nil, err
		}
		buildings = append(buildings, pb)
	}
	return buildings, rows.Err()
}

func (r *BaseLayoutRepository) PlaceBuilding(ctx context.Context, pb model.PlacedBuilding) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO base_layouts (user_id, building_id, x, y, level, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		pb.UserID, pb.BuildingID, pb.X, pb.Y, pb.Level, pb.Metadata,
	)
	return err
}

func (r *BaseLayoutRepository) RemoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, x, y int) error {
	_, err := r.DB.Exec(ctx,
		`DELETE FROM base_layouts
		 WHERE user_id = $1 AND building_id = $2 AND x = $3 AND y = $4`,
		userID, buildingID, x, y,
	)
	return err
}

func (r *BaseLayoutRepository) MoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, fromX, fromY, toX, toY int) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE base_layouts SET x = $5, y = $6, updated_at = now()
		 WHERE user_id = $1 AND building_id = $2 AND x = $3 AND y = $4`,
		userID, buildingID, fromX, fromY, toX, toY,
	)
	return err
}

func (r *BaseLayoutRepository) GetBuildingAtPosition(ctx context.Context, userID uuid.UUID, x, y int) (*model.PlacedBuilding, error) {
	var pb model.PlacedBuilding
	pb.UserID = userID
	err := r.DB.QueryRow(ctx,
		`SELECT building_id, x, y, level, metadata
		 FROM base_layouts
		 WHERE user_id = $1 AND x = $2 AND y = $3`,
		userID, x, y,
	).Scan(&pb.BuildingID, &pb.X, &pb.Y, &pb.Level, &pb.Metadata)
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func (r *BaseLayoutRepository) GetBuildingLevelStats(ctx context.Context, buildingID string, level int) (*model.BuildingLevel, error) {
	var l model.BuildingLevel
	err := r.DB.QueryRow(ctx,
		`SELECT building_id, level, hp, damage_per_second, production_rate,
		        storage_capacity, attack_range, upgrade_cost, upgrade_time
		 FROM building_levels
		 WHERE building_id = $1 AND level = $2`,
		buildingID, level,
	).Scan(&l.BuildingID, &l.Level, &l.HP, &l.DamagePerSec,
		&l.ProductionRate, &l.StorageCapacity, &l.AttackRange,
		&l.UpgradeCost, &l.UpgradeTime)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
