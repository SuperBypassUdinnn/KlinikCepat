package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"KlinikCepat/internal/middleware"
	"KlinikCepat/internal/models"
	"KlinikCepat/internal/services"

	"github.com/go-chi/chi/v5"
)

// ProcessTriage handles POST /api/v1/triage
// ProcessTriage handles POST /api/v1/triage
func (h *Handler) ProcessTriage(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Gunakan struct payload khusus agar JSON pasti
	// terbaca sesuai nama field yang dikirim frontend.
	var payload struct {
		KlinikID    string               `json:"klinik_id"`
		NamaPasien  string               `json:"nama_pasien"`
		EmailPasien string               `json:"email_pasien"`
		Gejala      []models.GejalaInput `json:"gejala"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		http.Error(
			w,
			`{"error":"Payload JSON tidak valid"}`,
			http.StatusBadRequest,
		)
		return
	}

	req := models.TriageRequest{
		KlinikID:    strings.TrimSpace(payload.KlinikID),
		NamaPasien:  strings.TrimSpace(payload.NamaPasien),
		EmailPasien: strings.TrimSpace(payload.EmailPasien),
		Gejala:      payload.Gejala,
	}

	res, err := h.TriageService.ProcessTriage(
		r.Context(),
		&req,
	)
	if err != nil {
		var validationError *services.TriageValidationError

		if errors.As(err, &validationError) {
			writeTicketError(
				w,
				http.StatusBadRequest,
				validationError.Error(),
			)
			return
		}

		log.Printf(
			"Gagal memproses triage: %v",
			err,
		)

		writeTicketError(
			w,
			http.StatusInternalServerError,
			"Gagal memproses pendaftaran antrean",
		)
		return
	}

	writeTicketJSON(
		w,
		http.StatusCreated,
		res,
	)
}

// GetAntreanByKlinikID handles GET /api/v1/admin/antrean
func (h *Handler) GetAntreanByKlinikID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(
		r.Context(),
	)
	if !ok {
		http.Error(
			w,
			`{"error": "Unauthorized: Informasi user tidak ditemukan"}`,
			http.StatusUnauthorized,
		)
		return
	}

	requestedKlinikID := strings.TrimSpace(
		r.URL.Query().Get("klinik_id"),
	)

	var klinikID string

	switch claims.Role {
	case "klinik_admin":
		if claims.KlinikID == nil ||
			strings.TrimSpace(*claims.KlinikID) == "" {
			http.Error(
				w,
				`{"error": "Forbidden: Akun admin belum terhubung ke klinik"}`,
				http.StatusForbidden,
			)
			return
		}

		klinikID = strings.TrimSpace(
			*claims.KlinikID,
		)

		// Jika frontend masih mengirim ID klinik berbeda,
		// request langsung ditolak.
		if requestedKlinikID != "" &&
			requestedKlinikID != klinikID {
			http.Error(
				w,
				`{"error": "Forbidden: Anda tidak boleh mengakses klinik lain"}`,
				http.StatusForbidden,
			)
			return
		}

	case "superadmin":
		if requestedKlinikID == "" {
			http.Error(
				w,
				`{"error": "Parameter 'klinik_id' wajib diisi untuk superadmin"}`,
				http.StatusBadRequest,
			)
			return
		}

		klinikID = requestedKlinikID

	default:
		http.Error(
			w,
			`{"error": "Forbidden: Role tidak diizinkan"}`,
			http.StatusForbidden,
		)
		return
	}

	status := strings.TrimSpace(
		r.URL.Query().Get("status"),
	)

	if status == "" {
		status = "Menunggu"
	}

	if status != "Menunggu" &&
		status != "Selesai" &&
		status != "Dilewati" {
		http.Error(
			w,
			`{"error": "Status harus berupa 'Menunggu', 'Selesai', atau 'Dilewati'"}`,
			http.StatusBadRequest,
		)
		return
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

	claims, ok := middleware.GetClaimsFromContext(
		r.Context(),
	)
	if !ok {
		http.Error(
			w,
			`{"error": "Unauthorized: Informasi user tidak ditemukan"}`,
			http.StatusUnauthorized,
		)
		return
	}

	var scopedKlinikID *string

	switch claims.Role {
	case "klinik_admin":
		if claims.KlinikID == nil ||
			strings.TrimSpace(*claims.KlinikID) == "" {
			http.Error(
				w,
				`{"error": "Forbidden: Akun admin belum terhubung ke klinik"}`,
				http.StatusForbidden,
			)
			return
		}

		klinikID := strings.TrimSpace(
			*claims.KlinikID,
		)

		scopedKlinikID = &klinikID

	case "superadmin":
		// nil berarti superadmin tidak dibatasi klinik tertentu.
		scopedKlinikID = nil

	default:
		http.Error(
			w,
			`{"error": "Forbidden: Role tidak diizinkan"}`,
			http.StatusForbidden,
		)
		return
	}

	updated, err := h.Repo.UpdateStatusAntrean(
		r.Context(),
		id,
		payload.Status,
		scopedKlinikID,
	)

	if err != nil {
		http.Error(
			w,
			`{"error": "Gagal memperbarui status antrean: `+
				err.Error()+`"}`,
			http.StatusInternalServerError,
		)
		return
	}

	if !updated {
		http.Error(
			w,
			`{"error": "Antrean tidak ditemukan"}`,
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Status antrean berhasil diperbarui menjadi ` + payload.Status + `"}`))
}
