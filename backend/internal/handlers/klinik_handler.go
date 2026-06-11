package handlers

import (
	"encoding/json"
	"net/http"
	"KlinikCepat/internal/models"

	"github.com/go-chi/chi/v5"
)

// CreateKlinik handles POST /api/v1/klinik
func (h *Handler) CreateKlinik(w http.ResponseWriter, r *http.Request) {
	var k models.Klinik
	if err := json.NewDecoder(r.Body).Decode(&k); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	if k.NamaKlinik == "" {
		http.Error(w, `{"error": "Nama klinik wajib diisi"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateKlinik(r.Context(), &k)
	if err != nil {
		http.Error(w, `{"error": "Gagal menyimpan data klinik: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(k)
}

// GetAllKlinik handles GET /api/v1/klinik
func (h *Handler) GetAllKlinik(w http.ResponseWriter, r *http.Request) {
	kliniks, err := h.Repo.GetAllKlinik(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Gagal mengambil data klinik: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(kliniks)
}

// GetKlinikByID handles GET /api/v1/klinik/{id}
func (h *Handler) GetKlinikByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	k, err := h.Repo.GetKlinikByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "Data klinik tidak ditemukan atau error: `+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(k)
}

// UpdateKlinik handles PUT /api/v1/klinik/{id}
func (h *Handler) UpdateKlinik(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	var k models.Klinik
	if err := json.NewDecoder(r.Body).Decode(&k); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}
	k.ID = id

	if k.NamaKlinik == "" {
		http.Error(w, `{"error": "Nama klinik wajib diisi"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.UpdateKlinik(r.Context(), &k)
	if err != nil {
		http.Error(w, `{"error": "Gagal memperbarui data klinik: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(k)
}

// DeleteKlinik handles DELETE /api/v1/klinik/{id}
func (h *Handler) DeleteKlinik(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.DeleteKlinik(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "Gagal menghapus data klinik: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Klinik berhasil dihapus"}`))
}
