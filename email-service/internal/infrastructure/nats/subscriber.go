package nats

import (
	"encoding/json"
	"log"

	"github.com/aitu/food-delivery/email-service/internal/usecase"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn *nats.Conn
	subs []*nats.Subscription
}

func NewSubscriber(url string, uc *usecase.EmailUsecase) (*Subscriber, error) {
	conn, err := nats.Connect(url, nats.Name("email-service"))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, _ = js.AddStream(&nats.StreamConfig{Name: "PAYMENT_EVENTS", Subjects: []string{"payment.>"}})

	s := &Subscriber{conn: conn}
	sub, err := js.QueueSubscribe("payment.completed", "email-service", func(msg *nats.Msg) {
		var event usecase.PaymentCompletedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			_ = msg.Nak()
			return
		}
		if err := uc.HandlePaymentCompleted(event); err != nil {
			log.Println("email send failed:", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	}, nats.Durable("email-payment-completed"), nats.ManualAck())
	if err != nil {
		conn.Close()
		return nil, err
	}
	s.subs = append(s.subs, sub)
	return s, nil
}

func (s *Subscriber) Close() {
	for _, sub := range s.subs {
		_ = sub.Drain()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
