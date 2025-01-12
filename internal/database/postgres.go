package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PsqlConnection interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type Transactional[Repository any] interface {
	WithTx(db PsqlConnection) Repository
}
