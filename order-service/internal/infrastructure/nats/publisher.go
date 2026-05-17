package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url, nats.Name("order-service"))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "ORDER_EVENTS",
		Subjects: []string{"order.>"},
	})
	return &Publisher{conn: conn, js: js}, nil
}

func (p *Publisher) Publish(subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(subject, data)
	return err
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
