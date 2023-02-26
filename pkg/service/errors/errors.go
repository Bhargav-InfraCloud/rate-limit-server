package errors

import (
	"net/http"
)

type Error struct {
	Message     string   `json:"message"`
	Code        string   `json:"service_code"`
	RelatedLogs []string `json:"logs"`

	// For internal use
	statusCode int
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) StatusCode() int {
	return e.statusCode
}

var (
	RateLimitedError = &Error{
		Message:    "rate limit reached for the specific ID",
		Code:       "1001",
		statusCode: http.StatusTooManyRequests,
	}
	InternalServerError = &Error{
		Message:    "internal server error",
		Code:       "1002",
		statusCode: http.StatusInternalServerError,
	}
	InvalidInputsError = &Error{
		Message:    "missing or invalid inputs error",
		Code:       "1003",
		statusCode: http.StatusBadRequest,
	}
)

func CopyOfServiceError(err *Error) *Error {
	servErr := *err
	return &servErr
}
