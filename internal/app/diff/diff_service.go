package diff

import (
	"synthema/internal/observability"
	"synthema/internal/ports/repository"
)

type Service struct {
	logger *observability.Logger
	repo   repository.DiffRepository
}

func NewService(logger *observability.Logger, repo repository.DiffRepository) *Service {
	return &Service{logger: logger, repo: repo}
}
