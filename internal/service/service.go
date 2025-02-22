package service

import (
	"api-gateway/internal/repository"
	"api-gateway/internal/service/match"
	"api-gateway/internal/service/middlewares"
	"context"

	"github.com/gofreego/goutils/logger"
)

type Config struct {
	Middlewares middlewares.Config
}

type Service struct {
	match       *match.MatchService
	middlewares *middlewares.MiddlewareService
	repo        repository.Repository
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
	return &Service{
		match:       matchService,
		middlewares: middlewareService,
		repo:        repo,
	}
}
