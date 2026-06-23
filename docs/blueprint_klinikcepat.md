# Blueprint Sistem KlinikCepat

Dokumen ini menjadi blueprint teknis dan fungsional KlinikCepat berdasarkan implementasi aktual.

**Pembaruan terakhir:** 24 Juni 2026

---

## 1. Ringkasan Produk

KlinikCepat adalah sistem antrean klinik berbasis triage digital.

Pasien tidak hanya ditempatkan berdasarkan urutan kedatangan, tetapi juga berdasarkan tingkat urgensi kondisi.

Sistem memiliki tiga aktor utama:

1. Pasien
2. Admin Klinik
3. Superadmin

---

## 2. Tujuan Sistem

KlinikCepat bertujuan untuk:

* membantu pasien menemukan klinik;
* melakukan penyaringan awal kondisi pasien;
* menentukan prioritas antrean;
* membantu klinik mengelola antrean;
* memisahkan data antar-klinik;
* menyediakan pengelolaan klinik dan gejala secara terpusat.

---

## 3. Status Implementasi

### Implemented

* pencarian klinik;
* kalkulasi jarak menggunakan Haversine;
* katalog gejala;
* form triage;
* klasifikasi Merah, Kuning, dan Hijau;
* pembuatan antrean;
* login admin;
* role `klinik_admin`;
* role `superadmin`;
* dashboard antrean;
* update status antrean;
* CRUD klinik;
* CRUD katalog gejala;
* isolasi data antar-klinik;
* role-based routing frontend;
* JWT Supabase ES256;
* verifikasi JWT melalui JWKS.

### Partial

* tiket hasil triage;
* konfigurasi deployment;
* statistik dashboard;
* administrasi akun admin.

### Planned

* akun pasien;
* tiket permanen;
* live tracking posisi antrean;
* estimasi waktu tunggu;
* status `Dipanggil`;
* notifikasi pasien;
* analitik global;
* laporan agregat;
* WebSocket;
* deployment production.

---

## 4. Arsitektur Sistem

```text
Pasien
  │
  ▼
React Frontend
  │
  ├── Supabase Auth
  │       └── Admin Klinik dan Superadmin
  │
  ▼
Go REST API
  │
  ├── AuthMiddleware
  ├── RequireRole
  ├── Triage Service
  ├── Repository
  │
  ▼
PostgreSQL Supabase
```

---

## 5. Teknologi

### Frontend

* React
* Vite
* React Router
* Supabase JavaScript Client
* CSS
* React Icons

### Backend

* Go
* Chi Router
* pgxpool
* JWT ES256
* JWKS
* REST API

### Database dan Auth

* Supabase PostgreSQL
* Supabase Auth

---

## 6. Aktor Sistem

### 6.1 Pasien

Pasien menggunakan fitur publik tanpa login.

Pasien dapat:

* melihat klinik;
* melihat jarak klinik;
* mengisi triage;
* mendapatkan tiket hasil triage.

Pasien tidak dapat:

* melihat seluruh antrean;
* mengubah status antrean;
* mengakses dashboard admin.

### 6.2 Admin Klinik

Admin Klinik terikat pada satu klinik.

```text
role      = klinik_admin
klinik_id = UUID klinik
```

Admin Klinik dapat:

* login;
* melihat antrean kliniknya;
* memfilter status antrean;
* menyelesaikan antrean;
* melewati antrean;
* logout.

Admin Klinik tidak dapat:

* memilih klinik lain;
* melihat antrean klinik lain;
* mengubah antrean klinik lain;
* mengakses fitur Superadmin.

### 6.3 Superadmin

Superadmin tidak terikat pada satu klinik.

```text
role      = superadmin
klinik_id = NULL
```

Superadmin dapat:

* mengelola klinik;
* mengelola katalog gejala;
* mengakses antrean berdasarkan klinik;
* mengakses halaman Superadmin.

---

## 7. Alur Pasien

```text
Buka aplikasi
→ izinkan lokasi
→ ambil daftar klinik
→ hitung jarak
→ pilih klinik
→ pilih gejala
→ isi skala keparahan
→ kirim triage
→ backend hitung skor
→ backend tentukan status
→ backend buat antrean
→ frontend tampilkan tiket
```

---

## 8. Pencarian Klinik

Frontend mengambil daftar klinik melalui:

```http
GET /api/v1/klinik
```

Lokasi pengguna diperoleh melalui browser Geolocation API.

Jarak dihitung menggunakan formula Haversine.

Formula:

```text
a =
sin²(Δlat / 2)
+
cos(lat1)
× cos(lat2)
× sin²(Δlng / 2)

c =
2 × atan2(√a, √(1-a))

d =
R × c
```

Keterangan:

```text
d = jarak
R = radius bumi
```

---

## 9. Triage Digital

Pasien memilih gejala dan tingkat keparahan.

Setiap gejala memiliki bobot urgensi.

