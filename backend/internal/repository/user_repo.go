package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// GetUserRole mengambil role pengguna berdasarkan ID pengguna dari database
func (r *RepositoryWrapper) GetUserRole(ctx context.Context, userID string) (string, error) {
	var role string
	query := `SELECT role FROM user_roles WHERE user_id = $1`
	err := r.Pool.QueryRow(ctx, query, userID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Jika belum di-assign, kita bisa return error atau default role
			// Untuk sekarang, kita kembalikan error "role tidak ditemukan"
			return "", nil
		}
		return "", err
	}
	return role, nil
}
