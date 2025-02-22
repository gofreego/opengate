package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
)

type GatewayController struct {
}

func NewGatewayController() *GatewayController {
	return &GatewayController{}
}

func (c *GatewayController) Register(ctx context.Context, router gin.IRouter) {

}
