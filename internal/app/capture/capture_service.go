package capture

import (
	"context"

	"synthema/internal/observability"
	"synthema/internal/ports/queue"
	"synthema/internal/ports/repository"
)

type Service struct {
	logger *observability.Logger

	repo  repository.TrafficRepository
	queue queue.TrafficQueue
}

func NewService(logger *observability.Logger, repo repository.TrafficRepository, q queue.TrafficQueue) *Service {
	return &Service{logger: logger, repo: repo, queue: q}
}

func (s *Service) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
