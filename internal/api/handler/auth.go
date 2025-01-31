package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"github.com/roadmap-thesis/backend/pkg/render"
)

func (h *Handler) Auth(c echo.Context) error {
	var input io.AuthInput

	if err := c.Bind(&input); err != nil {
		return apperrors.InvalidData()
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	output, err := h.backend.Auth(c.Request().Context(), input)
	if err != nil {
		return err
	}

	if output.Created {
		return render.Created(c, "Successfully registered.", output)
	} else {
		return render.OK(c, "Successfully logged in.", output)
	}
}
