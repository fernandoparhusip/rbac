package jwt_test

import (
	"testing"
	"time"

	"rbac/pkg/jwt"
)

func TestTokenManager_Lifecycle(t *testing.T) {
	secret := "test-secret-key-12345"
	manager := jwt.NewManager(secret, 5*time.Minute, 1*time.Hour)

	userID := "user-123"
	role := "Admin"
	sessionUUID := "sess-abc-789"

	// 1. Uji pembuatan Access Token
	token, jti, expiresAt, err := manager.NewAccessToken(userID, role, sessionUUID)
	if err != nil {
		t.Fatalf("Gagal membuat access token: %v", err)
	}
	if token == "" || jti == "" {
		t.Fatalf("Token atau JTI kosong")
	}
	if expiresAt.Before(time.Now()) {
		t.Fatalf("ExpiresAt sudah lewat")
	}

	// 2. Uji validasi Access Token yang sah
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Gagal memvalidasi token yang valid: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("Subject tidak cocok: expected %s, got %s", userID, claims.Subject)
	}
	if claims.Role != role {
		t.Errorf("Role tidak cocok: expected %s, got %s", role, claims.Role)
	}
	if claims.SessionUUID != sessionUUID {
		t.Errorf("SessionUUID tidak cocok: expected %s, got %s", sessionUUID, claims.SessionUUID)
	}
	if claims.TokenType != "access" {
		t.Errorf("TokenType tidak cocok: expected access, got %s", claims.TokenType)
	}

	// 3. Uji validasi token dengan secret berbeda (palsu / tampered)
	otherManager := jwt.NewManager("wrong-secret-key-99999", 5*time.Minute, 1*time.Hour)
	_, err = otherManager.ValidateToken(token)
	if err == nil {
		t.Errorf("Harusnya gagal validasi dengan secret yang berbeda")
	}
}
