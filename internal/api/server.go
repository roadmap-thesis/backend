package api

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/api/handler"
	"github.com/HotPotatoC/roadmap_gen/internal/backend"
	"github.com/HotPotatoC/roadmap_gen/pkg/server"
)

type Server struct {
	instance *server.Server
	backend  backend.Backend
	handler  *handler.Handler
}

func NewServer(port string, backend backend.Backend) *Server {
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
