package service

import (
	"context"

	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/internal/service/auth"
	changedetector "github.com/gofreego/opengate/internal/service/change_detector"
	routemanager "github.com/gofreego/opengate/internal/service/route_manager"
)

type Config struct {
	Auth                  auth.Config           `yaml:"Auth"`
	ChangeDetector        changedetector.Config `yaml:"ChangeDetector"`
	InitialRoutes         []models.ServiceRoute `yaml:"InitialRoutes"`
	EnablePermissionCheck bool                  `yaml:"EnablePermissionCheck"`
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
	// Seed initial routes from config
	service.seedInitialRoutes(ctx)
	go changedetector.New(repo, service.routeManager, &cfg.ChangeDetector).DetectChanges(ctx)
	return service
}

// seedInitialRoutes seeds initial routes from config if they don't exist
func (s *Service) seedInitialRoutes(ctx context.Context) {
	if len(s.cfg.InitialRoutes) == 0 {
		return
	}

	// Get existing routes
	existingConfigs, _, err := s.repo.ListConfigs(ctx, &models.ConfigFilter{Limit: 1000})
	if err != nil {
		logger.Error(ctx, "failed to list existing configs for seeding: %v", err)
		return
	}

	// Create a map of existing route names for quick lookup
	existingNames := make(map[string]bool)
	for _, cfg := range existingConfigs {
		existingNames[cfg.Name] = true
	}

	// Seed routes that don't exist
	for _, route := range s.cfg.InitialRoutes {
		if existingNames[route.Name] {
			logger.Debug(ctx, "route %s already exists, skipping seed", route.Name)
			continue
		}

		config := &models.Config{
			Name:           route.Name,
			PathPrefix:     route.PathPrefix,
			TargetURL:      route.TargetURL,
			StripPrefix:    route.StripPrefix,
			Authentication: route.Authentication,
			Middleware:     route.Middleware,
			Timeout:        route.Timeout,
		}

		_, err := s.repo.CreateConfig(ctx, config)
		if err != nil {
			logger.Error(ctx, "failed to seed route %s: %v", route.Name, err)
			continue
		}
		logger.Info(ctx, "seeded initial route: %s", route.Name)
	}
}
