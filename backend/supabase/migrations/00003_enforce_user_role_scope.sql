BEGIN;

ALTER TABLE user_roles
DROP CONSTRAINT IF EXISTS user_roles_role_check;

ALTER TABLE user_roles
ADD CONSTRAINT user_roles_role_check
CHECK (
    role IN (
        'superadmin',
        'klinik_admin'
    )
);

ALTER TABLE user_roles
DROP CONSTRAINT IF EXISTS user_roles_klinik_scope_check;

ALTER TABLE user_roles
ADD CONSTRAINT user_roles_klinik_scope_check
CHECK (
    (
        role = 'superadmin'
        AND klinik_id IS NULL
    )
    OR
    (
        role = 'klinik_admin'
        AND klinik_id IS NOT NULL
    )
);

ALTER TABLE user_roles
DROP CONSTRAINT IF EXISTS user_roles_klinik_id_fkey;

ALTER TABLE user_roles
ADD CONSTRAINT user_roles_klinik_id_fkey
FOREIGN KEY (klinik_id)
REFERENCES klinik(id)
ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_user_roles_klinik_id
ON user_roles(klinik_id);

COMMIT;