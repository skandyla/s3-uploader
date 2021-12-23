package service

import (
	"context"

	"github.com/skandyla/s3-uploader/internal/models"
)

type HealthService interface {
	Ping(ctx context.Context) error
	Info(ctx context.Context) (models.InfoDependencyItem, error)
}

type Health struct {
	repo HealthService
}

func NewHealth(repo HealthService) *Health {
	return &Health{
		repo: repo,
	}
}

func (s *Health) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *Health) Info(ctx context.Context) (models.InfoDependencyItem, error) {
	return s.repo.Info(ctx)
}
