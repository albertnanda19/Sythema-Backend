package config

import (
	"fmt"
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
	Host     string
	Port     int
	Addr     string
	Password string
	DB       int

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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

	redisHost := os.Getenv("REDIS_HOST")
	redisPortRaw := os.Getenv("REDIS_PORT")
	var redisCfg RedisConfig
	if redisHost != "" || redisPortRaw != "" {
		if redisHost == "" || redisPortRaw == "" {
			return Config{}, fmt.Errorf("redis requires both REDIS_HOST and REDIS_PORT")
		}
		p, err := strconv.Atoi(redisPortRaw)
		if err != nil {
			return Config{}, err
		}
		db := 0
		if v := os.Getenv("REDIS_DB"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return Config{}, err
			}
			db = n
		}

		dialTimeout := 5 * time.Second
		if v := os.Getenv("REDIS_DIAL_TIMEOUT"); v != "" {
			d, err := time.ParseDuration(v)
			if err != nil {
				return Config{}, err
			}
			dialTimeout = d
		}
		readTimeout := 1 * time.Second
		if v := os.Getenv("REDIS_READ_TIMEOUT"); v != "" {
			d, err := time.ParseDuration(v)
			if err != nil {
				return Config{}, err
			}
			readTimeout = d
		}
		writeTimeout := 1 * time.Second
		if v := os.Getenv("REDIS_WRITE_TIMEOUT"); v != "" {
			d, err := time.ParseDuration(v)
			if err != nil {
				return Config{}, err
			}
			writeTimeout = d
		}

		redisCfg = RedisConfig{
			Host:         redisHost,
			Port:         p,
			Addr:         fmt.Sprintf("%s:%d", redisHost, p),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           db,
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		}
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
		Redis:               redisCfg,
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
