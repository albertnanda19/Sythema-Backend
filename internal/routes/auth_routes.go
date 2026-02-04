package routes

import (
	"github.com/gofiber/fiber/v2"

	authhandlers "synthema/internal/handlers/auth"
)

func RegisterAuthRoutes(v1 fiber.Router, authMW fiber.Handler, meHandler *authhandlers.MeHandler, logoutMW fiber.Handler, logoutHandler *authhandlers.LogoutHandler) {
	auth := v1.Group("/auth")
	auth.Get("/me", authMW, meHandler.Me)
	auth.Post("/logout", logoutMW, logoutHandler.Logout)
}
