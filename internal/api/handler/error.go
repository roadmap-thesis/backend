package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"github.com/roadmap-thesis/backend/pkg/render"
)

func (h *Handler) ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *apperrors.AppError
	var httpErr *echo.HTTPError
	code := http.StatusInternalServerError
	switch {
	case errors.As(err, &appErr):
		code = appErr.Code()
	case errors.As(err, &httpErr):
		if httpErr.Code == http.StatusNotFound {
			err = apperrors.ResourceNotFound("Path")
		}
		code = httpErr.Code
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
		if len(validationErrMsgs) > 0 {
			err = render.Error(c, code, "Validation failed.", validationErrMsgs)
		} else {
			err = render.Error(c, code, err.Error(), nil)
		}
	}
}

type validationErrMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getValidationErrMsg(err validator.FieldError) validationErrMsg {
	errMsg := validationErrMsg{
		Field: strings.ToLower(err.Field()),
	}

	errMsg.Message = map[string]string{
		"required": err.Field() + " is required.",
		"email":    "Must be a valid email address.",
		"min":      err.Field() + " must be at least " + err.Param() + " characters long.",
		"max":      err.Field() + " must not exceed " + err.Param() + " characters.",
	}[err.Tag()]

	return errMsg
}
