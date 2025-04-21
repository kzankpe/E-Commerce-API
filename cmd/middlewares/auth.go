package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal"
	"github.com/kzankpe/e-commerce-api/models"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Missing authorization header"})
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Missing authorization header"})
			return
		}
		fmt.Print(accessToken)
		userSub, err := internal.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid Token"})
			return
		}
		var user models.User
		result := models.DB.First(&user, "email=?", fmt.Sprint(userSub))
		if result.Error != nil {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"status": "fail", "message": `the token is invalid`},
			)
			return
		}
		c.Set("currentUser", userSub)
		c.Next()
	}
}
