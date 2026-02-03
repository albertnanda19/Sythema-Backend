package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"synthema/internal/seed"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		log.Fatal().Msg("refusing to run seeder in production environment")
	}

	databaseURL := os.Getenv("SYNTHEMA_POSTGRES_DSN")
	if databaseURL == "" {
		log.Fatal().Msg("SYNTHEMA_POSTGRES_DSN environment variable is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}

	fmt.Println("Successfully connected to the database.")

	seeder := seed.NewSeeder(pool)

	if err := seeder.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to run seeder")
	}

	log.Info().Msg("database seeding completed successfully")
}
