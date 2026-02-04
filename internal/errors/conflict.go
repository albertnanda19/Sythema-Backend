package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeConflict = "resource.conflict"
	MsgConflict  = "Conflict"
)

func Conflict(message string) Error {
	if message == "" {
		message = MsgConflict
	}
	return New(CodeConflict, fiber.StatusConflict, message)
}
