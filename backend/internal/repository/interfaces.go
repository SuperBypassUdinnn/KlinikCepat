package repository

import (
	"KlinikCepat/internal/models"
	"context"
)

// KlinikRepository mendefinisikan kontrak akses data untuk Klinik
type KlinikRepository interface {
	CreateKlinik(ctx context.Context, k *models.Klinik) error
	GetAllKlinik(ctx context.Context) ([]models.Klinik, error)
	GetKlinikByID(ctx context.Context, id string) (*models.Klinik, error)
	UpdateKlinik(ctx context.Context, k *models.Klinik) error
	DeleteKlinik(ctx context.Context, id string) error
}

// GejalaRepository mendefinisikan kontrak akses data untuk Katalog Gejala
type GejalaRepository interface {
	CreateGejala(ctx context.Context, g *models.Gejala) error
	GetAllGejala(ctx context.Context) ([]models.Gejala, error)
	GetGejalaByID(ctx context.Context, id string) (*models.Gejala, error)
	UpdateGejala(ctx context.Context, g *models.Gejala) error
	DeleteGejala(ctx context.Context, id string) error
}

// AntreanRepository mendefinisikan kontrak akses data untuk Antrean
type AntreanRepository interface {
	CreateAntrean(ctx context.Context, a *models.Antrean) error
	GetAntreanByID(ctx context.Context, id string) (*models.Antrean, error)
	GetAntreanByKlinikID(ctx context.Context, klinikID string, status string) ([]models.Antrean, error)
	UpdateStatusAntrean(
		ctx context.Context,
		id string,
		status string,
		klinikID *string,
	) (bool, error)
}

// UserRepository mendefinisikan kontrak akses data untuk User Roles
type UserRepository interface {
	GetUserAccess(ctx context.Context, userID string) (*models.UserAccess, error)
}

// ClinicAdminManagementRepository mendefinisikan
// kebutuhan database untuk pembuatan Admin Klinik.
type ClinicAdminManagementRepository interface {
	GetKlinikByID(
		ctx context.Context,
		id string,
	) (*models.Klinik, error)

	CreateUserAccess(
		ctx context.Context,
		access *models.UserAccess,
	) error
}

// RepositoryInterface menggabungkan seluruh fungsionalitas repository database
type RepositoryInterface interface {
	KlinikRepository
	GejalaRepository
	AntreanRepository
	UserRepository
}
