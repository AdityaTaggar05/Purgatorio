package postgres

import (
	"context"
	"fmt"

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
            prev.production_rate AS previous_rate,
            bl.metadata
        FROM base_layout bl
				JOIN buildings b
						ON b.id = bl.building_id
        JOIN building_levels curr
            ON curr.building_id = bl.building_id
           AND curr.level = bl.level
        LEFT JOIN building_levels prev
            ON prev.building_id = bl.building_id
           AND prev.level = bl.level - 1
        WHERE bl.user_id = $1
          AND b.category = resource 
    `

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("querying resource buildings: %w", err)
	}
	defer rows.Close()

	var buildings []model.ResourceGenerationInfo

	for rows.Next() {
		var b model.ResourceGenerationInfo

		if err := rows.Scan(
			&b.BuildingID,
			&b.BuildingLevel,
			&b.CurrentRate,
			&b.PreviousRate,
			&b.Metadata,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return buildings, nil
}
