package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"rbac/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.Load()

	// 1. Buat koneksi ke database default bawaan postgres (database: postgres)
	// Tujuannya agar kita bisa menjalankan perintah CREATE DATABASE rbac
	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBSSLMode)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, defaultDSN)
	if err != nil {
		log.Fatalf("Gagal menyambung ke default database (postgres): %v", err)
	}

	// 2. Eksekusi CREATE DATABASE
	log.Printf("Memeriksa dan membuat database '%s' jika belum ada...\n", cfg.DBName)
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		// Jika error karena database sudah ada, kita abaikan saja.
		if strings.Contains(err.Error(), "already exists") {
			log.Println("Database sudah ada, melewati pembuatan database.")
		} else {
			log.Fatalf("Gagal membuat database: %v", err)
		}
	} else {
		log.Println("Berhasil membuat database baru!")
	}
	pool.Close()

	// 3. Jalankan Migrasi ke database yang dituju (rbac)
	log.Println("Menjalankan migrasi tabel...")
	
	// Format DSN khusus untuk golang-migrate
	migrateDSN := cfg.DSN()
	
	m, err := migrate.New(
		"file://migration", // lokasi folder migrasi
		migrateDSN,
	)
	if err != nil {
		log.Fatalf("Gagal inisialisasi alat migrasi: %v", err)
	}

	// Mengeksekusi file .up.sql
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Semua tabel sudah versi terbaru (tidak ada migrasi baru).")
		} else {
			log.Fatalf("Gagal menjalankan migrasi: %v", err)
		}
	} else {
		log.Println("Migrasi tabel selesai sukses!")
	}
}
