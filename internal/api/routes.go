package api

import "github.com/HotPotatoC/roadmap_gen/internal/api/middleware"

func (s *Server) setupRoutes() {
	s.instance.GET("/health", s.api.HealthCheck)

	s.instance.POST("/auth", s.api.Auth)
	s.instance.GET("/profile", s.api.Profile, middleware.Auth)
}
