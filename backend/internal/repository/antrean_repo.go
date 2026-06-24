package repository

import (
	"KlinikCepat/internal/models"
	"context"
)

// CreateAntrean menambahkan antrean baru ke database
func (r *RepositoryWrapper) CreateAntrean(
	ctx context.Context,
	a *models.Antrean,
) error {
	query := `
		INSERT INTO antrean (
			klinik_id,
			nama_pasien,
			email_pasien,
			kode_tiket,
			total_skor,
			status_triage,
			status_antrean
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			public_token::text,
			created_at
	`

	return r.Pool.QueryRow(
		ctx,
		query,
		a.KlinikID,
		a.NamaPasien,
		a.EmailPasien,
		a.KodeTiket,
		a.TotalSkor,
		a.StatusTriage,
		a.StatusAntrean,
	).Scan(
		&a.ID,
		&a.PublicToken,
		&a.CreatedAt,
	)
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
func (r *RepositoryWrapper) UpdateStatusAntrean(
	ctx context.Context,
	id string,
	status string,
	klinikID *string,
) (bool, error) {
	if klinikID == nil {
		commandTag, err := r.Pool.Exec(
			ctx,
			`
				UPDATE antrean
				SET status_antrean = $1
				WHERE id = $2
			`,
			status,
			id,
		)
		if err != nil {
			return false, err
		}

		return commandTag.RowsAffected() > 0, nil
	}

	commandTag, err := r.Pool.Exec(
		ctx,
		`
			UPDATE antrean
			SET status_antrean = $1
			WHERE id = $2
			  AND klinik_id = $3
		`,
		status,
		id,
		*klinikID,
	)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (r *RepositoryWrapper) GetPublicTicketByToken(
	ctx context.Context,
	publicToken string,
) (*models.PublicTicket, error) {
	query := `
		SELECT
			a.public_token::text,
			COALESCE(a.kode_tiket, ''),
			a.nama_pasien,
			k.nama_klinik,
			a.total_skor,
			a.status_triage,
			a.status_antrean,
			a.created_at
		FROM antrean AS a
		JOIN klinik AS k
			ON k.id = a.klinik_id
		WHERE a.public_token::text = $1
	`

	var ticket models.PublicTicket

	err := r.Pool.QueryRow(
		ctx,
		query,
		publicToken,
	).Scan(
		&ticket.PublicToken,
		&ticket.KodeTiket,
		&ticket.NamaPasien,
		&ticket.NamaKlinik,
		&ticket.TotalSkor,
		&ticket.StatusTriage,
		&ticket.StatusAntrean,
		&ticket.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *RepositoryWrapper) CheckPublicTicket(
	ctx context.Context,
	kodeTiket string,
	email string,
) (*models.PublicTicket, error) {
	query := `
		SELECT
			a.public_token::text,
			COALESCE(a.kode_tiket, ''),
			a.nama_pasien,
			k.nama_klinik,
			a.total_skor,
			a.status_triage,
			a.status_antrean,
			a.created_at
		FROM antrean AS a
		JOIN klinik AS k
			ON k.id = a.klinik_id
		WHERE UPPER(a.kode_tiket) = UPPER($1)
		  AND LOWER(a.email_pasien) = LOWER($2)
		LIMIT 1
	`

	var ticket models.PublicTicket

	err := r.Pool.QueryRow(
		ctx,
		query,
		kodeTiket,
		email,
	).Scan(
		&ticket.PublicToken,
		&ticket.KodeTiket,
		&ticket.NamaPasien,
		&ticket.NamaKlinik,
		&ticket.TotalSkor,
		&ticket.StatusTriage,
		&ticket.StatusAntrean,
		&ticket.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
