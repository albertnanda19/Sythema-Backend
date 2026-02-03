package queue

import (
	"context"

	"synthema/internal/domain/traffic"
)

type TrafficQueue interface {
	EnqueueCapturedTraffic(ctx context.Context, t traffic.CapturedTraffic) error
}
