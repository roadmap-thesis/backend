package commonerrors

import "net/http"

func NotFound(resource string) *AppError {
	return &AppError{
		code:    http.StatusNotFound,
		message: resource + " not found.",
	}
}
