package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/internal/repository"
	"github.com/roadmap-thesis/backend/pkg/openai"
)

type Backend interface {
	Auth(ctx context.Context, input io.AuthInput) (io.AuthOutput, error)
	Profile(ctx context.Context) (io.ProfileOutput, error)

	GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error)
	GenerateRoadmap(ctx context.Context, input io.GenerateRoadmapInput) (io.GenerateRoadmapOutput, error)
}

type backend struct {
	openai     *openai.Client
	repository *repository.Repository
}

func New(openai *openai.Client, repository *repository.Repository) Backend {
	return &backend{
		openai:     openai,
		repository: repository,
	}
}
