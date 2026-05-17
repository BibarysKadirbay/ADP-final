package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aitu/food-delivery/payment-service/internal/usecase"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	js nats.JetStreamContext
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url, nats.Name("payment-service"))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, _ = js.AddStream(&nats.StreamConfig{Name: "ORDER_EVENTS", Subjects: []string{"order.>"}})
	_, _ = js.AddStream(&nats.StreamConfig{Name: "PAYMENT_EVENTS", Subjects: []string{"payment.>"}})
	return &Publisher{js: js}, nil
}

func (p *Publisher) Publish(subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(subject, data)
	return err
}

type Subscriber struct {
	subs []*nats.Subscription
}

func NewSubscriber(url string, uc *usecase.PaymentUsecase) (*Subscriber, error) {
	conn, err := nats.Connect(url, nats.Name("payment-service-subscriber"))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, _ = js.AddStream(&nats.StreamConfig{Name: "ORDER_EVENTS", Subjects: []string{"order.>"}})

	s := &Subscriber{}
	sub, err := js.QueueSubscribe("order.created", "payment-service", func(msg *nats.Msg) {
		var event usecase.OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			_ = msg.Nak()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := uc.ProcessOrderCreated(ctx, event); err != nil {
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	}, nats.Durable("payment-order-created"), nats.ManualAck())
	if err != nil {
		return nil, err
	}
	s.subs = append(s.subs, sub)
	return s, nil
}

func (s *Subscriber) Close() {
	for _, sub := range s.subs {
		_ = sub.Drain()
	}
}
