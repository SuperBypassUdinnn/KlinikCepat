# Arsitektur Pengguna dan Hak Akses KlinikCepat

Dokumen ini menjelaskan aktor, batas akses, dan alur kerja aktual pada sistem KlinikCepat.

**Pembaruan terakhir:** 24 Juni 2026

KlinikCepat memiliki tiga kelompok pengguna:

1. Pasien
2. Admin Klinik
3. Superadmin

Fitur yang belum diimplementasikan diberi status **Planned** dan tidak dianggap sebagai kemampuan aktif sistem.

---

## 1. Pasien

### 1.1 Deskripsi

Pasien merupakan pengguna publik yang mencari klinik, mengisi formulir triage, dan menerima tiket hasil pendaftaran antrean.

Pada versi MVP saat ini, pasien tidak memerlukan akun dan tidak melakukan login.

### 1.2 Hak Akses

Pasien dapat:

- melihat daftar klinik;
- melihat detail klinik;
- melihat katalog gejala;
- memilih klinik tujuan;
- mengisi formulir triage;
- mengirim pendaftaran antrean;
- melihat tiket hasil triage.

Pasien tidak dapat:

- mengakses dashboard admin;
- melihat seluruh antrean klinik;
- mengubah status antrean;
- mengelola klinik;
- mengelola katalog gejala.

### 1.3 Alur Kerja

1. Pasien membuka halaman utama.
2. Aplikasi meminta izin lokasi dari browser.
3. Frontend mengambil daftar klinik dari backend.
4. Frontend menghitung jarak pasien ke klinik menggunakan formula Haversine.
5. Pasien memilih klinik.
6. Frontend mengambil katalog gejala dari backend.
7. Pasien memilih gejala dan tingkat keparahan.
8. Frontend mengirim payload triage ke backend.
9. Backend menghitung skor urgensi.
10. Backend menetapkan status triage:
    - Merah
    - Kuning
    - Hijau

11. Backend menyimpan data antrean.
12. Frontend menampilkan tiket hasil triage.

### 1.4 Endpoint Pasien

Endpoint berikut bersifat publik:

```http
GET /api/v1/klinik
GET /api/v1/klinik/{id}
GET /api/v1/gejala
GET /api/v1/gejala/{id}
POST /api/v1/triage
```

Pasien tidak mengirimkan JWT saat menggunakan endpoint tersebut.

### 1.5 Keterbatasan Saat Ini

Tiket pasien saat ini masih bergantung pada state navigasi frontend.

Akibatnya:

- tiket dapat hilang setelah browser di-refresh;
- tiket belum memiliki URL permanen;
- posisi antrean belum dapat dipantau;
- estimasi waktu tunggu belum tersedia;
- belum ada riwayat antrean pasien.

Fitur tersebut termasuk dalam status **Planned**.

---

## 2. Admin Klinik

### 2.1 Deskripsi

Admin Klinik merupakan staf klinik yang bertugas memantau dan mengelola antrean pada satu klinik tertentu.

Setiap Admin Klinik wajib terhubung dengan tepat satu `klinik_id`.

```text
role      = klinik_admin
klinik_id = UUID klinik
```

Admin Klinik tidak dapat memilih atau mengganti kliniknya melalui frontend.

### 2.2 Pembuatan Akun

Admin Klinik tidak memiliki registrasi mandiri.

Pembuatan akun dilakukan melalui:

1. pembuatan user pada Supabase Auth;
2. penambahan data pada tabel `user_roles`;
3. pemberian role `klinik_admin`;
4. pengaitan akun dengan satu `klinik_id`.

Contoh struktur data:

```text
user_id   = UUID user Supabase
role      = klinik_admin
klinik_id = UUID klinik
```

Database menolak akun `klinik_admin` yang tidak memiliki `klinik_id`.

### 2.3 Hak Akses

Admin Klinik dapat:

- login menggunakan email dan password;
- melihat profil autentikasinya;
- melihat antrean kliniknya sendiri;
- memfilter antrean berdasarkan status;
- mengubah status antrean kliniknya;
- melihat statistik triage dari antrean yang sedang ditampilkan;
- logout.

Admin Klinik tidak dapat:

