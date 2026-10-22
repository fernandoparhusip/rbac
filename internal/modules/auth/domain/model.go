package domain

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan entitas pengguna di dalam sistem
type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password,omitempty"`
	FullName  string     `json:"full_name"`
	Role      string     `json:"role"`             // Role utama user (misal: "Admin", "User")
	Roles     []string   `json:"roles,omitempty"` // Daftar seluruh role jika user memiliki banyak role
	Status    int        `json:"status"`          // 1: Aktif, 0: Suspend
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Role merepresentasikan entitas peran hak akses di sistem
type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Session merepresentasikan riwayat sesi token user di database
type Session struct {
	ID          uuid.UUID  `json:"id"`
	SessionUUID uuid.UUID  `json:"session_uuid"`
	UserID      uuid.UUID  `json:"user_id"`
	AccessJTI   string     `json:"access_jti"`
	RefreshJTI  string     `json:"refresh_jti"`
	RefreshHash string     `json:"-"`
	ExpiredAt   time.Time  `json:"expired_at"`
	LogoutAt    *time.Time `json:"logout_at,omitempty"`
	Revoked     bool       `json:"revoked"`
	CreatedAt   time.Time  `json:"created_at"`
}

// LoginRequest DTO untuk input formulir login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest DTO untuk input pendaftaran user baru
type RegisterRequest struct {
	Username string     `json:"username" binding:"required,min=3"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	FullName string     `json:"full_name" binding:"required"`
	RoleID   *uuid.UUID `json:"role_id,omitempty"` // Opsional: memasukkan Role ID UUID langsung
	Role     string     `json:"role,omitempty"`    // Opsional: memasukkan nama role (misal: "Admin" atau "User")
}

// AssignRoleRequest DTO untuk mengaitkan role ke user tertentu
type AssignRoleRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// TokenResponse DTO untuk payload token JWT yang dikembalikan ke client
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Durasi kedaluwarsa dalam detik
}

// LoginResponse DTO kembalian saat login atau register berhasil
type LoginResponse struct {
	User   User          `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}
