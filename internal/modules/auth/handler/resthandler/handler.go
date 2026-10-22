package handler

import (
	"errors"
	"net/http"
	"rbac/internal/httpserver"
	"rbac/internal/middleware"
	"rbac/internal/modules/auth/domain"
	"rbac/internal/modules/auth/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc usecase.AuthUsecase
}

func NewHandler(uc usecase.AuthUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// Register menangani pendaftaran akun baru melalui HTTP POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New("format data registrasi tidak valid: " + err.Error()))
		return
	}

	res, err := h.uc.Register(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data":    res,
	})
}

// Login menangani autentikasi akun melalui HTTP POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New("format data login tidak valid: " + err.Error()))
		return
	}

	res, err := h.uc.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data":    res,
	})
}

// AssignRole menangani penugasan role ke user (khusus Admin) melalui HTTP POST /auth/assign-role
func (h *Handler) AssignRole(c *gin.Context) {
	var req domain.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New("format data assign-role tidak valid: " + err.Error()))
		return
	}

	updatedUser, err := h.uc.AssignRole(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role berhasil ditugaskan ke user",
		"data":    updatedUser,
	})
}

// Mount mendaftarkan seluruh endpoint modul auth ke routing server
func (h *Handler) Mount(public *httpserver.RouterGroup, protected *httpserver.RouterGroup) {
	public.POST("/auth/register", h.Register)
	public.POST("/auth/login", h.Login)

	// Rute penugasan role diproteksi (harus login & khusus role Admin)
	protected.POST("/auth/assign-role", h.AssignRole, middleware.RequireRole("Admin"))
}
