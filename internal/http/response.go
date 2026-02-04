package http

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *fiber.Ctx, status int, message string, data any) error {
	resp := APIResponse{Status: status, Message: message, Data: data}
	return c.Status(status).JSON(resp)
}

func Fail(c *fiber.Ctx, status int, message string) error {
	resp := APIResponse{Status: status, Message: message, Data: nil}
	return c.Status(status).JSON(resp)
}
