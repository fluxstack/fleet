package usero

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/lynx-go/lynx/hook"
	"github.com/weflux/fastapi/jsonapi"
	"github.com/weflux/fastapi/module/user/api"
	"github.com/weflux/fastapi/module/user/config"
	"github.com/weflux/fastapi/module/user/repoimpl"
	"github.com/weflux/fastapi/module/user/usecase"
	"github.com/weflux/fastapi/storage/database"
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

func New(dbClient *database.Client, cfg config.Config) *Module {
	users := repoimpl.NewUsers(dbClient)
	uc := usecase.NewAuth(users, cfg)
	userAPI := api.NewAPI(uc)
	return &Module{
		api: userAPI,
	}
}

func (mod *Module) Mount(router chi.Router) {
	router.Mount("/api/v1/auth/login", jsonapi.H(mod.api.Login))
}
