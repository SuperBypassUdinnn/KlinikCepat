# Progres Pengembangan KlinikCepat

Dokumen ini mencatat kondisi implementasi aktual KlinikCepat agar pengembangan backend, frontend, database, dan dokumentasi tetap sinkron.

---

## 1. Status Saat Ini

**Pembaruan terakhir:** 24 Juni 2026

### Ringkasan

| Area                         | Status      |
| ---------------------------- | ----------- |
| Backend API Go               | Implemented |
| Database PostgreSQL Supabase | Implemented |
| Triage pasien                | Implemented |
| Autentikasi admin            | Implemented |
| Role-based access control    | Implemented |
| Isolasi data per klinik      | Implemented |
| Dashboard admin klinik       | Implemented |
| CRUD superadmin              | Implemented |
| Deployment production        | Belum       |
| Live tracking tiket          | Planned     |
| Estimasi waktu antrean       | Planned     |
| Analitik global              | Planned     |

---

## 2. Backend

### 2.1 Database dan Repository

- [x] Tabel `klinik`, `katalog_gejala`, `antrean`, dan `user_roles`.
- [x] Repository menggunakan PostgreSQL connection pool melalui `pgxpool`.
- [x] CRUD klinik dan katalog gejala.
- [x] Penyimpanan hasil triage ke tabel antrean.
- [x] Pengambilan antrean berdasarkan klinik dan status.
- [x] Pengurutan antrean berdasarkan prioritas:
  - Merah
  - Kuning
  - Hijau
  - waktu pendaftaran paling awal

- [x] Repository user membaca `role` dan `klinik_id`.
- [x] Update status antrean dibatasi berdasarkan klinik admin yang login.

### 2.2 Triage Engine

- [x] Perhitungan skor urgensi berdasarkan bobot gejala dan skala keparahan.
- [x] Klasifikasi status triage menjadi Merah, Kuning, atau Hijau.
- [x] Aturan khusus untuk gejala dengan tingkat urgensi tinggi.
- [x] Pembuatan antrean setelah proses triage berhasil.

### 2.3 Autentikasi dan Otorisasi

- [x] Autentikasi admin menggunakan Supabase Auth.
- [x] Verifikasi JWT Supabase menggunakan ES256 dan public key dari JWKS.
- [x] Validasi issuer, audience, expiration, dan subject token.
- [x] Endpoint `GET /api/v1/auth/me`.
- [x] Role aplikasi:
  - `superadmin`
  - `klinik_admin`

- [x] Middleware `RequireRole`.
- [x] Role dan `klinik_id` dimasukkan ke request context.
- [x] Admin klinik hanya dapat membaca antrean kliniknya sendiri.
- [x] Admin klinik hanya dapat memperbarui antrean kliniknya sendiri.
- [x] Superadmin tidak terikat pada satu klinik.

### 2.4 Constraint Database

- [x] `superadmin` harus memiliki `klinik_id = NULL`.
- [x] `klinik_admin` wajib memiliki `klinik_id`.
- [x] Foreign key `user_roles.klinik_id` menggunakan `ON DELETE RESTRICT`.
- [x] Klinik yang masih memiliki admin tidak dapat dihapus.
- [x] Index untuk kolom `user_roles.klinik_id`.

### 2.5 Endpoint Backend

#### Publik

- [x] `GET /api/v1/klinik`
- [x] `GET /api/v1/klinik/{id}`
- [x] `GET /api/v1/gejala`
- [x] `GET /api/v1/gejala/{id}`
- [x] `POST /api/v1/triage`

#### User terautentikasi

- [x] `GET /api/v1/auth/me`

#### Admin Klinik

- [x] `GET /api/v1/admin/antrean?status=Menunggu`
- [x] `PUT /api/v1/admin/antrean/{id}/status`

`klinik_admin` tidak mengirimkan `klinik_id` secara bebas. Backend mengambil `klinik_id` dari akun yang terautentikasi.

#### Superadmin

- [x] Tambah, edit, dan hapus klinik.
- [x] Tambah, edit, dan hapus katalog gejala.
- [x] Mengakses data antrean berdasarkan klinik yang dipilih.

### 2.6 Testing Backend

- [x] Unit test handler.
- [x] Unit test triage service.
- [x] Unit test middleware autentikasi.
- [x] Unit test token ES256 menggunakan key lokal.
- [x] Unit test akses antrean berdasarkan tenant.
- [x] Unit test penolakan update antrean milik klinik lain.
- [x] Manual integration test menggunakan JWT Supabase asli.
- [x] Verifikasi request tanpa token menghasilkan `401 Unauthorized`.
- [x] Verifikasi akses tenant lain ditolak.
- [ ] Continuous Integration melalui GitHub Actions.

---

## 3. Frontend

### 3.1 Halaman Pasien

- [x] Halaman pencarian klinik.
- [x] Pengambilan daftar klinik dari backend.
- [x] Penggunaan lokasi pengguna dan kalkulasi jarak Haversine.
- [x] Form triage berdasarkan katalog gejala dari backend.
- [x] Pengiriman hasil triage ke backend.
- [x] Halaman tiket hasil triage.

### 3.2 Admin Klinik

- [x] Login menggunakan Supabase Auth.
- [x] Dashboard antrean.
- [x] Filter status:
  - Menunggu
  - Selesai
  - Dilewati

