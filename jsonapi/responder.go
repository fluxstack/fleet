package jsonapi

import "github.com/lynx-go/fastapi/apierror"

type M map[string]interface{}

type Response[T any] struct {
	apierror.APIError
	Data T `json:"data"`
}

func NewResponse[T any](data T) *Response[T] {
	return &Response[T]{}
}
