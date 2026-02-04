package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeInternal          = "internal"
	CodeServiceUnavailable = "service.unavailable"
	MsgInternal           = "Internal server error"
	MsgServiceUnavailable = "Service unavailable"
)

func Internal(err error) Error {
	return Wrap(CodeInternal, fiber.StatusInternalServerError, MsgInternal, err)
}

func ServiceUnavailable(err error) Error {
	return Wrap(CodeServiceUnavailable, fiber.StatusServiceUnavailable, MsgServiceUnavailable, err)
}
