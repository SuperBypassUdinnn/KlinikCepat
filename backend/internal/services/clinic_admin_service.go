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

// CreateClinicAdmin membuat akun awal Admin Klinik,
// menghasilkan password sementara, lalu menyimpan role.
func (s *ClinicAdminService) CreateClinicAdmin(
	ctx context.Context,
	email string,
	klinikID string,
) (*models.CreateClinicAdminResponse, error) {
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

	temporaryPassword, err :=
		generateTemporaryPassword(14)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal menghasilkan password sementara: %w",
			err,
		)
	}

	createdUser, err := s.supabaseAdmin.CreateUser(
		ctx,
		normalizedEmail,
		temporaryPassword,
	)
	if err != nil {
		return nil, err
	}

	access := &models.UserAccess{
		UserID:   createdUser.ID,
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
			createdUser.ID,
		)

		if cleanupErr != nil {
			return nil, fmt.Errorf(
				"gagal menyimpan role Admin Klinik: %v; "+
					"rollback akun Supabase juga gagal: %w",
				err,
				cleanupErr,
			)
		}

		return nil, fmt.Errorf(
			"gagal menyimpan role Admin Klinik: %w",
			err,
		)
	}

	return &models.CreateClinicAdminResponse{
		Message:           "Akun Admin Klinik berhasil dibuat",
		UserID:            createdUser.ID,
		Email:             createdUser.Email,
		TemporaryPassword: temporaryPassword,
		Role:              "klinik_admin",
		KlinikID:          klinikID,
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

const (
	passwordLowercase = "abcdefghijkmnopqrstuvwxyz"
	passwordUppercase = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	passwordDigits    = "23456789"
	passwordSymbols   = "!@#$%&*"
)

func generateTemporaryPassword(
	length int,
) (string, error) {
	if length < 12 {
		length = 12
	}

	characterGroups := []string{
		passwordLowercase,
		passwordUppercase,
		passwordDigits,
		passwordSymbols,
	}

	password := make([]byte, 0, length)

	// Pastikan minimal ada huruf kecil, huruf besar,
	// angka, dan simbol.
	for _, group := range characterGroups {
		character, err := secureRandomCharacter(group)
		if err != nil {
			return "", err
		}

		password = append(password, character)
	}

	allCharacters :=
		passwordLowercase +
			passwordUppercase +
			passwordDigits +
			passwordSymbols

	for len(password) < length {
		character, err :=
			secureRandomCharacter(allCharacters)
		if err != nil {
			return "", err
		}

		password = append(password, character)
	}

	// Acak posisi karakter agar pola karakter wajib
	// tidak selalu muncul di awal.
	for index := len(password) - 1; index > 0; index-- {
		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(index+1)),
		)
		if err != nil {
			return "", err
		}

		target := int(randomIndex.Int64())

		password[index], password[target] =
			password[target], password[index]
	}

	return string(password), nil
}

func secureRandomCharacter(
	characters string,
) (byte, error) {
	randomIndex, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(len(characters))),
	)
	if err != nil {
		return 0, err
	}

	return characters[randomIndex.Int64()], nil
}
