package handlers

import (
	"encoding/json"
	"net/http"

	"KlinikCepat/internal/middleware"
	"KlinikCepat/internal/models"
)

// GetCurrentUser menangani GET /api/v1/auth/me.
func (h *Handler) GetCurrentUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	claims, ok := r.Context().
		Value(middleware.UserContextKey).(*middleware.JWTClaims)

	if !ok || claims == nil {
		http.Error(
			w,
			`{"error": "Unauthorized: Informasi user tidak ditemukan"}`,
			http.StatusUnauthorized,
		)
		return
	}

	response := models.AuthMeResponse{
		ID:       claims.Sub,
		Email:    claims.Email,
		Role:     claims.Role,
		KlinikID: claims.KlinikID,
	}

	if response.KlinikID != nil {
		klinik, err := h.Repo.GetKlinikByID(
			r.Context(),
			*response.KlinikID,
		)
		if err != nil {
			http.Error(
				w,
				`{"error":"Gagal mengambil data klinik user"}`,
				http.StatusInternalServerError,
			)
			return
		}

		response.NamaKlinik = &klinik.NamaKlinik
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			`{"error": "Gagal membuat respons"}`,
			http.StatusInternalServerError,
		)
	}
}
