package internal

import (
	"context"

	_ "api-gateway/docs"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type InternalController struct {
}

func NewInternalController() *InternalController {
	return &InternalController{}
}

func (c *InternalController) Register(ctx context.Context, router gin.IRouter) {

	internalGroup := router.Group("/internal")

	internalGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
