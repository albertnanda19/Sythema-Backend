package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env string

	API APIConfig

	Postgres PostgresConfig
	Redis    RedisConfig

	ShutdownGracePeriod time.Duration
}

type APIConfig struct {
	Host string
	Port int
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadFromEnv() (Config, error) {
	port := 8080
	if v := os.Getenv("SYNTHEMA_API_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, err
		}
		port = p
	}

	redisDB := 0
	if v := os.Getenv("SYNTHEMA_REDIS_DB"); v != "" {
		db, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, err
		}
		redisDB = db
	}

	grace := 10 * time.Second
	if v := os.Getenv("SYNTHEMA_SHUTDOWN_GRACE"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, err
		}
		grace = d
	}

	cfg := Config{
		Env: os.Getenv("SYNTHEMA_ENV"),
		API: APIConfig{
			Host: getenvDefault("SYNTHEMA_API_HOST", "0.0.0.0"),
			Port: port,
		},
		Postgres: PostgresConfig{
			DSN: os.Getenv("SYNTHEMA_POSTGRES_DSN"),
		},
		Redis: RedisConfig{
			Addr:     getenvDefault("SYNTHEMA_REDIS_ADDR", "127.0.0.1:6379"),
			Password: os.Getenv("SYNTHEMA_REDIS_PASSWORD"),
			DB:       redisDB,
		},
		ShutdownGracePeriod: grace,
	}

	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	if cfg.API.Port <= 0 {
		return Config{}, errors.New("invalid SYNTHEMA_API_PORT")
	}

	return cfg, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
