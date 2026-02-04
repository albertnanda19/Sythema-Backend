package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeNotFound = "resource.not_found"
	MsgNotFound  = "Not found"
)

func NotFound() Error {
	return New(CodeNotFound, fiber.StatusNotFound, MsgNotFound)
}
