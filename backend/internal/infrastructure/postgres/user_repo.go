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
		SET penitence=$2, grace=$3, max_penitence=$4, collector_pending_penitence=$5, collector_reset_at=$6, updated_at=$6
		WHERE user_id=$1
	`

	_, err := r.DB.Exec(ctx, query, eco.ID, eco.Penitence, eco.Grace, eco.MaxPenitence, eco.CollectorPendingPenitence, time.Now())

	return err
}
