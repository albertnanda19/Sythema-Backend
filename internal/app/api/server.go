package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	httpapi "synthema/internal/app/api/http"
	"synthema/internal/app/health"
	"synthema/internal/config"
	apihttp "synthema/internal/http"
	"synthema/internal/observability"
)

type Server struct {
	app           *fiber.App
	listenAddr    string
	shutdownGrace time.Duration
	logger        *observability.Logger
}

func NewServer(cfg config.Config, logger *observability.Logger, healthHandler *health.Handler) *Server {
	app := fiber.New(fiber.Config{ErrorHandler: apihttp.FiberErrorHandler(logger)})
	app.Use(recover.New())

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
