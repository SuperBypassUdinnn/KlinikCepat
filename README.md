# KlinikCepat
### B2B Micro-SaaS Triage & Manajemen Antrean Faskes
> **Tema SDGs 3:** Good Health and Well-being (Kehidupan Sehat dan Sejahtera)  
> **Target Proyek:** UAS Praktikum Pemrograman Berbasis Web  

---

## Deskripsi Proyek
**KlinikCepat** adalah platform manajemen antrean cerdas berbasis *multi-tenant* yang dirancang untuk Fasilitas Kesehatan (Faskes) Tingkat Pertama independen atau poliklinik institusi. Berbeda dengan sistem antrean konvensional yang mengandalkan model *First-In-First-Out* (FIFO) murni, aplikasi ini bertindak sebagai gerbang *triage* digital otomatis.

Sistem mengkalkulasi skor urgensi pasien berdasarkan input kuesioner gejala secara *real-time* sebelum mereka tiba di lokasi, lalu memprioritaskan penanganan medis darurat (**Status Merah**) di atas keluhan ringan (**Status Hijau**). Pendekatan ini secara langsung mereduksi penumpukan pasien berisiko tinggi di ruang tunggu fisik dan mengoptimalkan manajemen waktu tenaga medis.

---

## 📂 Struktur Proyek
Aplikasi ini dikembangkan dalam repositori tunggal (*monorepo*) dengan struktur berikut:
```text
KlinikCepat/
├── backend/                  # RESTful API Go (go-chi, pgxpool)
│   ├── cmd/api/main.go       # Entrypoint aplikasi backend
│   └── internal/             # Domain logic (Models, Repositories, Services, Handlers)
├── docs/                     # Berkas panduan, blueprint, dan arsitektur proyek
│   ├── blueprint_klinikcepat.md  # Konsep dasar & rumus matematis triage
│   ├── backend_architecture.md   # Panduan teknis lengkap & daftar API endpoint
│   ├── progress.md               # Pelacakan progres fitur
│   └── user_roles_klinikcepat.md # Pembagian batasan peran pengguna
└── supabase/                 # Konfigurasi database remote
    └── migrations/           # File migrasi SQL untuk skema tabel & seed data
```

---

## Memulai (Panduan Instalasi Backend)

### 1. Prasyarat
Pastikan Anda sudah menginstal:
- **Golang** (v1.26+)
- Akun dan Proyek aktif di **Supabase**

### 2. Konfigurasi Database & Migrasi
Jika Anda ingin menerapkan ulang skema basis data di Supabase:
1. Pastikan Supabase CLI sudah terpasang. Tautkan dengan proyek Supabase Anda:
   ```bash
   supabase link --project-ref <PROJECT-ID-ANDA>
   ```
2. Jalankan perintah migrasi untuk membuat tabel (`klinik`, `katalog_gejala`, `antrean`) beserta data seed awal:
   ```bash
   supabase db push
   ```
   *Alternatif:* Salin kode di dalam berkas [00001_initial_schema.sql](file:///home/superbypassudin/.clone/Github/KlinikCepat/supabase/migrations/00001_initial_schema.sql) dan tempel di **SQL Editor** pada dashboard web Supabase Anda, lalu klik **Run**.

### 3. Jalankan Backend
1. Masuk ke direktori `backend/`:
   ```bash
   cd backend
   ```
2. Buat file `.env` (salin dari `.env.example`) dan isi variabel environment Anda:
   ```env
   DATABASE_URL="postgresql://postgres.<project-id>:<password>@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
   PORT=8080
   SUPABASE_JWT_SECRET="JWT_SECRET_DARI_API_SETTINGS_SUPABASE"
   ```
3. Unduh dependensi modul Go:
   ```bash
   go mod tidy
   ```
4. Jalankan server backend lokal Anda:
   ```bash
   go run cmd/api/main.go
   ```
   Server backend sekarang aktif di `http://localhost:8080`. Anda dapat mengakses healthcheck di `GET http://localhost:8080/health`.

---

## Dokumentasi Tambahan
Untuk panduan pengembangan lebih lanjut, Anda dapat merujuk ke berkas dokumentasi internal kami:
- [Struktur Organisasi Proyek](file:///home/superbypassudin/.clone/Github/KlinikCepat/docs/ProjectStructure.md)
- [Cetak Biru Proyek (Blueprint)](file:///home/superbypassudin/.clone/Github/KlinikCepat/docs/blueprint_klinikcepat.md)
- [Arsitektur & Spesifikasi API Backend](file:///home/superbypassudin/.clone/Github/KlinikCepat/docs/backend_architecture.md)
- [Arsitektur Peran & Alur Pengguna (User Roles)](file:///home/superbypassudin/.clone/Github/KlinikCepat/docs/user_roles_klinikcepat.md)
- [Status Progress Pengembangan](file:///home/superbypassudin/.clone/Github/KlinikCepat/docs/progress.md)
