# Struktur Proyek KlinikCepat

Dokumen ini menjelaskan struktur direktori KlinikCepat berdasarkan implementasi aktual.

**Pembaruan terakhir:** 24 Juni 2026

Struktur di bawah menampilkan file dan direktori utama. File pengujian dan stylesheet tertentu dapat bertambah seiring pengembangan.

---

## 1. Struktur Utama

```text
KlinikCepat/
├── backend/
├── frontend/
├── supabase/
├── docs/
├── .gitignore
└── README.md
```

Keterangan:

- `backend/`
  REST API Go, business logic, autentikasi, dan akses database.

- `frontend/`
  Aplikasi React untuk pasien, Admin Klinik, dan Superadmin.

- `supabase/`
  Migration database PostgreSQL.

- `docs/`
  Dokumentasi arsitektur, integrasi, role, blueprint, dan progres proyek.

---

## 2. Backend

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

File `.env` bersifat lokal dan tidak boleh dimasukkan ke Git.

---

## 3. `backend/cmd/api`

```text
backend/cmd/api/
└── main.go
```

`main.go` merupakan entry point backend.

Tanggung jawabnya meliputi:

- membaca environment variable;
- membuat koneksi PostgreSQL;
- membuat repository;
- membuat service;
- membuat handler;
- mendaftarkan route;
- memasang middleware;
- menjalankan HTTP server.

Route utama dikelompokkan menjadi:

- route publik;
- route user terautentikasi;
- route Admin Klinik;
- route Superadmin.

---

## 4. `backend/internal/handlers`

Contoh struktur:

```text
backend/internal/handlers/
├── auth_handler.go
├── antrean_handler.go
├── gejala_handler.go
├── klinik_handler.go
├── handler.go
├── auth_handler_test.go
├── antrean_handler_test.go
└── mock_test.go
```

Nama file dapat berbeda sedikit sesuai perkembangan repository.

### Tanggung Jawab

Handler bertugas:

- membaca parameter URL;
- membaca query parameter;
- membaca body JSON;
- melakukan validasi HTTP-level;
- mengambil claims dari request context;
- memanggil service atau repository;
- menentukan status HTTP;
- mengembalikan JSON response.

### Handler Utama

#### `auth_handler.go`

Menangani:

```http
GET /api/v1/auth/me
```

Mengembalikan:

- user ID;
- email;
- role aplikasi;
- `klinik_id`.

#### `antrean_handler.go`

Menangani:

```http
POST /api/v1/triage
GET /api/v1/admin/antrean
PUT /api/v1/admin/antrean/{id}/status
```

Handler ini juga menerapkan tenant scope berdasarkan role dan klinik user.

#### `klinik_handler.go`

Menangani pembacaan dan CRUD klinik.

#### `gejala_handler.go`

Menangani pembacaan dan CRUD katalog gejala.

---

## 5. `backend/internal/middleware`

```text
backend/internal/middleware/
├── auth.go
├── role.go
├── auth_test.go
└── role_test.go
```

### `auth.go`

Berisi:

- pembacaan header `Authorization`;
- parsing Bearer token;
- verifikasi JWT Supabase;
- verifikasi ES256 melalui JWKS;
- validasi issuer;
- validasi audience;
- validasi expiration;
- penyimpanan claims ke request context.

### `role.go`

Berisi middleware `RequireRole`.

Tanggung jawabnya:

- mengambil user ID dari claims;
- membaca role dan `klinik_id` dari database;
- memeriksa role yang diperbolehkan;
- memastikan `klinik_admin` memiliki `klinik_id`;
- menyimpan role dan `klinik_id` ke request context.

---

## 6. `backend/internal/models`

```text
backend/internal/models/
└── models.go
```

Model utama:

```text
Klinik
Gejala
Antrean
UserAccess
AuthMeResponse
GejalaInput
TriageRequest
TriageResponse
```

### `UserAccess`

Merepresentasikan hasil pembacaan tabel `user_roles`.

```go
type UserAccess struct {
	UserID   string
	Role     string
	KlinikID *string
}
```

### `AuthMeResponse`

Merepresentasikan respons:

```http
GET /api/v1/auth/me
```

---

## 7. `backend/internal/repository`

Contoh struktur:

```text
backend/internal/repository/
├── interfaces.go
├── repository.go
├── klinik_repo.go
├── gejala_repo.go
├── antrean_repo.go
└── user_repo.go
```

Nama file wrapper repository dapat berbeda sesuai implementasi aktual.

### `interfaces.go`

Mendefinisikan kontrak repository, seperti:

- `KlinikRepository`
- `GejalaRepository`
- `AntreanRepository`
- `UserRepository`
- `RepositoryInterface`

### `klinik_repo.go`

Berisi query:

- membuat klinik;
- mengambil semua klinik;
- mengambil klinik berdasarkan ID;
- memperbarui klinik;
- menghapus klinik.

### `gejala_repo.go`

