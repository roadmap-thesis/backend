package main

import (
	"context"
	"os"

	"github.com/HotPotatoC/roadmap_gen/api"
	"github.com/HotPotatoC/roadmap_gen/backend"
	"github.com/HotPotatoC/roadmap_gen/clients"
	"github.com/HotPotatoC/roadmap_gen/config"
	"github.com/HotPotatoC/roadmap_gen/repository"
	_ "github.com/joho/godotenv/autoload"
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
	ctx := context.Background()

	clients, err := clients.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize clients")
	}

	repository := repository.New(clients.DB)
	backend := backend.New(repository)
	server := api.New(backend)

	exit := server.Listen(config.Port())

	signal := <-exit
	server.Shutdown(ctx, signal)
}
