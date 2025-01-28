package io

import (
	"time"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/domain/object"
)

type GetRoadmapOutput struct {
	ID                   int                                    `json:"id"`
	Title                string                                 `json:"title"`
	Slug                 string                                 `json:"slug"`
	Description          string                                 `json:"description"`
	Creator              GetRoadmapOutputCreator                `json:"creator"`
	Topics               []GetRoadmapOutputTopics               `json:"topics"`
	TotalTopics          int                                    `json:"total_topics"`
	PersonalizationOpts  GetRoadmapOutputPersonalizationOptions `json:"personalization_options"`
	CompletionPercentage float64                                `json:"completion_percentage"`
	CreatedAt            time.Time                              `json:"created_at"`
	UpdatedAt            time.Time                              `json:"updated_at"`
}

func (o *GetRoadmapOutput) Make(roadmap domain.Roadmap, personalizationOptions domain.PersonalizationOptions) {
	o.ID = roadmap.ID
	o.Title = roadmap.Title
	o.Slug = roadmap.Slug
	o.Description = roadmap.Description
	o.Creator = GetRoadmapOutputCreator{
		ID:     roadmap.Account.ID,
		Name:   roadmap.Account.Profile.Name,
		Avatar: roadmap.Account.Profile.Avatar,
	}

	o.PersonalizationOpts = GetRoadmapOutputPersonalizationOptions{
		DailyTimeAvailability: object.NewIntervalFromDuration(personalizationOptions.DailyTimeAvailability),
		TotalDuration:         object.NewIntervalFromDuration(personalizationOptions.TotalDuration),
		SkillLevel:            string(personalizationOptions.SkillLevel),
		AdditionalInfo:        personalizationOptions.AdditionalInfo,
	}
	o.TotalTopics = roadmap.TotalTopics()
	o.CompletionPercentage = roadmap.CompletionPercentage()
	o.CreatedAt = roadmap.CreatedAt
	o.UpdatedAt = roadmap.UpdatedAt

	topicMap := make(map[int][]GetRoadmapOutputTopics)

	for _, topic := range roadmap.Topics {
		outputTopic := GetRoadmapOutputTopics{
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

	var buildTopics func(parentID int) []GetRoadmapOutputTopics
	buildTopics = func(parentID int) []GetRoadmapOutputTopics {
		var outputTopics []GetRoadmapOutputTopics
		for _, topic := range topicMap[parentID] {
			topic.Subtopics = buildTopics(topic.ID)
			outputTopics = append(outputTopics, topic)
		}
		return outputTopics
	}

	o.Topics = buildTopics(0)
}

type GetRoadmapOutputCreator struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type GetRoadmapOutputTopics struct {
	ID          int                      `json:"id"`
	RoadmapID   int                      `json:"roadmap_id"`
	ParentID    int                      `json:"parent_id"`
	Title       string                   `json:"title"`
	Slug        string                   `json:"slug"`
	Description string                   `json:"description"`
	Order       int                      `json:"order"`
	Finished    bool                     `json:"finished"`
	Subtopics   []GetRoadmapOutputTopics `json:"subtopics"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type GetRoadmapOutputPersonalizationOptions struct {
	DailyTimeAvailability object.Interval `json:"daily_time_availability"`
	TotalDuration         object.Interval `json:"total_duration"`
	SkillLevel            string          `json:"skill_level"`
	AdditionalInfo        string          `json:"additional_info"`
}
