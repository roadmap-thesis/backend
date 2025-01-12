package handler

import "github.com/HotPotatoC/roadmap_gen/internal/backend"

type Handler struct {
	backend *backend.Backend
}

func New(backend *backend.Backend) *Handler {
	return &Handler{backend: backend}
}
