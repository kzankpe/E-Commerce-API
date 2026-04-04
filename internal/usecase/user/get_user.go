package user

import (
	"context"

	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
)

type GetUserUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetUserUseCase(userRepo repositories.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to fetch user", err.Error())
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
