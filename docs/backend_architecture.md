# Arsitektur Backend KlinikCepat

Dokumen ini menjelaskan arsitektur backend KlinikCepat berdasarkan implementasi aktual.

**Pembaruan terakhir:** 24 Juni 2026

---

## 1. Ringkasan

Backend KlinikCepat dibangun menggunakan Go dan menyediakan REST API untuk:

* data klinik;
* katalog gejala;
* proses triage;
* antrean pasien;
* autentikasi admin;
* otorisasi berdasarkan role;
* isolasi data berdasarkan klinik;
* CRUD master data oleh Superadmin.

Backend menjadi sumber kebenaran utama untuk:

* perhitungan skor triage;
* klasifikasi status triage;
* validasi JWT;
* pemeriksaan role;
* penentuan klinik Admin;
* pembatasan akses data antrean.

Frontend tidak dipercaya untuk menentukan sendiri role atau klinik yang boleh diakses.

---

## 2. Teknologi

Backend menggunakan:

* Go
* Chi Router
* PostgreSQL
* Supabase Database
* Supabase Auth
* `pgxpool`
* JWT ES256
* Supabase JWKS
* environment variable
* unit testing bawaan Go

---

## 3. Struktur Backend

Struktur utama:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   └── services/
├── go.mod
├── go.sum
└── .env
```

### `cmd/api`

Berisi entry point aplikasi:

* membaca environment variable;
* membuat koneksi database;
* membuat repository;
* membuat service;
* membuat handler;
* mendaftarkan route;
* menjalankan HTTP server.

### `internal/handlers`

Bertanggung jawab atas:

* membaca HTTP request;
* validasi input dasar;
* membaca path dan query parameter;
* memanggil service atau repository;
* mengubah hasil menjadi JSON response;
* menentukan HTTP status code.

### `internal/middleware`

Berisi:

* autentikasi JWT;
* verifikasi token Supabase;
* penyimpanan claims ke request context;
* pemeriksaan role aplikasi;
* penyimpanan `klinik_id` ke request context.

### `internal/models`

Berisi model data dan payload:

* `Klinik`
* `Gejala`
* `Antrean`
* `UserAccess`
* `AuthMeResponse`
* `TriageRequest`
* `TriageResponse`

### `internal/repository`

Bertanggung jawab atas akses database:

* query klinik;
* query gejala;
* query antrean;
* query user role;
* tenant-scoped update;
* mapping row database ke model.

### `internal/services`

Berisi business logic, terutama:

* perhitungan skor triage;
* klasifikasi Merah, Kuning, atau Hijau;
* validasi gejala;
* pembuatan antrean.

---

## 4. Alur Request

Alur request publik:

```text
Client
→ Router
→ Handler
→ Service atau Repository
→ PostgreSQL
→ JSON Response
```

Alur request terautentikasi:

```text
Client
→ Authorization Header
→ AuthMiddleware
→ Verifikasi JWT Supabase
→ RequireRole
→ Query user_roles
→ Role + klinik_id
→ Handler
→ Repository tenant-scoped
→ PostgreSQL
→ JSON Response
```

---

## 5. Konfigurasi Environment

Backend menggunakan environment variable seperti:

```env
PORT=8080
DATABASE_URL=
SUPABASE_URL=
```

Keterangan:

### `PORT`

Port HTTP server.

Contoh:

```env
PORT=8080
```

### `DATABASE_URL`

Connection string PostgreSQL Supabase.

### `SUPABASE_URL`

URL project Supabase.

Contoh:

```env
SUPABASE_URL=https://project-ref.supabase.co
```

Nilai ini digunakan untuk membentuk endpoint JWKS:

```text
https://project-ref.supabase.co/auth/v1/.well-known/jwks.json
```

File `.env` tidak boleh dimasukkan ke Git.

---

## 6. Database

Tabel utama:

```text
klinik
katalog_gejala
antrean
user_roles
```

---

## 7. Tabel `klinik`

Menyimpan data fasilitas kesehatan.

Atribut utama:

```text
id
nama_klinik
lat
lng
kapasitas_aktif
created_at
```

Digunakan untuk:

* pencarian klinik;
* perhitungan jarak;
* tujuan pendaftaran antrean;
* relasi dengan Admin Klinik.

---

## 8. Tabel `katalog_gejala`

Menyimpan gejala dan bobot urgensi.

Atribut utama:

```text
id
nama_gejala
bobot_urgensi
created_at
```

Bobot urgensi digunakan dalam perhitungan triage.

Perubahan bobot hanya memengaruhi triage baru.

Data antrean lama tidak dihitung ulang secara otomatis.

---

## 9. Tabel `antrean`

Menyimpan hasil pendaftaran pasien.

Atribut utama:

```text
id
klinik_id
nama_pasien
total_skor
status_triage
status_antrean
created_at
```

Nilai `status_triage`:

```text
Merah
Kuning
Hijau
```

Nilai `status_antrean`:

```text
Menunggu
Selesai
Dilewati
```

Status `Dipanggil` belum diimplementasikan.

---

## 10. Tabel `user_roles`

Menyimpan role aplikasi dan scope klinik.

Atribut utama:

```text
user_id
role
klinik_id
```

Role yang didukung:

```text
superadmin
klinik_admin
```

Aturan:

```text
superadmin
→ klinik_id harus NULL

