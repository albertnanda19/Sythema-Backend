package http

import (
	"github.com/gofiber/fiber/v2"
	"synthema/internal/repository"
)

// AuthMiddleware creates a middleware for authentication.
func AuthMiddleware(sessionRepo repository.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		session, err := sessionRepo.FindByID(c.Context(), sessionID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if session == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// Inject user ID into context for downstream handlers
		c.Locals("userID", session.UserID)

		return c.Next()
	}
}
