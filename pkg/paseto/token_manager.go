package paseto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kzankpe/e-commerce-api/internal/domain/models"
	"github.com/o1egl/paseto"
)

type TokenManager struct {
	symmetricKey []byte
	v2           *paseto.V2
}

func NewTokenManager(secretKey string) *TokenManager {
	// Ensure key is 32 bytes for encryption
	key := []byte(secretKey)
	if len(key) > 32 {
		key = key[:32]
	} else if len(key) < 32 {
		for len(key) < 32 {
			key = append(key, 0)
		}
	}

	return &TokenManager{
		symmetricKey: key,
		v2:           paseto.NewV2(),
	}
}

func (tm *TokenManager) GenerateTokenPair(user *models.User, accessDuration, refreshDuration time.Duration) (*models.TokenPair, error) {
	now := time.Now()

	// Access token payload
	accessPayload := models.TokenPayload{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      models.RoleCustomer,
		ExpiresAt: now.Add(accessDuration),
		IssuedAt:  now,
	}

	// Refresh token payload
	refreshPayload := models.TokenPayload{
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: now.Add(refreshDuration),
		IssuedAt:  now,
	}

	// Create tokens
	accessToken, err := tm.encryptToken(accessPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := tm.encryptToken(refreshPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// VerifyToken decrypts and validates a token
func (tm *TokenManager) VerifyToken(token string) (*models.TokenPayload, error) {
	var payload models.TokenPayload

	err := tm.v2.Decrypt(token, tm.symmetricKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	// Check expiration
	if time.Now().After(payload.ExpiresAt) {
		return nil, fmt.Errorf("token has expired")
	}

	return &payload, nil
}

// encryptToken encrypts the payload into a PASETO token
func (tm *TokenManager) encryptToken(payload interface{}) (string, error) {
	// Convert payload to JSON bytes
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a JSONToken wrapper
	jsonToken := paseto.JSONToken{}
	if err := json.Unmarshal(jsonData, &jsonToken); err != nil {
		return "", fmt.Errorf("failed to create JSON token: %w", err)
	}

	// Encrypt using V2 local (symmetric) encryption
	token, err := tm.v2.Encrypt(tm.symmetricKey, jsonToken, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}

	return token, nil
}
