package replay

import (
	"context"

	"synthema/internal/observability"
	"synthema/internal/ports/repository"
)

type Service struct {
	logger *observability.Logger

	repo repository.ReplayRepository
}

func NewService(logger *observability.Logger, repo repository.ReplayRepository) *Service {
	return &Service{logger: logger, repo: repo}
}

func (s *Service) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
