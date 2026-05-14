package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	client *goredis.Client
}

func New(addr string) (*Cache, error) {

	client := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		client: client,
	}, nil
}

func (c *Cache) Set(
	key string,
	value string,
	ttl time.Duration,
) error {

	return c.client.Set(
		ctx,
		key,
		value,
		ttl,
	).Err()
}

func (c *Cache) Get(
	key string,
) (string, error) {

	return c.client.Get(
		ctx,
		key,
	).Result()
}

func (c *Cache) Delete(
	key string,
) error {

	return c.client.Del(
		ctx,
		key,
	).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
