package middleware

import (
	"net/http"
	"strings"

	"rbac/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RequireAuth memverifikasi token JWT dari header Authorization (Bearer <token>)
func RequireAuth(tm jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "Authorization header diperlukan",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "Format Authorization header harus 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := tm.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "Token tidak valid atau sudah kedaluwarsa: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Pastikan token yang digunakan adalah Access Token, bukan Refresh Token
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "Tipe token tidak valid untuk mengakses resource ini",
			})
			c.Abort()
			return
		}

		// Simpan informasi user ke context request
		c.Set("userID", claims.Subject)
		c.Set("sessionUUID", claims.SessionUUID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
