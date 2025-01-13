package middleware

import (
	"strings"

	"github.com/HotPotatoC/roadmap_gen/pkg/auth"
	"github.com/HotPotatoC/roadmap_gen/pkg/commonerrors"
	"github.com/HotPotatoC/roadmap_gen/pkg/server"
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

		server.InjectEchoCtx(c, auth.AuthCtxKey, payload)
		return next(c)
	}
}
