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

func TestHandler_KlinikCRUD(t *testing.T) {
	repo := NewMockRepository()
	h := NewHandler(repo, nil)

	// Pre-seed a klinik
	kSeed := models.Klinik{
		ID:             "k-1",
		NamaKlinik:     "Klinik Awal",
		Lat:            1.23,
		Lng:            4.56,
		KapasitasAktif: 10,
	}
	repo.Kliniks[kSeed.ID] = &kSeed

	t.Run("Create Klinik - Success", func(t *testing.T) {
		kNew := models.Klinik{
			NamaKlinik:     "Klinik Baru",
			Lat:            -6.1,
			Lng:            106.8,
			KapasitasAktif: 20,
		}
		body, _ := json.Marshal(kNew)
		req := httptest.NewRequest("POST", "/api/v1/klinik", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreateKlinik(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusCreated)
		}

		var created models.Klinik
		_ = json.NewDecoder(rr.Body).Decode(&created)
		if created.ID != "mock-klinik-id" {
			t.Errorf("ID = %v, want 'mock-klinik-id'", created.ID)
		}
		if created.NamaKlinik != kNew.NamaKlinik {
			t.Errorf("NamaKlinik = %v, want %v", created.NamaKlinik, kNew.NamaKlinik)
		}
	})

	t.Run("Create Klinik - Invalid Name", func(t *testing.T) {
		kInvalid := models.Klinik{
			NamaKlinik: "",
		}
		body, _ := json.Marshal(kInvalid)
		req := httptest.NewRequest("POST", "/api/v1/klinik", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreateKlinik(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("Get All Klinik", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/klinik", nil)
		rr := httptest.NewRecorder()

		h.GetAllKlinik(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		var list []models.Klinik
		_ = json.NewDecoder(rr.Body).Decode(&list)
		if len(list) < 1 {
			t.Error("Harus mengembalikan minimal 1 klinik")
		}
	})

	t.Run("Get Klinik By ID - Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/klinik/k-1", nil)
		// Set URL parameter chi
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "k-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.GetKlinikByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		var k models.Klinik
		_ = json.NewDecoder(rr.Body).Decode(&k)
		if k.ID != "k-1" {
			t.Errorf("Klinik ID = %v, want k-1", k.ID)
		}
	})

	t.Run("Get Klinik By ID - Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/klinik/k-notfound", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "k-notfound")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.GetKlinikByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("Update Klinik - Success", func(t *testing.T) {
		kUpdate := models.Klinik{
			NamaKlinik:     "Klinik Awal Diperbarui",
			Lat:            1.23,
			Lng:            4.56,
			KapasitasAktif: 15,
		}
		body, _ := json.Marshal(kUpdate)
		req := httptest.NewRequest("PUT", "/api/v1/klinik/k-1", bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "k-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.UpdateKlinik(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		// Verifikasi perubahan di repo
		if repo.Kliniks["k-1"].NamaKlinik != kUpdate.NamaKlinik {
			t.Errorf("NamaKlinik setelah update = %v, want %v", repo.Kliniks["k-1"].NamaKlinik, kUpdate.NamaKlinik)
		}
	})

	t.Run("Delete Klinik - Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/klinik/k-1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "k-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		h.DeleteKlinik(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", rr.Code, http.StatusOK)
		}

		if _, exists := repo.Kliniks["k-1"]; exists {
			t.Error("Klinik 'k-1' harusnya sudah terhapus")
		}
	})
}
