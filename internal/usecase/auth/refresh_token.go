package auth

import (
	"context"
	"time"

	"github.com/kzankpe/e-commerce-api/internal/domain/models"
	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/paseto"
)

type RefreshTokenUseCase struct {
	userRepo        repositories.UserRepository
	tokenManager    *paseto.TokenManager
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewRefreshTokenUseCase(
	userRepo repositories.UserRepository,
	tokenManager *paseto.TokenManager,
	accessDuration, refreshDuration time.Duration,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:        userRepo,
		tokenManager:    tokenManager,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Verify refresh token
	payload, err := uc.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, payload.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.NewAppError(403, "Account is disabled", "")
	}

	// Generate new token pair
	tokenPair, err := uc.tokenManager.GenerateTokenPair(user, uc.accessDuration, uc.refreshDuration)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to generate tokens", err.Error())
	}

	return tokenPair, nil
}
