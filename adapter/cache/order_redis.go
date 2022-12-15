package cache

import (
	"context"
	"time"

	"github.com/Vractos/dolly/entity"
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
func (c *OrderRedis) SetOrder(o *entity.Order) error {
	return c.rdb.Set(context.Background(), o.MarketplaceID, o.Status.String(), time.Hour*10).Err()
}

// GetOrder implements order.Cache
func (c *OrderRedis) GetOrder(orderId string) (*entity.OrderStatus, error) {
	value, err := c.rdb.Get(context.Background(), orderId).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return (*entity.OrderStatus)(&value), nil
}
