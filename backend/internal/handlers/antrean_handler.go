package handlers

import (
	"encoding/json"
	"net/http"
	"KlinikCepat/internal/models"

	"github.com/go-chi/chi/v5"
)

// ProcessTriage handles POST /api/v1/triage
func (h *Handler) ProcessTriage(w http.ResponseWriter, r *http.Request) {
	var req models.TriageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	res, err := h.TriageService.ProcessTriage(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

// GetAntreanByKlinikID handles GET /api/v1/admin/antrean
func (h *Handler) GetAntreanByKlinikID(w http.ResponseWriter, r *http.Request) {
	klinikID := r.URL.Query().Get("klinik_id")
	if klinikID == "" {
		http.Error(w, `{"error": "Parameter 'klinik_id' wajib diisi"}`, http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "Menunggu" // Default status
	}

	antreans, err := h.Repo.GetAntreanByKlinikID(r.Context(), klinikID, status)
	if err != nil {
		http.Error(w, `{"error": "Gagal mengambil data antrean: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(antreans)
}

// UpdateStatusAntrean handles PUT /api/v1/admin/antrean/{id}/status
func (h *Handler) UpdateStatusAntrean(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "ID tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	// Payload struct untuk parsing request body
	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	if payload.Status == "" {
		http.Error(w, `{"error": "Status wajib diisi"}`, http.StatusBadRequest)
		return
	}

	// Validasi status sesuai enum yang ada di database
	if payload.Status != "Menunggu" && payload.Status != "Selesai" && payload.Status != "Dilewati" {
		http.Error(w, `{"error": "Status harus berupa 'Menunggu', 'Selesai', atau 'Dilewati'"}`, http.StatusBadRequest)
		return
	}

	err := h.Repo.UpdateStatusAntrean(r.Context(), id, payload.Status)
	if err != nil {
		http.Error(w, `{"error": "Gagal memperbarui status antrean: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Status antrean berhasil diperbarui menjadi ` + payload.Status + `"}`))
}
