package health

import "github.com/gofiber/fiber/v2"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Healthz(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
