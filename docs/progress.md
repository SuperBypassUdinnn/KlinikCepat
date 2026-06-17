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
*   [x] **Keamanan Tingkat Lanjut (RBAC & RLS)**:
    - **RBAC**: Implementasi kontrol akses berbasis peran (Superadmin & Admin Klinik) menggunakan tabel `user_roles` dan *middleware* `RequireRole`.
    - **RLS**: Proteksi data sisi database dengan mengaktifkan *Row Level Security* dan memberlakukan *Deny All* untuk menutup celah eksploitasi API langsung dari sisi klien (Supabase JS Anon Key).

### Sisi Frontend (Scaffolding Selesai)
*   [x] **Inisialisasi Proyek React + Vite**:
    - Proyek frontend di-bootstrap menggunakan Vite 6 dengan plugin React.
    - Konfigurasi proxy API (`/api/*` → `http://localhost:8080`) agar komunikasi ke backend Go berjalan tanpa masalah CORS di lingkungan development.
*   [x] **Struktur Direktori Sesuai Arsitektur**:
    - Direktori `src/` diorganisasi berdasarkan peran: `pages/Patient/`, `pages/AdminKlinik/`, `pages/SuperAdmin/`.
    - Folder terpisah untuk `components/` (UI Kit), `context/` (state management), `hooks/` (custom hooks), dan `services/` (API layer).
*   [x] **AuthContext (Skeleton)**:
    - `AuthProvider` dan hook `useAuth` siap digunakan, menunggu integrasi dengan Supabase Auth client SDK.
*   [x] **Custom Hook `useHaversine`**:
    - Hook kalkulasi jarak lokasi antara posisi GPS pengguna dan koordinat klinik menggunakan rumus Haversine, siap dipakai di halaman Cari Klinik.
*   [x] **API Service Layer**:
    - Modul `services/api.js` sebagai satu-satunya titik komunikasi HTTP ke Go Backend, lengkap dengan injeksi token auth otomatis dari `localStorage`.
    - Fungsi awal: `submitTriage()`, `getQueue()`, `getClinics()`, `getClinicById()`.
*   [x] **Routing Dasar**:
    - React Router v7 terpasang dengan route skeleton untuk tiga aktor: Pasien (`/`), Admin Klinik (`/admin/*`), Super Admin (`/superadmin/*`).
*   [x] **Dokumentasi README**:
    - Panduan instalasi dan menjalankan frontend ditambahkan ke README utama proyek.

---

## 2. Pekerjaan Berikutnya (Rekomendasi Langkah Selanjutnya)

### A. Sisi Frontend — Integrasi UI dengan Backend
Fondasi scaffolding sudah siap. Langkah selanjutnya adalah membangun halaman dan mengintegrasikan dengan endpoint backend:
1.  **Halaman Cari Klinik (Pasien)**:
    - Membangun UI daftar klinik dengan memanggil `GET /api/v1/klinik`.
    - Menggunakan `useHaversine` untuk mengurutkan klinik berdasarkan jarak terdekat dari posisi GPS pasien.
2.  **Kuesioner Triage Pasien**:
    - Membangun form input gejala dinamis dari `GET /api/v1/gejala` beserta pilihan keparahan (1-3).
    - Submit ke `POST /api/v1/triage`, lalu tampilkan status warna hasil triage, nomor antrean, dan estimasi waktu.
3.  **Dashboard Admin Klinik**:
    - Integrasi login admin via Supabase Auth client SDK → simpan session di `AuthContext`.
    - Ambil antrean aktif via `GET /api/v1/admin/antrean?klinik_id=<uuid>&status=Menunggu` dengan token JWT di header.
    - Tombol aksi "Panggil/Selesai" via `PUT /api/v1/admin/antrean/{id}/status`.
4.  **Komponen UI Reusable**:
    - Membangun `Navbar`, `Card`, `Button`, `Modal` di folder `components/`.

### B. Integrasi Lingkungan (Deployment)
1.  Menyediakan file `.env` di server deployment backend dengan konfigurasi `DATABASE_URL` dan `SUPABASE_JWT_SECRET` yang valid.
2.  Deploy backend Go (misal menggunakan Docker, fly.io, Render, atau Railway).
3.  Deploy frontend (Vercel, Netlify, atau sebagai static build di server yang sama).
