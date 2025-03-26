package messaging

import (
	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	conn *nats.Conn
}

func NewNATSBroker(url string) (*NATSBroker, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSBroker{conn: nc}, nil
}

func (b *NATSBroker) Publish(topic string, message []byte) error {
	return b.conn.Publish(topic, message)
}

func (b *NATSBroker) Subscribe(topic string, handler func(msg []byte)) error {
	_, err := b.conn.Subscribe(topic, func(m *nats.Msg) {
		handler(m.Data)
	})
	return err
}

func (b *NATSBroker) Close() error {
	b.conn.Close()
	return nil
}
