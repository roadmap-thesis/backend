package domain

import "time"

const (
	RoadmapTable = "roadmaps"
)

type Roadmap struct {
	ID    int
	Topic string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Roadmap) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
