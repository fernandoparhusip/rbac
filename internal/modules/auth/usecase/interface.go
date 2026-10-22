package usecase

import (
	"context"
	"rbac/internal/modules/auth/domain"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error)
	Register(ctx context.Context, req domain.RegisterRequest) (domain.LoginResponse, error)
	CreateUser(ctx context.Context, req domain.RegisterRequest, roleID *uuid.UUID) (domain.User, error)
	AssignRole(ctx context.Context, req domain.AssignRoleRequest) (domain.User, error)
}
