package usero

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/lynx-go/lynx/hook"
	"github.com/weflux/fastapi/jsonapi"
	"github.com/weflux/fastapi/module/user/api"
)

type Module struct {
	api *api.API
}

func (mod *Module) Name() string {
	return "user-module"
}

func (mod *Module) Start(ctx context.Context) error {
	return nil
}

func (mod *Module) Stop(ctx context.Context) error {
	return nil
}

func (mod *Module) Status() (hook.Status, error) {
	return hook.StatusStarted, nil
}

func New() *Module {
	return &Module{}
}

func (mod *Module) Mount(router chi.Router) {
	router.Mount("/auth/login", jsonapi.H(mod.api.Login))
}
