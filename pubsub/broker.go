package pubsub

import "gocloud.dev/pubsub"

type Broker struct {
	clients map[string]*pubsub.Topic
}
