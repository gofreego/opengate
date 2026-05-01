package configs

import (
	"context"
	"fmt"
	"time"

	repo "github.com/gofreego/opengate/internal/repository"
	"github.com/gofreego/opengate/internal/service"

	"github.com/gofreego/goutils/api/debug"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/configutils"
	"github.com/gofreego/goutils/logger"
)

type Configuration struct {
	LogConfig  bool           `yaml:"LogConfig"`
	Logger     logger.Config  `yaml:"Logger"`
	AppNames   []string       `yaml:"AppNames"`
	Server     Server         `yaml:"Server"`
	Repository repo.Config    `yaml:"Repository"`
	Service    service.Config `yaml:"Service"`
	Cache      cache.Config   `yaml:"Cache"`
}

func LoadConfig(ctx context.Context, path string, env string) *Configuration {
	filePath := fmt.Sprintf("%s/%s.yaml", path, env)
	var conf Configuration
	err := configutils.ReadConfig(ctx, filePath, &conf)
	if err != nil {
		logger.Panic(ctx, "failed to read configs : %v", err)
	}
	// logging config for debug
	if conf.LogConfig {
		configutils.LogConfig(ctx, conf)
	}
	return &conf
}

// Config represents admin server settings
type Server struct {
	AdminPort      int           `json:"adminPort" yaml:"AdminPort"`
	GatewayPort    int           `json:"gatewayPort" yaml:"GatewayPort"`
	GinMode        string        `json:"ginMode" yaml:"GinMode"`
	ReadTimeout    time.Duration `json:"readTimeout" yaml:"ReadTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout" yaml:"WriteTimeout"`
	IdleTimeout    time.Duration `json:"idleTimeout" yaml:"IdleTimeout"`
	MaxHeaderBytes int           `json:"maxHeaderBytes" yaml:"MaxHeaderBytes"`
	EnableCORS     bool          `json:"enableCors" yaml:"EnableCors"`
	Debug          debug.Config  `json:"debug" yaml:"Debug"`
}
