package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roadmap-thesis/backend/pkg/config"
)

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
