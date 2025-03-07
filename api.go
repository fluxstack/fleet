package fastapi

import (
	"context"
	"github.com/lynx-go/lynx/hook"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Option func(api *API)

//func WithRouter(r chi.Router) Option {
//	return func(api *API) {
//		api.router = r
//	}
//}

func WithModules(modules ...Module) Option {
	return func(api *API) {
		api.mods = append(api.mods, modules...)
	}
}

type API struct {
	router http.Handler
	mods   []Module
}

func (api *API) Name() string {
	return "fastapi"
}

func (api *API) Start(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)
	for _, mod := range api.mods {
		mod := mod
		wg.Go(func() error {
			return mod.Start(ctx)
		})
	}
	return wg.Wait()
}

func (api *API) Stop(ctx context.Context) error {
	var errs error
	for _, mod := range api.mods {
		mod := mod
		err := mod.Stop(ctx)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}

func (api *API) Status() (hook.Status, error) {
	return hook.StatusStarted, nil
}

func NewAPI(opts ...Option) *API {
	api := &API{}
	for _, opt := range opts {
		opt(api)
	}
	api.ensureDefaults()

	//for _, mod := range api.mods {
	//mod := mod
	//mod.Mount(api.router)
	//}
	return api
}

func (api *API) ensureDefaults() {
	//if api.router == nil {
	//	api.router = chi.NewRouter()
	//}
}
func (api *API) Router() http.Handler {
	return api.router
}

var _ hook.Hook = new(API)
