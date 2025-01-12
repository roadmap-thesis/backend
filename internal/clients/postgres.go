package clients

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgreSQLClient creates a new PostgreSQL pgxpool client
func NewPostgreSQLClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := ExtractDatabaseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func ExtractDatabaseConfig(connString string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = config.DatabaseMaxConns()
	cfg.MinConns = config.DatabaseMinConns()
	cfg.MaxConnLifetime = config.DatabaseMaxConnLifetime()
	cfg.MaxConnIdleTime = config.DatabaseMaxConnIdleTime()
	cfg.HealthCheckPeriod = config.DatabaseHealthCheckPeriod()
	cfg.ConnConfig.ConnectTimeout = config.DatabaseDefaultConnectionTimeout()

	return cfg, nil
}
