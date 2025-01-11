package handler

import (
	"net/http"

	"github.com/HotPotatoC/roadmap_gen/api/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidData = echo.NewHTTPError(http.StatusBadRequest, "Invalid Data")
)

func (h *Handler) ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}

	var validationErrMsgs []validationErrMsg
	if validationErrs, isValidationErr := err.(validator.ValidationErrors); isValidationErr {
		code = http.StatusUnprocessableEntity

		validationErrMsgs = make([]validationErrMsg, 0)
		for _, err := range validationErrs {
			validationErrMsgs = append(validationErrMsgs, getValidationErrMsg(err))
		}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(code)
	} else {
		var res response.Response
		if len(validationErrMsgs) > 0 {
			res = response.NewErrorResponse("validation error", validationErrMsgs)
		} else {
			res = response.NewErrorResponse(err.Error(), nil)
		}
		err = c.JSON(code, res)
	}
}

type validationErrMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getValidationErrMsg(err validator.FieldError) validationErrMsg {
	errMsg := validationErrMsg{
		Field: err.Field(),
	}

	errMsg.Message = map[string]string{
		"required": err.Field() + " is required.",
		"email":    "Must be a valid email address.",
		"min":      err.Field() + " must be at least " + err.Param() + " characters long.",
		"max":      err.Field() + " must not exceed " + err.Param() + " characters.",
	}[err.Tag()]

	return errMsg
}
