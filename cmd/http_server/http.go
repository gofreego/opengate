package http_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/configs"
	"github.com/gofreego/opengate/internal/service"
	"github.com/gofreego/opengate/pkg/utils"

	"github.com/gofreego/goutils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type HTTPServer struct {
	cfg       *configs.Server
	server    *http.Server
	service   *service.Service
	env       string
	uiHandler http.Handler
}

func (a *HTTPServer) Name() string {
	return "HTTP_Server"
}

func (a *HTTPServer) Shutdown(ctx context.Context) {
	if a.server == nil {
		return
	}
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Panic(ctx, "failed to shutdown %s : %v", a.Name(), err)
	}
}

func NewHTTPServer(cfg *configs.Server, service *service.Service, env string, uiHandler http.Handler) *HTTPServer {
	return &HTTPServer{
		cfg:       cfg,
		service:   service,
		env:       env,
		uiHandler: uiHandler,
	}
}

func (a *HTTPServer) Run(ctx context.Context) error {

	if a.cfg.AdminPort == 0 {
		logger.Panic(ctx, "http port is not provided")
	}

	// Create grpc-gateway mux for API routes
	grpcMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "X-User-Id", "X-User-Perms":
				return strings.ToLower(key), true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
	)

	// Register OpenGateService with grpc-gateway
	err := opengate_v1.RegisterOpenGateServiceHandlerServer(ctx, grpcMux, a.service)
	if err != nil {
		logger.Panic(ctx, "failed to register OpenGateService: %v", err)
	}

	// Register Swagger handler
	utils.RegisterSwaggerHandler(grpcMux, "/opengate/v1/swagger", "./api/docs/proto", "/opengate/v1/opengate.swagger.json")

	// Create combined handler that routes based on path
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Direct API requests to grpc-gateway mux
		if strings.HasPrefix(path, "/opengate/v1/") {
			grpcMux.ServeHTTP(w, r)
			return
		}

		// Direct UI requests or root to uiHandler
		if strings.HasPrefix(path, "/gateway") || path == "/" {
			a.uiHandler.ServeHTTP(w, r)
			return
		}

		// Return 404 for other paths (proxy is on gateway server)
		http.NotFound(w, r)
	})

	// Always apply CORS middleware using dynamic config from settings store
	handler := utils.CorsMiddleware(finalHandler, a.service.GetCORSConfig)

	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.AdminPort),
		Handler:        logger.WithRequestMiddleware(logger.WithRequestTimeMiddleware(handler)),
		ReadTimeout:    a.cfg.ReadTimeout,
		WriteTimeout:   a.cfg.WriteTimeout,
		IdleTimeout:    a.cfg.IdleTimeout,
		MaxHeaderBytes: a.cfg.MaxHeaderBytes,
	}

	if a.cfg.Debug.Enabled {
		logger.Info(ctx, "Debug dashboard available at `http://localhost:%d/opengate/v1/debug`", a.cfg.AdminPort)
	}
	logger.Info(ctx, "Started Admin HTTP server on port %d", a.cfg.AdminPort)
	logger.Info(ctx, "Admin UI available at `http://localhost:%d/gateway/`", a.cfg.AdminPort)
	logger.Info(ctx, "API endpoints available at `http://localhost:%d/opengate/v1/`", a.cfg.AdminPort)
	logger.Info(ctx, "Swagger UI available at `http://localhost:%d/opengate/v1/swagger`", a.cfg.AdminPort)

	// Start HTTP server
	err = a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic(ctx, "failed to start http server : %v", err)
	}
	logger.Info(ctx, "Admin HTTP server stopped")
	return nil
}
