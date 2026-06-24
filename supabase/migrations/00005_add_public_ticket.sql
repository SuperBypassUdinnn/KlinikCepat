-- Menambahkan informasi yang dibutuhkan
-- untuk tiket pasien tanpa akun.

ALTER TABLE antrean
ADD COLUMN email_pasien VARCHAR(255),
ADD COLUMN public_token UUID NOT NULL DEFAULT gen_random_uuid(),
ADD COLUMN kode_tiket VARCHAR(16);

-- Token publik wajib unik karena digunakan
-- sebagai identifier pada URL tiket.
ALTER TABLE antrean
ADD CONSTRAINT antrean_public_token_unique
UNIQUE (public_token);

-- Kode tiket hanya diwajibkan untuk antrean baru.
-- Data lama diperbolehkan tetap NULL.
CREATE UNIQUE INDEX idx_antrean_kode_tiket_unique
ON antrean (kode_tiket)
WHERE kode_tiket IS NOT NULL;

-- Mempercepat pencarian menggunakan
-- kombinasi kode tiket dan email.
CREATE INDEX idx_antrean_ticket_lookup
ON antrean (
    kode_tiket,
    LOWER(email_pasien)
);