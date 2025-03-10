package cloudevents

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/lynx-go/x/encoding/json"
	"gocloud.dev/pubsub"
)

func New(source string, typ string, data any) *pubsub.Message {
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

func Unmarshal(message *pubsub.Message, data any) error {
	event := cloudevents.NewEvent()
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}
	return event.DataAs(&data)
}
