package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindUserCredentials(ctx context.Context, email string) (User, string, error) {
	const q = `
SELECT u.id::text, u.email::text, u.full_name, u.status, u.session_version, c.password_hash
FROM users u
JOIN user_password_credentials c ON c.user_id = u.id
WHERE u.email = $1 AND u.deleted_at IS NULL`
	var user User
	var hash string
	if err := r.pool.QueryRow(ctx, q, email).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Status, &user.SessionVersion, &hash); err != nil {
		return User{}, "", err
	}
	return user, hash, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, s Session) error {
	const q = `
INSERT INTO auth_sessions (user_id, token_hash, csrf_hash, session_version, remember_me, created_at, last_seen_at, idle_expires_at, absolute_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, q, s.UserID, s.TokenHash, s.CSRFHash, s.SessionVersion, s.RememberMe, s.CreatedAt, s.LastSeenAt, s.IdleExpiresAt, s.AbsoluteExpiresAt)
	return err
}

func (r *PostgresRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = now(), updated_at = now() WHERE id = $1`, userID)
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}
