package http

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"synthema/internal/repository"
)

func AuthMiddleware(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionIDRaw := c.Cookies("session_id")
		if sessionIDRaw == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		sessionID, err := uuid.Parse(sessionIDRaw)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		session, err := sessionRepo.FindActiveByID(c.Context(), sessionID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if session == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		user, err := userRepo.FindByID(c.Context(), session.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		if !user.IsActive {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		roles, err := userRepo.ListRolesByUserID(c.Context(), user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		user.Roles = roles

		c.Locals("user", user)
		c.Locals("authSession", session)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}

var _ = sql.ErrNoRows
