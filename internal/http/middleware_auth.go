package http

import (
	appErrors "synthema/internal/errors"
	"synthema/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware for authentication.
func AuthMiddlewareLegacy(sessionRepo repository.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_ = sessionRepo
		return appErrors.NotFound()
	}
}
