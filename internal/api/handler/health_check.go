package handler

import (
	"github.com/HotPotatoC/roadmap_gen/pkg/render"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	return render.OK(c, "OK", nil)
}
