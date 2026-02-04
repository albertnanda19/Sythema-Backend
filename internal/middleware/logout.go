package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authctx "synthema/internal/context"
	appErrors "synthema/internal/errors"
)

func Logout(cookieName string) fiber.Handler {
	if cookieName == "" {
		return func(c *fiber.Ctx) error {
			return appErrors.Internal(errors.New("auth cookie name is not configured"))
		}
	}
	return func(c *fiber.Ctx) error {
		sessionIDRaw := c.Cookies(cookieName)
		if sessionIDRaw == "" {
			return c.Next()
		}

		sessionID, err := uuid.Parse(sessionIDRaw)
		if err != nil {
			return appErrors.Unauthorized()
		}

		ctx := authctx.WithAuthInfo(c.UserContext(), authctx.AuthInfo{SessionID: sessionID})
		c.SetUserContext(ctx)
		return c.Next()
	}
}
