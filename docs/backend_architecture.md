# Arsitektur Backend & Dokumentasi API: KlinikCepat

Backend KlinikCepat dibangun menggunakan bahasa **Golang (Go)** dengan router HTTP **go-chi** dan driver PostgreSQL **pgxpool**. Backend ini bertindak sebagai perantara logis antara Frontend (React) dan database (Supabase PostgreSQL), serta mengimplementasikan mesin kalkulasi triage dan verifikasi keamanan JWT.

---

## 1. Alur Otentikasi (JWT Supabase Auth)
Semua rute yang berada di bawah pengawasan admin diproteksi menggunakan **Supabase JWT Verification Middleware** lokal.

- **Mekanisme**:
  1. Frontend (React) melakukan login menggunakan Supabase Auth API dan mendapatkan JWT Token (HS256).
  2. Frontend melakukan permintaan API ke Go Backend dengan menyertakan token pada header HTTP:  
     `Authorization: Bearer <token_jwt>`
  3. Go Backend mengekstrak token, memverifikasi tanda tangannya menggunakan kunci `SUPABASE_JWT_SECRET` dari file `.env`, memeriksa waktu kedaluwarsa (`exp`), dan mengizinkan/menolak akses.

---

## 2. Daftar Endpoint API (v1)

### A. Rute Publik (Pasien / B2C)

#### 1. Get All Klinik
- **Endpoint**: `GET /api/v1/klinik`
- **Deskripsi**: Mengambil semua data faskes/klinik yang terdaftar. Digunakan oleh klien untuk menghitung jarak terdekat menggunakan rumus Haversine di sisi frontend.
- **Response (200 OK)**:
  ```json
  [
    {
      "id": "77e5d8a0-2fbe-4972-bb2d-f75a7c29fb91",
      "nama_klinik": "Klinik Sehat Selalu",
      "lat": -6.2,
      "lng": 106.816666,
      "kapasitas_aktif": 50,
      "created_at": "2026-06-11T10:11:50Z"
    }
  ]
  ```

#### 2. Get Klinik By ID
- **Endpoint**: `GET /api/v1/klinik/{id}`
- **Deskripsi**: Mengambil detail satu klinik tertentu berdasarkan UUID-nya.
- **Response (200 OK)**:
  ```json
  {
    "id": "77e5d8a0-2fbe-4972-bb2d-f75a7c29fb91",
    "nama_klinik": "Klinik Sehat Selalu",
    "lat": -6.2,
    "lng": 106.816666,
    "kapasitas_aktif": 50,
    "created_at": "2026-06-11T10:11:50Z"
  }
  ```

#### 3. Get All Gejala
- **Endpoint**: `GET /api/v1/gejala`
- **Deskripsi**: Mengambil katalog referensi gejala beserta bobot urgensinya. Digunakan untuk merender formulir kuesioner gejala pasien.
- **Response (200 OK)**:
  ```json
  [
    {
      "id": "a5e9bb2a-3c0c-4e8c-8cfa-298a28cc12ab",
      "nama_gejala": "Pendarahan Hebat",
      "bobot_urgensi": 10,
      "created_at": "2026-06-11T10:11:50Z"
    }
  ]
  ```

#### 4. Get Gejala By ID
- **Endpoint**: `GET /api/v1/gejala/{id}`
- **Deskripsi**: Mengambil detail satu gejala berdasarkan UUID-nya.

#### 5. Pendaftaran Triage (Submit Questionnaire)
- **Endpoint**: `POST /api/v1/triage`
- **Deskripsi**: Menerima pilihan gejala dan tingkat keparahan (skala 1-3) dari pasien, menghitung skor total triage, mengklasifikasikan prioritas antrean, membuat tiket antrean baru, dan mengembalikannya ke klien.
- **Request Body**:
  ```json
  {
    "klinik_id": "77e5d8a0-2fbe-4972-bb2d-f75a7c29fb91",
    "nama_pasien": "Budi Santoso",
    "gejala": [
      {
        "gejala_id": "a5e9bb2a-3c0c-4e8c-8cfa-298a28cc12ab",
        "skala_keparahan": 2
      }
    ]
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "antrean_id": "cb2a197c-bd3d-4c31-9f93-018274a2cb5a",
    "status_triage": "Merah",
    "total_skor": 20,
    "pesan": "Kondisi DARURAT MEDIS (Status Merah). Silakan langsung menuju faskes utama untuk penanganan prioritas."
  }
  ```

