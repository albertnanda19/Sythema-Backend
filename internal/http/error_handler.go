package http

import (
	"errors"
	"strings"

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

func MapErrorData(err error) any {
	if err == nil {
		return fiber.Map{"detail": domainerrors.MsgInternal}
	}

	var de domainerrors.Error
	if errors.As(err, &de) {
		code := strings.TrimSpace(de.Code())
		detail := strings.TrimSpace(de.Message())
		switch code {
		case domainerrors.CodeAuthInvalidCredentials:
			detail = "Email or password is incorrect"
		case domainerrors.CodeAuthUnauthorized:
			detail = "Unauthorized: missing or invalid session"
		case domainerrors.CodeAuthMissingSession:
			detail = "Unauthorized: session cookie is missing"
		case domainerrors.CodeAuthInvalidSession:
			detail = "Unauthorized: session cookie is invalid"
		case domainerrors.CodeAuthSessionExpired:
			detail = "Unauthorized: session has expired"
		case domainerrors.CodeAuthSessionRevoked:
			detail = "Unauthorized: session has been revoked"
		case domainerrors.CodeAuthUserInactive:
			detail = "Unauthorized: user is inactive"
		}
		if code == "" {
			return fiber.Map{"detail": detail}
		}
		return fiber.Map{"code": code, "detail": detail}
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		status, msg := MapError(err)
		_ = status
		return fiber.Map{"detail": msg}
	}

	return fiber.Map{"detail": domainerrors.MsgInternal}
}

func FiberErrorHandler(logger *observability.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status, msg := MapError(err)
		if status >= 500 && logger != nil {
			logger.ErrorContext(c.Context(), err.Error())
		}
		return FailWithData(c, status, msg, MapErrorData(err))
	}
}
