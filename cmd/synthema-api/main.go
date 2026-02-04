package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"synthema/internal/bootstrap"
)

func main() {
	app, err := bootstrap.BootstrapAPI()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("bootstrap complete")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	listenAddr := ":" + strconv.Itoa(app.Config.API.Port)
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.App.Listen(listenAddr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.Config.ShutdownGracePeriod)
		defer cancel()
		_ = app.App.ShutdownWithContext(shutdownCtx)
		if app.Redis != nil {
			_ = app.Redis.Close()
		}
		if app.DB != nil {
			_ = app.DB.Close()
		}
		_ = os.Stdout.Sync()
		t := time.NewTimer(100 * time.Millisecond)
		<-t.C
		return
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	}
}