- melihat antrean klinik lain;
- mengubah antrean klinik lain;
- memilih `klinik_id` secara bebas;
- menambah atau menghapus klinik;
- mengubah katalog gejala;
- mengakses halaman superadmin.

### 2.4 Dashboard Antrean

Dashboard Admin Klinik menyediakan filter:

- `Menunggu`
- `Selesai`
- `Dilewati`

Untuk antrean berstatus `Menunggu`, data diurutkan berdasarkan:

1. Merah
2. Kuning
3. Hijau
4. waktu pendaftaran paling awal

Dashboard melakukan pembaruan data otomatis setiap 10 detik menggunakan polling.

Polling 10 detik tidak dikategorikan sebagai komunikasi real-time berbasis WebSocket.

### 2.5 Aksi Antrean

Admin Klinik dapat melakukan aksi:

- **Selesai**
  Mengubah `status_antrean` menjadi `Selesai`.

- **Lewati**
  Mengubah `status_antrean` menjadi `Dilewati`.

Status `Dipanggil` dan tombol `Panggil` belum diimplementasikan.

### 2.6 Isolasi Data Klinik

Frontend tidak menentukan klinik yang boleh diakses oleh Admin Klinik.

Backend menentukan scope klinik melalui alur berikut:

```text
JWT Supabase
→ user_id
→ tabel user_roles
→ role + klinik_id
→ request context
→ query antrean berdasarkan klinik_id
```

Untuk membaca antrean, Admin Klinik menggunakan:

```http
GET /api/v1/admin/antrean?status=Menunggu
```

Admin Klinik tidak perlu mengirimkan `klinik_id`.

Jika frontend mencoba mengirimkan `klinik_id` klinik lain, backend menolak request tersebut.

Untuk memperbarui antrean, backend menggunakan pembatasan yang setara dengan:

```sql
UPDATE antrean
SET status_antrean = $1
WHERE id = $2
  AND klinik_id = $3;
```

Nilai `$3` berasal dari identitas Admin Klinik yang login, bukan dari input bebas frontend.

Jika ID antrean tidak termasuk klinik Admin tersebut, backend mengembalikan respons bahwa antrean tidak ditemukan.

### 2.7 Endpoint Admin Klinik

```http
GET /api/v1/auth/me
GET /api/v1/admin/antrean?status={status}
PUT /api/v1/admin/antrean/{id}/status
```

Seluruh endpoint tersebut membutuhkan JWT Supabase yang valid.

---

## 3. Superadmin

### 3.1 Deskripsi

Superadmin merupakan pengelola master data dan konfigurasi global KlinikCepat.

Superadmin tidak terikat pada satu klinik.

```text
role      = superadmin
klinik_id = NULL
```

Database menolak akun `superadmin` yang memiliki `klinik_id`.

### 3.2 Pembuatan Akun

Akun Superadmin dibuat melalui Supabase Auth dan diberikan role secara manual pada tabel `user_roles`.

Tidak tersedia registrasi Superadmin dari frontend publik.

### 3.3 Hak Akses

Superadmin dapat:

- login menggunakan email dan password;
- melihat profil autentikasinya;
- menambah klinik;
- mengubah klinik;
- menghapus klinik yang tidak masih terhubung dengan Admin Klinik;
- menambah katalog gejala;
- mengubah katalog gejala;
- menghapus katalog gejala;
- mengakses antrean berdasarkan klinik yang dipilih;
- logout.

Superadmin tidak terikat pada satu `klinik_id`.

### 3.4 Manajemen Klinik

Superadmin dapat mengelola data:

- nama klinik;
- latitude;
- longitude;
- kapasitas aktif.

Klinik yang masih terhubung dengan Admin Klinik tidak dapat dihapus karena foreign key menggunakan aturan:

```text
ON DELETE RESTRICT
```

Admin harus dipindahkan atau dihapus keterkaitannya terlebih dahulu sebelum klinik dapat dihapus.

### 3.5 Manajemen Katalog Gejala

Superadmin dapat mengelola:

- nama gejala;
- bobot urgensi.

Perubahan bobot urgensi memengaruhi proses triage berikutnya.

Perubahan tidak menghitung ulang data antrean lama yang sudah tersimpan.

### 3.6 Endpoint Superadmin

