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

	svc, cleanup, err := bootstrap.WireCapture(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup(context.Background())

	if err := svc.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
