package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Wrapped    error
}

func (h *HttpError) Error() string {
	return h.Wrapped.Error()
}

func NewError(statusCode int, err interface{}) error {
	var errVal error
	if ev, ok := err.(error); ok {
		errVal = ev
	} else if sv, ok := err.(string); ok {
		errVal = errors.New(sv)
	} else {
		errVal = errors.New(fmt.Sprint(err))
	}
	return &HttpError{
		StatusCode: statusCode,
		Wrapped:    errVal,
	}
}

func UnsupportedMethod(method string) error {
	return NewError(http.StatusBadRequest, "Unsupported method "+method)
}