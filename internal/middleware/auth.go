package middleware

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authctx "synthema/internal/context"
	appErrors "synthema/internal/errors"
	"synthema/internal/repositories"
)

func Auth(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository, cookieName string) fiber.Handler {
	if cookieName == "" {
		return func(c *fiber.Ctx) error {
			return appErrors.Internal(errors.New("auth cookie name is not configured"))
		}
	}
	return func(c *fiber.Ctx) error {
		sessionIDRaw := c.Cookies(cookieName)
		if sessionIDRaw == "" {
			return appErrors.MissingSessionCookie()
		}

		sessionID, err := uuid.Parse(sessionIDRaw)
		if err != nil {
			return appErrors.InvalidSessionCookie()
		}

		session, err := sessionRepo.FindByID(c.UserContext(), sessionID)
		if err != nil {
			return appErrors.Internal(err)
		}
		if session == nil {
			return appErrors.InvalidSessionCookie()
		}
		if session.RevokedAt != nil {
			return appErrors.SessionRevoked()
		}
		if !session.ExpiresAt.After(time.Now()) {
			return appErrors.SessionExpired()
		}

		user, err := userRepo.FindByID(c.UserContext(), session.UserID)
		if err != nil {
			return appErrors.Internal(err)
		}
		if user == nil {
			return appErrors.Unauthorized()
		}
		if !user.IsActive {
			return appErrors.UserInactive()
		}

		roles, err := userRepo.ListRolesByUserID(c.UserContext(), user.ID)
		if err != nil {
			return appErrors.Internal(err)
		}
		roleNames := make([]string, 0, len(roles))
		for _, r := range roles {
			if r.Name == "" {
				continue
			}
			roleNames = append(roleNames, r.Name)
		}

		ctx := authctx.WithAuthInfo(c.UserContext(), authctx.AuthInfo{
			UserID:    user.ID,
			Email:     user.Email,
			Roles:     roleNames,
			SessionID: session.ID,
		})
		c.SetUserContext(ctx)

		return c.Next()
	}
}
