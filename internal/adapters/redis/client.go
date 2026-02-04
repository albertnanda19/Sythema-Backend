package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"

	"synthema/internal/config"
)

func NewClient(cfg config.RedisConfig) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, ErrRedisNotConfigured
	}
	if cfg.DialTimeout < 0 || cfg.ReadTimeout < 0 || cfg.WriteTimeout < 0 {
		return nil, errors.New("redis timeouts must be >= 0")
	}
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}), nil
}

var ErrRedisNotConfigured = errors.New("redis not configured")
