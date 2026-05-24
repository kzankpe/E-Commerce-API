package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal/delivery/http/handler"
)

func NewRouter(h *handler.AuthHandler) *gin.Engine {
	router := gin.Default()

	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}

	// Add more routes as needed

	return router
}
