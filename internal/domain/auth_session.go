package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
