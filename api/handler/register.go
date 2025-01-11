package handler

import (
	"net/http"

	"github.com/HotPotatoC/roadmap_gen/api/response"
	"github.com/HotPotatoC/roadmap_gen/backend"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var input backend.RegisterInput

	if err := c.Bind(&input); err != nil {
		return ErrInvalidData
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	err := h.backend.Register(c.Request().Context(), input)
	if err != nil {
		return err
	}

	res := response.NewSuccessResponse("Successfully registered", nil)
	return c.JSON(http.StatusOK, res)
}
