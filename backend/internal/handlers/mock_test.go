package handlers

import (
	"KlinikCepat/internal/models"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// MockRepository adalah implementasi in-memory dari repository.RepositoryInterface
type MockRepository struct {
	Kliniks  map[string]*models.Klinik
	Gejalas  map[string]*models.Gejala
	Antreans map[string]*models.Antrean
}

// NewMockRepository membuat instance MockRepository baru
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Kliniks:  make(map[string]*models.Klinik),
		Gejalas:  make(map[string]*models.Gejala),
		Antreans: make(map[string]*models.Antrean),
	}
}

// CreateKlinik menambahkan klinik baru ke MockRepository
func (m *MockRepository) CreateKlinik(ctx context.Context, k *models.Klinik) error {
	k.ID = "mock-klinik-id"
	k.CreatedAt = time.Now()
	m.Kliniks[k.ID] = k
	return nil
}

// GetAllKlinik mengambil semua klinik
func (m *MockRepository) GetAllKlinik(ctx context.Context) ([]models.Klinik, error) {
	var list []models.Klinik
	for _, k := range m.Kliniks {
		list = append(list, *k)
	}
	return list, nil
}

// GetKlinikByID mengambil klinik berdasarkan ID
func (m *MockRepository) GetKlinikByID(ctx context.Context, id string) (*models.Klinik, error) {
	k, exists := m.Kliniks[id]
	if !exists {
		return nil, errors.New("klinik tidak ditemukan")
	}
	return k, nil
}

// UpdateKlinik memperbarui data klinik
func (m *MockRepository) UpdateKlinik(ctx context.Context, k *models.Klinik) error {
	if _, exists := m.Kliniks[k.ID]; !exists {
		return errors.New("klinik tidak ditemukan")
	}
	m.Kliniks[k.ID] = k
	return nil
}

// DeleteKlinik menghapus klinik berdasarkan ID
func (m *MockRepository) DeleteKlinik(ctx context.Context, id string) error {
	if _, exists := m.Kliniks[id]; !exists {
		return errors.New("klinik tidak ditemukan")
	}
	delete(m.Kliniks, id)
	return nil
}

// CreateGejala menambahkan gejala baru ke MockRepository
func (m *MockRepository) CreateGejala(ctx context.Context, g *models.Gejala) error {
	g.ID = "mock-gejala-id"
	g.CreatedAt = time.Now()
	m.Gejalas[g.ID] = g
	return nil
}

// GetAllGejala mengambil semua gejala
func (m *MockRepository) GetAllGejala(ctx context.Context) ([]models.Gejala, error) {
	var list []models.Gejala
	for _, g := range m.Gejalas {
		list = append(list, *g)
	}
	return list, nil
}

// GetGejalaByID mengambil gejala berdasarkan ID
func (m *MockRepository) GetGejalaByID(ctx context.Context, id string) (*models.Gejala, error) {
	g, exists := m.Gejalas[id]
	if !exists {
		return nil, errors.New("gejala tidak ditemukan")
	}
	return g, nil
}

// UpdateGejala memperbarui data gejala
func (m *MockRepository) UpdateGejala(ctx context.Context, g *models.Gejala) error {
	if _, exists := m.Gejalas[g.ID]; !exists {
		return errors.New("gejala tidak ditemukan")
	}
	m.Gejalas[g.ID] = g
	return nil
}

// DeleteGejala menghapus gejala berdasarkan ID
func (m *MockRepository) DeleteGejala(ctx context.Context, id string) error {
	if _, exists := m.Gejalas[id]; !exists {
		return errors.New("gejala tidak ditemukan")
	}
	delete(m.Gejalas, id)
	return nil
}

// CreateAntrean menambahkan antrean baru ke MockRepository
func (m *MockRepository) CreateAntrean(
	ctx context.Context,
	antrean *models.Antrean,
) error {
	antrean.ID = "mock-antrean-id"

	if antrean.PublicToken == "" {
		antrean.PublicToken =
			"11111111-1111-1111-1111-111111111111"
	}

	antrean.CreatedAt = time.Now()

	m.Antreans[antrean.ID] = antrean

	return nil
}

func (m *MockRepository) GetPublicTicketByToken(
	ctx context.Context,
	publicToken string,
) (*models.PublicTicket, error) {
	for _, antrean := range m.Antreans {
		if antrean.PublicToken != publicToken {
			continue
		}

		namaKlinik := ""

		if klinik, exists :=
			m.Kliniks[antrean.KlinikID]; exists {
			namaKlinik = klinik.NamaKlinik
		}

		return &models.PublicTicket{
			PublicToken:   antrean.PublicToken,
			KodeTiket:     antrean.KodeTiket,
			NamaPasien:    antrean.NamaPasien,
			NamaKlinik:    namaKlinik,
			TotalSkor:     antrean.TotalSkor,
			StatusTriage:  antrean.StatusTriage,
			StatusAntrean: antrean.StatusAntrean,
			CreatedAt:     antrean.CreatedAt,
		}, nil
	}

	return nil, pgx.ErrNoRows
}

func (m *MockRepository) CheckPublicTicket(
	ctx context.Context,
	kodeTiket string,
	email string,
) (*models.PublicTicket, error) {
	for _, antrean := range m.Antreans {
		kodeSesuai := strings.EqualFold(
			antrean.KodeTiket,
			kodeTiket,
		)

		emailSesuai := strings.EqualFold(
			antrean.EmailPasien,
			email,
		)

		if !kodeSesuai || !emailSesuai {
			continue
		}

		return m.GetPublicTicketByToken(
			ctx,
			antrean.PublicToken,
		)
	}

	return nil, pgx.ErrNoRows
}

// GetAntreanByID mengambil antrean berdasarkan ID
func (m *MockRepository) GetAntreanByID(ctx context.Context, id string) (*models.Antrean, error) {
	a, exists := m.Antreans[id]
	if !exists {
		return nil, errors.New("antrean tidak ditemukan")
	}
	return a, nil
}

// GetAntreanByKlinikID mengambil antrean berdasarkan klinikID dan status
func (m *MockRepository) GetAntreanByKlinikID(ctx context.Context, klinikID string, status string) ([]models.Antrean, error) {
	var list []models.Antrean
	for _, a := range m.Antreans {
		if a.KlinikID == klinikID && a.StatusAntrean == status {
			list = append(list, *a)
		}
	}
	return list, nil
}

// UpdateStatusAntrean memperbarui status antrean
func (m *MockRepository) UpdateStatusAntrean(
	ctx context.Context,
	id string,
	status string,
	klinikID *string,
) (bool, error) {
	antrean, exists := m.Antreans[id]
	if !exists {
		return false, nil
	}

	if klinikID != nil &&
		antrean.KlinikID != *klinikID {
		return false, nil
	}

	antrean.StatusAntrean = status

	return true, nil
}

// GetUserAccess mensimulasikan pengambilan role
// dan klinik user dari database.
func (m *MockRepository) GetUserAccess(
	ctx context.Context,
	userID string,
) (*models.UserAccess, error) {
	switch userID {
	case "superadmin-id":
		return &models.UserAccess{
			UserID: "superadmin-id",
			Role:   "superadmin",
		}, nil

	case "klinik-admin-id":
		klinikID := "mock-klinik-id"

		return &models.UserAccess{
			UserID:   "klinik-admin-id",
			Role:     "klinik_admin",
			KlinikID: &klinikID,
		}, nil

	default:
		return nil, nil
	}
}