Contoh:

```text
Gejala A
bobot urgensi = 8

Skala keparahan = 2

Skor gejala = 8 × 2 = 16
```

Total skor:

```text
total_skor =
jumlah seluruh skor gejala
```

---

## 10. Klasifikasi Triage

Backend mengelompokkan pasien menjadi:

### Merah

Kondisi paling darurat.

Pasien dengan status Merah ditempatkan pada prioritas paling tinggi.

### Kuning

Kondisi membutuhkan penanganan, tetapi tidak seberat Merah.

### Hijau

Kondisi relatif stabil dan memiliki prioritas lebih rendah.

Aturan nilai skor harus mengikuti implementasi pada `TriageService`.

Dokumen tidak boleh mendefinisikan threshold berbeda dari kode.

---

## 11. Data Antrean

Struktur utama antrean:

```text
id
klinik_id
nama_pasien
total_skor
status_triage
status_antrean
created_at
```

Nilai status triage:

```text
Merah
Kuning
Hijau
```

Nilai status antrean:

```text
Menunggu
Selesai
Dilewati
```

Status `Dipanggil` belum tersedia.

---

## 12. Pengurutan Antrean

Antrean diurutkan berdasarkan:

1. Merah
2. Kuning
3. Hijau
4. waktu pendaftaran paling awal

Query:

```sql
ORDER BY
  CASE status_triage
    WHEN 'Merah' THEN 1
    WHEN 'Kuning' THEN 2
    WHEN 'Hijau' THEN 3
    ELSE 4
  END ASC,
  created_at ASC;
```

Jangan menggunakan:

```sql
ORDER BY status_triage DESC;
```

karena urutan string atau enum tidak menjamin Merah berada paling awal.

---

## 13. Dashboard Admin Klinik

Dashboard menyediakan:

* total antrean;
* jumlah triage Merah;
* jumlah triage Kuning;
* jumlah triage Hijau;
* filter status;
* auto-refresh setiap 10 detik;
* aksi Selesai;
* aksi Lewati.

Filter status:

```text
Menunggu
Selesai
Dilewati
```

Dashboard tidak memiliki dropdown untuk memilih seluruh klinik.

---

## 14. Polling Antrean

Dashboard melakukan polling setiap 10 detik.

```text
Frontend
→ GET antrean
→ tunggu 10 detik
→ GET antrean kembali
```

Polling ini bukan WebSocket dan tidak dikategorikan sebagai real-time penuh.

---

## 15. Autentikasi

Admin login melalui Supabase Auth.

Alur:

```text
Email + password
→ Supabase Auth
→ session
→ access token
→ request ke backend
```

Frontend mengirim:

```http
Authorization: Bearer <access_token>
```

---

## 16. Verifikasi JWT

Backend memverifikasi token menggunakan:

* algoritma ES256;
* JWKS Supabase;
* issuer;
* audience `authenticated`;
* expiration;
* subject;
* signature.

Alur:

```text
JWT
→ baca kid
→ ambil public key
→ verifikasi ES256
→ validasi claims
→ request diteruskan
```

---

## 17. Role Aplikasi

Role Supabase seperti:

```text
authenticated
```

bukan role aplikasi KlinikCepat.

Role aplikasi disimpan pada tabel:

```text
user_roles
```

Role yang didukung:

```text
superadmin
klinik_admin
```

---

## 18. Endpoint Profil User

Frontend mengambil profil aplikasi melalui:

```http
GET /api/v1/auth/me
```

Respons Admin Klinik:

```json
{
  "id": "uuid-user",
  "email": "admin@klinik.com",
  "role": "klinik_admin",
  "klinik_id": "uuid-klinik"
}
```

Respons Superadmin:

```json
{
  "id": "uuid-user",
  "email": "superadmin@klinik.com",
  "role": "superadmin",
  "klinik_id": null
}
```

---

## 19. Multi-Tenant Isolation

Setiap Admin Klinik hanya dapat mengakses satu klinik.

Backend menjadi sumber kebenaran untuk tenant scope.

Alur:

```text
JWT
→ user_id
→ tabel user_roles
→ role + klinik_id
→ request context
→ query berdasarkan klinik_id
```

Frontend tidak menentukan tenant secara bebas.

---

## 20. Membaca Antrean

Admin Klinik menggunakan:

```http
GET /api/v1/admin/antrean?status=Menunggu
```

Admin Klinik tidak mengirimkan `klinik_id`.

Backend mengambil `klinik_id` dari akun yang login.

Jika Admin mencoba meminta klinik lain:

```text
403 Forbidden
```

---

## 21. Memperbarui Antrean

Endpoint:

```http
PUT /api/v1/admin/antrean/{id}/status
```

Payload:

```json
{
  "status": "Selesai"
}
```

Untuk Admin Klinik, query dibatasi:

```sql
UPDATE antrean
SET status_antrean = $1
WHERE id = $2
  AND klinik_id = $3;
```

