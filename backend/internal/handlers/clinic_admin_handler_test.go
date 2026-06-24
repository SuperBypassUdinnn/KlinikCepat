package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"KlinikCepat/internal/models"
	"KlinikCepat/internal/services"

	"github.com/jackc/pgx/v5"
)

type clinicAdminRepoStub struct {
	clinic        *models.Klinik
	createdAccess *models.UserAccess
	createErr     error
}

func (s *clinicAdminRepoStub) GetKlinikByID(
	ctx context.Context,
	id string,
) (*models.Klinik, error) {
	if s.clinic == nil || s.clinic.ID != id {
		return nil, pgx.ErrNoRows
	}

	return s.clinic, nil
}

func (s *clinicAdminRepoStub) CreateUserAccess(
	ctx context.Context,
	access *models.UserAccess,
) error {
	if s.createErr != nil {
		return s.createErr
	}

	s.createdAccess = access
	return nil
}

type supabaseAdminStub struct {
	invitedUser *models.CreatedAuthUser
	inviteErr   error
	deletedID   string
}

func (s *supabaseAdminStub) CreateUser(
	ctx context.Context,
	email string,
	password string,
) (*models.CreatedAuthUser, error) {
	if s.inviteErr != nil {
		return nil, s.inviteErr
	}

	return s.invitedUser, nil
}

func (s *supabaseAdminStub) DeleteUser(
	ctx context.Context,
	userID string,
) error {
	s.deletedID = userID
	return nil
}

func newClinicAdminHandlerForTest(
	repo *clinicAdminRepoStub,
	auth *supabaseAdminStub,
) *Handler {
	service := services.NewClinicAdminService(
		repo,
		auth,
	)

	return &Handler{
		ClinicAdminService: service,
	}
}

func TestCreateClinicAdminSuccess(
	t *testing.T,
) {
	repo := &clinicAdminRepoStub{
		clinic: &models.Klinik{
			ID:         "clinic-001",
			NamaKlinik: "Klinik Sehat",
		},
	}

	auth := &supabaseAdminStub{
		invitedUser: &models.CreatedAuthUser{
			ID:    "user-001",
			Email: "admin@klinik.com",
		},
	}

	handler := newClinicAdminHandlerForTest(
		repo,
		auth,
	)

	payload := models.CreateClinicAdminRequest{
		Email:    "admin@klinik.com",
		KlinikID: "clinic-001",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf(
			"gagal marshal payload: %v",
			err,
		)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/superadmin/admin-klinik/invite",
		bytes.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateClinicAdmin(
		recorder,
		request,
	)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if repo.createdAccess == nil {
		t.Fatal(
			"expected user access dibuat",
		)
	}

	if repo.createdAccess.Role != "klinik_admin" {
		t.Fatalf(
			"expected role klinik_admin, got %s",
			repo.createdAccess.Role,
		)
	}

	if repo.createdAccess.KlinikID == nil ||
		*repo.createdAccess.KlinikID !=
			"clinic-001" {
		t.Fatal(
			"expected klinik ID clinic-001",
		)
	}

	var response models.CreateClinicAdminResponse

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"gagal membaca response: %v",
			err,
		)
	}

	if response.TemporaryPassword == "" {
		t.Fatal(
			"expected temporary password dikembalikan",
		)
	}
}

func TestCreateClinicAdminInvalidEmail(
	t *testing.T,
) {
	repo := &clinicAdminRepoStub{
		clinic: &models.Klinik{
			ID: "clinic-001",
		},
	}

	auth := &supabaseAdminStub{}

	handler := newClinicAdminHandlerForTest(
		repo,
		auth,
	)

	body := bytes.NewBufferString(`
		{
			"email": "email-tidak-valid",
			"klinik_id": "clinic-001"
		}
	`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/superadmin/admin-klinik/invite",
		body,
	)

	recorder := httptest.NewRecorder()

	handler.CreateClinicAdmin(
		recorder,
		request,
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestCreateClinicAdminClinicNotFound(
	t *testing.T,
) {
	repo := &clinicAdminRepoStub{}

	auth := &supabaseAdminStub{}

	handler := newClinicAdminHandlerForTest(
		repo,
		auth,
	)

	body := bytes.NewBufferString(`
		{
			"email": "admin@klinik.com",
			"klinik_id": "clinic-tidak-ada"
		}
	`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/superadmin/admin-klinik/invite",
		body,
	)

	recorder := httptest.NewRecorder()

	handler.CreateClinicAdmin(
		recorder,
		request,
	)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
