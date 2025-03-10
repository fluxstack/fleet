package pubsub

import (
	"context"
	"github.com/lynx-go/x/json"
)

type HandlerFunc func(ctx context.Context, event []byte) error

func H[E any](h func(ctx context.Context, e E) error) HandlerFunc {
	return func(ctx context.Context, event []byte) error {
		var e E
		if err := json.Unmarshal(event, &e); err != nil {
			return err
		}
		return h(ctx, e)
	}
}
