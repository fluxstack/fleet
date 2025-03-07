package jsonapi

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HandlerFunc[I any, O any] func(ctx context.Context, in I) (O, error)

func H[I any, O any](h HandlerFunc[*I, *O]) echo.HandlerFunc {
	return func(ec echo.Context) error {
		in := new(I)
		if err := ec.Bind(in); err != nil {
			return err
		}
		out, err := h(ec.Request().Context(), in)
		if err != nil {
			return err
		}

		return ec.JSON(http.StatusOK, out)
	}
}
