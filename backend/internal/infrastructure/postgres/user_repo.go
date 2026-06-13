package postgres

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, hash, username string) (model.User, error) {
	var user model.User

	query := `
		WITH new_auth AS (
			INSERT INTO auth (email, password_hash)
			VALUES ($1, $2)
			RETURNING id
		)

		INSERT INTO users (id, username)
		SELECT id, $3
		FROM new_auth
		RETURNING id, username, xp, level, terrace_level
	`

	err := r.DB.QueryRow(ctx, query, email, hash, username).Scan(&user.ID, &user.Username, &user.XP, &user.Level, &user.TerraceLevel)
	return user, err
}

func (r *UserRepository) InitializeNewUser(ctx context.Context, userID uuid.UUID) error {
	buildings := []struct {
		id  string
		qty int
	}{
		{"bastion", 6},
		{"angel-spire", 1},
		{"lament-basin", 1},
		{"sanctum", 1},
		{"barracks", 1},
	}

	placements := []struct {
		buildingID string
		x          int
		y          int
	}{
		{"lament-basin", 4, 4},
		{"angel-spire", 14, 14},
		{"sanctum", 24, 5},
		{"barracks", 0, 20},
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO user_economy (user_id, penitence, grace, max_penitence) VALUES ($1, 500, 50, 5000)`,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO user_stats (user_id) VALUES ($1)`, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_combat (user_id, sin_meter) VALUES ($1, 0)`, userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO game_state (user_id, scene_id, fallback) VALUES ($1, 'base', 'base')`, userID,
	)
	if err != nil {
		return err
	}

	for _, b := range buildings {
		_, err = tx.Exec(ctx,
			`INSERT INTO user_buildings (user_id, building_id, quantity) VALUES ($1, $2, $3)`,
			userID, b.id, b.qty,
		)
		if err != nil {
			return err
		}
	}

	for _, p := range placements {
		_, err = tx.Exec(ctx,
			`INSERT INTO base_layouts (user_id, building_id, x, y, metadata) VALUES ($1, $2, $3, $4, '{}'::jsonb)`,
			userID, p.buildingID, p.x, p.y,
		)
		if err != nil {
			return err
		}
	}

	bastionPlacements := [][2]int{
		{13, 13}, {14, 13}, {15, 13}, {16, 13},
		{16, 14}, {16, 15}, {16, 16},
		{15, 16}, {14, 16}, {13, 16},
		{13, 15}, {13, 14},
	}
	for i, cell := range bastionPlacements {
		if i >= 6 {
			break
		}
		if cell[0] < 0 || cell[1] < 0 || cell[0] > 29 || cell[1] > 29 {
			continue
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO base_layouts (user_id, building_id, x, y, metadata) VALUES ($1, 'bastion', $2, $3, '{}'::jsonb)`,
			userID, cell[0], cell[1],
		)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_army (user_id) VALUES ($1)`, userID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepository) GetAuthAndUserByEmail(ctx context.Context, email string) (model.AuthAndUser, error) {
	var user model.AuthAndUser

	query := `
		SELECT auth.id, auth.password_hash, users.username, users.xp, users.level, users.terrace_level
		FROM auth
		INNER JOIN users ON users.id=auth.id 
		WHERE auth.email=$1
	`

	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.PasswordHash, &user.Username, &user.XP, &user.Level, &user.TerraceLevel,
	)

	return user, err
}

func (r *UserRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, exp time.Time) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, exp,
	)
	return err
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	var rt model.RefreshToken

	err := r.DB.QueryRow(
		ctx,
		`SELECT user_id, token, revoked, expires_at FROM refresh_tokens WHERE token=$1`,
		token,
	).Scan(
		&rt.UserID, &rt.Token, &rt.Revoked, &rt.ExpiresAt,
	)
	return rt, err
}

func (r *UserRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := r.DB.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked=true WHERE token=$1`,
		token,
	)
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	user.ID = id

	err := r.DB.QueryRow(
		ctx,
		`SELECT username, xp, level, terrace_level FROM users WHERE id=$1`,
		id,
	).Scan(&user.Username, &user.XP, &user.Level, &user.TerraceLevel)

	return user, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM auth WHERE id=$1`, id)
	return err
}

func (r *UserRepository) GetEconomy(ctx context.Context, id uuid.UUID) (model.UserEconomy, error) {
	var eco model.UserEconomy
	eco.ID = id

	query := `
		SELECT penitence, grace, max_penitence, collector_pending_penitence, collector_reset_at
		FROM user_economy
		WHERE user_id=$1
	`

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&eco.Penitence,
		&eco.Grace,
		&eco.MaxPenitence,
		&eco.CollectorPendingPenitence,
		&eco.CollectorResetAt,
	)

	return eco, err
}

func (r *UserRepository) UpdateEconomy(ctx context.Context, eco model.UserEconomy) error {
	query := `
		UPDATE user_economy
		SET penitence=$2, grace=$3, max_penitence=$4, collector_pending_penitence=$5, collector_reset_at=$6, updated_at=now()
		WHERE user_id=$1
	`

	_, err := r.DB.Exec(ctx, query, eco.ID, eco.Penitence, eco.Grace, eco.MaxPenitence, eco.CollectorPendingPenitence, eco.CollectorResetAt)

	return err
}

func (r *UserRepository) UpdateTerraceLevel(ctx context.Context, id uuid.UUID, level int) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE users SET terrace_level = $2, updated_at = now() WHERE id = $1`,
		id, level,
	)
	return err
}
