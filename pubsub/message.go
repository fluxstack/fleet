package pubsub

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/lynx-go/x/encoding/json"
	"github.com/samber/lo"
	"gocloud.dev/pubsub"
)

func NewMessage(data any) *pubsub.Message {

	return &pubsub.Message{
		LoggableID: uuid.NewString(),
		Body:       lo.Must(json.Marshal(data)),
	}
}

func NewCloudEvent(source string, typ string, data any) *pubsub.Message {
	id := uuid.New().String()
	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetSource(source)
	event.SetType(typ)
	_ = event.SetData(cloudevents.ApplicationJSON, data)
	msg, _ := event.MarshalJSON()
	return &pubsub.Message{
		LoggableID: id,
		Body:       msg,
	}
}
