package files

import (
	"api-gateway/internal/models/dao"
	"api-gateway/pkg/goutils"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Path               string
	RefreshIntervalSec int
}

type Repository struct {
	cfg            *Config
	eventPublisher *goutils.EventPublisher[*dao.RouteConfig]
	once           sync.Once
	watchDetails   map[string]os.FileInfo
}

func NewRepository(cfg *Config) *Repository {
	return &Repository{
		cfg:          cfg,
		watchDetails: make(map[string]os.FileInfo),
	}
}

// ReadFile reads the file and returns the route config
func (r *Repository) ReadFile(ctx context.Context, filename string, value any) error {
	v := viper.New()
	v.SetConfigFile(r.cfg.Path + "/" + filename)
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		v.SetConfigType("yaml")
	} else if strings.HasSuffix(filename, ".json") {
		v.SetConfigType("json")
	} else {
		logger.Error(ctx, "unsupported file type: %s", filename)
		return fmt.Errorf("unsupported file type: %s", filename)
	}

	if err := v.ReadInConfig(); err != nil {
		logger.Error(ctx, "failed to read file: %s", filename)
		return err
	}

	if err := v.Unmarshal(value); err != nil {
		logger.Error(ctx, "failed to unmarshal file: %s", filename)
		return err
	}
	return nil
}

func (r *Repository) readRouteConfig(ctx context.Context, filename string) (*dao.RouteConfig, error) {
	var routeConfig *dao.RouteConfig = new(dao.RouteConfig)
	err := r.ReadFile(ctx, filename, routeConfig)
	if err != nil {
		logger.Error(ctx, "failed to read file: %s", filename)
		return nil, err
	}
	var isJson bool = strings.HasSuffix(filename, ".json")
	err = routeConfig.Validate()
	if err != nil {
		logger.Error(ctx, "invalid route config file: %s, \n%s", filename, routeConfig.String(isJson))
		return nil, err
	}
	logger.Info(ctx, "route config for file :%s imported as \n%s", filename, routeConfig.String(isJson))
	return routeConfig, nil
}

// read all the files in the directory and return the list of route configs
func (r *Repository) GetRoutesConfig(ctx context.Context) ([]*dao.RouteConfig, error) {
	files, err := os.ReadDir(r.cfg.Path)
	if err != nil {
		logger.Error(ctx, "failed to read directory: %s", r.cfg.Path)
		return nil, err
	}

	var routeConfigs []*dao.RouteConfig
	for _, file := range files {
		if !file.IsDir() {
			routeConfig, err := r.readRouteConfig(ctx, file.Name())
			if err != nil {
				logger.Error(ctx, "failed to read route config: %s", file.Name())
				continue
			}
			routeConfigs = append(routeConfigs, routeConfig)
		}
	}

	return routeConfigs, nil
}

// WatchRoutesConfigChanges returns a channel that will receive route config changes if any
func (r *Repository) WatchRoutesConfigChanges(ctx context.Context) (<-chan *dao.RouteConfig, error) {
	r.once.Do(func() {
		r.eventPublisher = goutils.NewEventPublisher[*dao.RouteConfig]()
		go r.watch(ctx)
	})

	return r.eventPublisher.Subscribe(), nil
}

// watch for file changes and publish the changes
func (r *Repository) watch(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.cfg.RefreshIntervalSec) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logger.Debug(ctx, "checking for file changes")
		// check if files have added or modified time has changed
		files, err := os.ReadDir(r.cfg.Path)
		if err != nil {
			logger.Error(ctx, "failed to read directory: %s", r.cfg.Path)
			continue
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				logger.Error(ctx, "failed to get file info: %s", file.Name())
				continue
			}
			if !info.IsDir() {
				if fileInfo, ok := r.watchDetails[file.Name()]; !ok || (fileInfo.ModTime() != info.ModTime()) {
					logger.Debug(ctx, "file changed: %s", file.Name())
					routeConfig, err := r.readRouteConfig(ctx, file.Name())
					if err != nil {
						logger.Error(ctx, "failed to read route config file: %s", file.Name())
						continue
					}
					r.eventPublisher.Publish(routeConfig)
					r.watchDetails[file.Name()] = info
					logger.Info(ctx, "file changed: %s", file.Name())
				}
			}
		}
	}
}
