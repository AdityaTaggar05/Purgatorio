package postgres

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository struct {
	DB *pgxpool.Pool
}

func NewBaseRepository(db *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{
		DB: db,
	}
}

func (r *BaseRepository) GetResourceGenerationInfo(ctx context.Context, id uuid.UUID) ([]model.ResourceGenerationInfo, error) {
	query := `
        SELECT
            bl.building_id,
            bl.level,
            curr.production_rate AS current_rate,
						curr.storage_capacity,
            bl.metadata
        FROM base_layouts bl
				JOIN buildings b
						ON b.id = bl.building_id
        JOIN building_levels curr
            ON curr.building_id = bl.building_id
           AND curr.level = bl.level
        WHERE bl.user_id = $1
          AND b.category = 'resource'
    `

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.ResourceGenerationInfo

	for rows.Next() {
		var b model.ResourceGenerationInfo

		if err := rows.Scan(
			&b.BuildingID,
			&b.BuildingLevel,
			&b.CurrentRate,
			&b.StorageCapacity,
			&b.Metadata,
		); err != nil {
			return nil, err
		} else {
			buildings = append(buildings, b)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}

func (r *BaseRepository) RemoveUpgradeInfo(ctx context.Context, userID uuid.UUID, category model.BuildingCategory) error {
	query := `
		UPDATE base_layouts bl
		SET 
			metadata = bl.metadata - 'upgrade_ends_at'
		FROM buildings b
		WHERE
			b.id = bl.building_id
			AND b.category=$1
			AND bl.user_id=$2
			AND bl.metadata ? 'upgrade_ends_at'
			AND (bl.metadata->>'upgrade_ends_at')::timestamptz <= NOW()
	`

	_, err := r.DB.Exec(ctx, query, category, userID)

	return err
}
