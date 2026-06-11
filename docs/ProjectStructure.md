# Struktur Proyek: KlinikCepat

Dokumen ini mendeskripsikan struktur direktori dan organisasi kode dari repositori monorepo proyek **KlinikCepat**.

---

## Pohon Direktori

```text
KlinikCepat/
├── backend/                  # Porsi Kerja Rekan A (Go Engine)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go       # Entry point aplikasi Go
│   ├── internal/
│   │   ├── config/           # Konfigurasi Supabase URL & Service Role Key
│   │   ├── handlers/         # Controller HTTP (Triage, Queue, Clinic) & Unit Tests
│   │   ├── middleware/       # Validasi JWT Token dari Supabase Auth & Unit Tests
│   │   ├── models/           # Definisi Struct Data (Pasien, Klinik, Gejala)
│   │   ├── repository/       # Logika Kueri PostgreSQL & interfaces.go
│   │   └── services/         # Layanan Triage & Unit Tests
│   ├── go.mod
│   └── go.sum
├── frontend/                 # Porsi Kerja Rekan B (React App)
│   ├── src/
│   │   ├── assets/
│   │   ├── components/       # UI Kit Reusable (Navbar, Card, Button)
│   │   ├── context/          # AuthContext untuk menyimpan session user
│   │   ├── hooks/            # Custom hooks (e.g., useHaversine untuk lokasi)
│   │   ├── pages/
│   │   │   ├── Patient/      # Tampilan Cari Klinik, Form Triage, Token Antrean
│   │   │   ├── AdminKlinik/  # Dashboard Antrean (Merah, Kuning, Hijau)
│   │   │   └── SuperAdmin/   # Manajemen Tenant & Katalog Gejala
│   │   ├── services/         # API Fetcher ke Backend Go (bukan langsung ke Supabase)
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
└── supabase/                 # Konfigurasi & Migrasi Basis Data
    ├── migrations/           # File SQL untuk skema tabel & RLS (Row Level Security)
    └── seed.sql              # Data awal untuk katalog_gejala
```

---

## Deskripsi Direktori Utama

### 1. `backend/`
Berisi kode sumber untuk API server berbasis Go.
- **`cmd/api/main.go`**: Menginisialisasi koneksi database, router go-chi, mendaftarkan *middleware*, dan menjalankan peladen HTTP.
- **`internal/handlers/`**: Menangani permintaan HTTP, validasi payload JSON, dan mengirimkan respon JSON. Juga memuat berkas pengujian unit (`*_test.go`) dan `mock_test.go` untuk testing *in-memory*.
- **`internal/middleware/`**: Menyediakan middleware verifikasi JWT Supabase Auth untuk rute admin faskes.
- **`internal/models/`**: Berisi definisi objek Go struct yang dipetakan langsung ke tabel database dan request payload.
- **`internal/repository/`**: Abstraksi database menggunakan interface (`interfaces.go`) dan implementasi query SQL menggunakan connection pool `pgxpool`.
- **`internal/services/`**: Memuat triage engine bisnis logic untuk kalkulasi skor urgensi medis.

### 2. `frontend/`
Aplikasi web klien berbasis React.js (Vite).
- **`src/components/`**: Komponen UI yang dapat digunakan kembali secara konsisten (Buttons, Cards, Modals).
- **`src/pages/`**: Dibagi berdasarkan hak akses/aktor: Pasien (B2C), Admin Klinik (B2B), dan Super Admin.
- **`src/services/`**: Menangani seluruh komunikasi HTTP *fetch* ke Go Backend.

### 3. `supabase/`
Konfigurasi skema database relasional PostgreSQL.
- **`migrations/`**: Berisi file-file SQL migrasi beruntun untuk membangun skema tabel dan relasi secara bersih pada instance Supabase.
