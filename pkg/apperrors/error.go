package apperrors

type AppError struct {
	code    int
	message string
	err     error
}

func (e AppError) Code() int {
	return e.code
}

func (e AppError) Error() string {
	return e.message
}
