package jsonapi

import (
	"github.com/fluxstack/fleet/jsonapi/errors"
)

type M map[string]interface{}

type Response[T any] struct {
	errors.APIError
	Data T `json:"data"`
}

func OK[T any](data T) *Response[T] {
	return &Response[T]{
		APIError: errors.APIError{
			Message: "OK",
		},
		Data: data,
	}
}

func ERR[T any](err error, datas ...T) *Response[T] {
	var data T
	if len(datas) > 0 {
		data = datas[0]
	}
	var apiErr *errors.APIError
	switch e := err.(type) {
	case *errors.APIError:
		apiErr = e
	default:
		apiErr = &errors.APIError{
			Code:    -1,
			Message: e.Error(),
		}
	}
	return &Response[T]{
		APIError: *apiErr,
		Data:     data,
	}
}
