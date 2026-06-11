package handlers

import (
	"encoding/json"
	"net/http"
	"KlinikCepat/internal/models"

	"github.com/go-chi/chi/v5"
)

// CreateGejala handles POST /api/v1/gejala
func (h *Handler) CreateGejala(w http.ResponseWriter, r *http.Request) {
	var g models.Gejala
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	if g.NamaGejala == "" {
		http.Error(w, `{"error": "Nama gejala wajib diisi"}`, http.StatusBadRequest)
		return
	}
	if g.BobotUrgensi < 1 || g.BobotUrgensi > 10 {
		http.Error(w, `{"error": "Bobot urgensi harus berkisar antara 1 hingga 10"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateGejala(r.Context(), &g)
	if err != nil {
		http.Error(w, `{"error": "Gagal menyimpan data gejala: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(g)
}

// GetAllGejala handles GET /api/v1/gejala
func (h *Handler) GetAllGejala(w http.ResponseWriter, r *http.Request) {
	gejalas, err := h.Repo.GetAllGejala(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Gagal mengambil data gejala: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(gejalas)
}

// GetGejalaByID handles GET /api/v1/gejala/{id}
func (h *Handler) GetGejalaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	g, err := h.Repo.GetGejalaByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "Data gejala tidak ditemukan atau error: `+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(g)
}

// UpdateGejala handles PUT /api/v1/gejala/{id}
func (h *Handler) UpdateGejala(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	var g models.Gejala
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}
	g.ID = id

	if g.NamaGejala == "" {
		http.Error(w, `{"error": "Nama gejala wajib diisi"}`, http.StatusBadRequest)
		return
	}
	if g.BobotUrgensi < 1 || g.BobotUrgensi > 10 {
		http.Error(w, `{"error": "Bobot urgensi harus berkisar antara 1 hingga 10"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.UpdateGejala(r.Context(), &g)
	if err != nil {
		http.Error(w, `{"error": "Gagal memperbarui data gejala: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(g)
}

// DeleteGejala handles DELETE /api/v1/gejala/{id}
func (h *Handler) DeleteGejala(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.DeleteGejala(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "Gagal menghapus data gejala: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Gejala berhasil dihapus"}`))
}
