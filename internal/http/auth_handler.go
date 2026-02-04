package http

import (
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"synthema/internal/domain"
	"synthema/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService  service.AuthService
	cookieName   string
	cookieSecure bool
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService service.AuthService, cookieName string, cookieSecure bool) *AuthHandler {
	if cookieName == "" {
		cookieName = "session_id"
	}
	return &AuthHandler{authService: authService, cookieName: cookieName, cookieSecure: cookieSecure}
}

// LoginRequest represents the request body for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if len(req.Email) > 320 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userAgent := strings.TrimSpace(c.Get("User-Agent"))
	if len(userAgent) > 512 {
		userAgent = userAgent[:512]
	}
	ipAddress := strings.TrimSpace(c.IP())
	if len(ipAddress) > 64 {
		ipAddress = ipAddress[:64]
	}

	user, session, err := h.authService.Authenticate(c.Context(), req.Email, req.Password, service.SessionMeta{UserAgent: userAgent, IPAddress: ipAddress})
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		slog.Error("auth login failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
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
		SameSite: "Strict",
		MaxAge:   maxAgeSeconds,
	})

	roleNames := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		if r.Name == "" {
			continue
		}
		roleNames = append(roleNames, r.Name)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
		"roles": roleNames,
	})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionAny := c.Locals("authSession")
	if sessionAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	session, ok := sessionAny.(*domain.AuthSession)
	if !ok || session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.authService.RevokeSession(c.Context(), session.ID); err != nil {
		slog.Error("auth logout failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		Path:     "/",
		SameSite: "Strict",
		MaxAge:   -1,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout successful"})
}
