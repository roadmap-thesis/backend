package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Connection interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)

	InTx(ctx context.Context, fn TransactionFunc) error

	Close()
}

// New creates a new PostgreSQL pgxpool client
func New(ctx context.Context, connString string) (Connection, error) {
	cfg, err := ExtractDatabaseConfig(connString)
	if err != nil {
		return nil, err
	}

	cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.Pool.Query(ctx, query, args...)
}

func (db *DB) Close() {
	db.Pool.Close()
}
