package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	httpapi "synthema/internal/app/api/http"
	"synthema/internal/app/health"
	"synthema/internal/config"
	apihttp "synthema/internal/http"
	"synthema/internal/middleware"
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
	app.Use(middleware.RequestLogger(logger))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "*",
		MaxAge:           3600,
	}))
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if c.Get("Origin") == "" {
			c.Set("Access-Control-Allow-Origin", "*")
		}
		return err
	})
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
