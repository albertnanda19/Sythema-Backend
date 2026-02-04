package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	postgresCheck func(ctx context.Context) error
	redisCheck    func(ctx context.Context) error
	timeout       time.Duration
}

func NewHandler(postgresCheck func(ctx context.Context) error, redisCheck func(ctx context.Context) error, timeout time.Duration) *Handler {
	if timeout <= 0 {
		timeout = time.Second
	}
	return &Handler{postgresCheck: postgresCheck, redisCheck: redisCheck, timeout: timeout}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.timeout)
	defer cancel()

	services := map[string]string{}
	overall := "ok"
	statusCode := fiber.StatusOK

	if h.postgresCheck != nil {
		if err := h.postgresCheck(ctx); err != nil {
			services["postgres"] = "unhealthy"
			overall = "unhealthy"
			statusCode = fiber.StatusServiceUnavailable
		} else {
			services["postgres"] = "ok"
		}
	}
	if h.redisCheck != nil {
		if err := h.redisCheck(ctx); err != nil {
			services["redis"] = "unhealthy"
			overall = "unhealthy"
			statusCode = fiber.StatusServiceUnavailable
		} else {
			services["redis"] = "ok"
		}
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":   overall,
		"services": services,
	})
}
