package repository

import (
	"github.com/roadmap-thesis/backend/pkg/database"
)

type Repository struct {
	Account                *AccountRepository
	Roadmap                *RoadmapRepository
	Topic                  *TopicRepository
	PersonalizationOptions *PersonalizationOptionsRepository
}

func New(db database.Connection) *Repository {
	return &Repository{
		Account:                NewAccountRepository(db),
		Roadmap:                NewRoadmapRepository(db),
		Topic:                  NewTopicRepository(db),
		PersonalizationOptions: NewPersonalizationOptionsRepository(db),
	}
}
