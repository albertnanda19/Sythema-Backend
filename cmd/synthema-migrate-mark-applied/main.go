package main

import (
	"context"
	"log"
	"os"

	"synthema/internal/adapters/postgres"
	"synthema/internal/config"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run mark_applied.go <version>")
	}
	version := os.Args[1]

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

	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_versions (version TEXT NOT NULL PRIMARY KEY)`); err != nil {
		log.Fatal(err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO schema_versions (version) VALUES ($1) ON CONFLICT (version) DO NOTHING`, version); err != nil {
		log.Fatal(err)
	}

	log.Printf("marked version %s as applied", version)
}
