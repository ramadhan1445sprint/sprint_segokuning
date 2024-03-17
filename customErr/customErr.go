package customErr

type CustomError struct {
	message string
	status  int
}

func (e CustomError) Error() string {
	return e.message
}
func (e CustomError) Status() int {
	return e.status
}

func NewBadRequestError(message string) error {
	return CustomError{message: message, status: 400}
}

func NewUnauthorizedError(message string) error {
	return CustomError{message: message, status: 401}
}

func NewForbiddenError(message string) error {
	return CustomError{message: message, status: 403}
}

func NewNotFoundError(message string) error {
	return CustomError{message: message, status: 404}
}

func NewConflictError(message string) error {
	return CustomError{message: message, status: 409}
}

func NewInternalServerError(message string) error {
	return CustomError{message: message, status: 500}
}