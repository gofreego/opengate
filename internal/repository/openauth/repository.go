package openauth

import (
	"context"
	"encoding/json"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/openauth/pkg/clients/openauth"
	"github.com/gofreego/opengate/internal/models"
)

type Repository struct {
	client *openauth.OpenauthConfigFetcher
}

func NewRepository(ctx context.Context, cfg *openauth.ClientConfig) (*Repository, error) {
	client, err := openauth.NewOpenauthConfigFetcher(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Repository{
		client: client,
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}

func (r *Repository) GetRoutes(ctx context.Context) ([]*models.ServiceRoute, error) {

	configs, err := r.client.GetConfigsByEntityName(ctx, "opengate.routes")
	if err != nil {
		logger.Error(ctx, "Failed to fetch route configs from OpenAuth: %v", err)
		return nil, err
	}
	var routes []*models.ServiceRoute
	for _, config := range configs.GetConfigs() {
		var route models.ServiceRoute
		if err := json.Unmarshal([]byte(config.GetJsonValue()), &route); err != nil {
			logger.Error(ctx, "Failed to unmarshal route config: %v", err)
			continue
		}
		routes = append(routes, &route)
	}
	return routes, nil
}
