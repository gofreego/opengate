package service

import (
	"api-gateway/internal/repository"
	"api-gateway/internal/service/authentication"
	"api-gateway/internal/service/match"
	"api-gateway/internal/service/middlewares"
	"context"

	"github.com/gofreego/goutils/logger"
)

type Config struct {
	Authentication authentication.Config
	Middlewares    middlewares.Config
}

type Service struct {
	match       *match.MatchService
	middlewares *middlewares.MiddlewareService

	repo repository.Repository
}

func NewService(ctx context.Context, cfg *Config, repo repository.Repository) *Service {

	matchService := match.NewMatchService()
	middlewareService := middlewares.NewMiddlewareService(&cfg.Middlewares)

	routeConfigs, err := repo.GetRoutesConfig(ctx)
	if err != nil {
		logger.Panic(ctx, "failed to fetch route configs: Err:%v", err.Error())
		return nil
	}
	err = matchService.UpdateRoutesMatch(ctx, routeConfigs...)
	if err != nil {
		logger.Panic(ctx, "failed to update match, Err:%s", err.Error())
		return nil
	}
	err = middlewareService.UpdateMiddlewares(ctx, routeConfigs...)
	if err != nil {
		logger.Panic(ctx, "failed to update middlewares, Err:%s", err.Error())
	}
	srv := &Service{
		match:       matchService,
		middlewares: middlewareService,
		repo:        repo,
	}

	// Watch for route config changes
	go srv.WatchRoutesConfigChanges(ctx)
	return srv
}

func (s *Service) WatchRoutesConfigChanges(ctx context.Context) {
	routeConfigChan, err := s.repo.WatchRoutesConfigChanges(ctx)
	if err != nil {
		logger.Panic(ctx, "failed to watch route config changes: Err:%v", err.Error())
		return
	}
	for routeConfig := range routeConfigChan {

		err := s.match.UpdateRoutesMatch(ctx, routeConfig)
		if err != nil {
			logger.Error(ctx, "failed to update match: Err:%v", err.Error())
		}
		err = s.middlewares.UpdateMiddlewares(ctx, routeConfig)
		if err != nil {
			logger.Error(ctx, "failed to update middleware: Err:%v", err.Error())
		}

	}
}
