package api

import "github.com/HotPotatoC/roadmap_gen/internal/api/middleware"

func (s *Server) setupRoutes() {
	s.instance.GET("/health", s.api.HealthCheck)

	s.instance.POST("/register", s.api.Register)
	s.instance.GET("/profile", s.api.Profile, middleware.Auth)
}
