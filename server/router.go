package server

import (
	"rbac/internal/httpserver"
	"rbac/internal/middleware"
	"rbac/pkg/jwt"

	authHandler "rbac/internal/modules/auth/handler/resthandler"
	authUsecase "rbac/internal/modules/auth/usecase"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Auth         authUsecase.AuthUsecase
	TokenManager jwt.TokenManager
}

// SetupRouter mendaftarkan seluruh rute aplikasi ke dalam HTTP Server
func SetupRouter(srv *httpserver.Server, deps Dependencies) {
	// 1. BUAT ZONA KEAMANAN
	// Zona Public (URL berawalan /api/v1 yang bebas diakses tanpa token)
	public := srv.Group("/api/v1")

	// Zona Protected (URL berawalan /api/v1 yang WAJIB menyertakan Bearer Token valid)
	protected := srv.Group("/api/v1", middleware.RequireAuth(deps.TokenManager))

	// Mount auth routes (/api/v1/auth/register dan /api/v1/auth/login)
	authHandler.NewHandler(deps.Auth).Mount(public, protected)

	// Health check & ping
	srv.GET("/health", srv.HealthHandler())
	srv.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Rute yang membutuhkan otentikasi (RequireAuth)
	srv.GET("/profile", func(c *gin.Context) {
		userID := c.GetString("userID")
		role := c.GetString("role")
		sessionUUID := c.GetString("sessionUUID")

		c.JSON(200, gin.H{
			"message":      "Ini profil Anda",
			"user_id":      userID,
			"role":         role,
			"session_uuid": sessionUUID,
		})
	}, middleware.RequireAuth(deps.TokenManager))

	// Rute yang membutuhkan otentikasi + otorisasi role khusus (RequireAuth + RequireRole)
	srv.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User dengan ID " + id + " berhasil dihapus oleh Admin",
		})
	}, middleware.RequireAuth(deps.TokenManager), middleware.RequireRole("Admin"))
}
