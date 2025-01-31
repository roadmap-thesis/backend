package backend

import (
	"context"
	"errors"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/domain/object"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error) {
	ctx, span := tracer.Start(ctx, "(*backend.GetRoadmapBySlug)", trace.WithAttributes(attribute.String("slug", slug)))
	defer span.End()

	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, domain.ErrRoadmapNotFound) {
			return io.GetRoadmapOutput{}, apperrors.ResourceNotFound("roadmap")
		}
		return io.GetRoadmapOutput{}, err
	}

	account, err := b.repository.Account.GetByID(ctx, roadmap.AccountID)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}
	roadmap.SetCreator(account)

	output := io.GetRoadmapOutput{
		ID:          roadmap.ID,
		Title:       roadmap.Title,
		Slug:        roadmap.Slug,
		Description: roadmap.Description,
		Creator: io.GetRoadmapOutputCreator{
			ID:     roadmap.Account.ID,
			Name:   roadmap.Account.Profile.Name,
			Avatar: roadmap.Account.Profile.Avatar,
		},
		PersonalizationOpts: io.GetRoadmapOutputPersonalizationOptions{
			DailyTimeAvailability: object.NewIntervalFromDuration(roadmap.PersonalizationOptions.DailyTimeAvailability),
			TotalDuration:         object.NewIntervalFromDuration(roadmap.PersonalizationOptions.TotalDuration),
			SkillLevel:            roadmap.PersonalizationOptions.SkillLevel.String(),
			AdditionalInfo:        roadmap.PersonalizationOptions.AdditionalInfo,
		},
		TotalTopics:          roadmap.TotalTopics(),
		CompletionPercentage: roadmap.CompletionPercentage(),
		CreatedAt:            roadmap.CreatedAt,
		UpdatedAt:            roadmap.UpdatedAt,
	}

	topicMap := make(map[int][]io.GetRoadmapOutputTopics)

	for _, topic := range roadmap.Topics {
		outputTopic := io.GetRoadmapOutputTopics{
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
		}

		topicMap[topic.ParentID] = append(topicMap[topic.ParentID], outputTopic)
	}

	var buildTopics func(ctx context.Context, parentID int) []io.GetRoadmapOutputTopics
	buildTopics = func(ctx context.Context, parentID int) []io.GetRoadmapOutputTopics {
		traceCtx, span := tracer.Start(ctx, "buildTopics", trace.WithAttributes(attribute.Int("parentID", parentID)))
		defer span.End()

		outputTopics := topicMap[parentID]
		for i := range outputTopics {
			outputTopics[i].Subtopics = buildTopics(traceCtx, outputTopics[i].ID)
		}

		return outputTopics
	}

	output.Topics = buildTopics(ctx, 0)

	return output, nil
}
