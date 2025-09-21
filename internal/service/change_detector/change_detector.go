package changedetector

import (
	"context"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/internal/models"
	routemanager "github.com/gofreego/opengate/internal/service/route_manager"
)

type ChangeDetector interface {
	DetectChanges(ctx context.Context) error
}

type Repository interface {
	GetRoutes(ctx context.Context) ([]*models.ServiceRoute, error)
}

type Config struct {
	RouteUpdateInterval time.Duration `yaml:"RouteUpdateInterval"`
}

type changeDetector struct {
	repo         Repository
	routeManager routemanager.Manager
	cfg          *Config
}

func New(repo Repository, routeManager routemanager.Manager, cfg *Config) ChangeDetector {
	return &changeDetector{
		repo:         repo,
		routeManager: routeManager,
		cfg:          cfg,
	}
}

// detectChanges periodically checks for changes in the routes and updates the route manager accordingly
func (cd *changeDetector) DetectChanges(ctx context.Context) error {
	logger.Info(ctx, "Route change detector started")

	// Load initial routes
	if err := cd.checkAndUpdateRoutes(ctx); err != nil {
		logger.Error(ctx, "Failed to load initial routes: %v", err)
	}
	// Default update interval if not configured
	updateInterval := 30 * time.Second
	if cd.cfg.RouteUpdateInterval > 0 {
		updateInterval = cd.cfg.RouteUpdateInterval
	}

	logger.Info(ctx, "Starting route change detector with interval: %v", updateInterval)
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "Route change detector stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := cd.checkAndUpdateRoutes(ctx); err != nil {
				logger.Error(ctx, "Error checking for route changes: %v", err)
				// Continue monitoring despite errors
				continue
			}
		}
	}
}

// checkAndUpdateRoutes checks for changes and updates routes if necessary
func (cd *changeDetector) checkAndUpdateRoutes(ctx context.Context) error {
	// Fetch latest routes from repository
	latestRoutes, err := cd.repo.GetRoutes(ctx)
	if err != nil {
		return err
	}

	if cd.isDeletedAny(latestRoutes) {
		logger.Info(ctx, "Route deletion detected, reloading all routes")

		return nil
	}

	for _, route := range latestRoutes {
		existingRoute := cd.routeManager.GetRouteByName(route.Name)
		if existingRoute == nil {
			// New route added
			cd.routeManager.AddRoute(route)
			logger.Info(ctx, "New route added: %s, value: %+v", route.Name, route)
			continue
		}

		if route.UpdatedAt > existingRoute.UpdatedAt {
			// Route updated
			cd.routeManager.AddRoute(route)
			logger.Info(ctx, "Route updated: %s, value: %v", route.Name, route)
		}
	}
	return nil
}

func (cd *changeDetector) isDeletedAny(latestRoutes []*models.ServiceRoute) bool {
	existingRoutes := cd.routeManager.GetRoutes()
	latestRouteNames := make(map[string]struct{})
	for _, route := range latestRoutes {
		latestRouteNames[route.Name] = struct{}{}
	}

	for _, existingRoute := range existingRoutes {
		if _, exists := latestRouteNames[existingRoute.Name]; !exists {
			return true
		}
	}
	return false
}
