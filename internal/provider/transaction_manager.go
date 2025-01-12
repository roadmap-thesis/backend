package provider

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{db: db}
}

type TransactionManagerFunc func(ctx context.Context, tx pgx.Tx) error

func (m *TransactionManager) Execute(ctx context.Context, fn TransactionManagerFunc) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		return errors.Join(tx.Rollback(ctx), err)
	}

	return tx.Commit(ctx)
}
