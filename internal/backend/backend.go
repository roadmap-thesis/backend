package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/internal/repository"
	"github.com/roadmap-thesis/backend/pkg/openai"
)

type Backend interface {
	Auth(ctx context.Context, input io.AuthInput) (io.AuthOutput, error)
	GetProfile(ctx context.Context) (io.ProfileOutput, error)

	GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error)
	GenerateRoadmap(ctx context.Context, input io.GenerateRoadmapInput) (io.GenerateRoadmapOutput, error)
}

type backend struct {
	repository *repository.Repository
	openai     *openai.Client
}

func New(repository *repository.Repository, openai *openai.Client) Backend {
	return &backend{
		repository: repository,
		openai:     openai,
	}
}
