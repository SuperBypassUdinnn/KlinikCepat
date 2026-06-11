package services

import (
	"context"
	"testing"
	"KlinikCepat/internal/models"
)

// mockTriageRepo adalah mock minimal untuk TriageService
type mockTriageRepo struct {
	gejalas  []models.Gejala
	antreans []*models.Antrean
}

func (m *mockTriageRepo) GetAllGejala(ctx context.Context) ([]models.Gejala, error) {
	return m.gejalas, nil
}

func (m *mockTriageRepo) CreateAntrean(ctx context.Context, a *models.Antrean) error {
	a.ID = "test-antrean-id"
	m.antreans = append(m.antreans, a)
	return nil
}

// Implementasikan method lain dari RepositoryInterface sebagai dummy
func (m *mockTriageRepo) CreateKlinik(ctx context.Context, k *models.Klinik) error { return nil }
func (m *mockTriageRepo) GetAllKlinik(ctx context.Context) ([]models.Klinik, error) { return nil, nil }
func (m *mockTriageRepo) GetKlinikByID(ctx context.Context, id string) (*models.Klinik, error) { return nil, nil }
func (m *mockTriageRepo) UpdateKlinik(ctx context.Context, k *models.Klinik) error { return nil }
func (m *mockTriageRepo) DeleteKlinik(ctx context.Context, id string) error { return nil }
func (m *mockTriageRepo) CreateGejala(ctx context.Context, g *models.Gejala) error { return nil }
func (m *mockTriageRepo) GetGejalaByID(ctx context.Context, id string) (*models.Gejala, error) { return nil, nil }
func (m *mockTriageRepo) UpdateGejala(ctx context.Context, g *models.Gejala) error { return nil }
func (m *mockTriageRepo) DeleteGejala(ctx context.Context, id string) error { return nil }
func (m *mockTriageRepo) GetAntreanByID(ctx context.Context, id string) (*models.Antrean, error) { return nil, nil }
func (m *mockTriageRepo) GetAntreanByKlinikID(ctx context.Context, kID string, status string) ([]models.Antrean, error) { return nil, nil }
func (m *mockTriageRepo) UpdateStatusAntrean(ctx context.Context, id string, status string) error { return nil }

func TestTriageService_ProcessTriage(t *testing.T) {
	// Mock katalog gejala
	katalogGejala := []models.Gejala{
		{ID: "g1", NamaGejala: "Pendarahan Hebat", BobotUrgensi: 10},
		{ID: "g2", NamaGejala: "Sesak Napas Ekstrem", BobotUrgensi: 10},
		{ID: "g3", NamaGejala: "Nyeri Dada Kiri", BobotUrgensi: 8},
		{ID: "g4", NamaGejala: "Demam Tinggi (> 39C)", BobotUrgensi: 5},
		{ID: "g5", NamaGejala: "Batuk Pilek Biasa", BobotUrgensi: 1},
		{ID: "g6", NamaGejala: "Layanan Non-Darurat", BobotUrgensi: 0},
	}

	repo := &mockTriageRepo{gejalas: katalogGejala}
	service := NewTriageService(repo)

	tests := []struct {
		name         string
		req          *models.TriageRequest
		expectStatus string
		expectSkor   int
		expectError  bool
	}{
		{
			name: "Status Merah - Skor Tinggi (>=15)",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Merah 1",
				Gejala: []models.GejalaInput{
					{GejalaID: "g3", SkalaKeparahan: 2}, // Nyeri Dada Kiri (8) * 2 = 16
				},
			},
			expectStatus: "Merah",
			expectSkor:   16,
			expectError:  false,
		},
		{
			name: "Status Merah - Kondisi Fatal (Bobot 10)",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Fatal",
				Gejala: []models.GejalaInput{
					{GejalaID: "g1", SkalaKeparahan: 1}, // Pendarahan Hebat (10) * 1 = 10 (Fatal!)
				},
			},
			expectStatus: "Merah",
			expectSkor:   10,
			expectError:  false,
		},
		{
			name: "Status Kuning - Skor Sedang (7-14)",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Kuning",
				Gejala: []models.GejalaInput{
					{GejalaID: "g4", SkalaKeparahan: 2}, // Demam (5) * 2 = 10
				},
			},
			expectStatus: "Kuning",
			expectSkor:   10,
			expectError:  false,
		},
		{
			name: "Status Hijau - Skor Rendah (<7)",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Hijau",
				Gejala: []models.GejalaInput{
					{GejalaID: "g5", SkalaKeparahan: 3}, // Batuk (1) * 3 = 3
				},
			},
			expectStatus: "Hijau",
			expectSkor:   3,
			expectError:  false,
		},
		{
			name: "Status Hijau - Layanan Non-Darurat (Skor 0)",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Non-Darurat",
				Gejala: []models.GejalaInput{
					{GejalaID: "g6", SkalaKeparahan: 1}, // Layanan Non-Medis (0) * 1 = 0
				},
			},
			expectStatus: "Hijau",
			expectSkor:   0,
			expectError:  false,
		},
		{
			name: "Validasi - Klinik ID Kosong",
			req: &models.TriageRequest{
				KlinikID:   "",
				NamaPasien: "Pasien Tanpa Klinik",
				Gejala: []models.GejalaInput{
					{GejalaID: "g5", SkalaKeparahan: 1},
				},
			},
			expectError: true,
		},
		{
			name: "Validasi - Nama Pasien Kosong",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "",
				Gejala: []models.GejalaInput{
					{GejalaID: "g5", SkalaKeparahan: 1},
				},
			},
			expectError: true,
		},
		{
			name: "Validasi - Tanpa Gejala",
			req: &models.TriageRequest{
				KlinikID:   "klinik-1",
				NamaPasien: "Pasien Tanpa Gejala",
				Gejala:     []models.GejalaInput{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.ProcessTriage(context.Background(), tt.req)
			if (err != nil) != tt.expectError {
				t.Fatalf("ProcessTriage() error = %v, expectError = %v", err, tt.expectError)
			}
			if tt.expectError {
				return
			}
			if res.StatusTriage != tt.expectStatus {
				t.Errorf("StatusTriage = %v, want %v", res.StatusTriage, tt.expectStatus)
			}
			if res.TotalSkor != tt.expectSkor {
				t.Errorf("TotalSkor = %v, want %v", res.TotalSkor, tt.expectSkor)
			}
		})
	}
}
