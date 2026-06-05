# Arsitektur Pengguna & Alur Kerja: KlinikCepat

Dokumen ini mendefinisikan batasan akses, fungsionalitas, dan alur kerja untuk tiga kelompok aktor utama dalam sistem KlinikCepat. Jangan menambahkan peran baru di luar dokumen ini untuk mencegah pelebaran ruang lingkup (*scope creep*).

---

## 1. Pasien (Public End-User / B2C)
Aktor yang membutuhkan pelayanan medis atau non-medis di faskes tujuan.

### Fitur Utama:
* **Manajemen Akun:** Registrasi dan login standar.
* **Pencarian Klinik:** Pemilihan faskes tujuan (dengan kalkulasi jarak *client-side* berbasis formula Haversine).
* **Triage Digital:** Pengisian kuesioner gejala untuk penentuan status urgensi secara otomatis.
* **Live Tracking:** Pemantauan nomor antrean digital secara *real-time*.

### Alur Kerja:
1.  **Akses:** Mendaftar/Login menggunakan Supabase Auth (ditangani di sisi React).
2.  **Pemilihan:** Sistem meminta izin GPS peramban $\rightarrow$ merender daftar klinik terdekat $\rightarrow$ pasien memilih satu faskes tujuan.
3.  **Triage:** Pasien mengisi formulir gejala $\rightarrow$ menekan "Daftar".
4.  **Pemrosesan:** React mengirim token JWT dan payload gejala ke Backend Go. Backend Go memvalidasi token, menghitung skor total, menetapkan status warna urgensi (Merah/Kuning/Hijau), dan menyimpannya ke basis data Supabase.
5.  **Output:** Pasien dialihkan ke layar tiket antrean digital yang menampilkan estimasi waktu dan status urgensi mereka.

---

## 2. Admin Klinik (Tenant Admin / B2B)
Staf administrasi atau tenaga kesehatan (Nakes) di masing-masing faskes yang bertugas mengelola alur fisik antrean. Entitas ini diikat secara absolut pada satu `klinik_id`.

### Fitur Utama:
* **Dashboard Triage Terpusat:** Memantau antrean aktif yang sudah disortir otomatis oleh sistem (Merah di urutan teratas, Hijau di bawah).
* **Kontrol Antrean:** Aksi Panggil (*Call*), Lewati (*Skip*), dan Selesai (*Mark as Done*).
* **Statistik Harian:** Pantauan beban operasional klinik pada hari tersebut.

### Alur Kerja:
1.  **Akses:** Tidak ada registrasi mandiri. Akun dibuatkan oleh Super Admin. Admin Klinik login melalui portal khusus faskes.
2.  **Monitoring:** Setelah masuk, *dashboard* React meminta data dari Backend Go dengan filter mutlak `WHERE klinik_id = {admin_clinic_id} AND status_antrean = 'Menunggu' ORDER BY status_triage DESC, created_at ASC`.
3.  **Eksekusi:** Admin memanggil pasien pada baris teratas (terlepas dari jam kedatangan jika statusnya Merah).
4.  **Penyelesaian:** Admin menekan tombol "Selesai" $\rightarrow$ Backend Go mengubah status baris tersebut di Supabase $\rightarrow$ sisa antrean di *frontend* bergeser maju secara otomatis.

---

## 3. Super Admin (Platform Owner)
Pengembang sistem (kalian berdua) atau entitas pengawas yang memiliki otoritas penuh atas infrastruktur dan konfigurasi global platform.

### Fitur Utama:
* **Manajemen Tenant:** Registrasi faskes baru dan injeksi titik koordinat statis (`lat`, `lng`).
* **Manajemen Logika Medis:** Menyesuaikan nilai bobot pada tabel `katalog_gejala`.
* **Pemantauan Sistem:** Dasbor analitik global untuk seluruh operasi *tenant*.

### Alur Kerja:
1.  **Akses:** Akun disuntikkan secara manual melalui *database seeding* (Supabase SQL).
2.  **Operasional Eksternal:** Menambahkan data klinik baru ke dalam tabel `klinik` agar faskes tersebut bisa langsung dipilih oleh pasien di aplikasi utama.
3.  **Kalibrasi Sistem:** Mengubah bobot poin urgensi gejala jika diperlukan, yang secara instan akan mengubah standar kalkulasi algoritma di Backend Go untuk seluruh pendaftaran antrean berikutnya.
