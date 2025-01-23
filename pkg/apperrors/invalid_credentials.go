package apperrors

import (
	"errors"
	"net/http"
)

func InvalidCredentials() *AppError {
	return &AppError{
		code:    http.StatusUnauthorized,
		message: "Invalid Credentials.",
		err:     errors.New("invalid credentials"),
	}
}
