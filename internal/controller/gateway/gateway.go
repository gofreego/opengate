package gateway

import (
	"api-gateway/internal/service"
	"context"

	"github.com/gin-gonic/gin"
)

type GatewayController struct {
	service *service.Service
}

func NewGatewayController(ctx context.Context, service *service.Service) *GatewayController {
	return &GatewayController{
		service: service,
	}
}

func (c *GatewayController) Register(ctx context.Context, router gin.IRouter) {
	router.Use(c.service.Route)
}
