package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"rbac/internal/middleware"
	"rbac/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(tm jwt.TokenManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	protected := r.Group("/api", middleware.RequireAuth(tm))
	{
		protected.GET("/user/profile", func(c *gin.Context) {
			userID := c.GetString("userID")
			role := c.GetString("role")
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
		})

		protected.DELETE("/admin/users/:id", middleware.RequireRole("Admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
		})
	}

	return r
}

func TestMiddleware_RequireAuth_And_RequireRole(t *testing.T) {
	tm := jwt.NewManager("secret-key-test", 10*time.Minute, 1*time.Hour)
	router := setupTestRouter(tm)

	// Buat token user biasa dan token admin
	userToken, _, _, _ := tm.NewAccessToken("user-1", "User", "session-1")
	adminToken, _, _, _ := tm.NewAccessToken("admin-1", "Admin", "session-2")

	// 1. Request tanpa token -> harus 401
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/profile", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Harusnya 401 Unauthorized, tapi dapat %d", w.Code)
	}

	// 2. Request user biasa ke profile -> harus 200
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/user/profile", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Harusnya 200 OK, tapi dapat %d", w.Code)
	}

	// 3. Request user biasa ke rute admin -> harus 403 Forbidden
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/users/123", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Harusnya 403 Forbidden untuk role User, tapi dapat %d", w.Code)
	}

	// 4. Request admin ke rute admin -> harus 200 OK
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/users/123", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Harusnya 200 OK untuk role Admin, tapi dapat %d", w.Code)
	}
}