```http
GET /api/v1/auth/me

POST /api/v1/klinik
PUT /api/v1/klinik/{id}
DELETE /api/v1/klinik/{id}

POST /api/v1/gejala
PUT /api/v1/gejala/{id}
DELETE /api/v1/gejala/{id}
```

Superadmin juga dapat menggunakan endpoint antrean dengan menentukan klinik tujuan:

```http
GET /api/v1/admin/antrean?klinik_id={id}&status={status}
PUT /api/v1/admin/antrean/{id}/status
```

### 3.7 Keterbatasan Saat Ini

Dashboard analitik global belum tersedia.

Statistik lintas klinik, grafik performa, total kunjungan harian, dan laporan agregat masih berstatus **Planned**.

---

## 4. Autentikasi

### 4.1 Penyedia Autentikasi

Autentikasi menggunakan Supabase Auth.

Frontend menerima session yang berisi access token, kemudian mengirim token melalui header:

```http
Authorization: Bearer <access_token>
```

Frontend membaca token langsung dari session Supabase dan tidak menyimpan salinan token aplikasi secara manual.

### 4.2 Verifikasi Token

Backend memverifikasi JWT Supabase menggunakan:

- algoritma ES256;
- public signing key dari endpoint JWKS Supabase;
- issuer;
- audience `authenticated`;
- expiration;
- subject pengguna.

Setelah token valid, backend mengambil role aplikasi dan `klinik_id` dari tabel `user_roles`.

Role yang terdapat pada JWT Supabase seperti `authenticated` bukan role aplikasi KlinikCepat.

Role aplikasi tetap berasal dari database:

```text
superadmin
klinik_admin
```

---

## 5. Otorisasi

### 5.1 Backend sebagai Sumber Kebenaran

Frontend menyembunyikan menu dan halaman berdasarkan role untuk meningkatkan pengalaman pengguna.

Namun, backend tetap menjadi sumber kebenaran utama dalam otorisasi.

Pengguna tidak memperoleh akses hanya karena:

- mengubah URL;
- mengubah state React;
- mengedit request;
- mengirim `klinik_id` berbeda;
- membuka halaman melalui DevTools.

Backend selalu memeriksa:

1. validitas JWT;
2. keberadaan user pada `user_roles`;
3. role aplikasi;
4. `klinik_id` untuk Admin Klinik;
5. kepemilikan data antrean.

### 5.2 Role-Based Routing Frontend

Route Admin Klinik hanya dapat diakses oleh:

```text
klinik_admin
```

Route Superadmin hanya dapat diakses oleh:

```text
superadmin
```

Pengguna dengan role yang salah diarahkan ke dashboard sesuai rolenya.

---

## 6. Matriks Hak Akses

| Fitur                            | Pasien |     Admin Klinik |       Superadmin |
| -------------------------------- | -----: | ---------------: | ---------------: |
| Melihat daftar klinik            |     Ya |               Ya |               Ya |
| Mengisi triage                   |     Ya | Tidak diperlukan | Tidak diperlukan |
| Melihat tiket hasil triage       |     Ya |            Tidak |            Tidak |
| Login admin                      |  Tidak |               Ya |               Ya |
| Melihat antrean klinik sendiri   |  Tidak |               Ya |               Ya |
| Melihat antrean klinik lain      |  Tidak |            Tidak |               Ya |
| Mengubah status antrean          |  Tidak |   Klinik sendiri |               Ya |
| Mengelola klinik                 |  Tidak |            Tidak |               Ya |
| Mengelola katalog gejala         |  Tidak |            Tidak |               Ya |
| Mengakses dashboard admin klinik |  Tidak |               Ya |            Tidak |
| Mengakses halaman superadmin     |  Tidak |            Tidak |               Ya |
| Analitik global                  |  Tidak |            Tidak |          Planned |

---

## 7. Fitur Planned

Fitur berikut belum aktif:

- akun pasien;
- registrasi dan login pasien;
- riwayat kunjungan pasien;
- tiket permanen;
- live tracking posisi antrean;
- estimasi waktu tunggu;
- status `Dipanggil`;
- tombol `Panggil`;
- notifikasi pasien;
- statistik harian dari endpoint agregasi;
- dashboard analitik global;
- manajemen akun admin melalui frontend;
- pengaturan role melalui frontend.

Fitur Planned tidak boleh dipresentasikan sebagai fitur yang sudah tersedia.
