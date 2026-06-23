-- Tabel untuk memetakan user Supabase ke role tertentu di aplikasi kita
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'klinik_admin',
    klinik_id UUID REFERENCES klinik(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- Indeks untuk mempercepat pencarian role
CREATE INDEX idx_user_roles_role ON user_roles(role);

-- Komentar tabel untuk dokumentasi
COMMENT ON TABLE user_roles IS 'Menyimpan role dan asosiasi klinik untuk pengguna (admin/superadmin)';
COMMENT ON COLUMN user_roles.role IS 'Role pengguna: superadmin, klinik_admin';
COMMENT ON COLUMN user_roles.klinik_id IS 'ID Klinik yang dikelola oleh user (null untuk superadmin)';
