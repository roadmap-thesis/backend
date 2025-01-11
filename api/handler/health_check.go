package handler

import (
	"net/http"

	"github.com/HotPotatoC/roadmap_gen/api/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	res := response.NewSuccessResponse("OK", nil)
	return c.JSON(http.StatusOK, res)
}
