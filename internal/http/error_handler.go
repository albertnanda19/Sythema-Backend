package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	domainerrors "synthema/internal/errors"
	"synthema/internal/observability"
)

func MapError(err error) (int, string) {
	if err == nil {
		return fiber.StatusInternalServerError, domainerrors.MsgInternal
	}

	var de domainerrors.Error
	if errors.As(err, &de) {
		return de.Status(), de.Message()
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		s := fe.Code
		switch s {
		case fiber.StatusNotFound:
			return s, domainerrors.MsgNotFound
		case fiber.StatusUnauthorized:
			return s, domainerrors.MsgUnauthorized
		case fiber.StatusForbidden:
			return s, domainerrors.MsgForbidden
		case fiber.StatusBadRequest:
			return s, domainerrors.MsgInvalidRequest
		default:
			if s >= 500 {
				return s, domainerrors.MsgInternal
			}
			if fe.Message != "" {
				return s, fe.Message
			}
			return s, domainerrors.MsgInvalidRequest
		}
	}

	return fiber.StatusInternalServerError, domainerrors.MsgInternal
}

func FiberErrorHandler(logger *observability.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status, msg := MapError(err)
		if status >= 500 && logger != nil {
			logger.ErrorContext(c.Context(), err.Error())
		}
		return Fail(c, status, msg)
	}
}
