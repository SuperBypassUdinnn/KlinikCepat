package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"KlinikCepat/internal/models"

	"github.com/go-chi/chi/v5"
)

func TestHandler_GejalaCRUD(t *testing.T) {
	repo := NewMockRepository()
	h := NewHandler(repo, nil)

	// Pre-seed a gejala
	gSeed := models.Gejala{
		ID:           "g-1",
		NamaGejala:   "Sakit Kepala",
		BobotUrgensi: 3,
	}
	repo.Gejalas[gSeed.ID] = &gSeed

	t.Run("Create Gejala - Success", func(t *testing.T) {
		gNew := models.Gejala{
			NamaGejala:   "Demam Tinggi",
			BobotUrgensi: 6,
		}
		body, _ := json.Marshal(gNew)
		req := httptest.NewRequest("POST", "/api/v1/gejala", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreateGejala(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusCreated)
		}

		var created models.Gejala
		_ = json.NewDecoder(rr.Body).Decode(&created)
		if created.ID != "mock-gejala-id" {
			t.Errorf("ID = %v, want 'mock-gejala-id'", created.ID)
		}
		if created.BobotUrgensi != gNew.BobotUrgensi {
			t.Errorf("BobotUrgensi = %v, want %v", created.BobotUrgensi, gNew.BobotUrgensi)
		}
	})

	t.Run("Create Gejala - Invalid Bobot", func(t *testing.T) {
		gInvalid := models.Gejala{
			NamaGejala:   "Gejala Aneh",
			BobotUrgensi: 99, // Bobot harus 1-10 (atau 0-10)
		}
		body, _ := json.Marshal(gInvalid)
		req := httptest.NewRequest("POST", "/api/v1/gejala", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreateGejala(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("Get All Gejala", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/gejala", nil)
		rr := httptest.NewRecorder()

		h.GetAllGejala(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		var list []models.Gejala
		_ = json.NewDecoder(rr.Body).Decode(&list)
		if len(list) < 1 {
			t.Error("Harus mengembalikan minimal 1 gejala")
		}
	})

	t.Run("Get Gejala By ID - Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/gejala/g-1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.GetGejalaByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		var g models.Gejala
		_ = json.NewDecoder(rr.Body).Decode(&g)
		if g.ID != "g-1" {
			t.Errorf("Gejala ID = %v, want g-1", g.ID)
		}
	})

	t.Run("Get Gejala By ID - Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/gejala/g-notfound", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g-notfound")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.GetGejalaByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("Update Gejala - Success", func(t *testing.T) {
		gUpdate := models.Gejala{
			NamaGejala:   "Sakit Kepala Akut",
			BobotUrgensi: 5,
		}
		body, _ := json.Marshal(gUpdate)
		req := httptest.NewRequest("PUT", "/api/v1/gejala/g-1", bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.UpdateGejala(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		if repo.Gejalas["g-1"].NamaGejala != gUpdate.NamaGejala {
			t.Errorf("NamaGejala setelah update = %v, want %v", repo.Gejalas["g-1"].NamaGejala, gUpdate.NamaGejala)
		}
	})

	t.Run("Delete Gejala - Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/gejala/g-1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.DeleteGejala(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		if _, exists := repo.Gejalas["g-1"]; exists {
			t.Error("Gejala 'g-1' harusnya sudah terhapus")
		}
	})
}