klinik_admin
→ klinik_id wajib terisi
```

Constraint database memastikan kombinasi role dan `klinik_id` tetap valid.

Foreign key `klinik_id` menggunakan:

```text
ON DELETE RESTRICT
```

Klinik tidak dapat dihapus selama masih memiliki Admin Klinik.

---

## 11. Repository Pattern

Handler tidak menulis SQL secara langsung.

Handler berkomunikasi melalui interface repository.

Contoh kontrak:

```go
type UserRepository interface {
	GetUserAccess(
		ctx context.Context,
		userID string,
	) (*models.UserAccess, error)
}
```

`GetUserAccess` mengembalikan:

```go
type UserAccess struct {
	UserID   string
	Role     string
	KlinikID *string
}
```

`KlinikID` menggunakan pointer karena Superadmin memiliki nilai `NULL`.

---

## 12. Triage Service

Triage diproses oleh service layer.

Alur:

```text
TriageRequest
→ validasi klinik
→ validasi gejala
→ ambil bobot gejala
→ hitung skor
→ tentukan status triage
→ buat antrean
→ TriageResponse
```

Request:

```json
{
  "klinik_id": "uuid-klinik",
  "nama_pasien": "Nama Pasien",
  "gejala": [
    {
      "gejala_id": "uuid-gejala",
      "skala_keparahan": 2
    }
  ]
}
```

Respons:

```json
{
  "antrean_id": "uuid-antrean",
  "status_triage": "Kuning",
  "total_skor": 12,
  "pesan": "Pendaftaran antrean berhasil"
}
```

---

## 13. Pengurutan Antrean

Antrean diurutkan berdasarkan prioritas:

1. Merah
2. Kuning
3. Hijau
4. waktu pendaftaran paling awal

Query menggunakan ekspresi `CASE`.

Contoh:

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

Pengurutan tidak menggunakan:

```sql
ORDER BY status_triage DESC;
```

karena urutan string atau enum tidak menjamin Merah berada paling awal.

---

## 14. Autentikasi JWT

Supabase Auth menghasilkan access token JWT.

Frontend mengirim token melalui:

```http
Authorization: Bearer <access_token>
```

Backend membaca header tersebut melalui `AuthMiddleware`.

---

## 15. Verifikasi ES256 dan JWKS

Token Supabase diverifikasi menggunakan algoritma:

```text
ES256
```

Backend mengambil public signing key dari endpoint JWKS Supabase.

Alur:

```text
JWT
→ baca header dan kid
→ ambil public key dari JWKS
→ verifikasi signature ES256
→ validasi issuer
→ validasi audience
→ validasi expiration
→ validasi subject
```

Validasi mencakup:

* algoritma harus ES256;
* issuer harus sesuai project Supabase;
* audience harus `authenticated`;
* token belum expired;
* subject tersedia;
* signature valid.

Backend tidak menggunakan HMAC `SUPABASE_JWT_SECRET` untuk token ES256.

---

## 16. Claims Internal

Backend menyimpan claims dalam request context.

Contoh:

```go
type JWTClaims struct {
	Sub      string
	Email    string
	Role     string
	Exp      int64
	KlinikID *string
}
```

Pada tahap awal:

```text
Sub
Email
Exp
```

berasal dari JWT.

Setelah `RequireRole`:

```text
Role
KlinikID
```

diisi berdasarkan tabel `user_roles`.

Role `authenticated` pada JWT Supabase bukan role aplikasi KlinikCepat.

---

## 17. Middleware `RequireRole`

`RequireRole` menerima daftar role yang diperbolehkan.

Contoh:

```go
r.Use(
	middleware.RequireRole(
		repo,
		"superadmin",
		"klinik_admin",
	),
)
```

Alur:

```text
claims.Sub
→ GetUserAccess()
→ role + klinik_id
→ validasi role
→ validasi klinik_admin punya klinik_id
→ simpan ke context
→ lanjut handler
```

Jika user tidak terdapat pada `user_roles`, backend mengembalikan:

```text
403 Forbidden
```

Jika role tidak sesuai route:

```text
403 Forbidden
```

---

## 18. Endpoint Profil User

Endpoint:

```http
GET /api/v1/auth/me
```

Membutuhkan JWT valid.

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

Frontend menggunakan endpoint ini untuk:

* menyimpan role;
* menentukan redirect;
* melindungi route;
* menampilkan menu berdasarkan role.

---

## 19. Isolasi Tenant Admin Klinik

Admin Klinik hanya boleh mengakses satu klinik.

Backend tidak mempercayai `klinik_id` dari frontend.

### Membaca Antrean

Request Admin Klinik:

```http
GET /api/v1/admin/antrean?status=Menunggu
```

Backend mengambil `klinik_id` dari request context.

Jika frontend mencoba:

```http
GET /api/v1/admin/antrean?klinik_id=klinik-lain
```

backend membandingkan ID tersebut dengan scope Admin.

Jika berbeda:

```text
403 Forbidden
```

### Memperbarui Antrean

Repository menerima scope klinik:

```go
UpdateStatusAntrean(
	ctx,
	id,
	status,
	klinikID,
)
```

Untuk Admin Klinik, query menggunakan:

```sql
UPDATE antrean
SET status_antrean = $1
WHERE id = $2
  AND klinik_id = $3;
