package io

import "time"

type GetRoadmapOutput struct {
	ID          int                      `json:"id"`
	Title       string                   `json:"title"`
	Slug        string                   `json:"slug"`
	Description string                   `json:"description"`
	Creator     GetRoadmapOutputCreator  `json:"creator"`
	Topics      []GetRoadmapOutputTopics `json:"topics"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type GetRoadmapOutputCreator struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
