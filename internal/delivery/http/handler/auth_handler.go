package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal/delivery/request"
	"github.com/kzankpe/e-commerce-api/internal/delivery/response"
	"github.com/kzankpe/e-commerce-api/internal/usecase/auth"
	"github.com/kzankpe/e-commerce-api/pkg/errors"
	"github.com/kzankpe/e-commerce-api/pkg/logger"

	useruc "github.com/kzankpe/e-commerce-api/internal/usecase/user"
)

type AuthHandler struct {
	registerUC     *auth.RegisterUseCase
	loginUC        *auth.LoginUseCase
	refreshTokenUC *auth.RefreshTokenUseCase
	logoutUC       *auth.LogoutUseCase
	getUserUC      *useruc.GetUserUseCase
	logger         *logger.Logger
}

func NewAuthHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	refreshTokenUC *auth.RefreshTokenUseCase,
	logoutUC *auth.LogoutUseCase,
	getUserUC *useruc.GetUserUseCase,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUC:     registerUC,
		loginUC:        loginUC,
		refreshTokenUC: refreshTokenUC,
		logoutUC:       logoutUC,
		getUserUC:      getUserUC,
		logger:         logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid registration request", err)
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Success: false,
			Message: "Invalid request data",
			Errors: []response.ErrorDetail{
				{
					Field:   "request",
					Message: err.Error(),
				},
			},
		})
		return
	}

	user, err := h.registerUC.Execute(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Data: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", err)
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	tokenPair, err := h.loginUC.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Success: true,
		Message: "Login successful",
		Data: response.TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    tokenPair.ExpiresIn,
			TokenType:    tokenPair.TokenType,
		},
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh token request", err)
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	tokenPair, err := h.refreshTokenUC.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: response.TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    tokenPair.ExpiresIn,
			TokenType:    tokenPair.TokenType,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from Authorization header
	token := c.GetString("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, response.AuthResponse{
			Success: false,
			Message: "Token is required",
		})
		return
	}

	if err := h.logoutUC.Execute(c.Request.Context(), token); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// GetProfile retrieves authenticated user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.AuthResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	user, err := h.getUserUC.Execute(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// handleError handles different error types
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error", err)
		c.JSON(http.StatusInternalServerError, response.AuthResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(appErr.Code, response.AuthResponse{
		Success: false,
		Message: appErr.Message,
		Errors: []response.ErrorDetail{
			{
				Message: appErr.Details,
			},
		},
	})
}
