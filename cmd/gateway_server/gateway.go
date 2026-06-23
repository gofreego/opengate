package gateway_server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofreego/opengate/internal/configs"
	"github.com/gofreego/opengate/internal/service"
	"github.com/gofreego/opengate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/api"
	"github.com/gofreego/goutils/logger"
)

type GatewayServer struct {
	cfg     *configs.Server
	server  *http.Server
	service *service.Service
}

func (g *GatewayServer) Name() string {
	return "Gateway_Server"
}

func (g *GatewayServer) Shutdown(ctx context.Context) {
	if g.server == nil {
		return
	}
	if err := g.server.Shutdown(ctx); err != nil {
		logger.Panic(ctx, "failed to shutdown %s : %v", g.Name(), err)
	}
}

func NewGatewayServer(cfg *configs.Server, service *service.Service) *GatewayServer {
	return &GatewayServer{
		cfg:     cfg,
		service: service,
	}
}

func (g *GatewayServer) Run(ctx context.Context) error {
	if g.cfg.GatewayPort == 0 {
		logger.Panic(ctx, "gateway port is not provided")
	}

	// Create gin router for proxy routes
	gin.SetMode(g.cfg.GinMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(api.RequestTimeMiddleware)
	ginRouter.Use(api.RequestIDMiddleware)
	ginRouter.Use(api.OptionRequestMiddleware)

	// Catch-all route handler - forwards all requests to service.RouteRequest
	ginRouter.NoRoute(g.service.RouteRequest)

	// Apply CORS middleware using dynamic config from settings store
	handler := utils.CorsMiddleware(ginRouter, g.service.GetCORSConfig)

	g.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", g.cfg.GatewayPort),
		Handler:        logger.WithRequestMiddleware(logger.WithRequestTimeMiddleware(handler)),
		ReadTimeout:    g.cfg.ReadTimeout,
		WriteTimeout:   g.cfg.WriteTimeout,
		IdleTimeout:    g.cfg.IdleTimeout,
		MaxHeaderBytes: g.cfg.MaxHeaderBytes,
	}

	logger.Info(ctx, "Started Gateway server on port %d", g.cfg.GatewayPort)

	// Start HTTP server
	err := g.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic(ctx, "failed to start gateway server : %v", err)
	}
	logger.Info(ctx, "Gateway server stopped")
	return nil
}
