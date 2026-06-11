package repository

import (
	"context"
	"KlinikCepat/internal/models"
)

// CreateKlinik menambahkan klinik baru ke database
func (r *RepositoryWrapper) CreateKlinik(ctx context.Context, k *models.Klinik) error {
	query := `
		INSERT INTO klinik (nama_klinik, lat, lng, kapasitas_aktif)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.Pool.QueryRow(ctx, query, k.NamaKlinik, k.Lat, k.Lng, k.KapasitasAktif).Scan(&k.ID, &k.CreatedAt)
	return err
}

// GetAllKlinik mengambil semua data klinik dari database
func (r *RepositoryWrapper) GetAllKlinik(ctx context.Context) ([]models.Klinik, error) {
	query := `SELECT id, nama_klinik, lat, lng, kapasitas_aktif, created_at FROM klinik`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kliniks []models.Klinik
	for rows.Next() {
		var k models.Klinik
		err := rows.Scan(&k.ID, &k.NamaKlinik, &k.Lat, &k.Lng, &k.KapasitasAktif, &k.CreatedAt)
		if err != nil {
			return nil, err
		}
		kliniks = append(kliniks, k)
	}
	return kliniks, nil
}

// GetKlinikByID mengambil satu klinik berdasarkan ID
func (r *RepositoryWrapper) GetKlinikByID(ctx context.Context, id string) (*models.Klinik, error) {
	query := `SELECT id, nama_klinik, lat, lng, kapasitas_aktif, created_at FROM klinik WHERE id = $1`
	var k models.Klinik
	err := r.Pool.QueryRow(ctx, query, id).Scan(&k.ID, &k.NamaKlinik, &k.Lat, &k.Lng, &k.KapasitasAktif, &k.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// UpdateKlinik memperbarui data klinik di database
func (r *RepositoryWrapper) UpdateKlinik(ctx context.Context, k *models.Klinik) error {
	query := `
		UPDATE klinik
		SET nama_klinik = $1, lat = $2, lng = $3, kapasitas_aktif = $4
		WHERE id = $5
	`
	_, err := r.Pool.Exec(ctx, query, k.NamaKlinik, k.Lat, k.Lng, k.KapasitasAktif, k.ID)
	return err
}

// DeleteKlinik menghapus klinik dari database berdasarkan ID
func (r *RepositoryWrapper) DeleteKlinik(ctx context.Context, id string) error {
	query := `DELETE FROM klinik WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}
