package api

import "github.com/HotPotatoC/roadmap_gen/internal/api/middleware"

func (s *Server) setupRoutes() {
	s.srv.GET("/health", s.handler.HealthCheck)

	s.srv.POST("/auth", s.handler.Auth)
	s.srv.GET("/profile", s.handler.Profile, middleware.Auth)
}
