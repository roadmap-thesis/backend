package api

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HotPotatoC/roadmap_gen/internal/api/handler"
	"github.com/HotPotatoC/roadmap_gen/internal/backend"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	instance *echo.Echo
	backend  backend.Backend
	api      *handler.Handler
}

func New(backend backend.Backend) *Server {
	echoInstance := NewEchoInstance()

	api := handler.New(backend)
	srv := &Server{
		instance: echoInstance,
		backend:  backend,
		api:      api,
	}

	srv.setupMiddlewares()
	srv.setupRoutes()
	srv.instance.HTTPErrorHandler = api.ErrorHandler

	return srv
}

func (s *Server) Listen(port string) <-chan os.Signal {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		log.Info().Msgf("Listening on %s", port)
		s.instance.Start(":" + port)
	}()

	return exitSignal
}

// Shutdown gracefully shuts down the API server
func (s *Server) Shutdown(ctx context.Context, signal os.Signal) {
	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	shutdownChan := make(chan error, 1)

	go func() {
		log.Warn().Any("signal", signal.String()).Msg("received signal, shutting down...")
		shutdownChan <- s.instance.Shutdown(ctx)
	}()

	select {
	case <-timeout.Done():
		log.Warn().Msg("shutdown timed out, forcing exit")
		os.Exit(1)
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal().Err(err).Msg("there was an error shutting down")
		} else {
			log.Info().Msg("server shutdown complete")
		}
	}
}
