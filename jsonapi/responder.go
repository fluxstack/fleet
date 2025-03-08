package jsonapi

import (
	"github.com/fluxstack/flumix/jsonapi/errors"
)

type M map[string]interface{}

type Response[T any] struct {
	errors.APIError
	Data T `json:"data"`
}

func NewResponse[T any](data T) *Response[T] {
	return &Response[T]{}
}
