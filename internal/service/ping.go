package service

import (
	"context"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/internal/models"
)

func (s *Service) Ping(ctx context.Context) (*models.PingResponse, error) {
	logger.Debug(ctx, "Ping request received")
	err := s.repo.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &models.PingResponse{
		Message: "Its fine here...!",
	}, nil
}
