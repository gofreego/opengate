package http_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofreego/opengate/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/api"
	"github.com/gofreego/goutils/api/debug"
	"github.com/gofreego/goutils/logger"
)

// GlobalConfig represents global gateway settings
type Config struct {
	Port           int           `json:"port" yaml:"Port"`
	GinMode        string        `json:"ginMode" yaml:"GinMode"`
	ReadTimeout    time.Duration `json:"readTimeout" yaml:"ReadTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout" yaml:"WriteTimeout"`
	IdleTimeout    time.Duration `json:"idleTimeout" yaml:"IdleTimeout"`
	MaxHeaderBytes int           `json:"maxHeaderBytes" yaml:"MaxHeaderBytes"`
	EnableCORS     bool          `json:"enableCors" yaml:"EnableCors"`
	Debug          debug.Config  `json:"debug" yaml:"Debug"`
}

type HTTPServer struct {
	cfg     *Config
	server  *http.Server
	service *service.Service
	env     string
}

func (a *HTTPServer) Name() string {
	return "HTTP_Server"
}

func (a *HTTPServer) Shutdown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Panic(ctx, "failed to shutdown %s : %v", a.Name(), err)
	}
}

func NewHTTPServer(cfg *Config, service *service.Service, env string) *HTTPServer {
	return &HTTPServer{
		cfg:     cfg,
		service: service,
		env:     env,
	}
}

func (s *HTTPServer) registerRoutes(ctx context.Context, router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/opengate/v1")

	// Ping endpoint
	v1.GET("/ping", s.ping)

	// debug endpoint - register on the main router instead of group
	// since debug.RegisterDebugHandlers expects *http.ServeMux
	if s.cfg.Debug.Enabled {
		// Create a subrouter for debug endpoints to work with the debug package
		// We'll handle this differently since Gin and the debug package expect different router types
		logger.Info(ctx, "Debug endpoints will not be available!! since Gin framework is being used")
	}

	// Catch-all route handler - forwards all non-matching requests to service.RouteRequest
	router.NoRoute(s.service.RouteRequest)
}

func (a *HTTPServer) Run(ctx context.Context) error {

	if a.cfg.Port == 0 {
		logger.Panic(ctx, "http port is not provided")
	}

	gin.SetMode(a.cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(api.RequestTimeMiddleware)
	router.Use(api.RequestIDMiddleware)

	if a.cfg.EnableCORS {
		router.Use(api.OptionRequestMiddleware)
	}

	a.registerRoutes(ctx, router)

	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.Port),
		Handler:        router,
		ReadTimeout:    a.cfg.ReadTimeout,
		WriteTimeout:   a.cfg.WriteTimeout,
		IdleTimeout:    a.cfg.IdleTimeout,
		MaxHeaderBytes: a.cfg.MaxHeaderBytes,
	}

	if a.cfg.Debug.Enabled {
		logger.Info(ctx, "Debug dashboard available at `http://localhost:%d/opengate/v1/debug`", a.cfg.Port)
	}
	logger.Info(ctx, "Started HTTP server on port %d", a.cfg.Port)
	// Start HTTP server
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic(ctx, "failed to start http server : %v", err)
	}
	logger.Info(ctx, "HTTP server stopped")
	return nil
}
