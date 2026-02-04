package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	apihttp "synthema/internal/http"
	"synthema/internal/observability"
)

func RequestLogger(logger *observability.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		dur := time.Since(start)

		status := c.Response().StatusCode()
		if err != nil {
			mappedStatus, _ := apihttp.MapError(err)
			status = mappedStatus
		}

		msg := fmt.Sprintf(
			"request method=%s path=%s status=%d dur=%s ip=%s ua=%q origin=%q",
			c.Method(),
			c.Path(),
			status,
			dur,
			c.IP(),
			c.Get("User-Agent"),
			c.Get("Origin"),
		)

		if logger != nil {
			switch {
			case status >= 500:
				logger.ErrorContext(c.Context(), msg)
			case status >= 400:
				logger.WarnContext(c.Context(), msg)
			default:
				logger.InfoContext(c.Context(), msg)
			}
		}

		return err
	}
}
