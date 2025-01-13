package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*echo.Echo
	port string
}

func New(port string) *Server {
	instance := NewEchoInstance()

	srv := &Server{
		port: port,
		Echo: instance,
	}

	return srv
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) Listen() chan os.Signal {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		log.Info().Msgf("Listening on %s", s.port)
		s.Echo.Start(":" + s.port)
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
		shutdownChan <- s.Echo.Shutdown(ctx)
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