Berisi query CRUD katalog gejala.

### `antrean_repo.go`

Berisi query:

- membuat antrean;
- mengambil antrean;
- mengambil antrean berdasarkan klinik;
- mengurutkan triage;
- memperbarui status antrean;
- membatasi update berdasarkan `klinik_id`.

Untuk Admin Klinik, update dilakukan dengan pembatasan:

```sql
WHERE id = $2
  AND klinik_id = $3
```

### `user_repo.go`

Berisi query tabel `user_roles`.

Method utama:

```go
GetUserAccess(
	ctx context.Context,
	userID string,
) (*models.UserAccess, error)
```

---

## 8. `backend/internal/services`

```text
backend/internal/services/
├── triage_service.go
└── triage_service_test.go
```

### `triage_service.go`

Berisi business logic triage:

- validasi input;
- pengambilan bobot gejala;
- perhitungan total skor;
- penentuan status triage;
- pembuatan antrean;
- pembuatan response.

Business logic ditempatkan di service agar tidak bercampur dengan HTTP handler.

---

## 9. Frontend

```text
frontend/
├── public/
├── src/
│   ├── components/
│   ├── context/
│   ├── pages/
│   ├── services/
│   ├── App.jsx
│   ├── main.jsx
│   └── index.css
├── .env
├── .env.example
├── index.html
├── package.json
├── package-lock.json
└── vite.config.js
```

File `.env` tidak boleh dimasukkan ke Git.

---

## 10. `frontend/src/components`

Contoh struktur:

```text
frontend/src/components/
├── Badge.jsx
├── Button.jsx
├── Card.jsx
├── LoadingSpinner.jsx
├── Navbar.jsx
├── Navbar.css
└── ProtectedRoute.jsx
```

### `Navbar.jsx`

Menampilkan menu berdasarkan role.

Menu publik:

```text
Cari Klinik
Login Admin
```

Menu Admin Klinik:

```text
Dashboard
Halaman Pasien
Logout
```

Menu Superadmin:

```text
Kelola Klinik
Kelola Gejala
Halaman Pasien
Logout
```

### `ProtectedRoute.jsx`

Melindungi route berdasarkan:

- session;
- loading state;
- role;
- allowed roles.

Contoh:

```jsx
<ProtectedRoute allowedRoles={["superadmin"]}>
  <ManajemenKlinik />
</ProtectedRoute>
```

---

## 11. `frontend/src/context`

```text
frontend/src/context/
└── AuthContext.jsx
```

`AuthContext` mengelola:

```text
user
profile
role
clinicId
loading
authError
signIn
signOut
```

Alurnya:

```text
Supabase session
→ GET /api/v1/auth/me
→ profile
→ role
→ clinicId
```

Token tidak disalin secara manual ke `localStorage`.

---

## 12. `frontend/src/pages`

```text
frontend/src/pages/
├── Patient/
│   ├── CariKlinik.jsx
│   ├── TriageForm.jsx
│   └── TicketAntrean.jsx
├── AdminKlinik/
│   ├── LoginAdmin.jsx
│   ├── LoginAdmin.css
│   ├── DashboardAdmin.jsx
│   └── DashboardAdmin.css
└── SuperAdmin/
    ├── ManajemenKlinik.jsx
    └── ManajemenGejala.jsx
```

File CSS dapat berada di folder yang sama dengan masing-masing halaman.

---

## 13. Halaman Pasien

### `CariKlinik.jsx`

Tanggung jawab:

- meminta lokasi pengguna;
- mengambil daftar klinik;
- menghitung jarak;
- menampilkan klinik;
- mengarahkan pasien ke form triage.

### `TriageForm.jsx`

Tanggung jawab:

- mengambil katalog gejala;
- menerima skala keparahan;
- membuat payload triage;
- mengirim triage ke backend.

### `TicketAntrean.jsx`

Menampilkan hasil triage.

Keterbatasan saat ini:

- masih bergantung pada state navigasi;
- belum memiliki URL permanen;
- data dapat hilang setelah refresh.

---

## 14. Halaman Admin Klinik

### `LoginAdmin.jsx`

Digunakan oleh:

- Admin Klinik;
- Superadmin.

Setelah login, redirect berdasarkan role:

```text
klinik_admin → /admin/dashboard
superadmin   → /superadmin/klinik
```

### `DashboardAdmin.jsx`

Tanggung jawab:

- mengambil antrean klinik user;
- memfilter status;
- menampilkan statistik;
- auto-refresh setiap 10 detik;
- menyelesaikan antrean;
- melewati antrean.

Dashboard tidak mengambil seluruh klinik dan tidak menampilkan clinic selector.

---

## 15. Halaman Superadmin

### `ManajemenKlinik.jsx`

Menyediakan CRUD klinik.

### `ManajemenGejala.jsx`

Menyediakan CRUD katalog gejala.

Kedua halaman hanya dapat diakses oleh role:

```text
superadmin
```

---

