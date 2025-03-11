package cloudevents

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/lynx-go/x/encoding/json"
	"github.com/samber/lo"
	"gocloud.dev/pubsub"
)

func New(source string, typ string, data any) *pubsub.Message {
	id := uuid.New().String()
	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetSource(source)
	event.SetType(typ)
	lo.Must0(event.SetData(cloudevents.ApplicationJSON, data))
	return &pubsub.Message{
		Body: lo.Must1(event.MarshalJSON()),
	}
}

func Unmarshal(message *pubsub.Message) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DataAs(message *pubsub.Message, data any) error {
	event, err := Unmarshal(message)
	if err != nil {
		return err
	}
	return event.DataAs(&data)
}