- [x] Statistik antrean berdasarkan status triage.
- [x] Aksi menyelesaikan dan melewati antrean.
- [x] Auto-refresh antrean setiap 10 detik.
- [x] Dashboard tidak menampilkan dropdown seluruh klinik.
- [x] Dashboard hanya mengambil antrean klinik milik admin yang login.

### 3.3 Superadmin

- [x] Halaman manajemen klinik.
- [x] Tambah klinik.
- [x] Edit klinik.
- [x] Hapus klinik.
- [x] Halaman manajemen katalog gejala.
- [x] Tambah gejala.
- [x] Edit gejala.
- [x] Hapus gejala.

### 3.4 AuthContext dan API Layer

- [x] Session dikelola melalui Supabase Auth.
- [x] Access token dibaca langsung dari session Supabase.
- [x] Tidak menggunakan salinan token manual di `localStorage`.
- [x] `AuthContext` menyimpan:
  - user
  - profile
  - role
  - clinicId
  - authError

- [x] Profile aplikasi diambil melalui `GET /api/v1/auth/me`.
- [x] Redirect setelah login berdasarkan role.
- [x] Session tetap aktif setelah browser di-refresh.
- [x] Logout membersihkan state autentikasi.

### 3.5 Route dan Navigasi

- [x] Route publik untuk pasien.
- [x] Route admin klinik hanya dapat diakses oleh `klinik_admin`.
- [x] Route superadmin hanya dapat diakses oleh `superadmin`.
- [x] Pengguna tanpa session diarahkan ke halaman login.
- [x] Pengguna dengan role yang salah diarahkan ke dashboard miliknya.
- [x] Navbar menampilkan menu berdasarkan role.
- [x] Admin klinik tidak melihat menu superadmin.
- [x] Superadmin tidak melihat menu dashboard admin klinik.

---

## 4. Fitur yang Masih Partial

### 4.1 Tiket Antrean

Tiket hasil triage sudah dapat ditampilkan, tetapi datanya masih bergantung pada state navigasi frontend.

Keterbatasan:

- tiket dapat hilang setelah halaman di-refresh;
- tiket belum dapat dibuka kembali melalui URL permanen;
- belum tersedia endpoint publik untuk membaca status tiket;
- belum tersedia posisi antrean secara real-time.

### 4.2 Konfigurasi Production

Frontend sudah mendukung konfigurasi URL API melalui environment variable, tetapi deployment production belum dilakukan.

Masih diperlukan:

- konfigurasi domain frontend dan backend;
- konfigurasi CORS atau reverse proxy;
- environment variable production;
- HTTPS;
- logging production;
- health monitoring.

### 4.3 Administrasi Akun

Akun admin dibuat melalui Supabase Auth dan kemudian diberi role melalui tabel `user_roles`.

Belum tersedia:

- halaman pembuatan akun admin;
- halaman pengaturan role;
- halaman pengaitan admin dengan klinik;
- mekanisme reset password dari dashboard aplikasi.

---

## 5. Fitur Planned

Fitur berikut belum dianggap selesai dan tidak boleh didokumentasikan sebagai fitur aktif:

- [ ] Akun pasien.
- [ ] Registrasi pasien.
- [ ] Riwayat kunjungan pasien.
- [ ] Tiket antrean permanen.
- [ ] Live tracking posisi antrean.
- [ ] Estimasi waktu tunggu.
- [ ] Status antrean `Dipanggil`.
- [ ] Notifikasi pasien.
- [ ] Dashboard analitik global.
- [ ] Statistik harian dari endpoint agregasi.
- [ ] Search klinik sebagai fallback ketika izin lokasi ditolak.
- [ ] Continuous Integration.
- [ ] Deployment production.

---

## 6. Prioritas Pengerjaan Berikutnya

### Prioritas 1 — Tiket Persisten

- Membuat endpoint untuk mengambil tiket berdasarkan ID atau token publik.
- Mengubah URL tiket menjadi `/ticket/:id`.
- Memastikan tiket dapat dibuka kembali setelah refresh.
- Membatasi data pasien yang dikembalikan endpoint publik.

### Prioritas 2 — Kesiapan Deployment

- Menentukan platform deployment backend.
- Menentukan platform deployment frontend.
- Mengatur environment variable production.
- Mengatur CORS atau reverse proxy.
- Menambahkan health check dan logging.

### Prioritas 3 — Testing Otomatis

- Menambahkan GitHub Actions.
- Menjalankan `go test ./...`.
- Menjalankan `go vet ./...`.
- Menjalankan build frontend.
- Menggagalkan pull request apabila test atau build gagal.

### Prioritas 4 — Dokumentasi

- Sinkronisasi `user_roles_klinikcepat.md`.
- Sinkronisasi `frontend_integration.md`.
- Sinkronisasi `backend_architecture.md`.
- Sinkronisasi `blueprint_klinikcepat.md`.
- Sinkronisasi `ProjectStructure.md`.

---

## 7. Perintah Validasi

### Backend

```bash
cd backend
go test ./...
go vet ./...
```

### Frontend

```bash
cd frontend
npm run build
```

Seluruh perubahan harus melewati validasi tersebut sebelum digabungkan ke branch utama.