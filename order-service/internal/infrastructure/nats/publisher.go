// order-service/internal/infrastructure/nats/publisher.go

package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

type OrderCreatedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Amount  int64  `json:"amount"`
}

func NewPublisher(url string) (*Publisher, error) {

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
