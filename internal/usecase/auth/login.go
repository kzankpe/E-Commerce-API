package auth

import (
	"context"
	"time"

	"github.com/kzankpe/e-commerce-api/internal/domain/models"
	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/paseto"
	"github.com/kzankpe/e-commerce-api/pkg/password"
)

type LoginUseCase struct {
	userRepo        repositories.UserRepository
	hasher          *password.Hasher
	tokenManager    *paseto.TokenManager
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewLoginUseCase(
	userRepo repositories.UserRepository,
	hasher *password.Hasher,
	tokenManager *paseto.TokenManager,
	accessDuration, refreshDuration time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:        userRepo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*models.TokenPair, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.NewAppError(403, "Account is disabled", "")
	}

	// Verify password
	if !uc.hasher.Compare(user.Password, password) {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := uc.tokenManager.GenerateTokenPair(user, uc.accessDuration, uc.refreshDuration)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to generate tokens", err.Error())
	}

	return tokenPair, nil
}
