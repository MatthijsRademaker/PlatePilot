package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/platepilot/backend/internal/common/domain"
)

// Repository handles auth data persistence.
type Repository struct {
	pool *pgxpool.Pool
}

// RefreshTokenRecord represents a persisted refresh token.
type RefreshTokenRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// NewRepository creates a new auth repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateUser creates a user with password credentials.
func (r *Repository) CreateUser(ctx context.Context, email, displayName, passwordHash string) (*domain.User, error) {
	user := &domain.User{
		ID:          uuid.New(),
		Email:       email,
		DisplayName: displayName,
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, email, display_name) VALUES ($1, $2, $3)`,
		user.ID, user.Email, user.DisplayName,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_credentials (user_id, password_hash) VALUES ($1, $2)`,
		user.ID, passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("insert credentials: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, display_name, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.DisplayName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	return &user, nil
}

// GetPasswordHash retrieves the password hash for a user.
func (r *Repository) GetPasswordHash(ctx context.Context, userID uuid.UUID) (string, error) {
	var hash string
	err := r.pool.QueryRow(ctx, `SELECT password_hash FROM user_credentials WHERE user_id = $1`, userID).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("query credentials: %w", err)
	}
	return hash, nil
}

// CreateRefreshToken stores a refresh token for a user.
func (r *Repository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_refresh_tokens (id, user_id, token_hash, expires_at, user_agent, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), userID, tokenHash, expiresAt, userAgent, ipAddress,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// GetRefreshTokenByHash retrieves a refresh token by its hash.
func (r *Repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error) {
	query := `SELECT id, user_id, expires_at, revoked_at FROM user_refresh_tokens WHERE token_hash = $1`

	var record RefreshTokenRecord
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&record.ID, &record.UserID, &record.ExpiresAt, &record.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("query refresh token: %w", err)
	}

	return &record, nil
}

// RevokeRefreshToken revokes a refresh token.
func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE user_refresh_tokens SET revoked_at = NOW(), last_used_at = NOW() WHERE id = $1`,
		tokenID,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

// RotateRefreshToken revokes a token and inserts a new one atomically.
func (r *Repository) RotateRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, newHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx,
		`UPDATE user_refresh_tokens
		 SET revoked_at = NOW(), last_used_at = NOW()
		 WHERE id = $1 AND revoked_at IS NULL`,
		tokenID,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrInvalidRefreshToken
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_refresh_tokens (id, user_id, token_hash, expires_at, user_agent, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), userID, newHash, expiresAt, userAgent, ipAddress,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
