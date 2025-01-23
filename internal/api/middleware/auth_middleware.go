package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"github.com/roadmap-thesis/backend/pkg/auth"
	"github.com/roadmap-thesis/backend/pkg/server"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization, ok := c.Request().Header["Authorization"]
		if !ok {
			return apperrors.Unauthorized()
		}

		bearer := strings.Split(authorization[0], " ")
		if len(bearer) < 2 {
			return apperrors.Unauthorized()
		}

		token := bearer[1]
		payload, err := auth.VerifyToken(token)
		if err != nil {
			return apperrors.Unauthorized()
		}

		server.InjectEchoCtx(c, auth.AuthCtxKey, payload)
		return next(c)
	}
}
