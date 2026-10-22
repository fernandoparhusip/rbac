package usecase

import (
	"rbac/internal/modules/auth/repository"
	"rbac/pkg/jwt"
	"rbac/pkg/password"
)

type Usecase struct {
	repo         repository.AuthRepository
	password     password.Hasher
	tokenManager jwt.TokenManager
}

func NewUsecase(repo repository.AuthRepository, pw password.Hasher, tm jwt.TokenManager) *Usecase {
	return &Usecase{
		repo:         repo,
		password:     pw,
		tokenManager: tm,
	}
}
