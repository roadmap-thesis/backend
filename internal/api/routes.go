package api

import "github.com/roadmap-thesis/backend/internal/api/middleware"

func (s *Server) setupRoutes() {
	s.instance.GET("/health", s.handler.HealthCheck)

	s.instance.POST("/auth", s.handler.Auth)
	s.instance.GET("/profile", s.handler.GetProfile, middleware.Auth)

	s.instance.GET("/roadmaps", s.handler.ListUserRoadmaps, middleware.Auth)
	s.instance.GET("/roadmaps/:slug", s.handler.GetRoadmapBySlug, middleware.Auth)
	s.instance.POST("/roadmaps", s.handler.GenerateRoadmap, middleware.Auth)
}
