package files

import (
	"api-gateway/internal/models/dao"
	"context"
)

type Config struct {
}

type Repository struct {
	cfg *Config
}

func NewRepository(cfg *Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) GetRoutesConfig(ctx context.Context) ([]dao.RouteConfig, error) {
	return []dao.RouteConfig{
		{
			ID:          "1",
			Name:        "user service",
			Description: "user service route configurations",
			Match: dao.MatchConfig{
				Host:    "",
				Prefix:  "/photo",
				Regex:   "",
				Methods: []string{"GET", "POST"},
			},
			Target: "https://images.unsplash.com",
		},
	}, nil
}

// WatchRoutesConfig implements repository.Repository.
func (r *Repository) WatchRoutesConfig(ctx context.Context) (<-chan []dao.RouteConfig, error) {
	panic("unimplemented")
}
