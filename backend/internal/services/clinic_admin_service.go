package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"KlinikCepat/internal/models"
	"KlinikCepat/internal/repository"

	"github.com/jackc/pgx/v5"
)

var (
	ErrClinicAdminInvalidEmail = errors.New(
		"format email tidak valid",
	)

	ErrClinicAdminKlinikRequired = errors.New(
		"klinik ID wajib diisi",
	)

	ErrClinicAdminKlinikNotFound = errors.New(
		"klinik tidak ditemukan",
	)
)

// ClinicAdminService menangani alur pembuatan
// Admin Klinik dari Auth hingga user_roles.
type ClinicAdminService struct {
	repo          repository.ClinicAdminManagementRepository
	supabaseAdmin SupabaseAdminClient
}

func NewClinicAdminService(
	repo repository.ClinicAdminManagementRepository,
	supabaseAdmin SupabaseAdminClient,
) *ClinicAdminService {
	return &ClinicAdminService{
		repo:          repo,
		supabaseAdmin: supabaseAdmin,
	}
}

// InviteClinicAdmin mengundang user Supabase dan
// mengaitkannya dengan sebuah klinik.
func (s *ClinicAdminService) InviteClinicAdmin(
	ctx context.Context,
	email string,
	klinikID string,
) (*models.InviteClinicAdminResponse, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	klinikID = strings.TrimSpace(klinikID)
	if klinikID == "" {
		return nil, ErrClinicAdminKlinikRequired
	}

	_, err = s.repo.GetKlinikByID(
		ctx,
		klinikID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClinicAdminKlinikNotFound
		}

		return nil, fmt.Errorf(
			"gagal memeriksa klinik: %w",
			err,
		)
	}

	invitedUser, err := s.supabaseAdmin.InviteUser(
		ctx,
		normalizedEmail,
	)
	if err != nil {
		return nil, err
	}

	access := &models.UserAccess{
		UserID:   invitedUser.ID,
		Role:     "klinik_admin",
		KlinikID: &klinikID,
	}

	err = s.repo.CreateUserAccess(
		ctx,
		access,
	)
	if err != nil {
		cleanupErr := s.supabaseAdmin.DeleteUser(
			ctx,
			invitedUser.ID,
		)

		if cleanupErr != nil {
			return nil, fmt.Errorf(
				"gagal menyimpan role Admin Klinik: %v; "+
					"rollback user Supabase juga gagal: %w",
				err,
				cleanupErr,
			)
		}

		return nil, fmt.Errorf(
			"gagal menyimpan role Admin Klinik: %w",
			err,
		)
	}

	return &models.InviteClinicAdminResponse{
		Message:  "Undangan Admin Klinik berhasil dikirim",
		UserID:   invitedUser.ID,
		Email:    invitedUser.Email,
		Role:     "klinik_admin",
		KlinikID: klinikID,
	}, nil
}

func normalizeEmail(value string) (string, error) {
	value = strings.TrimSpace(value)

	address, err := mail.ParseAddress(value)
	if err != nil {
		return "", ErrClinicAdminInvalidEmail
	}

	// Tolak format seperti:
	// Nama Admin <admin@example.com>
	if address.Address != value {
		return "", ErrClinicAdminInvalidEmail
	}

	return strings.ToLower(address.Address), nil
}
