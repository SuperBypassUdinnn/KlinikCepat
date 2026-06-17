-- Mengaktifkan Row Level Security (RLS) untuk semua tabel publik
ALTER TABLE klinik ENABLE ROW LEVEL SECURITY;
ALTER TABLE katalog_gejala ENABLE ROW LEVEL SECURITY;
ALTER TABLE antrean ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;

-- Keterangan:
-- Kita sengaja TIDAK MENAMBAHKAN POLICY apapun untuk role 'anon' atau 'authenticated'.
-- Di Supabase, jika RLS aktif tetapi tidak ada policy yang cocok, maka aksi tersebut akan DITOLAK (DENY ALL).
-- Hal ini memastikan bahwa Frontend (yang menggunakan Supabase JS dengan anon key) 
-- TIDAK BISA mengakses tabel-tabel ini secara langsung.
-- 
-- Semua akses baca/tulis data WAJIB melalui Backend Go kita.
-- Backend Go menggunakan koneksi 'postgres' role (via DATABASE_URL) yang
-- memiliki privilege superuser, sehingga otomatis mem-bypass RLS ini.
