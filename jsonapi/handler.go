package jsonapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type Options struct {
	bindingErrorHandler ErrorHandler
	handleErrorHandler  ErrorHandler
}

type Option func(*Options)

func WithBindingErrorHandler(bindingErrorHandler ErrorHandler) Option {
	return func(o *Options) {
		o.bindingErrorHandler = bindingErrorHandler
	}
}

var defaultOptions = &Options{}

func SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(defaultOptions)
	}
}

func WithHandleErrorHandler(handleErrorHandler ErrorHandler) Option {
	return func(o *Options) {
		o.handleErrorHandler = handleErrorHandler
	}
}

type HandlerFunc[I any, O any] func(ctx context.Context, in I) (O, error)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func H[I any, O any](h HandlerFunc[*I, *O]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := new(I)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			defaultOptions.bindingErrorHandler(w, r, err)
			return
		}
		out, err := h(r.Context(), in)
		if err != nil {
			defaultOptions.handleErrorHandler(w, r, err)
			return
		}
		R(w).OK(out)
	}
}
