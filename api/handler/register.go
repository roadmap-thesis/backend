package handler

import (
	"github.com/HotPotatoC/roadmap_gen/api/res"
	"github.com/HotPotatoC/roadmap_gen/backend"
	"github.com/labstack/echo/v4"
)

type RegisterOutput struct {
	Token string `json:"token"`
}

func (h *Handler) Register(c echo.Context) error {
	var input backend.RegisterInput

	if err := c.Bind(&input); err != nil {
		return ErrInvalidData
	}

	if err := c.Validate(&input); err != nil {
		return err
	}

	token, err := h.backend.Register(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return res.Created(c, "Successfully registered", RegisterOutput{
		Token: token,
	})
}
