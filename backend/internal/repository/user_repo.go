package repository

import (
	"context"
	"database/sql"
	"errors"

	"KlinikCepat/internal/models"

	"github.com/jackc/pgx/v5"
)

// GetUserAccess mengambil role dan klinik
// yang diasosiasikan dengan user.
func (r *RepositoryWrapper) GetUserAccess(
	ctx context.Context,
	userID string,
) (*models.UserAccess, error) {
	var access models.UserAccess
	var klinikID sql.NullString

	query := `
		SELECT
			user_id::text,
			role,
			klinik_id::text
		FROM user_roles
		WHERE user_id = $1
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&access.UserID,
		&access.Role,
		&klinikID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if klinikID.Valid {
		value := klinikID.String
		access.KlinikID = &value
	}

	return &access, nil
}

// CreateUserAccess menyimpan role aplikasi dan klinik
// untuk user Supabase yang baru dibuat.
func (r *RepositoryWrapper) CreateUserAccess(
	ctx context.Context,
	access *models.UserAccess,
) error {
	if access == nil {
		return errors.New("user access tidak boleh nil")
	}

	if access.UserID == "" {
		return errors.New("user ID wajib diisi")
	}

	if access.Role == "" {
		return errors.New("role wajib diisi")
	}

	if access.KlinikID == nil || *access.KlinikID == "" {
		return errors.New(
			"klinik ID wajib diisi untuk Admin Klinik",
		)
	}

	query := `
		INSERT INTO user_roles (
			user_id,
			role,
			klinik_id
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.Pool.Exec(
		ctx,
		query,
		access.UserID,
		access.Role,
		*access.KlinikID,
	)

	return err
}
