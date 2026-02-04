package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	domainerrors "synthema/internal/errors"
	"synthema/internal/observability"
	"synthema/internal/service"
)

func MapError(err error) (int, string) {
	if err == nil {
		return fiber.StatusInternalServerError, MsgInternalError
	}

	if errors.Is(err, service.ErrInvalidCredentials) {
		return fiber.StatusUnauthorized, MsgInvalidCreds
	}

	var de *domainerrors.DomainError
	if errors.As(err, &de) {
		return de.Status(), de.Message()
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		s := fe.Code
		switch s {
		case fiber.StatusNotFound:
			return s, MsgNotFound
		case fiber.StatusUnauthorized:
			return s, MsgUnauthorized
		case fiber.StatusForbidden:
			return s, MsgForbidden
		case fiber.StatusBadRequest:
			return s, MsgInvalidRequest
		default:
			if s >= 500 {
				return s, MsgInternalError
			}
			if fe.Message != "" {
				return s, fe.Message
			}
			return s, MsgInvalidRequest
		}
	}

	return fiber.StatusInternalServerError, MsgInternalError
}

func FiberErrorHandler(logger *observability.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status, msg := MapError(err)
		if status >= 500 && logger != nil {
			logger.ErrorContext(c.Context(), msg)
		}
		return Fail(c, status, msg)
	}
}
