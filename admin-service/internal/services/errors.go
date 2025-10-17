package services

import (
	"net/http"
	"strings"
)

type serviceError struct {
	status  int
	message string
}

func (e *serviceError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *serviceError) StatusCode() int {
	if e == nil {
		return http.StatusInternalServerError
	}
	return e.status
}

func newServiceError(status int, message string) *serviceError {
	if strings.TrimSpace(message) == "" {
		message = http.StatusText(status)
		if message == "" {
			message = "unexpected error"
		}
	}

	return &serviceError{
		status:  status,
		message: message,
	}
}
