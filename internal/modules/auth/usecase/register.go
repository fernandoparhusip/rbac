package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rbac/internal/modules/auth/domain"

	"github.com/google/uuid"
)

// CreateUser memvalidasi data, meng-hash password, dan menyimpan user baru beserta relasi rolenya
func (u *Usecase) CreateUser(
	ctx context.Context,
	req domain.RegisterRequest,
	roleID *uuid.UUID,
) (domain.User, error) {

	fmt.Printf("RegisterRequest: %+v\n", req)
	// 1. Cek apakah email sudah terdaftar sebelumnya
	existingUser, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser.Email != "" {
		return domain.User{}, errors.New("email sudah digunakan")
	}

	// 2. Hash password menggunakan Bcrypt
	hash, err := u.password.Hash(req.Password)
	if err != nil {
		return domain.User{}, errors.New("gagal mengenkripsi password: " + err.Error())
	}

	// 3. Tentukan Role ID jika roleID belum ditentukan langsung
	if roleID == nil || *roleID == uuid.Nil {
		if req.RoleID != nil && *req.RoleID != uuid.Nil {
			roleID = req.RoleID
		} else if req.Role != "" {
			foundRole, errRole := u.repo.GetRoleByName(ctx, req.Role)
			if errRole == nil {
				roleID = &foundRole.ID
			}
		}
	}

	// 4. Siapkan entity domain.User
	newUser := domain.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hash,
		FullName:  req.FullName,
		Status:    1,
		CreatedAt: time.Now(),
	}

	// 5. Simpan user ke repository database
	createdUser, err := u.repo.CreateUser(ctx, newUser, roleID)
	if err != nil {
		return domain.User{}, err
	}

	// Bersihkan hash password dari return value
	createdUser.Password = req.Password
	fmt.Printf("=====User Yang Berhasil Register =====")
	fmt.Printf("Username: %+v\n", createdUser.Username)
	fmt.Printf("Email: %+v\n", createdUser.Email)
	fmt.Printf("Password: %+v\n", createdUser.Password)
	fmt.Printf("FullName: %+v\n", createdUser.FullName)

	return createdUser, nil
}

// Register menangani alur pendaftaran user dan langsung menerbitkan sepasang token sesi
func (u *Usecase) Register(ctx context.Context, req domain.RegisterRequest) (domain.LoginResponse, error) {
	createdUser, err := u.CreateUser(ctx, req, req.RoleID)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	// Buat Session UUID dan Token otomatis setelah registrasi sukses
	sessionUUID := uuid.NewString()
	accessToken, _, expiresAt, err := u.tokenManager.NewAccessToken(createdUser.ID.String(), createdUser.Role, sessionUUID)
	if err != nil {
		return domain.LoginResponse{}, errors.New("registrasi berhasil namun gagal membuat token: " + err.Error())
	}

	refreshToken, _, _, err := u.tokenManager.NewRefreshToken(createdUser.ID.String(), createdUser.Role, sessionUUID)
	if err != nil {
		return domain.LoginResponse{}, errors.New("registrasi berhasil namun gagal membuat refresh token: " + err.Error())
	}

	return domain.LoginResponse{
		User: createdUser,
		Tokens: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		},
	}, nil
}

// AssignRole menugaskan role tertentu kepada user
func (u *Usecase) AssignRole(ctx context.Context, req domain.AssignRoleRequest) (domain.User, error) {
	updatedUser, err := u.repo.AssignRoleToUser(ctx, req.UserID, req.RoleID)
	if err != nil {
		return domain.User{}, errors.New("gagal menetapkan role: " + err.Error())
	}

	updatedUser.Password = ""
	return updatedUser, nil
}
