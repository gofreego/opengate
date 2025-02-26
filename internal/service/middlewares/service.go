package middlewares

import (
	"api-gateway/internal/models/dao"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/logger"
)

type MiddlewareService struct {
	middlewares map[string][]Middleware
}

type Config struct {
}

func NewMiddlewareService(config *Config) *MiddlewareService {
	return &MiddlewareService{
		middlewares: make(map[string][]Middleware),
	}
}

func (s *MiddlewareService) UpdateMiddlewares(ctx context.Context, cfg ...*dao.RouteConfig) error {
	for _, c := range cfg {

		var middlewares []Middleware
		// default last middleware
		reverseProxyMiddleware, err := s.GetReverseProxyMiddleware(c.Target)
		if err != nil {
			logger.Error(ctx, "failed to get reverse proxy middleware for id: %s, Err: %s", c.ID, err.Error())
			continue
		}
		middlewares = append(middlewares, reverseProxyMiddleware)
		s.middlewares[c.ID] = middlewares
	}
	return nil
}

func (s *MiddlewareService) ExecuteMiddleware(ctx *gin.Context, id string) error {
	if middlewares, found := s.middlewares[id]; found {
		for _, m := range middlewares {
			m(ctx)
		}
		return nil
	}
	logger.Error(ctx, "middleware not configured for id:%s", id)
	return customerrors.ERROR_INTERNAL_SERVER_ERROR
}
