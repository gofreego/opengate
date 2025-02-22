package internal

import (
	"context"

	"github.com/gin-gonic/gin"
)

type InternalController struct {
}

func NewInternalController() *InternalController {
	return &InternalController{}
}

func (c *InternalController) Register(ctx context.Context, router gin.IRouter) {

}
