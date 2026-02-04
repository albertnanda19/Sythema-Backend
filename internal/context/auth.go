package authctx

import (
	"context"

	"github.com/google/uuid"
)

type authInfoKey struct{}

type AuthInfo struct {
	UserID    uuid.UUID
	Email     string
	Roles     []string
	SessionID uuid.UUID
}

func WithAuthInfo(ctx context.Context, info AuthInfo) context.Context {
	roles := make([]string, len(info.Roles))
	copy(roles, info.Roles)
	info.Roles = roles
	return context.WithValue(ctx, authInfoKey{}, info)
}

func GetAuthInfo(ctx context.Context) (AuthInfo, bool) {
	if ctx == nil {
		return AuthInfo{}, false
	}
	v := ctx.Value(authInfoKey{})
	info, ok := v.(AuthInfo)
	return info, ok
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return uuid.UUID{}, false
	}
	return info.UserID, true
}

func Email(ctx context.Context) (string, bool) {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return "", false
	}
	return info.Email, true
}

func Roles(ctx context.Context) ([]string, bool) {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return nil, false
	}
	roles := make([]string, len(info.Roles))
	copy(roles, info.Roles)
	return roles, true
}

func SessionID(ctx context.Context) (uuid.UUID, bool) {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return uuid.UUID{}, false
	}
	return info.SessionID, true
}
