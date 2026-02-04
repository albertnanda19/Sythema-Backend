package http

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	appErrors "synthema/internal/errors"
	"synthema/internal/repository"
)

func AuthMiddleware(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, cookieName string) fiber.Handler {
	if cookieName == "" {
		cookieName = "session_id"
	}
	return func(c *fiber.Ctx) error {
		sessionIDRaw := c.Cookies(cookieName)
		if sessionIDRaw == "" {
			return appErrors.Unauthorized()
		}

		sessionID, err := uuid.Parse(sessionIDRaw)
		if err != nil {
			return appErrors.Unauthorized()
		}

		session, err := sessionRepo.FindActiveByID(c.Context(), sessionID)
		if err != nil {
			return appErrors.Internal(err)
		}
		if session == nil {
			return appErrors.Unauthorized()
		}

		user, err := userRepo.FindByID(c.Context(), session.UserID)
		if err != nil {
			return appErrors.Internal(err)
		}
		if user == nil {
			return appErrors.Unauthorized()
		}
		if !user.IsActive {
			return appErrors.Unauthorized()
		}

		roles, err := userRepo.ListRolesByUserID(c.Context(), user.ID)
		if err != nil {
			return appErrors.Internal(err)
		}
		user.Roles = roles

		c.Locals("user", user)
		c.Locals("authSession", session)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}

var _ = sql.ErrNoRows
