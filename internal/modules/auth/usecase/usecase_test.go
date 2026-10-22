package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"rbac/internal/modules/auth/domain"
	"rbac/internal/modules/auth/usecase"
	"rbac/pkg/jwt"
	"rbac/pkg/password"

	"github.com/google/uuid"
)

// MockRepository untuk pengujian usecase tanpa perlu koneksi database riil
type mockAuthRepo struct {
	users map[string]domain.User
	roles map[string]domain.Role
}

func newMockRepo() *mockAuthRepo {
	userRoleID := uuid.New()
	adminRoleID := uuid.New()
	return &mockAuthRepo{
		users: make(map[string]domain.User),
		roles: map[string]domain.Role{
			"User":  {ID: userRoleID, Name: "User"},
			"Admin": {ID: adminRoleID, Name: "Admin"},
		},
	}
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return u, nil
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, user domain.User, roleID *uuid.UUID) (domain.User, error) {
	if _, ok := m.users[user.Email]; ok {
		return domain.User{}, errors.New("email sudah digunakan")
	}
	if roleID != nil {
		for _, r := range m.roles {
			if r.ID == *roleID {
				user.Role = r.Name
				break
			}
		}
	}
	if user.Role == "" {
		user.Role = "User"
	}
	m.users[user.Email] = user
	return user, nil
}

func (m *mockAuthRepo) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (domain.User, error) {
	for email, u := range m.users {
		if u.ID == userID {
			for _, r := range m.roles {
				if r.ID == roleID {
					u.Role = r.Name
					m.users[email] = u
					return u, nil
				}
			}
		}
	}
	return domain.User{}, errors.New("user not found")
}

func (m *mockAuthRepo) GetDefaultRoleID(ctx context.Context) (uuid.UUID, error) {
	return m.roles["User"].ID, nil
}

func (m *mockAuthRepo) GetRoleByName(ctx context.Context, name string) (domain.Role, error) {
	r, ok := m.roles[name]
	if !ok {
		return domain.Role{}, errors.New("role not found")
	}
	return r, nil
}

func (m *mockAuthRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	for _, r := range m.roles {
		if r.ID == id {
			return r, nil
		}
	}
	return domain.Role{}, errors.New("role not found")
}

func TestAuthUsecase_Register_And_Login(t *testing.T) {
	repo := newMockRepo()
	hasher := password.NewBcrypt()
	tokenMgr := jwt.NewManager("secret-test-key", 15*time.Minute, 1*time.Hour)

	uc := usecase.NewUsecase(repo, hasher, tokenMgr)
	ctx := context.Background()

	// 1. Uji Register
	regReq := domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
		FullName: "John Doe",
		Role:     "User",
	}

	resReg, err := uc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Register gagal: %v", err)
	}

	if resReg.User.Email != regReq.Email {
		t.Errorf("Email mismatch: %s vs %s", resReg.User.Email, regReq.Email)
	}
	if resReg.Tokens.AccessToken == "" {
		t.Errorf("Access token kosong")
	}

	// 2. Uji Login dengan password benar
	loginReq := domain.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	resLogin, err := uc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login gagal: %v", err)
	}
	if resLogin.User.Role != "User" {
		t.Errorf("Expected role User, got %s", resLogin.User.Role)
	}

	// 3. Uji Login dengan password salah
	badLogin := domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	_, err = uc.Login(ctx, badLogin)
	if err == nil {
		t.Errorf("Login harusnya gagal untuk password salah")
	}

	// 4. Uji AssignRole ke Admin
	adminRoleID := repo.roles["Admin"].ID
	assignReq := domain.AssignRoleRequest{
		UserID: resLogin.User.ID,
		RoleID: adminRoleID,
	}
	updatedUser, err := uc.AssignRole(ctx, assignReq)
	if err != nil {
		t.Fatalf("Assign role gagal: %v", err)
	}
	if updatedUser.Role != "Admin" {
		t.Errorf("Expected updated role to be Admin, got %s", updatedUser.Role)
	}
}
