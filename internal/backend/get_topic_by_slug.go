package backend

import (
	"context"
	"errors"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
)

func (b *backend) GetTopicBySlug(ctx context.Context, slug string) (io.GetTopicOutput, error) {
	topic, err := b.repository.Topic.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, domain.ErrTopicNotFound) {
			return io.GetTopicOutput{}, apperrors.ResourceNotFound("topic")
		}
		return io.GetTopicOutput{}, err
	}

	return io.GetTopicOutput{
		ID:          topic.ID,
		RoadmapID:   topic.RoadmapID,
		ParentID:    topic.ParentID,
		Title:       topic.Title,
		Slug:        topic.Slug,
		Description: topic.Description,
		Order:       topic.Order,
		Finished:    topic.Finished,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
	}, nil
}
