package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"net/http"

	gateway_server "github.com/gofreego/opengate/cmd/gateway_server"
	"github.com/gofreego/opengate/cmd/http_server"
	"github.com/gofreego/opengate/internal/configs"
	"github.com/gofreego/opengate/internal/repository"
	"github.com/gofreego/opengate/internal/service"

	"github.com/gofreego/goutils/apputils"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/logger"
)

var (
	env  string
	path string
)

//go:embed all:admin/dist
var uiDist embed.FS

func getUIFileSystem() http.FileSystem {
	// Re-map the embedded filesystem to the root of 'dist'
	fsys, err := fs.Sub(uiDist, "admin/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func getIndexHTML() []byte {
	data, err := fs.ReadFile(uiDist, "admin/dist/index.html")
	if err != nil {
		panic(err)
	}
	return data
}

func getUIHandler() http.Handler {
	return http_server.GetUIHandler(getUIFileSystem(), getIndexHTML())
}

func main() {
	flag.StringVar(&env, "env", "dev", "-env=dev")
	flag.StringVar(&path, "path", ".", "-path=./")
	flag.Parse()
	ctx := context.Background()

	conf := configs.LoadConfig(ctx, path, env)

	conf.Logger.InitiateLogger()
	logger.AddMiddleLayers(logger.RequestMiddleLayer)

	// Create repository instance
	repo := repository.GetInstance(ctx, &conf.Repository)
	var cacheInstance cache.Cache
	if conf.Cache.Name != "" {
		cacheInstance = cache.NewCache(ctx, &conf.Cache)
	}
	// Create service instance
	svc := service.NewService(ctx, &conf.Service, repo, cacheInstance)

	// Create Admin HTTP server (APIs + UI)
	httpServer := http_server.NewHTTPServer(&conf.Server, svc, env, getUIHandler())
	go httpServer.Run(ctx)

	// Create Gateway server (proxy/routing)
	gatewayServer := gateway_server.NewGatewayServer(&conf.Server, svc)
	go gatewayServer.Run(ctx)

	apputils.GracefulShutdown(ctx, httpServer, gatewayServer)
}
