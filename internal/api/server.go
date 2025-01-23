package api

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/api/handler"
	"github.com/roadmap-thesis/backend/internal/backend"
	"github.com/roadmap-thesis/backend/pkg/server"
)

type Server struct {
	instance *server.Server
	backend  backend.Backend
	handler  *handler.Handler
}

func New(port string, backend backend.Backend) *Server {
	instance := server.New(port)

	handler := handler.New(backend)
	api := &Server{
		backend:  backend,
		handler:  handler,
		instance: instance,
	}

	api.setupMiddlewares()
	api.setupRoutes()
	api.instance.HTTPErrorHandler = handler.ErrorHandler

	return api
}

func (s *Server) Start(ctx context.Context) {
	exit := s.instance.Listen()

	signal := <-exit
	s.instance.Shutdown(ctx, signal)
}
