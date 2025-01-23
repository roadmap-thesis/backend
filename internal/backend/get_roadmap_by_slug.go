package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapBySlugOutput, error) {
	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		return io.GetRoadmapBySlugOutput{}, err
	}

	account, err := b.repository.Account.GetByID(ctx, roadmap.AccountID)
	if err != nil {
		return io.GetRoadmapBySlugOutput{}, err
	}

	roadmap.SetCreator(account)

	return b.makeGetRoadmapBySlugOutput(roadmap), nil
}

func (b *backend) makeGetRoadmapBySlugOutput(roadmap domain.Roadmap) io.GetRoadmapBySlugOutput {
	output := io.GetRoadmapBySlugOutput{
		ID:          roadmap.ID,
		Title:       roadmap.Title,
		Slug:        roadmap.Slug,
		Description: roadmap.Description,
		Creator: io.GetRoadmapBySlugOutputCreator{
			ID:   roadmap.Account.ID,
			Name: roadmap.Account.Name,
		},
		TotalTopics:          roadmap.TotalTopics(),
		CompletionPercentage: roadmap.CompletionPercentage(),
		CreatedAt:            roadmap.CreatedAt,
		UpdatedAt:            roadmap.UpdatedAt,
	}

	// Map topics to output also map subtopics into topics
	topicMap := make(map[int][]io.GetRoadmapBySlugOutputTopics)

	for _, topic := range roadmap.Topics {
		outputTopic := io.GetRoadmapBySlugOutputTopics{
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

	var buildTopics func(parentID int) []io.GetRoadmapBySlugOutputTopics
	buildTopics = func(parentID int) []io.GetRoadmapBySlugOutputTopics {
		var outputTopics []io.GetRoadmapBySlugOutputTopics
		for _, topic := range topicMap[parentID] {
			topic.Subtopics = buildTopics(topic.ID)
			outputTopics = append(outputTopics, topic)
		}
		return outputTopics
	}

	output.Topics = buildTopics(0)

	return output
}
