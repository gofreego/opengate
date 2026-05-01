package http_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/service"
	"github.com/gofreego/opengate/pkg/utils"

	"github.com/gofreego/goutils/api/debug"
	"github.com/gofreego/goutils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Config represents admin server settings
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
	cfg       *Config
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

func NewHTTPServer(cfg *Config, service *service.Service, env string, uifs http.FileSystem, indexHTML []byte) *HTTPServer {
	return &HTTPServer{
		cfg:       cfg,
		service:   service,
		env:       env,
		uiHandler: getUIHandler(uifs, indexHTML),
	}
}

func (a *HTTPServer) Run(ctx context.Context) error {

	if a.cfg.Port == 0 {
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

	// Apply CORS middleware if enabled
	var handler http.Handler = finalHandler
	if a.cfg.EnableCORS {
		handler = corsMiddleware(finalHandler)
	}

	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.Port),
		Handler:        logger.WithRequestMiddleware(logger.WithRequestTimeMiddleware(handler)),
		ReadTimeout:    a.cfg.ReadTimeout,
		WriteTimeout:   a.cfg.WriteTimeout,
		IdleTimeout:    a.cfg.IdleTimeout,
		MaxHeaderBytes: a.cfg.MaxHeaderBytes,
	}

	if a.cfg.Debug.Enabled {
		logger.Info(ctx, "Debug dashboard available at `http://localhost:%d/opengate/v1/debug`", a.cfg.Port)
	}
	logger.Info(ctx, "Started Admin HTTP server on port %d", a.cfg.Port)
	logger.Info(ctx, "Admin UI available at `http://localhost:%d/gateway/`", a.cfg.Port)
	logger.Info(ctx, "API endpoints available at `http://localhost:%d/opengate/v1/`", a.cfg.Port)
	logger.Info(ctx, "Swagger UI available at `http://localhost:%d/opengate/v1/swagger`", a.cfg.Port)

	// Start HTTP server
	err = a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Panic(ctx, "failed to start http server : %v", err)
	}
	logger.Info(ctx, "Admin HTTP server stopped")
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
