package main

import (
	"context"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx-contrib/transport/http"
	"github.com/lynx-go/lynx/hook"
	"github.com/lynx-go/x/log"
	"github.com/weflux/fastapi"
	user "github.com/weflux/fastapi/module/user"
)

func main() {
	serverSetup := func(ctx context.Context, hooks *hook.Hooks, o Option, args []string) (lynx.RunFunc, error) {
		logger := log.FromContext(ctx)
		cfg := o.Config
		logger.Info("option parsed", "option", o, "config", cfg)
		userMod := user.New()
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
