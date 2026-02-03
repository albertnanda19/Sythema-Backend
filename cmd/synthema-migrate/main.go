package main

import (
	"context"
	"log"

	"os"
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

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := migration.ApplyAll(ctx, pool, "migrations"); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := migration.RollbackLast(ctx, pool, "migrations"); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command: %s", command)
	}
}
