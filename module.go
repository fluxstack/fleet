package fastapi

import (
	"github.com/lynx-go/lynx/hook"
)

type Module interface {
	hook.Hook
	//Mount(router chi.Router)
}
