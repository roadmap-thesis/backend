package repository

import (
	"github.com/HotPotatoC/roadmap_gen/internal/database"
	"github.com/HotPotatoC/roadmap_gen/internal/domain"
)

type Repository struct {
	Account domain.AccountRepository
}

func New(db database.PsqlConnection) *Repository {
	return &Repository{
		Account: NewAccountRepository(db),
	}
}
