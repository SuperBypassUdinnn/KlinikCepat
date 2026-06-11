package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// RepositoryWrapper menampung koneksi pool database dan menyediakan akses ke sub-repository
type RepositoryWrapper struct {
	Pool *pgxpool.Pool
}

// NewRepositoryWrapper membuat instansi baru dari RepositoryWrapper
func NewRepositoryWrapper(pool *pgxpool.Pool) *RepositoryWrapper {
	return &RepositoryWrapper{Pool: pool}
}
