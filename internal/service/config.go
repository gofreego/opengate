package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gofreego/goutils/utils"
	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/constants"
	"github.com/gofreego/opengate/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultTimeout = 30 * time.Second
	MaxLimit       = 100
	DefaultLimit   = 20
)

// CreateConfig creates a new route configuration
func (s *Service) CreateConfig(ctx context.Context, req *opengate_v1.CreateConfigRequest) (*opengate_v1.CreateConfigResponse, error) {
	// Check write permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_WRITE) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:write required")
	}

	// Validate request
	if err := validateCreateConfigRequest(req); err != nil {
		return nil, err
	}

	// Convert proto to model
	config := protoToModel(req)

	// Create in repository
	created, err := s.repo.CreateConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &opengate_v1.CreateConfigResponse{
		Config:  modelToProto(created),
		Message: "Config created successfully",
	}, nil
}

// GetConfig retrieves a config by ID
func (s *Service) GetConfig(ctx context.Context, req *opengate_v1.GetConfigRequest) (*opengate_v1.GetConfigResponse, error) {
	// Check read permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_READ) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:read required")
	}

	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid config id")
	}

	config, err := s.repo.GetConfigByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &opengate_v1.GetConfigResponse{
		Config:  modelToProto(config),
		Message: "Config retrieved successfully",
	}, nil
}

// ListConfigs retrieves configs with pagination
func (s *Service) ListConfigs(ctx context.Context, req *opengate_v1.ListConfigsRequest) (*opengate_v1.ListConfigsResponse, error) {
	// Check read permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_READ) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:read required")
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset := int(req.GetOffset())
	if offset < 0 {
		offset = 0
	}

	filter := &models.ConfigFilter{
		Search: req.GetSearch(),
		Limit:  limit,
		Offset: offset,
	}

	configs, total, err := s.repo.ListConfigs(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoConfigs := make([]*opengate_v1.Config, len(configs))
	for i, cfg := range configs {
		protoConfigs[i] = modelToProto(cfg)
	}

	return &opengate_v1.ListConfigsResponse{
		Configs: protoConfigs,
		Total:   int32(total),
		Message: "Configs retrieved successfully",
	}, nil
}

// UpdateConfig updates an existing config
func (s *Service) UpdateConfig(ctx context.Context, req *opengate_v1.UpdateConfigRequest) (*opengate_v1.UpdateConfigResponse, error) {
	// Check write permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_WRITE) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:write required")
	}

	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid config id")
	}

	// Validate request
	if err := validateUpdateConfigRequest(req); err != nil {
		return nil, err
	}

	// Convert proto to model
	config := updateProtoToModel(req)

	// Update in repository
	updated, err := s.repo.UpdateConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &opengate_v1.UpdateConfigResponse{
		Config:  modelToProto(updated),
		Message: "Config updated successfully",
	}, nil
}

// DeleteConfig deletes a config by ID
func (s *Service) DeleteConfig(ctx context.Context, req *opengate_v1.DeleteConfigRequest) (*opengate_v1.DeleteConfigResponse, error) {
	// Check write permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_WRITE) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:write required")
	}

	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid config id")
	}

	err := s.repo.DeleteConfig(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &opengate_v1.DeleteConfigResponse{
		Message: "Config deleted successfully",
	}, nil
}

// GetRoutes retrieves all routes for the routing manager
func (s *Service) GetRoutes(ctx context.Context, req *opengate_v1.GetRoutesRequest) (*opengate_v1.GetRoutesResponse, error) {
	// Check read permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_READ) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:read required")
	}

	routes, err := s.repo.GetRoutes(ctx)
	if err != nil {
		return nil, err
	}

	protoRoutes := make([]*opengate_v1.Route, len(routes))
	for i, route := range routes {
		protoRoutes[i] = serviceRouteToProto(route)
	}

	return &opengate_v1.GetRoutesResponse{
		Routes:  protoRoutes,
		Message: "Routes retrieved successfully",
	}, nil
}

