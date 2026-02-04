package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	appErrors "synthema/internal/errors"
	apihttp "synthema/internal/http"
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

	if h.postgresCheck != nil {
		if err := h.postgresCheck(ctx); err != nil {
			services["postgres"] = "unhealthy"
			overall = "unhealthy"
		} else {
			services["postgres"] = "ok"
		}
	}
	if h.redisCheck != nil {
		if err := h.redisCheck(ctx); err != nil {
			services["redis"] = "unhealthy"
			overall = "unhealthy"
		} else {
			services["redis"] = "ok"
		}
	}

	return apihttp.Success(c, fiber.StatusOK, apihttp.MsgHealthOK, fiber.Map{
		"overall":  overall,
		"services": services,
	})
}

func (h *Handler) Healthz(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.timeout)
	defer cancel()

	healthy := true
	if h.postgresCheck != nil {
		if err := h.postgresCheck(ctx); err != nil {
			healthy = false
		}
	}
	if h.redisCheck != nil {
		if err := h.redisCheck(ctx); err != nil {
			healthy = false
		}
	}

	if !healthy {
		return appErrors.ServiceUnavailable(nil)
	}
	return apihttp.Success(c, fiber.StatusOK, apihttp.MsgHealthOK, nil)
}
