package middleware

import (
	"context"
	"strings"

	"github.com/HotPotatoC/roadmap_gen/internal/auth"
	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization, ok := c.Request().Header["Authorization"]
		if !ok {
			return commonerrors.Unauthorized()
		}

		bearer := strings.Split(authorization[0], " ")
		if len(bearer) < 2 {
			return commonerrors.Unauthorized()
		}

		token := bearer[1]
		payload, err := auth.VerifyToken(token)
		if err != nil {
			return commonerrors.Unauthorized()
		}

		ctx := context.WithValue(c.Request().Context(), auth.AuthCtxKey, payload)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
