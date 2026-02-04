package http

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	authctx "synthema/internal/context"
	appErrors "synthema/internal/errors"
	"synthema/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService    service.AuthService
	cookieName     string
	cookieSecure   bool
	cookieSameSite string
	cookieDomain   string
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService service.AuthService, cookieName string, cookieSecure bool, cookieSameSite, cookieDomain string) *AuthHandler {
	return &AuthHandler{authService: authService, cookieName: cookieName, cookieSecure: cookieSecure, cookieSameSite: cookieSameSite, cookieDomain: cookieDomain}
}

// LoginRequest represents the request body for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	if h.cookieName == "" {
		return appErrors.Internal(errors.New("auth cookie name is not configured"))
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.InvalidRequest()
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		return appErrors.InvalidRequest()
	}
	if len(req.Email) > 320 {
		return appErrors.InvalidRequest()
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return appErrors.InvalidRequest()
	}

	userAgent := strings.TrimSpace(c.Get("User-Agent"))
	if len(userAgent) > 512 {
		userAgent = userAgent[:512]
	}
	ipAddress := strings.TrimSpace(c.IP())
	if len(ipAddress) > 64 {
		ipAddress = ipAddress[:64]
	}

	user, session, err := h.authService.Authenticate(c.UserContext(), req.Email, req.Password, service.SessionMeta{UserAgent: userAgent, IPAddress: ipAddress})
	if err != nil {
		return err
	}
	maxAgeSeconds := int(time.Until(session.ExpiresAt).Seconds())
	if maxAgeSeconds < 0 {
		maxAgeSeconds = 0
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cookieName,
		Value:    session.ID.String(),
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		Path:     "/",
		SameSite: h.cookieSameSite,
		Domain:   h.cookieDomain,
		MaxAge:   maxAgeSeconds,
	})

	roleNames := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		if r.Name == "" {
			continue
		}
		roleNames = append(roleNames, r.Name)
	}

	type loginResponse struct {
		ID    any      `json:"id"`
		Email string   `json:"email"`
		Roles []string `json:"roles"`
	}

	return Success(c, fiber.StatusOK, MsgLoginSuccessful, loginResponse{ID: user.ID, Email: user.Email, Roles: roleNames})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if h.cookieName == "" {
		return appErrors.Internal(errors.New("auth cookie name is not configured"))
	}

	sessionID, ok := authctx.SessionID(c.UserContext())
	if !ok {
		return appErrors.Internal(errors.New("auth context is missing"))
	}

	if err := h.authService.RevokeSession(c.UserContext(), sessionID); err != nil {
		return err
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
		MaxAge:   -1,
	})

	return Success(c, fiber.StatusOK, MsgLogoutSuccessful, nil)
}
