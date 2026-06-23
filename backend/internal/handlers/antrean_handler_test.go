package handlers

import (
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

func TestHandler_AntreanAndTriage(t *testing.T) {
	t.Run("Process Triage - Success (Merah)", func(t *testing.T) {
		repo := NewMockRepository()
		triageService := services.NewTriageService(repo)
		h := NewHandler(repo, triageService)

		repo.Gejalas["g-1"] = &models.Gejala{ID: "g-1", NamaGejala: "Pendarahan Hebat", BobotUrgensi: 10}

		reqPayload := models.TriageRequest{
			KlinikID:   "k-1",
			NamaPasien: "Budi Santoso",
			Gejala: []models.GejalaInput{
				{GejalaID: "g-1", SkalaKeparahan: 2}, // Bobot 10 * 2 = 20 (Merah)
			},
		}

		body, _ := json.Marshal(reqPayload)
		req := httptest.NewRequest("POST", "/api/v1/triage", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.ProcessTriage(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusCreated)
		}

		var res models.TriageResponse
		_ = json.NewDecoder(rr.Body).Decode(&res)
		if res.StatusTriage != "Merah" {
			t.Errorf("StatusTriage = %v, want 'Merah'", res.StatusTriage)
		}
		if res.TotalSkor != 20 {
			t.Errorf("TotalSkor = %v, want 20", res.TotalSkor)
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

		req := httptest.NewRequest("GET", "/api/v1/admin/antrean?klinik_id=k-1&status=Menunggu", nil)
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

	t.Run("Get Antrean - Missing Klinik ID", func(t *testing.T) {
		repo := NewMockRepository()
		h := NewHandler(repo, nil)

		req := httptest.NewRequest("GET", "/api/v1/admin/antrean", nil)
		rr := httptest.NewRecorder()

		h.GetAntreanByKlinikID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("Update Status Antrean - Success", func(t *testing.T) {
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
			Status: "Selesai",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/admin/antrean/a-2/status", bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "a-2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.UpdateStatusAntrean(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		// Verifikasi status berubah di repo
		if repo.Antreans["a-2"].StatusAntrean != "Selesai" {
			t.Errorf("StatusAntrean = %v, want 'Selesai'", repo.Antreans["a-2"].StatusAntrean)
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
}
