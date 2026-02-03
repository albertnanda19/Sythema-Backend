package main

import (
	"context"
	"log"

	"synthema/internal/adapters/postgres"
	"synthema/internal/config"
	"synthema/internal/migration"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Postgres.DSN == "" {
		log.Fatal(postgres.ErrPostgresNotConfigured)
	}

	pool, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = postgres.Close(pool)
	}()

	if err := migration.ApplyAll(ctx, pool, "migrations"); err != nil {
		log.Fatal(err)
	}
}
