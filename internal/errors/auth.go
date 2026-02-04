package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeAuthUnauthorized       = "auth.unauthorized"
	CodeAuthInvalidCredentials = "auth.invalid_credentials"
	MsgUnauthorized            = "Unauthorized"
	MsgInvalidCredentials      = "Invalid credentials"
)

func Unauthorized() Error {
	return New(CodeAuthUnauthorized, fiber.StatusUnauthorized, MsgUnauthorized)
}

func InvalidCredentials() Error {
	return New(CodeAuthInvalidCredentials, fiber.StatusUnauthorized, MsgInvalidCredentials)
}
