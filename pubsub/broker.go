package pubsub

import (
	"gocloud.dev/pubsub"
	"strings"
)

type Broker struct {
	clients map[string]*pubsub.Topic
}

type TopicURI string

func (uri TopicURI) String() string {
	return string(uri)
}

func (uri TopicURI) Namespace() string {
	uris := strings.Split(string(uri), ":")
	return uris[0]
}

func (uri TopicURI) Topic() string {
	uris := strings.Split(string(uri), ":")
	return uris[len(uris)-1]
}

func NewTopicURI(namespace string, topic string) TopicURI {
	return TopicURI(namespace + ":" + topic)
}
