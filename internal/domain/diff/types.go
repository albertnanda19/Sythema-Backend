package diff

import (
	"time"

	"synthema/internal/domain/replay"
)

type DiffID string

type DiffResult struct {
	ID        DiffID
	ReplayID  replay.ReplayID
	CreatedAt time.Time
}
