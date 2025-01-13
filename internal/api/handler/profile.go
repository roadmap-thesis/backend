package handler

import (
	"github.com/HotPotatoC/roadmap_gen/internal/api/render"
	"github.com/labstack/echo/v4"
)

type ProfileOutput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) Profile(c echo.Context) error {
	output, err := h.backend.Profile(c.Request().Context())
	if err != nil {
		return err
	}

	return render.OK(c, "Profile details.", ProfileOutput{
		ID:   output.ID,
		Name: output.Name,
	})
}
