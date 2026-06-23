package repository

import (
	"context"
	"database/sql"

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
