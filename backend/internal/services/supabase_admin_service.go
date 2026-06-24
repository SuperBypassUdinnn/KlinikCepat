package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"KlinikCepat/internal/models"
)

// SupabaseAdminClient mendefinisikan operasi administratif
// Supabase Auth yang digunakan aplikasi.
type SupabaseAdminClient interface {
	CreateUser(
		ctx context.Context,
		email string,
		password string,
	) (*models.CreatedAuthUser, error)

	DeleteUser(
		ctx context.Context,
		userID string,
	) error
}

// SupabaseAdminService memanggil Supabase Admin Auth API.
type SupabaseAdminService struct {
	supabaseURL string
	serviceKey  string
	httpClient  *http.Client
}

// SupabaseAdminAPIError menyimpan error yang
// dikembalikan Supabase Auth API.
type SupabaseAdminAPIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *SupabaseAdminAPIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf(
			"Supabase Auth error %d (%s): %s",
			e.StatusCode,
			e.Code,
			e.Message,
		)
	}

	return fmt.Sprintf(
		"Supabase Auth error %d: %s",
		e.StatusCode,
		e.Message,
	)
}

// NewSupabaseAdminService membuat Supabase Admin client.
func NewSupabaseAdminService(
	supabaseURL string,
	serviceKey string,
	httpClient *http.Client,
) *SupabaseAdminService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	return &SupabaseAdminService{
		supabaseURL: strings.TrimRight(
			supabaseURL,
			"/",
		),
		serviceKey: serviceKey,
		httpClient: httpClient,
	}
}

// CreateUser membuat akun Supabase Auth secara langsung
// tanpa mengirim email undangan.
func (s *SupabaseAdminService) CreateUser(
	ctx context.Context,
	email string,
	password string,
) (*models.CreatedAuthUser, error) {
	endpoint := s.supabaseURL +
		"/auth/v1/admin/users"

	body, err := json.Marshal(map[string]any{
		"email":         email,
		"password":      password,
		"email_confirm": true,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat payload user: %w",
			err,
		)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat request user: %w",
			err,
		)
	}

	s.setAdminHeaders(request)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal menghubungi Supabase Auth: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {
		return nil, decodeSupabaseAdminError(response)
	}

	var responsePayload struct {
		ID    string                  `json:"id"`
		Email string                  `json:"email"`
		User  *models.CreatedAuthUser `json:"user"`
	}

	err = json.NewDecoder(
		io.LimitReader(response.Body, 1<<20),
	).Decode(&responsePayload)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membaca respons Supabase: %w",
			err,
		)
	}

	var createdUser models.CreatedAuthUser

	if responsePayload.User != nil {
		createdUser = *responsePayload.User
	} else {
		createdUser.ID = responsePayload.ID
		createdUser.Email = responsePayload.Email
	}

	if createdUser.ID == "" {
		return nil, fmt.Errorf(
			"Supabase tidak mengembalikan user ID",
		)
	}

	if createdUser.Email == "" {
		createdUser.Email = email
	}

	return &createdUser, nil
}

// DeleteUser menghapus user Supabase Auth.
// Digunakan sebagai rollback apabila insert user_roles gagal.
func (s *SupabaseAdminService) DeleteUser(
	ctx context.Context,
	userID string,
) error {
	endpoint := fmt.Sprintf(
		"%s/auth/v1/admin/users/%s",
		s.supabaseURL,
		url.PathEscape(userID),
	)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		endpoint,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat request hapus user: %w",
			err,
		)
	}

	s.setAdminHeaders(request)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"gagal menghubungi Supabase Auth: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {
		return decodeSupabaseAdminError(response)
	}

	return nil
}

func (s *SupabaseAdminService) setAdminHeaders(
	request *http.Request,
) {
	request.Header.Set(
		"Authorization",
		"Bearer "+s.serviceKey,
	)
	request.Header.Set("apikey", s.serviceKey)
	request.Header.Set(
		"Content-Type",
		"application/json",
	)
}

func decodeSupabaseAdminError(
	response *http.Response,
) error {
	data, readErr := io.ReadAll(
		io.LimitReader(response.Body, 1<<20),
	)

	if readErr != nil {
		return &SupabaseAdminAPIError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}

	var payload struct {
		Message          string `json:"message"`
		Msg              string `json:"msg"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Code             string `json:"code"`
		ErrorCode        string `json:"error_code"`
	}

	_ = json.Unmarshal(data, &payload)

	message := payload.Message

	if message == "" {
		message = payload.Msg
	}

	if message == "" {
		message = payload.ErrorDescription
	}

	if message == "" {
		message = payload.Error
	}

	if message == "" {
		message = strings.TrimSpace(string(data))
	}

	if message == "" {
		message = response.Status
	}

	code := payload.Code
	if code == "" {
		code = payload.ErrorCode
	}

	return &SupabaseAdminAPIError{
		StatusCode: response.StatusCode,
		Code:       code,
		Message:    message,
	}
}
