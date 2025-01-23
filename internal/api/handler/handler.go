package handler

import "github.com/roadmap-thesis/backend/internal/backend"

type Handler struct {
	backend backend.Backend
}

func New(backend backend.Backend) *Handler {
	return &Handler{backend: backend}
}
