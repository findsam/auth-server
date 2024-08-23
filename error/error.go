package ge

import (
	"net/http"
)

type CustomError struct {
	Message    string
	StatusCode int
}

func (e *CustomError) Error() string {
	return e.Message
}

func New(message string, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	Internal             = New("Internal Server Error", http.StatusInternalServerError)
	NotFound             = New("Resource Not Found", http.StatusNotFound)
	BadRequest           = New("Bad Request", http.StatusBadRequest)
	IncorrectCredentials = New("No user matches those credentials", http.StatusBadRequest)
	EmailExists          = New("A user with that email already exists", http.StatusBadRequest)
	Unauthorized         = New("Unauthorized request", http.StatusUnauthorized)
	UserNotFound         = New("No user was found", http.StatusNoContent)
)
