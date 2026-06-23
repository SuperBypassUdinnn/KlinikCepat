package services

import (
	"KlinikCepat/internal/models"
	"KlinikCepat/internal/repository"
	"context"
	"errors"
)

// TriageService menangani logika bisnis kalkulasi triage
type TriageService struct {
	repo repository.RepositoryInterface
}

// NewTriageService membuat instance TriageService baru
func NewTriageService(repo repository.RepositoryInterface) *TriageService {
	return &TriageService{repo: repo}
}

// ProcessTriage memproses formulir gejala pasien, menghitung skor urgensi,
// menentukan status triage (Merah/Kuning/Hijau), dan memasukkan ke dalam antrean.
func (s *TriageService) ProcessTriage(ctx context.Context, req *models.TriageRequest) (*models.TriageResponse, error) {
	if req.KlinikID == "" {
		return nil, errors.New("klinik_id wajib diisi")
	}
	if req.NamaPasien == "" {
		return nil, errors.New("nama_pasien wajib diisi")
	}
	if len(req.Gejala) == 0 {
		return nil, errors.New("paling tidak satu gejala harus dipilih")
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
		TotalSkor:     totalSkor,
		StatusTriage:  statusTriage,
		StatusAntrean: "Menunggu",
	}

	err = s.repo.CreateAntrean(ctx, &antrean)
	if err != nil {
		return nil, err
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
		Pesan:        pesan,
	}, nil
}
