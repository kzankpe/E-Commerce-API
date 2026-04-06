package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
)

type TokenBlacklistRepository struct {
	db *sql.DB
}

func NewTokenBlacklistRepository(db *sql.DB) repositories.TokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

func (r *TokenBlacklistRepository) Add(ctx context.Context, token string, expiresAt int64) error {
	query := `
        INSERT INTO token_blacklist (token, expires_at, created_at)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		token,
		time.Unix(expiresAt, 0),
		time.Now(),
	)

	return err
}

func (r *TokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE token = $1 AND expires_at > NOW())`

    var isBlacklisted bool
    err := r.db.QueryRowContext(ctx, query, token).Scan(&isBlacklisted)
    if err != nil {
        return false, err
    }

    return isBlacklisted, nil
}

func (r *TokenBlacklistRepository) CleanupExpired(ctx context.Context) error {
    query := `DELETE FROM token_blacklist WHERE expires_at < NOW()`

    _, err := r.db.ExecContext(ctx, query)
    return err
}