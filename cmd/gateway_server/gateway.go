package gateway_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofreego/opengate/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/api"
	"github.com/gofreego/goutils/logger"
)

// Config represents gateway server settings
type Config struct {
	Port           int           `json:"port" yaml:"Port"`
	GinMode        string        `json:"ginMode" yaml:"GinMode"`
	ReadTimeout    time.Duration `json:"readTimeout" yaml:"ReadTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout" yaml:"WriteTimeout"`
	IdleTimeout    time.Duration `json:"idleTimeout" yaml:"IdleTimeout"`
	MaxHeaderBytes int           `json:"maxHeaderBytes" yaml:"MaxHeaderBytes"`
	EnableCORS     bool          `json:"enableCors" yaml:"EnableCors"`
}

type GatewayServer struct {
	cfg     *Config
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

func NewGatewayServer(cfg *Config, service *service.Service) *GatewayServer {
	return &GatewayServer{
		cfg:     cfg,
		service: service,
	}
}

func (g *GatewayServer) Run(ctx context.Context) error {
	if g.cfg.Port == 0 {
		logger.Panic(ctx, "gateway port is not provided")
	}

	// Create gin router for proxy routes
	gin.SetMode(g.cfg.GinMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(api.RequestTimeMiddleware)
	ginRouter.Use(api.RequestIDMiddleware)

	if g.cfg.EnableCORS {
		ginRouter.Use(api.OptionRequestMiddleware)
	}

	// Catch-all route handler - forwards all requests to service.RouteRequest
	ginRouter.NoRoute(g.service.RouteRequest)

	// Apply CORS middleware if enabled
	var handler http.Handler = ginRouter
	if g.cfg.EnableCORS {
		handler = corsMiddleware(ginRouter)
	}

	g.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", g.cfg.Port),
		Handler:        logger.WithRequestMiddleware(logger.WithRequestTimeMiddleware(handler)),
		ReadTimeout:    g.cfg.ReadTimeout,
		WriteTimeout:   g.cfg.WriteTimeout,
		IdleTimeout:    g.cfg.IdleTimeout,
		MaxHeaderBytes: g.cfg.MaxHeaderBytes,
	}

	logger.Info(ctx, "Started Gateway server on port %d", g.cfg.Port)

	// Start HTTP server
	err := g.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic(ctx, "failed to start gateway server : %v", err)
	}
	logger.Info(ctx, "Gateway server stopped")
	return nil
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-User-Id, X-User-Perms")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
