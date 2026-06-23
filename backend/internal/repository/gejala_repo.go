package repository

import (
	"KlinikCepat/internal/models"
	"context"
)

// CreateGejala menambahkan gejala baru ke database
func (r *RepositoryWrapper) CreateGejala(ctx context.Context, g *models.Gejala) error {
	query := `
		INSERT INTO katalog_gejala (nama_gejala, bobot_urgensi)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.Pool.QueryRow(ctx, query, g.NamaGejala, g.BobotUrgensi).Scan(&g.ID, &g.CreatedAt)
	return err
}

// GetAllGejala mengambil seluruh katalog gejala
func (r *RepositoryWrapper) GetAllGejala(ctx context.Context) ([]models.Gejala, error) {
	query := `SELECT id, nama_gejala, bobot_urgensi, created_at FROM katalog_gejala`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gejalas []models.Gejala
	for rows.Next() {
		var g models.Gejala
		err := rows.Scan(&g.ID, &g.NamaGejala, &g.BobotUrgensi, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		gejalas = append(gejalas, g)
	}
	return gejalas, nil
}

// GetGejalaByID mengambil satu gejala berdasarkan ID
func (r *RepositoryWrapper) GetGejalaByID(ctx context.Context, id string) (*models.Gejala, error) {
	query := `SELECT id, nama_gejala, bobot_urgensi, created_at FROM katalog_gejala WHERE id = $1`
	var g models.Gejala
	err := r.Pool.QueryRow(ctx, query, id).Scan(&g.ID, &g.NamaGejala, &g.BobotUrgensi, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// UpdateGejala memperbarui gejala di database
func (r *RepositoryWrapper) UpdateGejala(ctx context.Context, g *models.Gejala) error {
	query := `
		UPDATE katalog_gejala
		SET nama_gejala = $1, bobot_urgensi = $2
		WHERE id = $3
	`
	_, err := r.Pool.Exec(ctx, query, g.NamaGejala, g.BobotUrgensi, g.ID)
	return err
}

// DeleteGejala menghapus gejala berdasarkan ID
func (r *RepositoryWrapper) DeleteGejala(ctx context.Context, id string) error {
	query := `DELETE FROM katalog_gejala WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}
