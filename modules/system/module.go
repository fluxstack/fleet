package system

import (
	"context"
	"database/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/weflux/fastapi"
	"github.com/weflux/fastapi/modules/system/api/healthcheck"
	"log/slog"
	"net/http"
)

type Option func(module *Module)

func WithLogger(logger *slog.Logger) Option {
	return func(mod *Module) {
		mod.logger = logger
	}
}

func WithDB(db *sql.DB) Option {
	return func(mod *Module) {
		mod.db = db
	}
}

func WithCache(rdb *redis.Client) Option {
	return func(mod *Module) {
		mod.rdb = rdb
	}
}

func WithEventBus(bus *cqrs.EventBus) Option {
	return func(mod *Module) {
		mod.bus = bus
	}
}

func (mod *Module) ensureDefaults() {
	if mod.logger == nil {
		mod.logger = slog.Default()
	}
}

func (mod *Module) initRoutes() {

	health := healthcheck.NewHandler()
	mod.mux.Handle("/health", health)
}

func (mod *Module) initBus() {
	//mod.bus.Publish()
}

func NewModule(opts ...Option) *Module {
	mod := &Module{}
	for _, opt := range opts {
		opt(mod)
	}
	mod.ensureDefaults()

	mod.mux = chi.NewRouter()
	mod.initRoutes()
	return mod
}

type Module struct {
	mux    *chi.Mux
	logger *slog.Logger
	db     *sql.DB
	rdb    *redis.Client
	bus    *cqrs.EventBus
}

func (mod *Module) Status() (int, error) {
	return 200, nil
}

func (mod *Module) Name() string {
	return "system-module"
}

func (mod *Module) Start(ctx context.Context) error {
	return nil
}

func (mod *Module) Stop(ctx context.Context) error {
	return nil
}

func (mod *Module) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mod.mux.ServeHTTP(writer, request)
}

var _ fastapi.Module = new(Module)
