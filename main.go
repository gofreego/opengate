package main

import (
	"context"
	"flag"

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

	// Create HTTP server with proper config
	app := http_server.NewHTTPServer(&conf.Server, svc, env)
	go app.Run(ctx)
	apputils.GracefulShutdown(ctx, app)
}
