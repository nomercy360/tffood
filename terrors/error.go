package terrors

import (
	"net/http"
)

type Error struct {
	Code    int
	Err     error
	Message string
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NotFound(err error, message string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Err:     err,
		Message: message,
	}
}

func BadRequest(err error, message string) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Err:     err,
		Message: message,
	}
}

func Conflict(err error, message string) *Error {
	return &Error{
		Code:    http.StatusConflict,
		Err:     err,
		Message: message,
	}

}

func InternalServerError(err error, message string) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Err:     err,
		Message: message,
	}
}

func Unauthorized(err error, message string) *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Err:     err,
		Message: message,
	}
}
