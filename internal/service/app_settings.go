package service

import (
	"context"
	"encoding/json"

	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/pkg/utils"
)

// GetCORSConfig returns the live CORS configuration for use by middleware.
func (s *Service) GetCORSConfig() *utils.CORSConfig {
	return s.settingsMgr.GetCORSConfig()
}

// GetAppSettings returns all settings as a key → raw-JSON map.
func (s *Service) GetAppSettings() map[string]json.RawMessage {
	return s.settingsMgr.GetAll()
}

// UpsertAppSetting persists a setting and immediately refreshes the in-memory cache.
func (s *Service) UpsertAppSetting(ctx context.Context, key string, value json.RawMessage) error {
	setting := &models.AppSetting{Key: key, Value: string(value)}
	if err := s.repo.UpsertAppSetting(ctx, setting); err != nil {
		return err
	}
	s.settingsMgr.Refresh(ctx)
	return nil
}
