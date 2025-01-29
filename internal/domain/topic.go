package domain

import (
	"errors"
	"time"

	"github.com/roadmap-thesis/backend/pkg/slug"
	"github.com/roadmap-thesis/backend/pkg/str"
)

const (
	TopicTable = "topics"
)

var (
	ErrTopicNotFound = errors.New("topic not found")
)

type Topic struct {
	ID          int
	RoadmapID   int
	ParentID    int
	Title       string
	Slug        string
	Description string
	Order       int
	Finished    bool

	Subtopics []*Topic

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTopic(title, description string) *Topic {
	return &Topic{
		Title:       title,
		Slug:        slug.Make(title + " " + str.Random(5)),
		Description: description,
		Finished:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (e *Topic) IsZero() bool {
	return e.ID == 0 &&
		e.RoadmapID == 0 &&
		e.ParentID == 0 &&
		e.Title == "" &&
		e.Slug == "" &&
		e.Description == "" &&
		e.Order == 0 &&
		e.Finished == false &&
		(e.Subtopics == nil || len(e.Subtopics) == 0) &&
		e.CreatedAt.IsZero() &&
		e.UpdatedAt.IsZero()
}

func (e *Topic) IsParent() bool {
	return e.ParentID == 0
}

func (e *Topic) IsChild() bool {
	return e.ParentID != 0
}

func (e *Topic) HasSubtopics() bool {
	return len(e.Subtopics) > 0
}

func (e *Topic) GetSubtopic(id int) *Topic {
	for _, subtopic := range e.Subtopics {
		if subtopic.ID == id {
			return subtopic
		}
	}

	return nil
}

func (e *Topic) AddSubtopic(subtopic *Topic) {
	if e.Subtopics == nil {
		e.Subtopics = make([]*Topic, 0)
	}

	subtopic.Order = len(e.Subtopics) + 1

	subtopic.ParentID = e.ID
	e.Subtopics = append(e.Subtopics, subtopic)
}

func (e *Topic) Update(title, description, slug string) {
	e.Title = title
	e.Description = description
	e.Slug = slug
	e.UpdateChangelog()
}

func (e *Topic) MarkAsFinished() {
	e.Finished = true
	e.UpdateChangelog()
}

func (e *Topic) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
