package handler

import (
	"github.com/HotPotatoC/roadmap_gen/internal/api/res"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	return res.OK(c, "OK", nil)
}
