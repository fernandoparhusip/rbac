package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
	}
}
func RequireRole(allowedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole := "UserBiasa"

		if currentRole != allowedRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You do not have the required role"})
			c.Abort()
			return
		}
		c.Next()
	}

}


