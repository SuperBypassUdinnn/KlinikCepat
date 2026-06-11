package repository

import (
	"context"
	"KlinikCepat/internal/models"
)

// CreateAntrean menambahkan antrean baru ke database
func (r *RepositoryWrapper) CreateAntrean(ctx context.Context, a *models.Antrean) error {
	query := `
		INSERT INTO antrean (klinik_id, nama_pasien, total_skor, status_triage, status_antrean)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.Pool.QueryRow(
		ctx,
		query,
		a.KlinikID,
		a.NamaPasien,
		a.TotalSkor,
		a.StatusTriage,
		a.StatusAntrean,
	).Scan(&a.ID, &a.CreatedAt)
	return err
}

// GetAntreanByID mengambil detail antrean berdasarkan ID
func (r *RepositoryWrapper) GetAntreanByID(ctx context.Context, id string) (*models.Antrean, error) {
	query := `
		SELECT id, klinik_id, nama_pasien, total_skor, status_triage, status_antrean, created_at
		FROM antrean
		WHERE id = $1
	`
	var a models.Antrean
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.KlinikID,
		&a.NamaPasien,
		&a.TotalSkor,
		&a.StatusTriage,
		&a.StatusAntrean,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAntreanByKlinikID mengambil antrean berdasarkan klinik_id dan status_antrean,
// diurutkan berdasarkan prioritas triage: Merah -> Kuning -> Hijau.
func (r *RepositoryWrapper) GetAntreanByKlinikID(ctx context.Context, klinikID string, status string) ([]models.Antrean, error) {
	// Menggunakan CASE untuk memastikan prioritas: Merah (1), Kuning (2), Hijau (3)
	query := `
		SELECT id, klinik_id, nama_pasien, total_skor, status_triage, status_antrean, created_at
		FROM antrean
		WHERE klinik_id = $1 AND status_antrean = $2
		ORDER BY 
			CASE status_triage
				WHEN 'Merah' THEN 1
				WHEN 'Kuning' THEN 2
				WHEN 'Hijau' THEN 3
				ELSE 4
			END ASC,
			created_at ASC
	`
	rows, err := r.Pool.Query(ctx, query, klinikID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var antreans []models.Antrean
	for rows.Next() {
		var a models.Antrean
		err := rows.Scan(
			&a.ID,
			&a.KlinikID,
			&a.NamaPasien,
			&a.TotalSkor,
			&a.StatusTriage,
			&a.StatusAntrean,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		antreans = append(antreans, a)
	}
	return antreans, nil
}

// UpdateStatusAntrean memperbarui status antrean (Menunggu / Selesai / Dilewati)
func (r *RepositoryWrapper) UpdateStatusAntrean(ctx context.Context, id string, status string) error {
	query := `
		UPDATE antrean
		SET status_antrean = $1
		WHERE id = $2
	`
	_, err := r.Pool.Exec(ctx, query, status, id)
	return err
}
