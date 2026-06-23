package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"KlinikCepat/internal/middleware"
	"KlinikCepat/internal/models"
)

func TestGetCurrentUser(t *testing.T) {
	t.Run("berhasil mengembalikan user klinik admin", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		klinikID := "mock-klinik-id"

		claims := &middleware.JWTClaims{
			Sub:      "klinik-admin-id",
			Email:    "admin@klinik.com",
			Role:     "klinik_admin",
			KlinikID: &klinikID,
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/auth/me",
			nil,
		)

		ctx := req.Context()
		ctx = middleware.WithClaims(ctx, claims)
		req = req.WithContext(ctx)

		recorder := httptest.NewRecorder()

		h.GetCurrentUser(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"status = %d, want %d",
				recorder.Code,
				http.StatusOK,
			)
		}

		var response models.AuthMeResponse

		if err := json.NewDecoder(
			recorder.Body,
		).Decode(&response); err != nil {
			t.Fatalf(
				"gagal decode response: %v",
				err,
			)
		}

		if response.Role != "klinik_admin" {
			t.Errorf(
				"role = %s, want klinik_admin",
				response.Role,
			)
		}

		if response.KlinikID == nil {
			t.Fatal(
				"klinik_id seharusnya tidak nil",
			)
		}

		if *response.KlinikID != "mock-klinik-id" {
			t.Errorf(
				"klinik_id = %s, want mock-klinik-id",
				*response.KlinikID,
			)
		}
	})

	t.Run("ditolak jika claims tidak tersedia", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/auth/me",
			nil,
		)

		recorder := httptest.NewRecorder()

		h.GetCurrentUser(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf(
				"status = %d, want %d",
				recorder.Code,
				http.StatusUnauthorized,
			)
		}
	})
}