## 16. `frontend/src/services`

```text
frontend/src/services/
├── api.js
└── supabaseClient.js
```

### `supabaseClient.js`

Membuat Supabase client menggunakan:

```text
VITE_SUPABASE_URL
VITE_SUPABASE_ANON_KEY
```

### `api.js`

Memusatkan komunikasi backend.

Tanggung jawab:

- menentukan base URL API;
- membaca session Supabase;
- memasang Bearer token;
- menangani response error;
- menyediakan fungsi endpoint.

Contoh fungsi:

```text
getCurrentUser
getClinics
getClinicById
getGejala
getGejalaById
submitTriage
getQueue
updateStatusAntrean
createKlinik
updateKlinik
deleteKlinik
createGejala
updateGejala
deleteGejala
```

---

## 17. `frontend/src/App.jsx`

Mendaftarkan route aplikasi.

### Route Publik

```text
/
/triage/:klinikId
/ticket
/admin/login
```

### Route Admin Klinik

```text
/admin/dashboard
```

Hanya untuk:

```text
klinik_admin
```

### Route Superadmin

```text
/superadmin/klinik
/superadmin/gejala
```

Hanya untuk:

```text
superadmin
```

---

## 18. Supabase

```text
supabase/
└── migrations/
    ├── 00001_initial_schema.sql
    ├── 00002_user_roles.sql
    └── 00003_enforce_user_role_scope.sql
```

Nama file migration harus disesuaikan dengan file aktual repository apabila terdapat perbedaan penamaan.

---

## 19. Migration Awal

### `00001_initial_schema.sql`

Membuat struktur dasar:

- klinik;
- katalog gejala;
- antrean;
- enum atau constraint status;
- relasi utama;
- data awal jika disertakan dalam migration.

Tidak terdapat file terpisah `seed.sql` pada struktur saat ini.

---

## 20. Migration User Roles

### `00002_user_roles.sql`

Membuat tabel:

```text
user_roles
```

Atribut utama:

```text
user_id
role
klinik_id
```

---

## 21. Migration Tenant Constraint

### `00003_enforce_user_role_scope.sql`

Menerapkan aturan:

```text
superadmin
→ klinik_id harus NULL

klinik_admin
→ klinik_id wajib terisi
```

Migration ini juga mengatur:

```text
user_roles.klinik_id
→ klinik.id
→ ON DELETE RESTRICT
```

dan menambahkan index `klinik_id`.

---

## 22. Dokumentasi

```text
docs/
├── ProjectStructure.md
├── backend_architecture.md
├── blueprint_klinikcepat.md
├── frontend_integration.md
├── progress.md
└── user_roles_klinikcepat.md
```

### `ProjectStructure.md`

Menjelaskan struktur repository.

### `backend_architecture.md`

Menjelaskan backend, JWT, repository, middleware, dan tenant authorization.

### `blueprint_klinikcepat.md`

Menjelaskan blueprint sistem dan status fitur.

### `frontend_integration.md`

Menjelaskan integrasi React, Supabase Auth, API layer, dan role-based routing.

### `progress.md`

Mencatat status implementasi aktual dan roadmap.

### `user_roles_klinikcepat.md`

Menjelaskan aktor dan hak akses.

---

## 23. File yang Tidak Digunakan

Struktur aktual tidak mengandalkan:

```text
backend/internal/config/
supabase/seed.sql
```

Konfigurasi backend dibaca melalui environment variable.

Data awal dapat dimasukkan dalam migration apabila diperlukan.

Dokumentasi tidak boleh mencantumkan file tersebut sebagai file aktif kecuali benar-benar dibuat dan digunakan.

---

## 24. Environment Files

File lokal:

```text
backend/.env
frontend/.env
```

Tidak boleh dimasukkan ke repository.

File contoh:

```text
frontend/.env.example
```

boleh disimpan karena tidak berisi kredensial asli.

Backend dapat memiliki `.env.example` apabila ingin mendokumentasikan variable yang diperlukan.

---

## 25. Perintah Menjalankan Backend

```bash
cd backend
go run ./cmd/api
```

Validasi:

```bash
go test ./...
go vet ./...
```

---

## 26. Perintah Menjalankan Frontend

```bash
cd frontend
npm install
npm run dev
```

Build:

```bash
npm run build
```

---

## 27. Prinsip Struktur Proyek

Struktur proyek mengikuti pemisahan tanggung jawab:

```text
Handler
→ menangani HTTP

Middleware
→ autentikasi dan otorisasi

Service
→ business logic

Repository
→ akses database

Model
→ struktur data

Context
→ state autentikasi frontend

Service API
→ komunikasi frontend-backend

Page
→ tampilan dan interaksi pengguna
```

Pemisahan ini membuat kode:

- lebih mudah diuji;
- lebih mudah dipelihara;
- lebih mudah dikembangkan;
- lebih aman dari pencampuran logika;
- lebih mudah dipahami anggota tim.
