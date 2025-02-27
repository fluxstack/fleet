package main

import (
	"context"
	"emperror.dev/emperror"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx-contrib/transport/http"
	"github.com/lynx-go/lynx/hook"
	"github.com/lynx-go/x/log"
	"github.com/spf13/viper"
	"github.com/weflux/fastapi"
	"github.com/weflux/fastapi/errors"
	"github.com/weflux/fastapi/jsonapi"
	usermod "github.com/weflux/fastapi/module/user"
	userconfig "github.com/weflux/fastapi/module/user/config"
	"github.com/weflux/fastapi/storage/database"
	gohttp "net/http"
)

type Config struct {
	userconfig.Config
	Data database.Config `json:"data"`
}

func main() {
	jsonapi.SetOptions(jsonapi.WithHandleErrorHandler(func(w gohttp.ResponseWriter, r *gohttp.Request, err error) {
		jsonapi.R(w).OK(errors.APIError{
			Code:     -1,
			Message:  err.Error(),
			Reason:   err.Error(),
			Metadata: nil,
		})
	}))

	serverSetup := func(ctx context.Context, hooks *hook.Hooks, o Option, args []string) (lynx.RunFunc, error) {
		logger := log.FromContext(ctx)
		logger.Info("option parsed", "option", o, "config", o.Config)
		viper.SetConfigFile(o.Config)
		emperror.Panic(viper.ReadInConfig())
		cfg := Config{}
		emperror.Panic(viper.Unmarshal(&cfg))
		dbClient, err := database.NewClient(ctx, cfg.Data)
		emperror.Panic(err)
		userMod := usermod.New(dbClient, cfg.Config)
		api := fastapi.NewAPI(fastapi.WithModules(userMod))
		hooks.Register(api)

		hs := http.New(http.WithHandler(api.Router()), http.WithAddr(o.Addr))
		hooks.Register(hs)
		hooks.OnStart(func(ctx context.Context) error {
			log.InfoContext(ctx, "onstart")
			return nil
		})
		hooks.OnStop(func(ctx context.Context) error {
			log.InfoContext(ctx, "onstop")
			return nil
		})
		return lynx.RunWaitSignal(), nil
	}

	cli := lynx.NewCLI[Option](
		lynx.CMD[Option](
			lynx.New(
				lynx.WithName[Option]("fastapi"),
				lynx.WithVersion[Option]("0.0.1"),
				lynx.WithSetup[Option](serverSetup),
			),
		),
	)

	cli.Run()
}

type Option struct {
	Addr   string `json:"addr"`
	Config string `json:"config"`
}
