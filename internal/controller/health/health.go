package health

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/response"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Register(ctx context.Context, router gin.IRouter) {
	healthGroup := router.Group("/health")

	healthGroup.GET("/liveness", c.liveness)
	healthGroup.GET("/readiness", c.readiness)
}

func (c *HealthController) liveness(ctx *gin.Context) {
	response.WriteSuccess(ctx, "OK")
}

func (c *HealthController) readiness(ctx *gin.Context) {
	response.WriteSuccess(ctx, "OK")
}
