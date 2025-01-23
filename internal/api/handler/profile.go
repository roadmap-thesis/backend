package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/roadmap-thesis/backend/pkg/render"
)

func (h *Handler) Profile(c echo.Context) error {
	output, err := h.backend.Profile(c.Request().Context())
	if err != nil {
		return err
	}

	return render.OK(c, "Profile details.", output)
}
