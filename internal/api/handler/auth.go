package handler

import (
	"github.com/HotPotatoC/roadmap_gen/internal/api/render"
	"github.com/HotPotatoC/roadmap_gen/internal/backend"
	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/labstack/echo/v4"
)

type AuthOutput struct {
	Token string `json:"token"`
}

func (h *Handler) Auth(c echo.Context) error {
	var input backend.AuthInput

	if err := c.Bind(&input); err != nil {
		return commonerrors.InvalidData()
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	output, err := h.backend.Auth(c.Request().Context(), input)
	if err != nil {
		return err
	}

	if output.Created {
		return render.Created(c, "Successfully registered.", AuthOutput{
			Token: output.Token,
		})
	} else {
		return render.OK(c, "Successfully logged in.", AuthOutput{
			Token: output.Token,
		})
	}
}
