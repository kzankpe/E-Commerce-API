package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal/delivery/response"
	"github.com/kzankpe/e-commerce-api/internal/domain/repositories"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/logger"
	"github.com/kzankpe/e-commerce-api/pkg/paseto"
)

type AuthMiddleware struct {
	tokenManager             *paseto.TokenManager
	tokenBlacklistRepository repositories.TokenBlacklistRepository
	logger                   *logger.Logger
}

func NewAuthMiddleware(
	tokenManager *paseto.TokenManager,
	tokenBlacklistRepository repositories.TokenBlacklistRepository,
	logger *logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager:             tokenManager,
		tokenBlacklistRepository: tokenBlacklistRepository,
		logger:                   logger,
	}
}

// Authenticate verifies the PASETO token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, response.AuthResponse{
				Success: false,
				Message: errors.ErrMissingToken.Message,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, response.AuthResponse{
				Success: false,
				Message: errors.ErrInvalidToken.Message,
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Check if token is blacklisted
		isBlacklisted, err := m.tokenBlacklistRepository.IsBlacklisted(c.Request.Context(), token)
		if err != nil {
			m.logger.Error("Failed to check token blacklist", err)
			c.JSON(http.StatusInternalServerError, response.AuthResponse{
				Success: false,
				Message: "Internal server error",
			})
			c.Abort()
			return
		}
		if isBlacklisted {
			m.logger.Warn("Token is blacklisted", "token", token[:10]+"...")
			c.JSON(http.StatusUnauthorized, response.AuthResponse{
				Success: false,
				Message: errors.ErrTokenBlacklisted.Message,
			})
			c.Abort()
			return
		}

		// Verify and parse the PASETO token
		payload, err := m.tokenManager.VerifyToken(token)
		if err != nil {
			m.logger.Warn("Invalid or expired token", "error", err)
			c.JSON(http.StatusUnauthorized, response.AuthResponse{
				Success: false,
				Message: errors.ErrInvalidToken.Message,
			})
			c.Abort()
			return
		}
		// Check if token has expired
		if time.Now().Unix() > payload.ExpiresAt.Unix() {
			m.logger.Warn("Token has expired")
			c.JSON(http.StatusUnauthorized, response.AuthResponse{
				Success: false,
				Message: errors.ErrTokenExpired.Message,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", payload.UserID)
		c.Set("email", payload.Email)
		c.Set("role", payload.Role)
		c.Set("token", token)

		m.logger.Info("Token authenticated successfully", "userID", payload.UserID)

		// Continue to next handler
		c.Next()
	}
}
