package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"

	"KlinikCepat/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type TriageValidationError struct {
	Message string
}

func (e *TriageValidationError) Error() string {
	return e.Message
}

func newTriageValidationError(
	message string,
) error {
	return &TriageValidationError{
		Message: message,
	}
}

type TriageRepository interface {
	GetAllGejala(
		ctx context.Context,
	) ([]models.Gejala, error)

	CreateAntrean(
		ctx context.Context,
		antrean *models.Antrean,
	) error
}

// TriageService menangani logika bisnis kalkulasi triage
type TriageService struct {
	repo TriageRepository
}

// NewTriageService membuat instance TriageService baru
func NewTriageService(
	repo TriageRepository,
) *TriageService {
	return &TriageService{
		repo: repo,
	}
}

// ProcessTriage memproses formulir gejala pasien, menghitung skor urgensi,
// menentukan status triage (Merah/Kuning/Hijau), dan memasukkan ke dalam antrean.
func (s *TriageService) ProcessTriage(
	ctx context.Context,
	req *models.TriageRequest,
) (*models.TriageResponse, error) {
	if req == nil {
		return nil, newTriageValidationError(
			"payload triage wajib diisi",
		)
	}

	req.KlinikID = strings.TrimSpace(
		req.KlinikID,
	)

	req.NamaPasien = strings.TrimSpace(
		req.NamaPasien,
	)

	req.EmailPasien = strings.TrimSpace(
		req.EmailPasien,
	)

	if req.KlinikID == "" {
		return nil, newTriageValidationError(
			"klinik_id wajib diisi",
		)
	}

	if req.NamaPasien == "" {
		return nil, newTriageValidationError(
			"nama_pasien wajib diisi",
		)
	}

	normalizedEmail, err := normalizePatientEmail(
		req.EmailPasien,
	)
	if err != nil {
		return nil, newTriageValidationError(
			"format email_pasien tidak valid",
		)
	}

	req.EmailPasien = normalizedEmail

	if len(req.Gejala) == 0 {
		return nil, newTriageValidationError(
			"paling tidak satu gejala harus dipilih",
		)
	}

	// 1. Ambil seluruh katalog gejala dari database untuk mendapatkan bobot
	katalog, err := s.repo.GetAllGejala(ctx)
	if err != nil {
		return nil, err
	}

	// Buat map untuk mempermudah pencarian bobot berdasarkan ID gejala
	bobotMap := make(map[string]int)
	for _, g := range katalog {
		bobotMap[g.ID] = g.BobotUrgensi
	}

	totalSkor := 0
	hasFatalCondition := false

	// 2. Hitung total skor berdasarkan input pasien
	for _, input := range req.Gejala {
		bobot, exists := bobotMap[input.GejalaID]
		if !exists {
			// Jika ID gejala tidak valid, abaikan atau return error
			continue
		}

		// Validasi skala keparahan (wajib 1-3)
		skala := input.SkalaKeparahan
		if skala < 1 {
			skala = 1
		} else if skala > 3 {
			skala = 3
		}

		totalSkor += bobot * skala

		// Kondisi fatal jika gejala memiliki bobot ekstrem (misal Bobot = 10)
		if bobot == 10 && skala >= 1 {
			hasFatalCondition = true
		}
	}

	// 3. Tentukan klasifikasi status triage berdasarkan parameter di blueprint
	var statusTriage string
	if totalSkor >= 15 || hasFatalCondition {
		statusTriage = "Merah"
	} else if totalSkor >= 7 {
		statusTriage = "Kuning"
	} else {
		statusTriage = "Hijau"
	}

	// 4. Daftarkan antrean baru ke database
	antrean := models.Antrean{
		KlinikID:      req.KlinikID,
		NamaPasien:    req.NamaPasien,
		EmailPasien:   req.EmailPasien,
		TotalSkor:     totalSkor,
		StatusTriage:  statusTriage,
		StatusAntrean: "Menunggu",
	}

	var createErr error

	for attempt := 0; attempt < 5; attempt++ {
		kodeTiket, err := generateTicketCode()
		if err != nil {
			return nil, fmt.Errorf(
				"gagal menghasilkan kode tiket: %w",
				err,
			)
		}

		antrean.KodeTiket = kodeTiket

		createErr = s.repo.CreateAntrean(
			ctx,
			&antrean,
		)

		if createErr == nil {
			break
		}

		if !isTicketCodeConflict(createErr) {
			return nil, createErr
		}
	}

	if createErr != nil {
		return nil, fmt.Errorf(
			"gagal menghasilkan kode tiket unik: %w",
			createErr,
		)
	}

	// 5. Kembalikan respons sukses
	pesan := "Pendaftaran antrean berhasil."
	if statusTriage == "Merah" {
		pesan = "Kondisi DARURAT MEDIS (Status Merah). Silakan langsung menuju faskes utama untuk penanganan prioritas."
	}

	return &models.TriageResponse{
		AntreanID:    antrean.ID,
		StatusTriage: statusTriage,
		TotalSkor:    totalSkor,
		KodeTiket:    antrean.KodeTiket,
		PublicToken:  antrean.PublicToken,
		Pesan:        pesan,
	}, nil
}

const ticketCharacters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateTicketCode() (string, error) {
	const length = 6

	result := make([]byte, length)

	for index := range result {
		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(
				int64(len(ticketCharacters)),
			),
		)
		if err != nil {
			return "", err
		}

		result[index] =
			ticketCharacters[randomIndex.Int64()]
	}

	return "KC-" + string(result), nil
}

func normalizePatientEmail(
	value string,
) (string, error) {
	value = strings.TrimSpace(value)

	address, err := mail.ParseAddress(value)
	if err != nil {
		return "", err
	}

	if address.Address != value {
		return "", errors.New(
			"format email tidak valid",
		)
	}

	return strings.ToLower(address.Address), nil
}

func isTicketCodeConflict(err error) bool {
	var pgError *pgconn.PgError

	if !errors.As(err, &pgError) {
		return false
	}

	return pgError.Code == "23505" &&
		pgError.ConstraintName ==
			"idx_antrean_kode_tiket_unique"
}
