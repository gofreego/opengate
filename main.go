package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/controller"
	"context"
	"flag"

	"github.com/gofreego/goutils/apputils"
	"github.com/gofreego/goutils/logger"
	"gopkg.in/yaml.v3"
)

func main() {
	flag.Parse()
	var ctx context.Context = context.Background()
	cfg := config.GetConfig(ctx)
	cfg.Logger.InitiateLogger()
	// logging the config
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		logger.Panic(ctx, "Failed to marshal config: %v", err)
		return
	}
	logger.Debug(ctx, "Configurations: \n%s", bytes)

	ctrl := controller.New(&cfg.Ctrl)
	// starting the controller
	go ctrl.Run(ctx)

	// Graceful shutdown goroutine
	apputils.GracefulShutdown(ctx, ctrl)
}
