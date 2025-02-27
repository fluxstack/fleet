package jsonapi

import (
	"context"
	"github.com/ggicci/httpin"
	"github.com/weflux/fastapi/errors"
	"net/http"
)

type HandlerFunc[I any, O any] func(ctx context.Context, in I) (O, error)

var bindingErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	R(w).JSON(http.StatusBadRequest, errors.APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Reason:  err.Error(),
	})
}

var handleErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	R(w).JSON(http.StatusInternalServerError, errors.APIError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Reason:  err.Error(),
	})
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func H[I any, O any](h HandlerFunc[I, O]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := httpin.Decode[I](r)
		if err != nil {
			bindingErrorHandler(w, r, err)
			return
		}
		out, err := h(r.Context(), *in)
		if err != nil {
			handleErrorHandler(w, r, err)
			return
		}
		R(w).OK(out)
	}
}
