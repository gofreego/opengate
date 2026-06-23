package openauth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/openauth/pkg/clients/openauth"
	"github.com/gofreego/opengate/internal/models"
)

// ErrNotImplemented is returned for operations not supported by openauth repository
var ErrNotImplemented = errors.New("operation not implemented for openauth repository")

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
		route.UpdatedAt = config.UpdatedAt
		routes = append(routes, &route)
	}
	return routes, nil
}

// CreateConfig is not implemented for openauth repository
func (r *Repository) CreateConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	return nil, ErrNotImplemented
}

// GetConfigByID is not implemented for openauth repository
func (r *Repository) GetConfigByID(ctx context.Context, id int64) (*models.Config, error) {
	return nil, ErrNotImplemented
}

// ListConfigs is not implemented for openauth repository
func (r *Repository) ListConfigs(ctx context.Context, filter *models.ConfigFilter) ([]*models.Config, int, error) {
	return nil, 0, ErrNotImplemented
}

// UpdateConfig is not implemented for openauth repository
func (r *Repository) UpdateConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	return nil, ErrNotImplemented
}

// DeleteConfig is not implemented for openauth repository
func (r *Repository) DeleteConfig(ctx context.Context, id int64) error {
	return ErrNotImplemented
}

// GetAppSettings is not implemented for openauth repository
func (r *Repository) GetAppSettings(ctx context.Context) ([]*models.AppSetting, error) {
	return nil, nil
}

// UpsertAppSetting is not implemented for openauth repository
func (r *Repository) UpsertAppSetting(ctx context.Context, setting *models.AppSetting) error {
	return ErrNotImplemented
}
