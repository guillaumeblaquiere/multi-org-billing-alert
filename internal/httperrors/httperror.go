package httperrors

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	httpCode int
	Err      error
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("error %v: %v", e.httpCode, e.Err)
}

func New(err error, statusCode int) error {
	return &HttpError{
		httpCode: statusCode,
		Err:      err,
	}
}

func GetHttpCode(err error) int {
	e, ok := err.(*HttpError)
	if !ok {
		return http.StatusInternalServerError
	}
	return e.httpCode
}
