package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"synthema/internal/config"
)

func Connect(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	if cfg.DSN == "" {
		return nil, nil
	}

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Close(pool *pgxpool.Pool) error {
	if pool == nil {
		return nil
	}
	pool.Close()
	return nil
}

var ErrPostgresNotConfigured = errors.New("postgres not configured")
