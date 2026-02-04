package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeInvalidRequest = "validation.invalid_request"
	MsgInvalidRequest  = "Invalid request"
)

func InvalidRequest() Error {
	return New(CodeInvalidRequest, fiber.StatusBadRequest, MsgInvalidRequest)
}

func Validation(code string, message string) Error {
	if code == "" {
		code = CodeInvalidRequest
	}
	if message == "" {
		message = MsgInvalidRequest
	}
	return New(code, fiber.StatusUnprocessableEntity, message)
}
