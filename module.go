package fastapi

import (
	"github.com/lynx-go/lynx/integration"
	"net/http"
)

type Module interface {
	integration.Integration
	ServeHTTP(http.ResponseWriter, *http.Request)
}
