package restidl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Request struct {
	IdToken  string
	ClientId string
	ReqMsg   interface{}
}

func (rr *Request) ReadTo(reqStruct interface{}) {
	jsonText, err := json.Marshal(rr.ReqMsg)
	if err != nil {
		panic(NewError(http.StatusInternalServerError, err))
	}
	err = json.Unmarshal(jsonText, reqStruct)
	if err != nil {
		panic(NewError(http.StatusBadRequest, err))
	}
}

type httpError struct {
	StatusCode int
	Wrapped    error
}

func (h *httpError) Error() string {
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
	return &httpError{
		StatusCode: statusCode,
		Wrapped:    errVal,
	}
}
