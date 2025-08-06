// Filename: cache/redis.go
package cache

import (
	"context"
	"fmt"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%v", addr.Host, addr.Port)})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &RedisCache{Client: client}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.Client.Get(ctx, key).Bytes()
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}
