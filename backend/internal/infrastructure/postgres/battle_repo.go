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

