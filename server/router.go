package server

import (
	"rbac/internal/httpserver"
	"rbac/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter sekarang menerima bungkusan Server buatan kita
func SetupRouter(srv *httpserver.Server) {

	// Route publik: Menggunakan fungsi Health bawaan bungkusan kita
	srv.GET("/health", srv.HealthHandler())

	srv.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Route yang butuh Auth (Perhatikan penempatan middleware-nya ada di belakang koma)
	srv.GET("/profile", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(200, gin.H{"message": "Ini profil Anda", "user_id": userID})
	}, middleware.RequireAuth())

	// Route dengan 2 Middleware
	srv.DELETE("/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "User dihapus"})
	}, middleware.RequireAuth(), middleware.RequireRole("Admin"))
}
