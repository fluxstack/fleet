package pubsub

import "context"

type HandlerFunc func(ctx context.Context, event []byte) error
