package repository

import (
	"context"

	"synthema/internal/domain/diff"
	"synthema/internal/domain/replay"
	"synthema/internal/domain/traffic"
)

type TrafficRepository interface {
	SaveCapturedTraffic(ctx context.Context, t traffic.CapturedTraffic) error
}

type ReplayRepository interface {
	SaveReplayJob(ctx context.Context, j replay.ReplayJob) error
}

type DiffRepository interface {
	SaveDiffResult(ctx context.Context, r diff.DiffResult) error
}
