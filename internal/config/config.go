package config

import (
	"api-gateway/internal/controller"
	"api-gateway/internal/repository"
	"api-gateway/internal/service"
	"context"
	"flag"
	"os"
	"sync"

	"github.com/gofreego/goutils/configutils"
	"github.com/gofreego/goutils/logger"
)

var (
	cfgFile = flag.String("config", "dev.yaml", "config file")
	once    = sync.Once{}
)

type Configurations struct {
	Name       string
	Logger     logger.Config
	Ctrl       controller.Config
	Service    service.Config
	Repository repository.Config
}

func GetConfig(ctx context.Context) *Configurations {
	flag.Parse()
	var cfg Configurations
	// Load configuration
	once.Do(func() {
		err := configutils.ReadConfig(ctx, *cfgFile, &cfg)
		if err != nil {
			logger.Panic(ctx, "Failed to read configuration: %v", err)
			os.Exit(1)
		}
	})
	return &cfg
}
