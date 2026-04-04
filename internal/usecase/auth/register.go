package auth

import (
	"context"

	"github.com/kzankpe/e-commerce-api/internal/domain/models"
	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/password"
	"github.com/kzankpe/e-commerce-api/pkg/validator"
)

type RegisterUseCase struct {
	userRepo repositories.UserRepository
	hasher   *password.Hasher
}

func NewRegisterUseCase(userRepo repositories.UserRepository, hasher *password.Hasher) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	// Validate input
	v := validator.New().
		ValidateEmail(email).
		ValidatePassword(password).
		ValidateRequired("first_name", firstName).
		ValidateRequired("last_name", lastName)

	if v.HasErrors() {
		return nil, errors.NewAppError(400, "Validation failed", "")
	}

	// Check if user already exists
	exists, err := uc.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.hasher.Hash(password)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to process password", err.Error())
	}

	// Create user
	user := &models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewAppError(500, "Failed to create user", err.Error())
	}

	return user, nil
}
