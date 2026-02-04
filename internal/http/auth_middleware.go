package http

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"synthema/internal/repository"
)

func AuthMiddleware(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, cookieName string) fiber.Handler {
	if cookieName == "" {
		cookieName = "session_id"
	}
	return func(c *fiber.Ctx) error {
		sessionIDRaw := c.Cookies(cookieName)
		if sessionIDRaw == "" {
			return Fail(c, fiber.StatusUnauthorized, MsgUnauthorized)
		}

		sessionID, err := uuid.Parse(sessionIDRaw)
		if err != nil {
			return Fail(c, fiber.StatusUnauthorized, MsgUnauthorized)
		}

		session, err := sessionRepo.FindActiveByID(c.Context(), sessionID)
		if err != nil {
			return err
		}
		if session == nil {
			return Fail(c, fiber.StatusUnauthorized, MsgUnauthorized)
		}

		user, err := userRepo.FindByID(c.Context(), session.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return Fail(c, fiber.StatusUnauthorized, MsgUnauthorized)
		}
		if !user.IsActive {
			return Fail(c, fiber.StatusUnauthorized, MsgUnauthorized)
		}

		roles, err := userRepo.ListRolesByUserID(c.Context(), user.ID)
		if err != nil {
			return err
		}
		user.Roles = roles

		c.Locals("user", user)
		c.Locals("authSession", session)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}

var _ = sql.ErrNoRows
