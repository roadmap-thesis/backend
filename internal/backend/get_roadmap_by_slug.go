package backend

import (
	"context"
	"errors"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/domain/object"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error) {
	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, domain.ErrRoadmapNotFound) {
			return io.GetRoadmapOutput{}, apperrors.ResourceNotFound("roadmap")
		}
		return io.GetRoadmapOutput{}, err
	}

	personalizationOptions, err := b.repository.PersonalizationOptions.GetByRoadmapID(ctx, roadmap.ID)
	if err != nil {
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
			DailyTimeAvailability: object.NewIntervalFromDuration(personalizationOptions.DailyTimeAvailability),
			TotalDuration:         object.NewIntervalFromDuration(personalizationOptions.TotalDuration),
			SkillLevel:            personalizationOptions.SkillLevel.String(),
			AdditionalInfo:        personalizationOptions.AdditionalInfo,
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

	var buildTopics func(parentID int) []io.GetRoadmapOutputTopics
	buildTopics = func(parentID int) []io.GetRoadmapOutputTopics {
		outputTopics := topicMap[parentID]
		for i := range outputTopics {
			outputTopics[i].Subtopics = buildTopics(outputTopics[i].ID)
		}
		return outputTopics
	}

	output.Topics = buildTopics(0)

	return output, nil
}
