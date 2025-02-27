package jsonapi

import (
	"fmt"
	"github.com/lynx-go/x/json"
	"log/slog"
	"net/http"
)

type M map[string]interface{}

func R(w http.ResponseWriter) *Responder {
	return &Responder{ResponseWriter: w}
}

type Responder struct {
	http.ResponseWriter
}

func (w *Responder) OK(body interface{}) {
	w.JSON(http.StatusOK, body)
}

func (w *Responder) JSON(status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, json.SafeMarshalString(body)); err != nil {
		slog.Warn("could not write output", "error", err)
	}
}

func (w *Responder) TEXT(status int, body interface{}) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, body); err != nil {
		slog.Warn("could not write output", "error", err)
	}
}
