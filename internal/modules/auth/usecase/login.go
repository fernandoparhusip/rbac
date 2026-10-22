package usecase

import (
	"context"
	"errors"
	"time"

	"rbac/internal/modules/auth/domain"

	"github.com/google/uuid"
)

// Login memverifikasi email & password user, kemudian menerbitkan sepasang token JWT (Access & Refresh)
func (u *Usecase) Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginResponse{}, errors.New("email atau password salah")
	}

	// 2. Verifikasi kecocokan password menggunakan Bcrypt Hasher
	if err := u.password.Compare(user.Password, req.Password); err != nil {
		return domain.LoginResponse{}, errors.New("email atau password salah")
	}

	// 3. Pastikan akun dalam status aktif (status = 1)
	if user.Status == 0 {
		return domain.LoginResponse{}, errors.New("akun Anda sedang dinonaktifkan, silakan hubungi admin")
	}

	// 4. Default role jika kosong
	if user.Role == "" {
		user.Role = "User"
	}

	// 5. Buat Session UUID unik untuk sesi login ini
	sessionUUID := uuid.NewString()

	// 6. Buat Access Token dan Refresh Token
	accessToken, _, expiresAt, err := u.tokenManager.NewAccessToken(user.ID.String(), user.Role, sessionUUID)
	if err != nil {
		return domain.LoginResponse{}, errors.New("gagal membuat access token: " + err.Error())
	}

	refreshToken, _, _, err := u.tokenManager.NewRefreshToken(user.ID.String(), user.Role, sessionUUID)
	if err != nil {
		return domain.LoginResponse{}, errors.New("gagal membuat refresh token: " + err.Error())
	}

	// 7. Bersihkan hash password dari memori sebelum dikembalikan ke client
	user.Password = ""

	return domain.LoginResponse{
		User: user,
		Tokens: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		},
	}, nil
}
