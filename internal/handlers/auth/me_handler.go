package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	authctx "synthema/internal/context"
	appErrors "synthema/internal/errors"
	"synthema/internal/http"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

func (h *MeHandler) Me(c *fiber.Ctx) error {
	info, ok := authctx.GetAuthInfo(c.UserContext())
	if !ok {
		return appErrors.Internal(errors.New("auth context is missing"))
	}

	type meResponse struct {
		ID    any      `json:"id"`
		Email string   `json:"email"`
		Roles []string `json:"roles"`
	}

	roles := make([]string, len(info.Roles))
	copy(roles, info.Roles)
	return http.Success(c, fiber.StatusOK, "Authenticated user", meResponse{ID: info.UserID, Email: info.Email, Roles: roles})
}
