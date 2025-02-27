package fastapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/lynx-go/lynx/hook"
)

type Module interface {
	hook.Hook
	Mount(router chi.Router)
}
