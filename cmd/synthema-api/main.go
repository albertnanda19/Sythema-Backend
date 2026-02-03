package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"synthema/internal/bootstrap"
	"synthema/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	app, cleanup, err := bootstrap.WireAPI(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup(context.Background())

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
