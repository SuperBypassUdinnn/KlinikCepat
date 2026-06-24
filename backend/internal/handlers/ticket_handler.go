package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"KlinikCepat/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// GetPublicTicket menangani:
// GET /api/v1/ticket/{publicToken}
func (h *Handler) GetPublicTicket(
	w http.ResponseWriter,
	r *http.Request,
) {
	publicToken := strings.TrimSpace(
		chi.URLParam(r, "publicToken"),
	)

	if publicToken == "" {
		writeTicketError(
			w,
			http.StatusBadRequest,
			"Token tiket wajib diisi",
		)
		return
	}

	ticket, err := h.Repo.GetPublicTicketByToken(
		r.Context(),
		publicToken,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeTicketError(
				w,
				http.StatusNotFound,
				"Tiket tidak ditemukan",
			)
			return
		}

		log.Printf(
			"Gagal mengambil tiket: %v",
			err,
		)

		writeTicketError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil tiket",
		)
		return
	}

	writeTicketJSON(
		w,
		http.StatusOK,
		ticket,
	)
}

// CheckPublicTicket menangani:
// POST /api/v1/ticket/check
func (h *Handler) CheckPublicTicket(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		1<<20,
	)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload models.CheckTicketRequest

	if err := decoder.Decode(&payload); err != nil {
		writeTicketError(
			w,
			http.StatusBadRequest,
			"Payload JSON tidak valid",
		)
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeTicketError(
			w,
			http.StatusBadRequest,
			"Payload hanya boleh berisi satu objek JSON",
		)
		return
	}

	kodeTiket := strings.ToUpper(
		strings.TrimSpace(payload.KodeTiket),
	)

	email := strings.ToLower(
		strings.TrimSpace(payload.Email),
	)

	if kodeTiket == "" {
		writeTicketError(
			w,
			http.StatusBadRequest,
			"Kode tiket wajib diisi",
		)
		return
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		writeTicketError(
			w,
			http.StatusBadRequest,
			"Format email tidak valid",
		)
		return
	}

	ticket, err := h.Repo.CheckPublicTicket(
		r.Context(),
		kodeTiket,
		email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeTicketError(
				w,
				http.StatusNotFound,
				"Kode tiket atau email tidak sesuai",
			)
			return
		}

		log.Printf(
			"Gagal mengecek tiket: %v",
			err,
		)

		writeTicketError(
			w,
			http.StatusInternalServerError,
			"Gagal mengecek tiket",
		)
		return
	}

	writeTicketJSON(
		w,
		http.StatusOK,
		ticket,
	)
}

func writeTicketError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeTicketJSON(
		w,
		status,
		map[string]string{
			"error": message,
		},
	)
}

func writeTicketJSON(
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
			"Gagal menulis respons tiket: %v",
			err,
		)
	}
}
