package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/openauth/pkg/clients/openauth"
)

type Config struct {
	Name     string                `yaml:"Name"`
	OpenAuth openauth.ClientConfig `yaml:"OpenAuth"`
}

type AuthManager interface {
	Authenticate(ctx *gin.Context) error
}

type manager struct {
	strategy Strategy
}

func NewAuthManager(ctx context.Context, config *Config, cache cache.Cache) (AuthManager, error) {
	strategy, err := NewOpenAuthStrategy(ctx, &config.OpenAuth, cache)
	if err != nil {
		return nil, err
	}
	return &manager{
		strategy: strategy,
	}, nil
}

func (m *manager) Authenticate(ctx *gin.Context) error {
	return m.strategy.Authenticate(ctx)
}
