package server

import (
	"context"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewEchoInstance() *echo.Echo {
	instance := echo.New()
	instance.HideBanner = true
	instance.HidePort = true
	validator := validator.New()
	instance.Validator = &CustomValidator{validator: validator}
	return instance
}

func InjectEchoCtx(c echo.Context, key, val any) {
	ctx := context.WithValue(c.Request().Context(), key, val)
	c.SetRequest(c.Request().WithContext(ctx))
}
