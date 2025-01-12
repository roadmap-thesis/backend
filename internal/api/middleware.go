package api

import (
	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func (s *Server) setupMiddlewares() {
	s.instance.Use(middleware.CORS())
	s.instance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Debug().
					Str("uri", v.URI).
					Int("status", v.Status).
					Send()
			} else {
				if config.AppEnv() != "production" || v.Status >= 500 {
					log.Error().
						Err(v.Error).
						Str("uri", v.URI).
						Int("status", v.Status).
						Send()
				}
			}
			return nil
		},
	}))
}
