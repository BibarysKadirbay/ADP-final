package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/config"
	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Subscriber struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	subs []*nats.Subscription
	log  *zap.Logger
}

func NewSubscriber(cfg config.NATSConfig, uc *usecase.DeliveryUsecase, log *zap.Logger) (*Subscriber, error) {
	conn, err := nats.Connect(cfg.URL, nats.Name("delivery-service-subscriber"))
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
		Subjects: []string{"order.*"},
		Storage:  nats.FileStorage,
	})
	s := &Subscriber{conn: conn, js: js, log: log}
	if err := s.subscribe(uc); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

func (s *Subscriber) subscribe(uc *usecase.DeliveryUsecase) error {
	confirmed, err := s.js.QueueSubscribe("order.confirmed", "delivery-service", func(msg *nats.Msg) {
		var event usecase.OrderConfirmedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			_ = msg.Nak()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := uc.HandleOrderConfirmed(ctx, event); err != nil {
			s.log.Warn("order.confirmed handling failed", zap.Error(err))
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	}, nats.Durable("delivery-order-confirmed"), nats.ManualAck())
	if err != nil {
		return err
	}
	cancelled, err := s.js.QueueSubscribe("order.cancelled", "delivery-service", func(msg *nats.Msg) {
		var event struct {
			OrderID uuid.UUID `json:"order_id"`
		}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			_ = msg.Nak()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := uc.HandleOrderCancelled(ctx, event.OrderID); err != nil {
			s.log.Warn("order.cancelled handling failed", zap.Error(err))
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	}, nats.Durable("delivery-order-cancelled"), nats.ManualAck())
	if err != nil {
		return err
	}
	s.subs = append(s.subs, confirmed, cancelled)
	return nil
}

func (s *Subscriber) Close() error {
	for _, sub := range s.subs {
		_ = sub.Drain()
	}
	if s.conn != nil {
		s.conn.Drain()
		s.conn.Close()
	}
	return nil
}
