package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error) {
	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}

	account, err := b.repository.Account.GetByID(ctx, roadmap.AccountID)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}

	roadmap.SetCreator(account)

	return b.makeGetRoadmapBySlugOutput(roadmap), nil
}

func (b *backend) makeGetRoadmapBySlugOutput(roadmap domain.Roadmap) io.GetRoadmapOutput {
	output := io.GetRoadmapOutput{
		ID:          roadmap.ID,
		Title:       roadmap.Title,
		Slug:        roadmap.Slug,
		Description: roadmap.Description,
		Creator: io.GetRoadmapOutputCreator{
			ID:   roadmap.Account.ID,
			Name: roadmap.Account.Name,
		},
		TotalTopics:          roadmap.TotalTopics(),
		CompletionPercentage: roadmap.CompletionPercentage(),
		CreatedAt:            roadmap.CreatedAt,
		UpdatedAt:            roadmap.UpdatedAt,
	}

	// Map topics to output also map subtopics into topics
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
		var outputTopics []io.GetRoadmapOutputTopics
		for _, topic := range topicMap[parentID] {
			topic.Subtopics = buildTopics(topic.ID)
			outputTopics = append(outputTopics, topic)
		}
		return outputTopics
	}

	output.Topics = buildTopics(0)

	return output
}
