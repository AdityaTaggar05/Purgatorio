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

