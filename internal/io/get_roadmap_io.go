package io

import (
	"time"

	"github.com/roadmap-thesis/backend/internal/domain"
)

type GetRoadmapBySlugOutput struct {
	ID                   int                            `json:"id"`
	Title                string                         `json:"title"`
	Slug                 string                         `json:"slug"`
	Description          string                         `json:"description"`
	Creator              GetRoadmapBySlugOutputCreator  `json:"creator"`
	Topics               []GetRoadmapBySlugOutputTopics `json:"topics"`
	TotalTopics          int                            `json:"total_topics"`
	CompletionPercentage float64                        `json:"completion_percentage"`
	CreatedAt            time.Time                      `json:"created_at"`
	UpdatedAt            time.Time                      `json:"updated_at"`
}

func (o *GetRoadmapBySlugOutput) Make(roadmap domain.Roadmap) {
	o.ID = roadmap.ID
	o.Title = roadmap.Title
	o.Slug = roadmap.Slug
	o.Description = roadmap.Description
	o.Creator = GetRoadmapBySlugOutputCreator{
		ID:   roadmap.Account.ID,
		Name: roadmap.Account.Name,
	}
	o.TotalTopics = roadmap.TotalTopics()
	o.CompletionPercentage = roadmap.CompletionPercentage()
	o.CreatedAt = roadmap.CreatedAt
	o.UpdatedAt = roadmap.UpdatedAt

	topicMap := make(map[int][]GetRoadmapBySlugOutputTopics)

	for _, topic := range roadmap.Topics {
		outputTopic := GetRoadmapBySlugOutputTopics{
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

	var buildTopics func(parentID int) []GetRoadmapBySlugOutputTopics
	buildTopics = func(parentID int) []GetRoadmapBySlugOutputTopics {
		var outputTopics []GetRoadmapBySlugOutputTopics
		for _, topic := range topicMap[parentID] {
			topic.Subtopics = buildTopics(topic.ID)
			outputTopics = append(outputTopics, topic)
		}
		return outputTopics
	}

	o.Topics = buildTopics(0)
}

type GetRoadmapBySlugOutputCreator struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetRoadmapBySlugOutputTopics struct {
	ID          int                            `json:"id"`
	RoadmapID   int                            `json:"roadmap_id"`
	ParentID    int                            `json:"parent_id"`
	Title       string                         `json:"title"`
	Slug        string                         `json:"slug"`
	Description string                         `json:"description"`
	Order       int                            `json:"order"`
	Finished    bool                           `json:"finished"`
	Subtopics   []GetRoadmapBySlugOutputTopics `json:"subtopics"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}
