package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	authctx "synthema/internal/context"
	appErrors "synthema/internal/errors"
	"synthema/internal/http"
	"synthema/internal/service"
)

type LogoutHandler struct {
	authService    service.AuthService
	cookieName     string
	cookieSecure   bool
	cookieSameSite string
	cookieDomain   string
}

func NewLogoutHandler(authService service.AuthService, cookieName string, cookieSecure bool, cookieSameSite, cookieDomain string) *LogoutHandler {
	return &LogoutHandler{authService: authService, cookieName: cookieName, cookieSecure: cookieSecure, cookieSameSite: cookieSameSite, cookieDomain: cookieDomain}
}

func (h *LogoutHandler) Logout(c *fiber.Ctx) error {
	if h.cookieName == "" {
		return appErrors.Internal(errors.New("auth cookie name is not configured"))
	}

	sessionID, ok := authctx.SessionID(c.UserContext())
	if ok {
		if err := h.authService.Logout(c.UserContext(), sessionID); err != nil {
			return err
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		Path:     "/",
		SameSite: h.cookieSameSite,
		Domain:   h.cookieDomain,
		MaxAge:   0,
	})

	return http.Success(c, fiber.StatusOK, "Logged out successfully", nil)
}
