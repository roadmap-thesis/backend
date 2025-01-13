package repository

import (
	"github.com/HotPotatoC/roadmap_gen/internal/database"
)

type Repository struct {
	Account AccountRepository
}

func New(db database.Connection) *Repository {
	return &Repository{
		Account: NewAccountRepository(db),
	}
}
