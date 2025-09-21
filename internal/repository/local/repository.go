package local

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/internal/models"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// path to folder containing route definitions in JSON or YAML format
	// each file represents a route
	RoutesFolderPath string `yaml:"RoutesFolderPath"`
}

type Repository struct {
	cfg *Config
}

func NewRepository(ctx context.Context, cfg *Config) (*Repository, error) {
	return &Repository{
		cfg: cfg,
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}

func (r *Repository) GetRoutes(ctx context.Context) ([]*models.ServiceRoute, error) {
	var routes []*models.ServiceRoute

	if r.cfg.RoutesFolderPath == "" {
		logger.Warn(ctx, "Routes folder path is not configured")
		return routes, nil
	}

	// Check if directory exists
	if _, err := os.Stat(r.cfg.RoutesFolderPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("routes folder does not exist: %s", r.cfg.RoutesFolderPath)
	}

	// Walk through the directory
	err := filepath.WalkDir(r.cfg.RoutesFolderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error(ctx, "Error walking route files: %v", err)
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Process only JSON and YAML files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			logger.Debug(ctx, "Skipping non-route file: %s", path)
			return nil
		}

		// Parse the route file
		route, err := r.parseRouteFile(ctx, path)
		if err != nil {
			logger.Error(ctx, "Failed to parse route file %s: %v", path, err)
			return nil
		}

		// get the updated timestamp of the file
		if info, err := d.Info(); err == nil {
			route.UpdatedAt = info.ModTime().UnixMilli()
		}

		if route != nil {
			routes = append(routes, route)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read routes from folder %s: %w", r.cfg.RoutesFolderPath, err)
	}
	return routes, nil
}

// parseRouteFile reads and parses a single route configuration file
func (r *Repository) parseRouteFile(ctx context.Context, filePath string) (*models.ServiceRoute, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Skip empty files
	if len(data) == 0 {
		logger.Debug(ctx, "Skipping empty file: %s", filePath)
		return nil, nil
	}

	var route models.ServiceRoute
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &route); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &route); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s, only JSON and YAML files are supported", ext)
	}

	// Validate required fields
	if route.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if route.PathPrefix == "" {
		return nil, fmt.Errorf("path_prefix is required")
	}
	if route.TargetURL == "" {
		return nil, fmt.Errorf("target_url is required")
	}

	return &route, nil
}
