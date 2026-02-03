package replay

import (
	"time"

	"synthema/internal/domain/traffic"
)

type ReplayID string

type ReplayJob struct {
	ID        ReplayID
	CaptureID traffic.CaptureID
	CreatedAt time.Time
}
