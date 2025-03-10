package pubsub

import (
	"context"
	"github.com/lynx-go/x/encoding/json"
	"gocloud.dev/pubsub"
)

type HandlerFunc func(ctx context.Context, msg *pubsub.Message) error

//func H[E any](h func(ctx context.Context, e E) error) HandlerFunc {
//	return func(ctx context.Context, msg *pubsub.Message) error {
//		var e E
//		if err := cloudevents.Unmarshal(msg, &e); err != nil {
//			return err
//		}
//		return h(ctx, e)
//	}
//}

//
//type Handler interface {
//	TopicID() TopicID
//	Unmarshal(data []byte, out any) error
//	HandlerFunc() HandlerFunc
//}
//
//func WrapH(topic TopicID, h HandlerFunc) Handler {
//	return &handler{topic: topic, handler: h}
//}
//
//type handler struct {
//	topic   TopicID
//	handler HandlerFunc
//}
//
//func (h *handler) TopicID() TopicID {
//	return h.topic
//}
//
//func (h *handler) Unmarshal(data []byte, out any) error {
//	return json.Unmarshal(data, out)
//}
//
//func (h *handler) HandlerFunc() HandlerFunc {
//	return func(ctx context.Context, event []byte) error {
//
//	}
//	return H(h.handler)
//}

type Marshaller interface {
	Marshal(data any) ([]byte, error)
	Unmarshal(data []byte, out any) error
}

type JSONMarshaller struct{}

func (ms *JSONMarshaller) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (ms *JSONMarshaller) Unmarshal(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

var _ Marshaller = new(JSONMarshaller)
