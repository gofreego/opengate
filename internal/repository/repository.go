package repository

import (
	"context"
	"sync"

	"github.com/gofreego/openauth/pkg/clients/openauth"
	"github.com/gofreego/opengate/internal/repository/local"
	openauthRepo "github.com/gofreego/opengate/internal/repository/openauth"
	"github.com/gofreego/opengate/internal/service"
)

type Name string

const (
	Local    Name = "Local"
	OpenAuth Name = "OpenAuth"
)

type Config struct {
	Name     Name                  `yaml:"Name"`
	Local    local.Config          `yaml:"Local"`
	OpenAuth openauth.ClientConfig `yaml:"OpenAuth"`
}

var (
	instance service.Repository
	once     sync.Once
	mu       sync.RWMutex
)

// GetInstance returns the singleton instance of the repository
func GetInstance(ctx context.Context, cfg *Config) service.Repository {
	mu.RLock()
	if instance != nil {
		defer mu.RUnlock()
		return instance
	}
	mu.RUnlock()

	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			switch cfg.Name {
			case Local:
				repo, err := local.NewRepository(ctx, &cfg.Local)
				if err != nil {
					panic("failed to create repository: " + err.Error())
				}
				instance = repo
			case OpenAuth:
				repo, err := openauthRepo.NewRepository(ctx, &cfg.OpenAuth)
				if err != nil {
					panic("failed to create repository: " + err.Error())
				}
				instance = repo
			default:
				panic("unsupported repository: " + string(cfg.Name))
			}
		}
	})

	return instance
}
