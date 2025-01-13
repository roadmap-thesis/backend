package api

import "github.com/HotPotatoC/roadmap_gen/internal/api/middleware"

func (s *Server) setupRoutes() {
	s.instance.GET("/health", s.handler.HealthCheck)

	s.instance.POST("/auth", s.handler.Auth)
	s.instance.GET("/profile", s.handler.Profile, middleware.Auth)
}
