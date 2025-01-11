package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/HotPotatoC/roadmap_gen/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Stack().Logger()
	if os.Getenv("APP_ENV") == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Bootstrapping config and clients")
	config.Init()
}

func main() {
	var command *string = flag.String("command", "up", "migration command (up/down/reset)")

	flag.Parse()
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

	dir := "./db"
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
