package repository

import (
	"context"
	"rbac/internal/modules/auth/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User, roleID *uuid.UUID) (domain.User, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (domain.User, error)
	GetDefaultRoleID(ctx context.Context) (uuid.UUID, error)
	GetRoleByName(ctx context.Context, name string) (domain.Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error)
}
