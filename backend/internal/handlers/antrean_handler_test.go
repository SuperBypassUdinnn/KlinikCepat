package handlers

import (
	"KlinikCepat/internal/middleware"
	"KlinikCepat/internal/models"
	"KlinikCepat/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func withTestClaims(
	req *http.Request,
	role string,
	klinikID *string,
) *http.Request {
	claims := &middleware.JWTClaims{
		Sub:      "test-user-id",
		Email:    "test@example.com",
		Role:     role,
		KlinikID: klinikID,
	}

	ctx := middleware.WithClaims(
		req.Context(),
		claims,
	)

	return req.WithContext(ctx)
}

func TestHandler_AntreanAndTriage(t *testing.T) {
	t.Run("Process Triage - Success (Merah)", func(t *testing.T) {
		repo := NewMockRepository()

		repo.Gejalas["g-1"] = &models.Gejala{
			ID:           "g-1",
			NamaGejala:   "Pendarahan Hebat",
			BobotUrgensi: 10,
		}

		triageService := services.NewTriageService(repo)
		h := NewHandler(repo, triageService)

		body := []byte(`{
		"klinik_id": "k-1",
		"nama_pasien": "Budi Santoso",
		"email_pasien": "pasien@example.com",
		"gejala": [
			{
				"gejala_id": "g-1",
				"skala_keparahan": 2
			}
		]
	}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/triage",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

		rr := httptest.NewRecorder()

		h.ProcessTriage(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf(
				"StatusCode = %d, want %d; response body = %s",
				rr.Code,
				http.StatusCreated,
				rr.Body.String(),
			)
		}

		var res models.TriageResponse

		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatalf(
				"gagal membaca response: %v",
				err,
			)
		}

		if res.StatusTriage != "Merah" {
			t.Errorf(
				"StatusTriage = %q, want %q",
				res.StatusTriage,
				"Merah",
			)
		}

		if res.TotalSkor != 20 {
			t.Errorf(
				"TotalSkor = %d, want 20",
				res.TotalSkor,
			)
		}

		if res.KodeTiket == "" {
			t.Error("KodeTiket tidak boleh kosong")
		}

		if res.PublicToken == "" {
			t.Error("PublicToken tidak boleh kosong")
		}
	})

	t.Run("Get Antrean - Success", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		// Seed antrean
		repo.Antreans["a-1"] = &models.Antrean{
			ID:            "a-1",
			KlinikID:      "k-1",
			NamaPasien:    "Pasien Antre",
			TotalSkor:     5,
			StatusTriage:  "Hijau",
			StatusAntrean: "Menunggu",
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/admin/antrean?status=Menunggu",
			nil,
		)

		klinikID := "k-1"

		req = withTestClaims(
			req,
			"klinik_admin",
			&klinikID,
		)
		rr := httptest.NewRecorder()

		h.GetAntreanByKlinikID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		var list []models.Antrean
		_ = json.NewDecoder(rr.Body).Decode(&list)
		if len(list) != 1 {
			t.Errorf("Jumlah antrean = %v, want 1", len(list))
		}
	})

	t.Run(
		"Superadmin wajib menentukan klinik ID",
		func(t *testing.T) {
			repo := NewMockRepository()
			h := NewHandler(repo, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/admin/antrean",
				nil,
			)

			req = withTestClaims(
				req,
				"superadmin",
				nil,
			)

			rr := httptest.NewRecorder()

			h.GetAntreanByKlinikID(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf(
					"StatusCode = %v, want %v",
					rr.Code,
					http.StatusBadRequest,
				)
			}
		},
	)

	t.Run(
		"Klinik admin ditolak mengakses klinik lain",
		func(t *testing.T) {
			repo := NewMockRepository()
			h := NewHandler(repo, nil)

			klinikID := "k-1"

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/admin/antrean?klinik_id=k-2",
				nil,
			)

			req = withTestClaims(
				req,
				"klinik_admin",
				&klinikID,
			)

			rr := httptest.NewRecorder()

			h.GetAntreanByKlinikID(rr, req)

			if rr.Code != http.StatusForbidden {
				t.Errorf(
					"StatusCode = %v, want %v",
					rr.Code,
					http.StatusForbidden,
				)
			}
		},
	)

	t.Run(
		"Get antrean menolak status tidak valid",
		func(t *testing.T) {
			repo := NewMockRepository()
			h := NewHandler(repo, nil)

			klinikID := "k-1"

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/admin/antrean?status=StatusNgawur",
				nil,
			)

			req = withTestClaims(
				req,
				"klinik_admin",
				&klinikID,
			)

			rr := httptest.NewRecorder()

			h.GetAntreanByKlinikID(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf(
					"StatusCode = %v, want %v",
					rr.Code,
					http.StatusBadRequest,
				)
			}
		},
	)

	t.Run("Update Status Antrean - Success", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		repo.Antreans["a-1"] = &models.Antrean{
			ID:            "a-1",
			KlinikID:      "k-1",
			NamaPasien:    "Budi",
			TotalSkor:     10,
			StatusTriage:  "Kuning",
			StatusAntrean: "Menunggu",
		}

		payload := map[string]string{
			"status": "Selesai",
		}

		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(
			http.MethodPut,
			"/api/v1/admin/antrean/a-1/status",
			bytes.NewBuffer(body),
		)

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "a-1")

		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				routeContext,
			),
		)

		klinikID := "k-1"

		req = withTestClaims(
			req,
			"klinik_admin",
			&klinikID,
		)

		rr := httptest.NewRecorder()

		h.UpdateStatusAntrean(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf(
				"StatusCode = %v, want %v",
				rr.Code,
				http.StatusOK,
			)
		}

		if repo.Antreans["a-1"].StatusAntrean != "Selesai" {
			t.Errorf(
				"StatusAntrean = %v, want 'Selesai'",
				repo.Antreans["a-1"].StatusAntrean,
			)
		}
	})

	t.Run("Update Status Antrean - Invalid Status", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		repo.Antreans["a-2"] = &models.Antrean{
			ID:            "a-2",
			KlinikID:      "k-1",
			NamaPasien:    "Pasien B",
			StatusAntrean: "Menunggu",
		}

		payload := struct {
			Status string `json:"status"`
		}{
			Status: "StatusNgawur",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/admin/antrean/a-2/status", bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "a-2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.UpdateStatusAntrean(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run(
		"Klinik admin tidak bisa update antrean klinik lain",
		func(t *testing.T) {
			repo := NewMockRepository()
			h := NewHandler(repo, nil)

			repo.Antreans["a-klinik-lain"] = &models.Antrean{
				ID:            "a-klinik-lain",
				KlinikID:      "k-2",
				NamaPasien:    "Pasien Klinik B",
				StatusAntrean: "Menunggu",
			}

			payload := map[string]string{
				"status": "Selesai",
			}

			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(
				http.MethodPut,
				"/api/v1/admin/antrean/a-klinik-lain/status",
				bytes.NewBuffer(body),
			)

			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add(
				"id",
				"a-klinik-lain",
			)

			req = req.WithContext(
				context.WithValue(
					req.Context(),
					chi.RouteCtxKey,
					routeContext,
				),
			)

			// Admin terikat ke Klinik A.
			klinikID := "k-1"

			req = withTestClaims(
				req,
				"klinik_admin",
				&klinikID,
			)

			rr := httptest.NewRecorder()

			h.UpdateStatusAntrean(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Errorf(
					"StatusCode = %v, want %v",
					rr.Code,
					http.StatusNotFound,
				)
			}

			if repo.Antreans["a-klinik-lain"].
				StatusAntrean != "Menunggu" {
				t.Error(
					"Status antrean klinik lain tidak boleh berubah",
				)
			}
		},
	)
}
