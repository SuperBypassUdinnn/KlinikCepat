package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"KlinikCepat/internal/models"
	"KlinikCepat/internal/services"
)

// InviteClinicAdmin menangani:
// POST /api/v1/superadmin/admin-klinik/invite
func (h *Handler) InviteClinicAdmin(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h.ClinicAdminService == nil {
		writeClinicAdminError(
			w,
			http.StatusInternalServerError,
			"Service Admin Klinik belum dikonfigurasi",
		)
		return
	}

	// Batasi ukuran payload agar request tidak berlebihan.
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		1<<20,
	)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload models.InviteClinicAdminRequest

	if err := decoder.Decode(&payload); err != nil {
		writeClinicAdminError(
			w,
			http.StatusBadRequest,
			"Payload JSON tidak valid",
		)
		return
	}

	// Pastikan body hanya berisi satu objek JSON.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeClinicAdminError(
			w,
			http.StatusBadRequest,
			"Payload hanya boleh berisi satu objek JSON",
		)
		return
	}

	result, err := h.ClinicAdminService.InviteClinicAdmin(
		r.Context(),
		payload.Email,
		payload.KlinikID,
	)
	if err != nil {
		handleClinicAdminInviteError(w, err)
		return
	}

	writeClinicAdminJSON(
		w,
		http.StatusCreated,
		result,
	)
}

func handleClinicAdminInviteError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(
		err,
		services.ErrClinicAdminInvalidEmail,
	):
		writeClinicAdminError(
			w,
			http.StatusBadRequest,
			"Format email tidak valid",
		)

	case errors.Is(
		err,
		services.ErrClinicAdminKlinikRequired,
	):
		writeClinicAdminError(
			w,
			http.StatusBadRequest,
			"Klinik wajib dipilih",
		)

	case errors.Is(
		err,
		services.ErrClinicAdminKlinikNotFound,
	):
		writeClinicAdminError(
			w,
			http.StatusNotFound,
			"Klinik tidak ditemukan",
		)

	default:
		handleSupabaseInviteError(w, err)
	}
}

func handleSupabaseInviteError(
	w http.ResponseWriter,
	err error,
) {
	var apiError *services.SupabaseAdminAPIError

	if errors.As(err, &apiError) {
		message := strings.ToLower(apiError.Message)

		if apiError.StatusCode == http.StatusConflict ||
			apiError.StatusCode ==
				http.StatusUnprocessableEntity ||
			strings.Contains(message, "already registered") ||
			strings.Contains(message, "already been registered") ||
			strings.Contains(message, "already exists") {
			writeClinicAdminError(
				w,
				http.StatusConflict,
				"Email tersebut sudah terdaftar",
			)
			return
		}

		if apiError.StatusCode ==
			http.StatusTooManyRequests {
			writeClinicAdminError(
				w,
				http.StatusTooManyRequests,
				"Terlalu banyak permintaan undangan. Coba lagi nanti",
			)
			return
		}

		log.Printf(
			"Supabase Admin API error: %v",
			apiError,
		)

		writeClinicAdminError(
			w,
			http.StatusBadGateway,
			"Gagal mengirim undangan melalui layanan autentikasi",
		)
		return
	}

	log.Printf(
		"Gagal mengundang Admin Klinik: %v",
		err,
	)

	writeClinicAdminError(
		w,
		http.StatusInternalServerError,
		"Gagal membuat akun Admin Klinik",
	)
}

func writeClinicAdminError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeClinicAdminJSON(
		w,
		status,
		map[string]string{
			"error": message,
		},
	)
}

func writeClinicAdminJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf(
			"Gagal menulis JSON response: %v",
			err,
		)
	}
}
