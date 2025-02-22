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

// liveness godoc
// @Summary Check if the service is alive
// @Description Check if the service is alive
// @Produce  json
// @Success 200 {string} string "OK"
// @Router /health/liveness [get]
func (c *HealthController) liveness(ctx *gin.Context) {
	response.WriteSuccess(ctx, "OK")
}

// readiness godoc
// @Summary Check if the service is ready
// @Description Check if the service is ready
// @Produce  json
// @Success 200 {string} string "OK"
// @Router /health/readiness [get]
func (c *HealthController) readiness(ctx *gin.Context) {
	response.WriteSuccess(ctx, "OK")
}
