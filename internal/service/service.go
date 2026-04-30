package service

import (
	"context"

	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/internal/service/auth"
	changedetector "github.com/gofreego/opengate/internal/service/change_detector"
	routemanager "github.com/gofreego/opengate/internal/service/route_manager"
)

type Config struct {
	Auth           auth.Config           `yaml:"Auth"`
	ChangeDetector changedetector.Config `yaml:"ChangeDetector"`
}

type Repository interface {
	Ping(ctx context.Context) error
	GetRoutes(ctx context.Context) ([]*models.ServiceRoute, error)

	// Config CRUD operations
	CreateConfig(ctx context.Context, config *models.Config) (*models.Config, error)
	GetConfigByID(ctx context.Context, id int64) (*models.Config, error)
	ListConfigs(ctx context.Context, filter *models.ConfigFilter) ([]*models.Config, int, error)
	UpdateConfig(ctx context.Context, config *models.Config) (*models.Config, error)
	DeleteConfig(ctx context.Context, id int64) error
}

type Service struct {
	repo         Repository
	routeManager routemanager.Manager
	authManager  auth.AuthManager
	cfg          *Config
	opengate_v1.UnimplementedOpenGateServiceServer
}

func NewService(ctx context.Context, cfg *Config, repo Repository, cache cache.Cache) *Service {
	authManager, err := auth.NewAuthManager(ctx, &cfg.Auth, cache)
	if err != nil {
		panic("failed to create AuthManager: " + err.Error())
	}
	service := &Service{
		cfg:          cfg,
		repo:         repo,
		routeManager: routemanager.New(),
		authManager:  authManager,
	}
	go changedetector.New(repo, service.routeManager, &cfg.ChangeDetector).DetectChanges(ctx)
	return service
}