---

### B. Rute Terproteksi (Admin Klinik / Super Admin)
*Memerlukan Header: `Authorization: Bearer <JWT-token>`*

#### 1. CRUD Klinik (Kelola Tenant)
- **POST `/api/v1/klinik`**: Membuat klinik baru (Super Admin).
- **PUT `/api/v1/klinik/{id}`**: Memperbarui informasi koordinat atau kapasitas klinik.
- **DELETE `/api/v1/klinik/{id}`**: Menghapus klinik penyewa.

#### 2. CRUD Katalog Gejala (Kelola Logika Triage Medis)
- **POST `/api/v1/gejala`**: Menambahkan gejala baru ke kamus sistem.
- **PUT `/api/v1/gejala/{id}`**: Menyesuaikan deskripsi atau menaikkan/menurunkan bobot urgensi.
- **DELETE `/api/v1/gejala/{id}`**: Menghapus gejala dari katalog.

#### 3. Dapatkan Antrean Aktif Faskes
- **Endpoint**: `GET /api/v1/admin/antrean?klinik_id={uuid}&status={Menunggu/Selesai/Dilewati}`
- **Deskripsi**: Mengambil data antrean untuk dashboard admin faskes. Diurutkan secara berlapis: Pasien **Merah** di urutan teratas, lalu **Kuning**, lalu **Hijau**, serta diurutkan berdasarkan waktu pendaftaran terkecil (tertua).
- **Response (200 OK)**:
  ```json
  [
    {
      "id": "cb2a197c-bd3d-4c31-9f93-018274a2cb5a",
      "klinik_id": "77e5d8a0-2fbe-4972-bb2d-f75a7c29fb91",
      "nama_pasien": "Budi Santoso",
      "total_skor": 20,
      "status_triage": "Merah",
      "status_antrean": "Menunggu",
      "created_at": "2026-06-11T10:12:00Z"
    }
  ]
  ```

#### 4. Perbarui Status Antrean Pasien
- **Endpoint**: `PUT /api/v1/admin/antrean/{id}/status`
- **Deskripsi**: Mengubah status antrean (dipanggil/selesai/dilewati) oleh staf klinik nakes.
- **Request Body**:
  ```json
  {
    "status": "Selesai"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Status antrean berhasil diperbarui menjadi Selesai"
  }
  ```

---

## 3. Unit Testing & Mocking (Pengujian Mandiri)
Untuk memfasilitasi pengujian unit yang terisolasi dan mandiri (tanpa koneksi database fisik ataupun koneksi internet), arsitektur database dibungkus menggunakan interface Go.

### A. Abstraksi Repositori
Seluruh panggilan database didefinisikan melalui interface di `internal/repository/interfaces.go`:
- `KlinikRepository` (CRUD Klinik)
- `GejalaRepository` (CRUD Gejala)
- `AntreanRepository` (CRUD & prioritas Antrean)

Kedua modul utama (`Handler` dan `TriageService`) bergantung pada interface `RepositoryInterface`, bukan pada *concrete struct* pool database langsung.

### B. Mock Repository (*In-Memory*)
Di dalam berkas pengujian [mock_test.go](file:///home/superbypassudin/.clone/Github/KlinikCepat/backend/internal/handlers/mock_test.go), diimplementasikan struct `MockRepository` yang menggunakan struktur data `map` Go untuk menampung data sementara di memori selama unit test berjalan. Hal ini memungkinkan simulasi operasi database secara instan dan andal.

### C. Cara Menjalankan Tes
Seluruh pengujian unit (Middleware JWT, Triage Engine, CRUD Handlers) dapat dijalankan dengan perintah Go standar dari direktori `backend/`:
```bash
go test -v ./...
```
Untuk menguji linter kode:
```bash
golangci-lint run
```
