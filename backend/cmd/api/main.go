package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment system")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL ERROR: DATABASE_URL tidak disetel")
	}

	// 2. Inisialisasi Koneksi Connection Pool ke Supabase PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("FATAL ERROR: Gagal terhubung ke database: %v\n", err)
	}
	defer dbPool.Close()

	// Tes Ping ke Database
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("FATAL ERROR: Ping ke database gagal: %v\n", err)
	}
	fmt.Println("STATUS: Berhasil terhubung ke Supabase PostgreSQL")

	// 3. Inisialisasi Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Endpoint Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("API KlinikCepat Aktif dan Terhubung ke Database"))
	})

	// 4. Jalankan Peladen
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("STATUS: Peladen berjalan di port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}