// GetStats retrieves dashboard statistics
func (s *Service) GetStats(ctx context.Context, req *opengate_v1.GetStatsRequest) (*opengate_v1.GetStatsResponse, error) {
	// Check read permission
	if !utils.HasPermission(ctx, constants.PERMISSION_ROUTES_READ) {
		return nil, status.Error(codes.PermissionDenied, "permission denied: routes:read required")
	}

	// Get total count of configs using ListConfigs with limit 1
	_, total, err := s.repo.ListConfigs(ctx, &models.ConfigFilter{
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	return &opengate_v1.GetStatsResponse{
		TotalRoutes: int32(total),
		Message:     "Stats retrieved successfully",
	}, nil
}

// validateCreateConfigRequest validates the create config request
func validateCreateConfigRequest(req *opengate_v1.CreateConfigRequest) error {
	if req.GetName() == "" {
		return fmt.Errorf("name is required")
	}
	if req.GetPathPrefix() == "" {
		return fmt.Errorf("path_prefix is required")
	}
	if req.GetTargetUrl() == "" {
		return fmt.Errorf("target_url is required")
	}

	// Validate target URL format
	if _, err := url.Parse(req.GetTargetUrl()); err != nil {
		return fmt.Errorf("invalid target_url format: %w", err)
	}

	return nil
}

// validateUpdateConfigRequest validates the update config request
func validateUpdateConfigRequest(req *opengate_v1.UpdateConfigRequest) error {
	if req.GetName() == "" {
		return fmt.Errorf("name is required")
	}
	if req.GetPathPrefix() == "" {
		return fmt.Errorf("path_prefix is required")
	}
	if req.GetTargetUrl() == "" {
		return fmt.Errorf("target_url is required")
	}

	// Validate target URL format
	if _, err := url.Parse(req.GetTargetUrl()); err != nil {
		return fmt.Errorf("invalid target_url format: %w", err)
	}

	return nil
}

// protoToModel converts a CreateConfigRequest to a Config model
func protoToModel(req *opengate_v1.CreateConfigRequest) *models.Config {
	timeout := time.Duration(req.GetTimeout())
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	config := &models.Config{
		Name:        req.GetName(),
		PathPrefix:  req.GetPathPrefix(),
		TargetURL:   req.GetTargetUrl(),
		StripPrefix: req.GetStripPrefix(),
		Middleware:  req.GetMiddleware(),
		Timeout:     timeout,
	}

	if req.GetAuthentication() != nil {
		config.Authentication = protoAuthToModel(req.GetAuthentication())
	}

	return config
}

// updateProtoToModel converts an UpdateConfigRequest to a Config model
func updateProtoToModel(req *opengate_v1.UpdateConfigRequest) *models.Config {
	timeout := time.Duration(req.GetTimeout())
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	config := &models.Config{
		ID:          req.GetId(),
		Name:        req.GetName(),
		PathPrefix:  req.GetPathPrefix(),
		TargetURL:   req.GetTargetUrl(),
		StripPrefix: req.GetStripPrefix(),
		Middleware:  req.GetMiddleware(),
		Timeout:     timeout,
	}

	if req.GetAuthentication() != nil {
		config.Authentication = protoAuthToModel(req.GetAuthentication())
	}

	return config
}

// modelToProto converts a Config model to a proto Config
func modelToProto(config *models.Config) *opengate_v1.Config {
	if config == nil {
		return nil
	}

	protoConfig := &opengate_v1.Config{
		Id:          config.ID,
		Name:        config.Name,
		PathPrefix:  config.PathPrefix,
		TargetUrl:   config.TargetURL,
		StripPrefix: config.StripPrefix,
		Middleware:  config.Middleware,
		Timeout:     int64(config.Timeout),
		CreatedAt:   config.CreatedAt.Unix(),
		UpdatedAt:   config.UpdatedAt.Unix(),
	}

	if config.Authentication != nil {
		protoConfig.Authentication = modelAuthToProto(config.Authentication)
	}

	return protoConfig
}

// serviceRouteToProto converts a ServiceRoute to a proto Route
func serviceRouteToProto(route *models.ServiceRoute) *opengate_v1.Route {
	if route == nil {
		return nil
	}

	protoRoute := &opengate_v1.Route{
		Name:        route.Name,
		PathPrefix:  route.PathPrefix,
		TargetUrl:   route.TargetURL,
		StripPrefix: route.StripPrefix,
		Middleware:  route.Middleware,
		Timeout:     int64(route.Timeout),
		UpdatedAt:   route.UpdatedAt,
	}

	if route.Authentication != nil {
		protoRoute.Authentication = modelAuthToProto(route.Authentication)
	}

	return protoRoute
}

// protoAuthToModel converts proto Authentication to model Authentication
func protoAuthToModel(auth *opengate_v1.Authentication) *models.Authentication {
	if auth == nil {
		return nil
	}

	modelAuth := &models.Authentication{
		Required: auth.GetRequired(),
	}

	for _, except := range auth.GetExcept() {
		modelAuth.Except = append(modelAuth.Except, struct {
			Path    string   `json:"path" yaml:"Path"`
			Methods []string `json:"methods" yaml:"Methods"`
		}{
			Path:    except.GetPath(),
			Methods: except.GetMethods(),
		})
	}

	return modelAuth
}

// modelAuthToProto converts model Authentication to proto Authentication
func modelAuthToProto(auth *models.Authentication) *opengate_v1.Authentication {
	if auth == nil {
		return nil
	}

	protoAuth := &opengate_v1.Authentication{
		Required: auth.Required,
	}

	for _, except := range auth.Except {
		protoAuth.Except = append(protoAuth.Except, &opengate_v1.AuthenticationException{
			Path:    except.Path,
			Methods: except.Methods,
		})
	}

	return protoAuth
}
