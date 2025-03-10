package pubsub

import (
	"github.com/lynx-go/x/encoding/json"
	"github.com/samber/lo"
	"gocloud.dev/pubsub"
)

func NewMessage(data any) *pubsub.Message {
	return &pubsub.Message{
		Body: lo.Must(json.Marshal(data)),
	}
}
