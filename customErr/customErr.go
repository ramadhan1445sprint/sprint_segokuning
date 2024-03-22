package customErr

type CustomError struct {
	Message    string  `json:"message"`
	StatusCode int	   `json:"status"`
}

func (e CustomError) Error() string {
	return e.Message
}
func (e CustomError) Status() int {
	return e.StatusCode
}

func NewBadRequestError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 400}
}

func NewUnauthorizedError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 401}
}

func NewForbiddenError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 403}
}

func NewNotFoundError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 404}
}

func NewConflictError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 409}
}

func NewInternalServerError(message string) CustomError {
	return CustomError{Message: message, StatusCode: 500}
}

