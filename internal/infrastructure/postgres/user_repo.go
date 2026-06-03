package postgres

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
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

func (r *UserRepository) GetAuthAndUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

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

func (r *UserRepository) CreateRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, exp,
	)
	return err
}
