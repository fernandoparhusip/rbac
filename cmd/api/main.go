package main

import (
	"context"
	"log"
	"time"

	"rbac/config"
	"rbac/internal/httpserver"
	"rbac/server"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Memuat konfigurasi (termasuk file .env)
	cfg := config.Load()

	// 2. Setup Context Timeout (10 detik maksimal untuk nyambung)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Membuka Koneksi Database (menggunakan pgxpool)
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Gagal membuat connection pool ke database: %v", err)
	}
	defer pool.Close() // Pastikan ditutup saat aplikasi mati

	// 4. Memastikan Database Benar-Benar Nyambung (Ping)
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	// 5. Inisialisasi HTTP Server (Gin)
	// ... (kode koneksi db ke atas tetap sama)

	// 5. Inisialisasi HTTP Server (Menggunakan bungkusan kita)
	srv := httpserver.New()

	// Daftarkan rute-rute aplikasi ke dalam server tersebut
	server.SetupRouter(srv)

	log.Printf("Server berjalan di port %s...\n", cfg.AppPort)
	if err := srv.Start(":" + cfg.AppPort); err != nil { // Pakai srv.Start
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
