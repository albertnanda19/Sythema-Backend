package errors

import "github.com/gofiber/fiber/v2"

const (
	CodeAuthUnauthorized       = "auth.unauthorized"
	CodeAuthMissingSession     = "auth.missing_session_cookie"
	CodeAuthInvalidSession     = "auth.invalid_session_cookie"
	CodeAuthSessionExpired     = "auth.session_expired"
	CodeAuthSessionRevoked     = "auth.session_revoked"
	CodeAuthUserInactive       = "auth.user_inactive"
	CodeAuthInvalidCredentials = "auth.invalid_credentials"
	MsgUnauthorized            = "Unauthorized"
	MsgMissingSession          = "Missing session cookie"
	MsgInvalidSession          = "Invalid session cookie"
	MsgSessionExpired          = "Session expired"
	MsgSessionRevoked          = "Session revoked"
	MsgUserInactive            = "User inactive"
	MsgInvalidCredentials      = "Invalid credentials"
)

func Unauthorized() Error {
	return New(CodeAuthUnauthorized, fiber.StatusUnauthorized, MsgUnauthorized)
}

func MissingSessionCookie() Error {
	return New(CodeAuthMissingSession, fiber.StatusUnauthorized, MsgMissingSession)
}

func InvalidSessionCookie() Error {
	return New(CodeAuthInvalidSession, fiber.StatusUnauthorized, MsgInvalidSession)
}

func SessionExpired() Error {
	return New(CodeAuthSessionExpired, fiber.StatusUnauthorized, MsgSessionExpired)
}

func SessionRevoked() Error {
	return New(CodeAuthSessionRevoked, fiber.StatusUnauthorized, MsgSessionRevoked)
}

func UserInactive() Error {
	return New(CodeAuthUserInactive, fiber.StatusUnauthorized, MsgUserInactive)
}

func InvalidCredentials() Error {
	return New(CodeAuthInvalidCredentials, fiber.StatusUnauthorized, MsgInvalidCredentials)
}
