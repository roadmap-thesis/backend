package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/repository"
)

type Backend interface {
	Auth(ctx context.Context, input AuthInput) (AuthOutput, error)
	Profile(ctx context.Context) (ProfileOutput, error)
}

type backend struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) Backend {
	return &backend{
		repository: repository,
	}
}
