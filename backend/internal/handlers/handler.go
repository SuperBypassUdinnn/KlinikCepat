package handlers

import (
	"KlinikCepat/internal/repository"
	"KlinikCepat/internal/services"
)

// Handler menampung semua dependencies controller API
type Handler struct {
	Repo          repository.RepositoryInterface
	TriageService *services.TriageService
}

// NewHandler membuat instance Handler baru
func NewHandler(repo repository.RepositoryInterface, triage *services.TriageService) *Handler {
	return &Handler{
		Repo:          repo,
		TriageService: triage,
	}
}
