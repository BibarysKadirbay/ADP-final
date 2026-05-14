package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/config"
	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/metrics"
	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client  *goredis.Client
	metrics *metrics.Metrics
}

func New(cfg config.RedisConfig, m *metrics.Metrics) *Cache {
	return &Cache{
		client:  goredis.NewClient(&goredis.Options{Addr: cfg.Addr, Password: cfg.Password, DB: cfg.DB}),
		metrics: m,
	}
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) Get(ctx context.Context, key string, dest any) (bool, error) {
	raw, err := c.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		c.metrics.CacheEvents.WithLabelValues("miss").Inc()
		return false, nil
	}
	if err != nil {
		return false, err
	}
	c.metrics.CacheEvents.WithLabelValues("hit").Inc()
	return true, json.Unmarshal(raw, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, raw, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		if next == 0 {
			return nil
		}
		cursor = next
	}
}
