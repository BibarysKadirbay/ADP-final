package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/config"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewPublisher(cfg config.NATSConfig) (*Publisher, error) {
	conn, err := nats.Connect(cfg.URL, nats.Name("delivery-service"))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "DELIVERY_EVENTS",
		Subjects: []string{"delivery.*"},
		Storage:  nats.FileStorage,
	})
	return &Publisher{conn: conn, js: js}, nil
}

func (p *Publisher) Publish(ctx context.Context, subject string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ackCh := make(chan error, 1)
	go func() {
		_, err := p.js.Publish(subject, raw)
		ackCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ackCh:
		return err
	case <-time.After(3 * time.Second):
		return context.DeadlineExceeded
	}
}

func (p *Publisher) Close() error {
	if p.conn != nil {
		p.conn.Drain()
		p.conn.Close()
	}
	return nil
}
