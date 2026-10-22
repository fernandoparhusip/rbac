package main

import (
	"context"
	"log"
	"time"

	"rbac/config"
	"rbac/internal/httpserver"
	authRepo "rbac/internal/modules/auth/repository"
	authUsecase "rbac/internal/modules/auth/usecase"
	"rbac/pkg/jwt"
	"rbac/pkg/password"
	"rbac/server"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Memuat konfigurasi (termasuk file .env)
	cfg := config.Load()

	// 2. Setup Context Timeout (10 detik maksimal untuk menyambung)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Membuka Koneksi Database (menggunakan pgxpool)
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Gagal membuat connection pool ke database: %v", err)
	}
	defer pool.Close()

	// 4. Memastikan Database Benar-Benar Nyambung (Ping)
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	// 5. Inisialisasi Services, Token Manager, & Hasher
	passwordHasher := password.NewBcrypt()
	tokenManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)

	repoAuth := authRepo.NewRepository(pool)
	usecaseAuth := authUsecase.NewUsecase(repoAuth, passwordHasher, tokenManager)

	// 6. Inisialisasi HTTP Server
	srv := httpserver.New()

	deps := server.Dependencies{
		Auth:         usecaseAuth,
		TokenManager: tokenManager,
	}

	// 7. Daftarkan rute-rute aplikasi ke server
	server.SetupRouter(srv, deps)

	log.Printf("Server berjalan di port %s...\n", cfg.AppPort)
	if err := srv.Start(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