Jika antrean bukan milik klinik Admin:

```text
404 Not Found
```

---

## 22. Akses Superadmin

Superadmin dapat membaca antrean klinik tertentu melalui:

```http
GET /api/v1/admin/antrean
  ?klinik_id={uuid-klinik}
  &status=Menunggu
```

Superadmin wajib menyebut klinik yang ingin dibaca.

Jika tidak:

```text
400 Bad Request
```

---

## 23. Role-Based Routing Frontend

Route Admin Klinik:

```text
/admin/dashboard
```

Hanya dapat diakses oleh:

```text
klinik_admin
```

Route Superadmin:

```text
/superadmin/klinik
/superadmin/gejala
```

Hanya dapat diakses oleh:

```text
superadmin
```

Jika role salah, frontend mengarahkan user ke dashboard sesuai role.

Backend tetap memeriksa role pada setiap endpoint.

---

## 24. Navigasi

### Publik

```text
Cari Klinik
Login Admin
```

### Admin Klinik

```text
Dashboard
Halaman Pasien
Logout
```

### Superadmin

```text
Kelola Klinik
Kelola Gejala
Halaman Pasien
Logout
```

---

## 25. Database Constraint

Aturan tabel `user_roles`:

```text
superadmin
→ klinik_id harus NULL

klinik_admin
→ klinik_id wajib diisi
```

Foreign key:

```text
user_roles.klinik_id
→ klinik.id
→ ON DELETE RESTRICT
```

Klinik tidak dapat dihapus jika masih memiliki Admin Klinik.

---

## 26. Endpoint Publik

```http
GET /api/v1/klinik
GET /api/v1/klinik/{id}
GET /api/v1/gejala
GET /api/v1/gejala/{id}
POST /api/v1/triage
```

---

## 27. Endpoint Terautentikasi

```http
GET /api/v1/auth/me
GET /api/v1/admin/antrean
PUT /api/v1/admin/antrean/{id}/status
```

---

## 28. Endpoint Superadmin

```http
POST /api/v1/klinik
PUT /api/v1/klinik/{id}
DELETE /api/v1/klinik/{id}

POST /api/v1/gejala
PUT /api/v1/gejala/{id}
DELETE /api/v1/gejala/{id}
```

---

## 29. Error Handling

Status utama:

```text
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
500 Internal Server Error
```

### 401

Token tidak tersedia atau tidak valid.

### 403

User valid tetapi tidak memiliki hak akses.

### 404

Data tidak ditemukan atau berada di luar tenant scope.

---

## 30. Testing

Backend diuji menggunakan:

```bash
go test ./...
go vet ./...
```

Frontend diuji menggunakan:

```bash
npm run build
```

Manual test mencakup:

* login Admin Klinik;
* login Superadmin;
* pembatasan route;
* pembatasan antrean klinik;
* update antrean tenant lain;
* request tanpa token;
* refresh session;
* logout.

---

## 31. Batasan Sistem Saat Ini

### Tiket Pasien

Tiket belum persisten.

Masalah:

* hilang setelah refresh;
* tidak memiliki URL permanen;
* tidak dapat dicek ulang;
* belum memiliki posisi antrean.

### Statistik

Statistik dashboard dihitung dari data antrean yang sedang dimuat.

Belum tersedia endpoint agregasi harian.

### Real-Time

Sistem masih menggunakan polling 10 detik.

Belum menggunakan WebSocket.

### Akun Admin

Pembuatan akun dan role masih dilakukan melalui Supabase dan database.

Belum ada UI manajemen akun.

---

## 32. Roadmap

### Fase 1 — MVP Dasar

Status: Implemented

* pencarian klinik;
* triage;
* antrean;
* dashboard admin;
* CRUD superadmin.

### Fase 2 — Keamanan dan Multi-Tenant

Status: Implemented

* JWT ES256;
* JWKS;
* role-based authorization;
* tenant isolation;
* constraint database;
* role-based frontend routing.

### Fase 3 — Tiket Persisten

Status: Planned

* endpoint tiket;
* URL permanen;
* status antrean pasien;
* proteksi data pasien.

### Fase 4 — Real-Time

Status: Planned

* posisi antrean;
* estimasi waktu;
* WebSocket;
* notifikasi.

### Fase 5 — Deployment

Status: Planned

* deployment backend;
* deployment frontend;
* CORS;
* HTTPS;
* logging;
* monitoring;
* CI/CD.

---

## 33. Prinsip Arsitektur

KlinikCepat mengikuti prinsip:

```text
Frontend untuk UX
Backend untuk otorisasi
Database untuk integritas
```

Frontend dapat menyembunyikan menu, tetapi tidak menentukan izin.

Backend menentukan:

* siapa pengguna;
* role pengguna;
* klinik pengguna;
* data yang boleh diakses.

Database menjaga:

* relasi;
* constraint;
* integritas tenant.
