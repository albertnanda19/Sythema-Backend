package httpapi

import (
	"github.com/gofiber/fiber/v2"

	"synthema/internal/app/health"
)

func RegisterRoutes(app *fiber.App, healthHandler *health.Handler) {
	app.Get("/healthz", healthHandler.Healthz)
}
