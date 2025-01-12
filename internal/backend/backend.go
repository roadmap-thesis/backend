package backend

import (
	"github.com/HotPotatoC/roadmap_gen/internal/repository"
)

type Backend struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *Backend {
	return &Backend{
		repository: repository,
	}
}
