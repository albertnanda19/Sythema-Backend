package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Ping(ctx context.Context, client *redis.Client, timeout time.Duration) error {
	if client == nil {
		return ErrRedisNotConfigured
	}
	if timeout <= 0 {
		timeout = time.Second
	}
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return client.Ping(pingCtx).Err()
}
