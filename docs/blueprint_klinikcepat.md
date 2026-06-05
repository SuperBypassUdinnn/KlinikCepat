# Cetak Biru Proyek: KlinikCepat
## B2B Micro-SaaS Triage & Manajemen Antrean Faskes

**Tema:** SDGs 3 (Good Health and Well-being)  
**Target Proyek:** UAS Praktikum Pemrograman Berbasis Web  
**Struktur Tim:** 2 Orang (Backend & Frontend)

---

### 1. Deskripsi Proyek
KlinikCepat adalah platform manajemen antrean cerdas berbasis *multi-tenant* yang dirancang untuk Fasilitas Kesehatan (Faskes) Tingkat Pertama independen atau poliklinik institusi. Berbeda dengan sistem antrean konvensional yang mengandalkan model *First-In-First-Out* (FIFO) murni, aplikasi ini bertindak sebagai gerbang *triage* digital otomatis. 

Sistem mengkalkulasi skor urgensi pasien berdasarkan input kuesioner gejala secara *real-time* sebelum mereka tiba di lokasi, lalu memprioritaskan penanganan medis darurat (Status Merah) di atas keluhan ringan (Status Hijau). Pendekatan ini secara langsung mereduksi penumpukan pasien berisiko tinggi di ruang tunggu fisik dan mengoptimalkan manajemen waktu tenaga medis.

---

### 2. Arsitektur Basis Data (Shared Schema Multi-Tenant)
Isolasi data dilakukan secara absolut menggunakan *Foreign Key* `klinik_id` dalam satu basis data relasional tunggal untuk menghemat waktu konfigurasi infrastruktur.

| Nama Tabel | Kolom Esensial & Relasi | Fungsi Utama |
| :--- | :--- | :--- |
| **`klinik`** | `id` (PK), `nama_klinik`, `lat` (Float), `lng` (Float), `kapasitas_aktif` (Int) | Menyimpan entitas penyewa (*tenant*). Titik koordinat wajib statis. |
| **`katalog_gejala`** | `id` (PK), `nama_gejala` (Varchar), `bobot_urgensi` (Int, 1-10) | Kamus referensi statis untuk kalkulasi sistem. Menghindari *hardcode* pada backend. |
| **`antrean`** | `id` (PK), `klinik_id` (FK), `nama_pasien` (Varchar), `total_skor` (Int), `status_triage` (Enum: Merah, Kuning, Hijau), `status_antrean` (Enum: Menunggu, Selesai), `created_at` (Timestamp) | Jantung operasional aplikasi. Semua kueri admin wajib difilter dengan klausa `WHERE klinik_id = ?`. |

---

### 3. Mesin Logika Triage (Backend Engine)
*Backend* menerima *payload* kuesioner gejala dari pasien dan menghitung skor urgensi total menggunakan rumusan matematis berikut:

$$S_{urgensi} = \sum_{i=1}^{n} (W_i \cdot V_i)$$

*Keterangan:* * $W_i$ = Bobot dasar gejala dari tabel `katalog_gejala`.  
* $V_i$ = Nilai keparahan/skala (1-3) yang diinput oleh pengguna di *frontend*.

**Parameter Klasifikasi Urgensi:**
1. **Prioritas 1 (MERAH):** $S_{urgensi} \ge 15$ ATAU terdapat kondisi fatal (misal: Pendarahan Hebat = Ya / Sesak Napas = Ekstrem). Secara otomatis memotong antrean ke urutan paling atas.
2. **Prioritas 2 (KUNING):** $7 \le S_{urgensi} < 15$. Mengantre tepat di bawah kelompok pasien berstatus Merah.
3. **Prioritas 3 (HIJAU):** $S_{urgensi} < 7$ ATAU Layanan Non-Darurat (seperti Cek Darah / Surat Sehat). Mengantre di urutan paling bawah berdasarkan kronologi waktu pendaftaran.

> **Aturan Mutlak Backend:** Kueri pemanggilan pasien pada *dashboard* admin klinik wajib dieksekusi dengan perintah pengurutan berlapis: `ORDER BY status_triage DESC, created_at ASC`.

---

### 4. Offloading Kalkulasi Spasial (Frontend Logic)
Komputasi penentuan jarak tidak boleh membebani peladen (*server*). *Backend* hanya menyuplai *array* koordinat statis dari seluruh klinik yang aktif, lalu *frontend* (JavaScript) meminta akses lokasi pengguna via HTML5 Geolocation dan menghitung jarak radius terdekat menggunakan *Haversine formula*:

$$d = 2r \cdot rcsin\left(\sqrt{\sin^2\left(rac{\Delta\phi}{2}ight) + \cos\phi_1\cos\phi_2\sin^2\left(rac{\Delta\lambda}{2}ight)}ight)$$

Setelah nilai jarak ($d$) didapatkan, *frontend* mengurutkan *array* objek tersebut dari nilai terkecil ke terbesar sebelum merendernya ke antarmuka pengguna. Sediakan fitur *Search Bar* teks biasa sebagai mekanisme *fallback* mutlak jika izin GPS ditolak oleh pengguna.

---

### 5. Delegasi Tugas & Batasan Kerja Eksekusi

#### Porsi Backend (Infrastruktur & API)
* Perancangan skema *database* relasional.
* Implementasi autentikasi akses berbasis token (JWT) untuk memisahkan hak akses pasien umum dan manajemen klinik.
* Penyusunan *endpoint* RESTful yang aman.
* Validasi matematis pada algoritma *sorting* antrean berlapis.

#### Porsi Frontend (UI/UX & Klien)
* Implementasi alur kuesioner gejala pasien yang interaktif.
* Eksekusi perhitungan formula Haversine di sisi klien.
* Pembuatan *dashboard* admin klinik yang memvisualisasikan perubahan indikator warna antrean (Merah, Kuning, Hijau) secara instan.

---

### 6. Daftar Hitam Fitur (Scope Creep Warning)
Mengerjakan fitur di bawah ini sebelum fungsionalitas inti (Poin 1-5) selesai 100% akan langsung menggagalkan penyelesaian proyek UAS kalian:
* **Rekam Medis Elektronik (EMR):** Jangan membuat sistem riwayat penyakit jangka panjang; setelah status menjadi 'Selesai', data cukup diarsipkan.
* **Gateway Pembayaran (Payment Gateway):** Proses transaksi diasumsikan berjalan secara fisik di kasir faskes. Jangan buang waktu mengurus integrasi API Midtrans/Xendit.
* **Chat Real-time:** Penggunaan *WebSockets* untuk interaksi teks antara dokter dan pasien akan menghancurkan fokus kalian dalam menyempurnakan algoritma antrean.
