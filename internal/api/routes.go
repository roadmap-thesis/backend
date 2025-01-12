package api

func (s *Server) setupRoutes() {
	s.instance.GET("/health", s.api.HealthCheck)

	s.instance.POST("/register", s.api.Register)
}
