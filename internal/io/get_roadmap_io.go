package io

import "time"

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
