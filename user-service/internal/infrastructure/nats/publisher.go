package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func NewPublisher(
	url string,
) (*Publisher, error) {

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn: nc,
	}, nil
}

func (p *Publisher) Publish(
	subject string,
	data interface{},
) error {

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, body)
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
