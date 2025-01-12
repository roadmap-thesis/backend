package handler

import (
	"github.com/HotPotatoC/roadmap_gen/api/res"
	"github.com/HotPotatoC/roadmap_gen/backend"
	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/labstack/echo/v4"
)

type RegisterOutput struct {
	Token string `json:"token"`
}

func (h *Handler) Register(c echo.Context) error {
	var input backend.RegisterInput

	if err := c.Bind(&input); err != nil {
		return commonerrors.InvalidData()
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	output, err := h.backend.Register(c.Request().Context(), input)
	if err != nil {
		return err
	}

	if output.Created {
		return res.Created(c, "Successfully registered.", RegisterOutput{
			Token: output.Token,
		})
	} else {
		return res.OK(c, "Successfully logged in.", RegisterOutput{
			Token: output.Token,
		})
	}
}
