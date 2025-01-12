package repository

import (
	"context"
	"errors"

	"github.com/HotPotatoC/roadmap_gen/internal/database"
	"github.com/jackc/pgx/v5"
)

// transactionManager manages database transaction exclusive to repository
type transactionManager struct {
	db database.PsqlConnection
}

func newTransactionManager(db database.PsqlConnection) *transactionManager {
	return &transactionManager{db: db}
}

type transactionFunc func(ctx context.Context, tx pgx.Tx) error

// WithTransaction starts a new transaction then calls fn(). Rollback() the transaction if
// an error was found from fn(), Commit() otherwise.
func (t *transactionManager) WithTransaction(ctx context.Context, fn transactionFunc) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		return errors.Join(tx.Rollback(ctx), err)
	}

	return tx.Commit(ctx)
}
