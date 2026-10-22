package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 1. TUJUAN: Menentukan bentuk Formulir Data (Payload) yang akan masuk ke dalam perut token.
// Struct ini mewakili isi JSON yang akan dibaca oleh Middleware nantinya.
type Claims struct {
	SessionUUID          string `json:"session_uuid"`
	TokenType            string `json:"type"`           // Membedakan apakah ini token "access" atau "refresh"
	Role                 string `json:"role,omitempty"` // Role user untuk otorisasi RBAC
	jwt.RegisteredClaims        // Bawaan library: Menyediakan wadah untuk 'sub' (UserID), 'jti', dan 'exp'
}

// 2. TUJUAN: Membuat "Kontrak Kerja" (Interface) agar Usecase Login Anda tahu
// tombol apa saja yang bisa ditekan dari mesin pabrik ini.
type TokenManager interface {
	NewAccessToken(userID, role, sessionUUID string) (string, string, time.Time, error)
	NewRefreshToken(userID, role, sessionUUID string) (string, string, time.Time, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// 3. TUJUAN: Struct ini bertindak sebagai "Mesin Pabrik" yang menyimpan 3 bahan baku utama:
// Tinta Rahasia (Secret), Umur Access Token (pendek), dan Umur Refresh Token (panjang).
type manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// 4. TUJUAN: Tombol Start untuk menghidupkan mesin pabrik pertama kali.
// Fungsi ini akan dipanggil dari file `main.go` saat aplikasi baru menyala.
func NewManager(secret string, accessTTL, refreshTTL time.Duration) TokenManager {
	return &manager{[]byte(secret), accessTTL, refreshTTL}
}

// 5. TUJUAN: Ini adalah Jantung Mesinnya! Fungsi rahasia ini menempelkan "Stempel Digital".
// Ia mengambil data Claims (JSON) di atas, lalu menyegelnya menggunakan Algoritma HS256 dan Tinta Rahasia (Secret).
func (m *manager) sign(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret) // <- Disinilah token panjang "eyJhbG..." tercipta
}

// 6. TUJUAN: Fungsi khusus untuk memproduksi Access Token (Tiket Utama).
// Fungsi ini akan dipanggil oleh Usecase Login.
func (m *manager) NewAccessToken(userID, role, sessionUUID string) (string, string, time.Time, error) {
	// TUJUAN: Menciptakan Nomor Seri Unik (JTI) untuk token ini
	jti := uuid.NewString()

	// TUJUAN: Menghitung kapan token ini mati (Jam Sekarang + Umur Access Token)
	expires := time.Now().Add(m.accessTTL)

	// TUJUAN: Mengisi Formulir Data (Payload) yang akan dibungkus
	claims := Claims{
		SessionUUID: sessionUUID,
		TokenType:   "access",
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,                      // Menyisipkan ID User
			ID:        jti,                         // Menyisipkan Nomor Seri JTI
			ExpiresAt: jwt.NewNumericDate(expires), // Menyisipkan Waktu Kedaluwarsa
		},
	}

	token, err := m.sign(claims)
	return token, jti, expires, err
}

// TUJUAN: Fungsi khusus untuk memproduksi Refresh Token (Tiket Cadangan berumur panjang).
func (m *manager) NewRefreshToken(userID, role, sessionUUID string) (string, string, time.Time, error) {
	jti := uuid.NewString()
	expires := time.Now().Add(m.refreshTTL)

	claims := Claims{
		SessionUUID: sessionUUID,
		TokenType:   "refresh",
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token, err := m.sign(claims)
	return token, jti, expires, err
}

// 7. TUJUAN: Mesin Pemindai Tiket (Validasi & Verifikasi Token).
// Fungsi ini menerima string token dari header Authorization,
// memverifikasi stempel digital (signature HS256) menggunakan Secret Key,
// dan mengembalikan data Claims jika token asli serta belum kedaluwarsa.
func (m *manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode penandatanganan token tidak valid")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	return claims, nil
}
