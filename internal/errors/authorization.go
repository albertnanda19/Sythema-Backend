package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeForbidden = "auth.forbidden"
	MsgForbidden  = "Forbidden"
)

func Forbidden() Error {
	return New(CodeForbidden, fiber.StatusForbidden, MsgForbidden)
}
