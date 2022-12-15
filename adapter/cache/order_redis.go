package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type OrderRedis struct {
	rdb *redis.Client
}

func NewOrderRedis(rdb *redis.Client) *OrderRedis {
	return &OrderRedis{
		rdb: rdb,
	}
}

// SetOrder implements order.Cache
func (c *OrderRedis) SetOrder(orderId string) error {
	return c.rdb.Set(context.Background(), "order", orderId, time.Hour*10).Err()
}
