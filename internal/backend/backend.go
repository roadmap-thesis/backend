package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/internal/repository"
	"github.com/roadmap-thesis/backend/pkg/llm"
)

type Backend interface {
	Auth(ctx context.Context, input io.AuthInput) (io.AuthOutput, error)
	GetProfile(ctx context.Context) (io.GetProfileOutput, error)

	GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapBySlugOutput, error)
	GenerateRoadmap(ctx context.Context, input io.GenerateRoadmapInput) (io.GenerateRoadmapOutput, error)

	// ListUserRoadmaps(ctx context.Context) (io.ListUserRoadmapsOutput, error)
	// DeleteUserRoadmap(ctx context.Context, input io.DeleteUserRoadmapInput) (io.DeleteUserRoadmapOutput, error)
	// RegenerateRoadmap(ctx context.Context, input io.RegenerateRoadmapInput) (io.RegenerateRoadmapOutput, error)
	// GetTopic(ctx context.Context, input io.GetTopicInput) (io.GetTopicOutput, error)
	// GetTopicResources(ctx context.Context, input io.GetTopicResourcesInput) (io.GetTopicResourcesOutput, error)
	// MarkTopicAsFinish(ctx context.Context, input io.MarkTopicAsFinishInput) (io.TopicFinishOutput, error)
	// MarkTopicAsIncomplete(ctx context.Context, input io.MarkTopicAsIncompleteInput) (io.TopicFinishOutput, error)
}

type backend struct {
	repository *repository.Repository
	llm        llm.Client
}

func New(repository *repository.Repository, llm llm.Client) Backend {
	return &backend{
		repository: repository,
		llm:        llm,
	}
}
