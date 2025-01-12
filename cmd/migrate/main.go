package main

import (
	"database/sql"
	"flag"

	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/HotPotatoC/roadmap_gen/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func main() {
	var command *string = flag.String("command", "up", "migration command (up/down/reset)")

	flag.Parse()

	config.Init()
	logger.Init()

	log.Info().Msg("Initialized config and clients")

	err := goose.SetDialect("pgx")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize goose")
	}

	goose.SetTableName("schema_migrations")
	db, err := sql.Open("pgx", config.DatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize clients")
	}
	defer db.Close()

	dir := "./migrations"
	switch *command {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "reset":
		err = goose.Reset(db, dir)
	default:
		log.Fatal().Err(err).Msg("Unknown migration command")
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	log.Info().Msg("Migrations applied successfully!")
}
