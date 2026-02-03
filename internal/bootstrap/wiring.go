package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	redisv9 "github.com/redis/go-redis/v9"

	"synthema/internal/adapters/postgres"
	redisadapter "synthema/internal/adapters/redis"
	"synthema/internal/app/api"
	"synthema/internal/app/capture"
	"synthema/internal/app/diff"
	"synthema/internal/app/health"
	"synthema/internal/app/replay"
	"synthema/internal/app/transform"
	"synthema/internal/config"
	"synthema/internal/observability"
)

type Cleanup func(context.Context) error

func WireAPI(ctx context.Context, cfg config.Config) (*api.Server, Cleanup, error) {
	logger := observability.NewLogger()
	_ = observability.NewMetrics()

	pgPool, redisClient, cleanup, err := wireInfra(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	_ = pgPool
	_ = redisClient

	h := health.NewHandler()
	server := api.NewServer(cfg, logger, h)

	return server, cleanup, nil
}

func WireCapture(ctx context.Context, cfg config.Config) (*capture.Service, Cleanup, error) {
	logger := observability.NewLogger()
	_ = observability.NewMetrics()

	pgPool, redisClient, cleanup, err := wireInfra(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	_ = pgPool
	_ = redisClient

	svc := capture.NewService(logger, nil, nil)
	return svc, cleanup, nil
}

func WireWorker(ctx context.Context, cfg config.Config) (*replay.Service, Cleanup, error) {
	logger := observability.NewLogger()
	_ = observability.NewMetrics()

	pgPool, redisClient, cleanup, err := wireInfra(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	_ = pgPool
	_ = redisClient

	_ = transform.NewEngine(logger)
	_ = diff.NewService(logger, nil)

	svc := replay.NewService(logger, nil)
	return svc, cleanup, nil
}

func wireInfra(ctx context.Context, cfg config.Config) (*pgxpool.Pool, *redisv9.Client, Cleanup, error) {
	pgPool, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, nil, nil, err
	}

	redisClient := redisadapter.NewClient(cfg.Redis)

	cleanup := func(context.Context) error {
		if err := redisClient.Close(); err != nil {
			_ = postgres.Close(pgPool)
			return err
		}
		return postgres.Close(pgPool)
	}

	return pgPool, redisClient, cleanup, nil
}
