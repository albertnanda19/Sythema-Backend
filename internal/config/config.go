package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName     string
	Environment string
	LogLevel    string

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

	grace := 10 * time.Second
	if v := os.Getenv("SYNTHEMA_SHUTDOWN_GRACE"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, err
		}
		grace = d
	}

	cfg := Config{
		AppName:     getenvDefault("SYNTHEMA_APP_NAME", "synthema"),
		Environment: getenvDefault("SYNTHEMA_ENV", "dev"),
		LogLevel:    getenvDefault("SYNTHEMA_LOG_LEVEL", "info"),
		API: APIConfig{
			Host: getenvDefault("SYNTHEMA_API_HOST", "0.0.0.0"),
			Port: port,
		},
		Postgres:            PostgresConfig{},
		Redis:               RedisConfig{},
		ShutdownGracePeriod: grace,
	}

	return cfg, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
