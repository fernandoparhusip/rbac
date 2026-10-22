package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config menyimpan semua konfigurasi dari environment
type Config struct {
	AppName       string
	AppPort       string
	DBHost        string
	DBName        string
	DBUser        string
	DBPassword    string
	DBSSLMode     string
	DBPort        string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
}

// Load membaca file .env dan mengisinya ke dalam struct Config
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, mengambil dari OS environment")
	}

	return &Config{
		AppName:       getEnv("APP_NAME", "RBAC"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBName:        getEnv("DB_NAME", "rbac"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASWORD", ""), // Sesuai dengan penamaan di file .env Anda
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		DBPort:        getEnv("DB_PORT", "5432"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-rbac-key-change-in-production"),
		JWTAccessTTL:  15 * time.Minute,
		JWTRefreshTTL: 7 * 24 * time.Hour,
	}
}

// getEnv adalah helper function untuk mencegah nilai kosong
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// DSN (Data Source Name) menghasilkan connection string PostgreSQL
func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}
