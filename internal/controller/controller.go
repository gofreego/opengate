package controller

import (
	"api-gateway/internal/controller/gateway"
	"api-gateway/internal/controller/health"
	"api-gateway/internal/controller/internal"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/logger"
)

type Config struct {
	Port  int
	Debug bool
}

type Controller struct {
	cfg    *Config
	server *http.Server
}

func New(c *Config) *Controller {
	return &Controller{
		cfg: c,
	}
}

func (c *Controller) Run(ctx context.Context) error {
	if !c.cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// registering health controller
	healthController := health.NewHealthController()
	healthController.Register(ctx, router)

	// registering internal controller
	internalController := internal.NewInternalController()
	internalController.Register(ctx, router)

	// registering gateway controller
	gatewayController := gateway.NewGatewayController()
	gatewayController.Register(ctx, router)

	c.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", c.cfg.Port),
		Handler: router,
	}

	logger.Info(ctx, "Starting API Gateway on port %d", c.cfg.Port)
	if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(ctx, "Failed to start server: %v", err)
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func (c *Controller) Shutdown(ctx context.Context) {
	if err := c.server.Shutdown(ctx); err != nil {
		logger.Error(ctx, "Failed to shutdown server: %v", err)
	}
}

func (c *Controller) Name() string {
	return "API Gateway"
}