```

Jika antrean berada di klinik lain:

```text
RowsAffected = 0
```

Backend mengembalikan:

```text
404 Not Found
```

Respons `404` mencegah Admin mengetahui apakah antrean klinik lain benar-benar ada.

---

## 20. Akses Superadmin

Superadmin tidak memiliki `klinik_id`.

Untuk melihat antrean, Superadmin menentukan klinik melalui query:

```http
GET /api/v1/admin/antrean
  ?klinik_id={uuid-klinik}
  &status=Menunggu
```

Jika `klinik_id` tidak diberikan:

```text
400 Bad Request
```

Untuk update status antrean, Superadmin dapat memperbarui berdasarkan ID tanpa tenant restriction.

Endpoint tetap dilindungi oleh role `superadmin`.

---

## 21. Daftar Endpoint

### Health Check

```http
GET /health
```

### Publik

```http
GET /api/v1/klinik
GET /api/v1/klinik/{id}
GET /api/v1/gejala
GET /api/v1/gejala/{id}
POST /api/v1/triage
```

### User Terautentikasi

```http
GET /api/v1/auth/me
```

### Admin Klinik dan Superadmin

```http
GET /api/v1/admin/antrean
PUT /api/v1/admin/antrean/{id}/status
```

### Superadmin

```http
POST /api/v1/klinik
PUT /api/v1/klinik/{id}
DELETE /api/v1/klinik/{id}

POST /api/v1/gejala
PUT /api/v1/gejala/{id}
DELETE /api/v1/gejala/{id}
```

---

## 22. HTTP Status Code

### `200 OK`

Request berhasil.

### `201 Created`

Data berhasil dibuat, seperti hasil triage.

### `400 Bad Request`

Input atau query parameter tidak valid.

Contoh:

* status antrean tidak valid;
* Superadmin tidak memberikan `klinik_id`;
* payload JSON tidak valid.

### `401 Unauthorized`

JWT tidak tersedia atau tidak valid.

Contoh:

* header Authorization tidak ada;
* format Bearer salah;
* token expired;
* signature tidak valid;
* issuer salah.

### `403 Forbidden`

User valid tetapi tidak memiliki akses.

Contoh:

* role tidak sesuai;
* user tidak memiliki role aplikasi;
* Admin Klinik belum memiliki `klinik_id`;
* Admin mencoba membaca klinik lain.

### `404 Not Found`

Data tidak ditemukan atau berada di luar tenant scope.

### `500 Internal Server Error`

Terjadi error database atau server.

---

## 23. Testing

Backend memiliki unit test untuk:

* triage service;
* handler klinik;
* handler gejala;
* handler antrean;
* autentikasi JWT;
* auth middleware;
* role middleware;
* endpoint `/auth/me`;
* pembatasan baca antrean;
* pembatasan update antrean.

JWT unit test menggunakan private dan public key ES256 lokal.

Unit test tidak bergantung pada:

* internet;
* Supabase live;
* `.env`;
* endpoint JWKS production.

Production tetap menggunakan JWKS Supabase.

---

## 24. Manual Integration Test

Pengujian aktual telah mencakup:

```text
Admin Klinik login
→ berhasil

Admin Klinik membaca kliniknya
→ berhasil

Admin Klinik membaca klinik lain
→ 403

Admin Klinik update antreannya
→ berhasil

Admin Klinik update antrean klinik lain
→ 404

Request tanpa token
→ 401

Superadmin login
→ berhasil

Superadmin mengakses route dan endpoint miliknya
→ berhasil
```

---

## 25. Perintah Validasi

Dari folder backend:

```bash
go test ./...
```

Lalu:

```bash
go vet ./...
```

Backend tidak boleh digabungkan ke branch utama jika salah satu perintah gagal.

---

## 26. Keamanan

Aturan keamanan utama:

* backend tidak percaya role frontend;
* backend tidak percaya `klinik_id` frontend;
* backend memverifikasi JWT;
* backend mengambil role dari database;
* backend membatasi query berdasarkan tenant;
* database memvalidasi kombinasi role dan klinik;
* foreign key mencegah klinik terhapus ketika masih memiliki Admin.

Route protection frontend hanya meningkatkan pengalaman pengguna.

Otorisasi sebenarnya tetap dilakukan backend.

---

## 27. Fitur yang Belum Tersedia

Backend belum menyediakan:

* akun pasien;
* endpoint tiket permanen;
* endpoint posisi antrean pasien;
* estimasi waktu tunggu;
* WebSocket;
* status `Dipanggil`;
* notifikasi real-time;
* statistik agregat;
* analitik global;
* manajemen akun admin;
* deployment production;
* CI GitHub Actions.