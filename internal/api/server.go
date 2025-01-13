package api

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/api/handler"
	"github.com/HotPotatoC/roadmap_gen/internal/backend"
	"github.com/HotPotatoC/roadmap_gen/pkg/server"
)

type Server struct {
	srv     *server.Server
	backend backend.Backend
	handler *handler.Handler
}

func NewServer(port string, backend backend.Backend) *Server {
	srv := server.New(port)

	handler := handler.New(backend)
	api := &Server{
		backend: backend,
		handler: handler,
		srv:     srv,
	}

	api.setupMiddlewares()
	api.setupRoutes()
	api.srv.HTTPErrorHandler = handler.ErrorHandler

	return api
}

func (s *Server) Start(ctx context.Context) {
	exit := s.srv.Listen()

	signal := <-exit
	s.srv.Shutdown(ctx, signal)
}
