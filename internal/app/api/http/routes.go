package httpapi

import (
	"github.com/gofiber/fiber/v2"

	"synthema/internal/app/health"
)

func RegisterRoutes(app *fiber.App, healthHandler *health.Handler) {
	v1 := app.Group("/api/v1")
	v1.Get("/health", healthHandler.Health)
	app.Get("/healthz", healthHandler.Healthz)
}
