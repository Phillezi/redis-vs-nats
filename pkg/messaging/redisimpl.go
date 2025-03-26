package messaging

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client *redis.Client
}

func NewRedisBroker(addr string) *RedisBroker {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisBroker{client: client}
}

func (b *RedisBroker) Publish(topic string, message []byte) error {
	return b.client.Publish(context.Background(), topic, message).Err()
}

func (b *RedisBroker) Subscribe(topic string, handler func(msg []byte)) error {
	pubsub := b.client.Subscribe(context.Background(), topic)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			handler([]byte(msg.Payload))
		}
	}()
	return nil
}

func (b *RedisBroker) Close() error {
	return b.client.Close()
}
