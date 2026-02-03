package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	httpapi "synthema/internal/app/api/http"
	"synthema/internal/app/health"
	"synthema/internal/config"
	"synthema/internal/observability"
)

type Server struct {
	app           *fiber.App
	listenAddr    string
	shutdownGrace time.Duration
	logger        *observability.Logger
}

func NewServer(cfg config.Config, logger *observability.Logger, healthHandler *health.Handler) *Server {
	app := fiber.New()

	httpapi.RegisterRoutes(app, healthHandler)

	return &Server{
		app:           app,
		listenAddr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		shutdownGrace: cfg.ShutdownGracePeriod,
		logger:        logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("api server listening on " + s.listenAddr)
		errCh <- s.app.Listen(s.listenAddr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownGrace)
		defer cancel()
		_ = s.app.ShutdownWithContext(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}
