package repositories

import "context"

type TokenBlacklistRepository interface {
	Add(ctx context.Context, token string, expireAt int64) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpired(ctx context.Context) error
}
