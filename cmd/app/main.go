package main

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/api"
	"github.com/HotPotatoC/roadmap_gen/internal/backend"
	"github.com/HotPotatoC/roadmap_gen/internal/clients"
	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/HotPotatoC/roadmap_gen/internal/logger"
	"github.com/HotPotatoC/roadmap_gen/internal/repository"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	config.Init()
	logger.Init()

	log.Info().Msg("Initialized config and clients")

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
