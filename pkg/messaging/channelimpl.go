package messaging

import (
	"errors"
	"sync"
)

type ChannelBroker struct {
	subscribers map[string][]chan []byte
	mu          sync.RWMutex
	closed      bool
	closeCh     chan struct{}
}

func NewChannelBroker() Broker {
	return &ChannelBroker{
		subscribers: make(map[string][]chan []byte),
		closeCh:     make(chan struct{}),
	}
}

func (b *ChannelBroker) Publish(topic string, message []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return errors.New("broker is closed")
	}

	subs, exists := b.subscribers[topic]
	if !exists {
		return nil
	}

	for _, ch := range subs {
		select {
		case ch <- message:
		case <-b.closeCh:
			return errors.New("broker is closing")
		}
	}
	return nil
}

func (b *ChannelBroker) Subscribe(topic string, handler func(msg []byte)) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return errors.New("broker is closed")
	}

	ch := make(chan []byte, 10)
	b.subscribers[topic] = append(b.subscribers[topic], ch)

	go func() {
		for {
			select {
			case msg := <-ch:
				handler(msg)
			case <-b.closeCh:
				return
			}
		}
	}()

	return nil
}

func (b *ChannelBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	close(b.closeCh)

	for _, subs := range b.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}

	b.closed = true
	return nil
}
