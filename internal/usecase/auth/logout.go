package auth

import (
	"context"

	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/paseto"
)

type LogoutUseCase struct {
	tokenManager             *paseto.TokenManager
	tokenBlacklistRepository repositories.TokenBlacklistRepository
}

func NewLogoutUseCase(
	tokenManager *paseto.TokenManager,
	tokenBlacklistRepository repositories.TokenBlacklistRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenManager:             tokenManager,
		tokenBlacklistRepository: tokenBlacklistRepository,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, accessToken string) error {
	// Verify token to get expiration
	payload, err := uc.tokenManager.VerifyToken(accessToken)
	if err != nil {
		return errors.ErrInvalidToken
	}

	// Add token to blacklist
	expiresAt := payload.ExpiresAt.Unix()
	if err := uc.tokenBlacklistRepository.Add(ctx, accessToken, expiresAt); err != nil {
		return errors.NewAppError(500, "Failed to logout", err.Error())
	}

	return nil
}
