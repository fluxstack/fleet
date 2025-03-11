package pubsub

import (
	"context"
	"errors"
	"fmt"
	"github.com/lynx-go/lynx/hook"
	"github.com/lynx-go/x/log"
	"github.com/oklog/run"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/kafkapubsub"
	"gocloud.dev/pubsub/mempubsub"
	"sync"
	"time"
)

type TopicID string

func (t TopicID) String() string {
	return string(t)
}

type Broker interface {
	hook.Hook
	Topic(id TopicID) (*pubsub.Topic, error)
	Subscription(id TopicID) (*pubsub.Subscription, error)
	RegisterHandlerFunc(id TopicID, handler HandlerFunc)
}

var _ Broker = (*broker)(nil)

type TopicOption struct {
	Provider string            `json:"provider"`
	TopicID  string            `json:"topic_id"`
	Kafka    *KafkaTopicOption `json:"kafka,omitempty"`
}

type KafkaTopicOption struct {
	Servers      []string           `json:"servers"`
	Topic        string             `json:"topic"`
	Subscription *KafkaSubscription `json:"subscription"`
}

type KafkaSubscription struct {
	Group string `json:"group"`
}

type Option struct {
	Topics map[TopicID]TopicOption
}

func NewBroker(o Option) Broker {
	b := &broker{
		o:             o,
		mu:            sync.RWMutex{},
		g:             &run.Group{},
		topics:        make(map[TopicID]*pubsub.Topic),
		subscriptions: make(map[TopicID]*pubsub.Subscription),
		handlers:      map[TopicID]HandlerFunc{},
	}

	return b
}

type broker struct {
	o             Option
	mu            sync.RWMutex
	topics        map[TopicID]*pubsub.Topic
	handlers      map[TopicID]HandlerFunc
	subscriptions map[TopicID]*pubsub.Subscription
	g             *run.Group
	memTopics     map[TopicID]*pubsub.Topic
	memMu         sync.RWMutex
}

func (b *broker) Name() string {
	return "pubsub-broker"
}

func (b *broker) addTopic(id TopicID, topic *pubsub.Topic) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.topics[id] = topic
}

func (b *broker) addHandler(id TopicID, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[id] = handler
}

func (b *broker) getTopic(id TopicID) (*pubsub.Topic, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topic, ok := b.topics[id]
	return topic, ok
}

func (b *broker) getHandler(id TopicID) HandlerFunc {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.handlers[id]
}

func (b *broker) Start(ctx context.Context) error {
	for id := range b.o.Topics {
		topic, err := b.openTopic(id)
		if err != nil {
			return err
		}
		b.addTopic(id, topic)
	}

	for id := range b.handlers {
		sub, err := b.openSubscription(id)
		if err != nil {
			return err
		}
		h := b.getHandler(id)
		log.InfoContext(ctx, "starting subscription for topic", "topic_id", id)
		rctx, cancel := context.WithCancel(ctx)
		b.g.Add(func() error {
			var err error
			for {
				var msg *pubsub.Message
				msg, err = sub.Receive(rctx)
				if err != nil {
					break
				}
				if err = h(rctx, msg); err != nil {
					break
				}
				msg.Ack()
			}
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}, func(err error) {
			cancel()
			_ = sub.Shutdown(ctx)
		})
	}
	return b.g.Run()
}

func (b *broker) Stop(ctx context.Context) error {
	for id := range b.topics {
		topic, _ := b.getTopic(id)

		if err := topic.Shutdown(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
	return nil
}

func (b *broker) Status() (hook.Status, error) {
	return hook.StatusStarted, nil
}

func (b *broker) Subscription(id TopicID) (*pubsub.Subscription, error) {
	return b.openSubscription(id)
}

func (b *broker) RegisterHandlerFunc(id TopicID, handler HandlerFunc) {
	b.addHandler(id, handler)
}

func (b *broker) Topic(id TopicID) (*pubsub.Topic, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topic, ok := b.getTopic(id)
	if !ok {
		var err error
		topic, err = b.openTopic(id)
		if err != nil {
			return nil, err
		}
		b.addTopic(id, topic)
	}

	return topic, nil
}

func (b *broker) openTopic(id TopicID) (*pubsub.Topic, error) {
	o, ok := b.o.Topics[id]
	if !ok {
		return nil, fmt.Errorf("no such topic: %s", id)
	}

	switch o.Provider {
	case "kafka":
		return b.openKafkaTopic(o.Kafka)
	case "mem":
		return b.openMemTopic(id)
	default:
		return nil, fmt.Errorf("unknown provider: %s", o.Provider)
	}
}

func (b *broker) openKafkaTopic(o *KafkaTopicOption) (*pubsub.Topic, error) {
	config := kafkapubsub.MinimalConfig()
	topic, err := kafkapubsub.OpenTopic(o.Servers, config, o.Topic, nil)
	return topic, err
}

func (b *broker) getMemTopic(id TopicID) (*pubsub.Topic, bool) {
	b.memMu.RLock()
	defer b.memMu.RUnlock()

	topic, ok := b.memTopics[id]
	return topic, ok
}

func (b *broker) addMemTopic(id TopicID, topic *pubsub.Topic) {
	b.memMu.Lock()
	defer b.memMu.Unlock()
	b.memTopics[id] = topic
}

func (b *broker) openMemTopic(id TopicID) (*pubsub.Topic, error) {
	topic, ok := b.getMemTopic(id)
	if ok {
		return topic, nil
	}
	topic = mempubsub.NewTopic()
	b.addMemTopic(id, topic)
	return topic, nil
}

func (b *broker) openSubscription(topic TopicID) (*pubsub.Subscription, error) {
	o, ok := b.o.Topics[topic]
	if !ok {
		return nil, fmt.Errorf("no such topic: %s", topic)
	}
	switch o.Provider {
	case "kafka":
		return b.openKafkaSubscription(o.Kafka)
	case "mem":
		return b.openMemSubscription(topic)
	default:
		return nil, fmt.Errorf("unknown provider: %s", o.Provider)
	}
}

func (b *broker) openKafkaSubscription(o *KafkaTopicOption) (*pubsub.Subscription, error) {
	if o.Subscription == nil || o.Subscription.Group == "" {
		return nil, fmt.Errorf("no subscription.group specified")
	}
	config := kafkapubsub.MinimalConfig()
	sub, err := kafkapubsub.OpenSubscription(o.Servers, config, o.Subscription.Group, []string{o.Topic}, nil)
	return sub, err
}

func (b *broker) openMemSubscription(id TopicID) (*pubsub.Subscription, error) {
	topic, _ := b.openMemTopic(id)
	return mempubsub.NewSubscription(topic, 1*time.Minute), nil
}
