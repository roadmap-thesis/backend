package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

type Repository struct {
	Account *AccountRepository
}

func New(db DB) *Repository {
	return &Repository{
		Account: NewAccountRepository(db),
	}
}
