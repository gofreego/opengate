package repository

import (
	"api-gateway/internal/models/dao"
	"api-gateway/internal/repository/files"
	"context"
)

type Repository interface {
	GetRoutesConfig(ctx context.Context) ([]*dao.RouteConfig, error)
	WatchRoutesConfigChanges(ctx context.Context) (<-chan *dao.RouteConfig, error)
}

type Config struct {
	Name  string
	Files files.Config
}

const (
	FILES_REPOSITORY    = "files"
	MONGO_REPOSITORY    = "mongo"
	POSTGRES_REPOSITORY = "postgres"
	CONSUL_REPOSITORY   = "consul"
)

func NewRepository(cfg *Config) Repository {
	switch cfg.Name {
	case FILES_REPOSITORY:
		return files.NewRepository(&cfg.Files)
	}
	panic("unknown repository name")
}
