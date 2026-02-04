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

	API  APIConfig
	Auth AuthConfig

	Postgres PostgresConfig
	Redis    RedisConfig

	ShutdownGracePeriod time.Duration
}

type APIConfig struct {
	Host string
	Port int
}

type AuthConfig struct {
	SessionTTL   time.Duration
	CookieName   string
	CookieSecure bool
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
	if err := loadDotEnvIfPresent(".env"); err != nil {
		return Config{}, err
	}

	port := 8080
	if v := os.Getenv("SYNTHEMA_API_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, err
		}
		port = p
	}

	dsn := os.Getenv("SYNTHEMA_POSTGRES_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	grace := 10 * time.Second
	if v := os.Getenv("SYNTHEMA_SHUTDOWN_GRACE"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, err
		}
		grace = d
	}

	sessionTTL := 7 * 24 * time.Hour
	if v := os.Getenv("SYNTHEMA_AUTH_SESSION_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, err
		}
		sessionTTL = d
	}

	cookieSecure := true
	if v := os.Getenv("SYNTHEMA_AUTH_COOKIE_SECURE"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, err
		}
		cookieSecure = b
	}

	cfg := Config{
		AppName:     getenvDefault("SYNTHEMA_APP_NAME", "synthema"),
		Environment: getenvDefault("SYNTHEMA_ENV", "dev"),
		LogLevel:    getenvDefault("SYNTHEMA_LOG_LEVEL", "info"),
		API: APIConfig{
			Host: getenvDefault("SYNTHEMA_API_HOST", "0.0.0.0"),
			Port: port,
		},
		Auth: AuthConfig{
			SessionTTL:   sessionTTL,
			CookieName:   getenvDefault("SYNTHEMA_AUTH_COOKIE_NAME", "session_id"),
			CookieSecure: cookieSecure,
		},
		Postgres:            PostgresConfig{DSN: dsn},
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
