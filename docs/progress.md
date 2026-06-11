# Progres Pengembangan: KlinikCepat

Dokumen ini melacak kemajuan pengerjaan aplikasi KlinikCepat secara berkala untuk mempermudah kolaborasi antara tim Backend (Rekan A) dan Frontend (Rekan B).

---

## 1. Kemajuan Saat Ini (11 Juni 2026)

### Sisi Backend (Selesai 100% untuk Tahap Awal)
*   [x] **Setup Basis Data & Migrasi**: 
    - Tabel `klinik`, `katalog_gejala`, dan `antrean` berhasil dibuat di remote PostgreSQL Supabase.
    - Indeks optimasi antrean `idx_antrean_admin_view` telah di-apply.
    - Data awal (*seeding*) untuk referensi gejala dan beberapa klinik contoh telah di-injeksi.
*   [x] **Model & Akses Data (Repository)**: 
    - Model data Go struct siap digunakan di `internal/models`.
    - Modul koneksi pooling `pgxpool` diimplementasikan di `internal/repository`.
    - CRUD database handler untuk `Klinik` dan `Gejala` selesai.
    - Query antrean berlapis `status_triage ASC, created_at ASC` selesai.
*   [x] **Mesin Logika Triage (Triage Engine)**: 
    - Kode kalkulasi $S_{urgensi} = \sum (W_i \cdot V_i)$ di `internal/services/triage_service.go` selesai.
    - Logika otomatisasi "Merah" jika ada gejala fatal bernilai 10 (Pendarahan Hebat / Sesak Napas) selesai.
*   [x] **Middleware & Keamanan (JWT)**: 
    - Custom verification middleware untuk Supabase Auth JWT di `internal/middleware/auth.go` selesai tanpa membutuhkan dependensi luar.
*   [x] **API Route & Controller**: 
    - Integrasi endpoint publik (klinik, gejala, submit triage) dan terproteksi admin selesai.
    - Kode berhasil di-build tanpa error dan lulus audit linter `golangci-lint run`.
*   [x] **Unit Testing & Abstraksi**:
    - Refactoring repositori menjadi interface (`repository.RepositoryInterface`) untuk pengujian mandiri tanpa koneksi database.
    - Implementasi `MockRepository` *in-memory* di `internal/handlers/mock_test.go`.
    - Pengujian unit lengkap untuk Auth Middleware (`auth_test.go`), Triage Service (`triage_service_test.go`), dan seluruh Handlers API (`klinik_handler_test.go`, `gejala_handler_test.go`, `antrean_handler_test.go`).
    - Lulus pengujian unit dengan status **PASS** (total 31 subtests).

---

## 2. Pekerjaan Berikutnya (Rekomendasi Langkah Selanjutnya)

### A. Sisi Frontend (React App)
Rekan B dapat mulai mengintegrasikan aplikasi web React dengan endpoint backend yang sudah ada:
1.  **Integrasi Cari Klinik**: 
    - Memanggil `GET /api/v1/klinik`.
    - Meminta hak akses GPS browser, jalankan formula Haversine di client untuk mengurutkan jarak klinik, dan menampilkannya kepada pasien.
2.  **Kuesioner Triage Pasien**:
    - Memanggil `GET /api/v1/gejala` untuk menampilkan kuesioner input gejala dinamis beserta pilihan keparahan (1-3).
    - Melakukan POST ke `POST /api/v1/triage` saat menekan tombol daftar, menampilkan status warna hasil triage, nomor antrean, serta estimasi.
3.  **Dashboard Admin Klinik**:
    - Melakukan otentikasi login admin via Supabase Auth (di client).
    - Mengirimkan token ke backend pada header `Authorization: Bearer <token>` dan mengambil antrean aktif lewat `GET /api/v1/admin/antrean?klinik_id=<uuid>&status=Menunggu`.
    - Membuat tombol aksi "Panggil/Selesai" untuk mengubah status via `PUT /api/v1/admin/antrean/{id}/status`.

### B. Integrasi Lingkungan (Deployment)
1.  Menyediakan file `.env` di server deployment backend dengan konfigurasi `DATABASE_URL` dan `SUPABASE_JWT_SECRET` yang valid.
2.  Deploy backend Go (misal menggunakan Docker, fly.io, Render, atau Railway).
