package http

import (
	"synthema/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware for authentication.
func AuthMiddlewareLegacy(sessionRepo repository.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_ = sessionRepo
		return Fail(c, fiber.StatusNotImplemented, MsgNotFound)
	}
}
