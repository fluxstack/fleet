package healthcheck

import (
	"github.com/weflux/fastapi/pkg/restful"
	"net/http"
)

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restful.R(w).OK(restful.M{"status": "OK"})
	}
}
