package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi hak akses endpoint berdasarkan role pengguna.
// Dapat menerima satu atau beberapa role yang diizinkan (misal: RequireRole("Admin", "Superadmin")).
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "forbidden",
				"message": "Akses ditolak: role pengguna tidak ditemukan di session/token",
			})
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok || userRoleStr == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "forbidden",
				"message": "Akses ditolak: role pengguna tidak valid",
			})
			c.Abort()
			return
		}

		isAllowed := false
		for _, 	role := range allowedRoles {
			if strings.EqualFold(userRoleStr, role) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"status":         "forbidden",
				"message":        "Akses ditolak: Anda tidak memiliki role yang diizinkan untuk mengakses resource ini",
				"required_roles": allowedRoles,
				"your_role":      userRoleStr,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
