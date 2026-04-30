package service

import (
	"context"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/api/opengate_v1"
)

// Ping implements OpenGateServiceServer.Ping
func (s *Service) Ping(ctx context.Context, req *opengate_v1.PingRequest) (*opengate_v1.PingResponse, error) {
	logger.Debug(ctx, "Ping request received: %v", req.GetMessage())
	err := s.repo.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &opengate_v1.PingResponse{
		Message: "Pong",
	}, nil
}
