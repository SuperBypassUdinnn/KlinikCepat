-- Tabel untuk menyimpan entitas penyewa (Klinik)
CREATE TABLE klinik (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_klinik VARCHAR(255) NOT NULL,
    lat FLOAT NOT NULL,
    lng FLOAT NOT NULL,
    kapasitas_aktif INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- Tabel untuk katalog referensi bobot urgensi gejala
CREATE TABLE katalog_gejala (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_gejala VARCHAR(255) NOT NULL,
    bobot_urgensi INT NOT NULL CHECK (bobot_urgensi >= 0 AND bobot_urgensi <= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- Enum untuk membatasi status triage dan status antrean
CREATE TYPE status_triage_enum AS ENUM ('Merah', 'Kuning', 'Hijau');
CREATE TYPE status_antrean_enum AS ENUM ('Menunggu', 'Selesai', 'Dilewati');

-- Tabel jantung operasional aplikasi: Antrean Pasien
CREATE TABLE antrean (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    klinik_id UUID NOT NULL REFERENCES klinik(id) ON DELETE CASCADE,
    nama_pasien VARCHAR(255) NOT NULL,
    total_skor INT NOT NULL,
    status_triage status_triage_enum NOT NULL,
    status_antrean status_antrean_enum NOT NULL DEFAULT 'Menunggu',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- Membuat indeks untuk mempercepat kueri Admin Klinik
CREATE INDEX idx_antrean_admin_view ON antrean (klinik_id, status_antrean, status_triage DESC, created_at ASC);

-- --------------------------------------------------------
-- Data Seeding (Data Dummy Awal)
-- --------------------------------------------------------

-- 1. Insert beberapa klinik awal
INSERT INTO klinik (nama_klinik, lat, lng, kapasitas_aktif) VALUES
('Klinik Sehat Selalu', -6.200000, 106.816666, 50),
('Klinik Bakti Medika', -6.210000, 106.826666, 30);

-- 2. Insert referensi gejala
INSERT INTO katalog_gejala (nama_gejala, bobot_urgensi) VALUES
('Pendarahan Hebat', 10),
('Sesak Napas Ekstrem', 10),
('Nyeri Dada Kiri', 8),
('Demam Tinggi (> 39C)', 5),
('Mual Muntah Terus Menerus', 4),
('Sakit Kepala Ringan', 2),
('Batuk Pilek Biasa', 1),
('Layanan Non-Darurat (Cek Darah/Surat Sehat)', 0);
