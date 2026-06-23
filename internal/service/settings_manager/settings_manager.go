package settingsmanager

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/pkg/utils"
)

const (
	KeyCORSConfig = "cors_config"

	defaultRefreshInterval = 30 * time.Second
)

type Repository interface {
	GetAppSettings(ctx context.Context) ([]*models.AppSetting, error)
	UpsertAppSetting(ctx context.Context, setting *models.AppSetting) error
}

type Config struct {
	RefreshInterval time.Duration `yaml:"RefreshInterval"`
}

type Manager struct {
	repo     Repository
	cfg      *Config
	mu       sync.RWMutex
	settings map[string]string // key -> raw JSON value
}

func New(repo Repository, cfg *Config) *Manager {
	return &Manager{
		repo:     repo,
		cfg:      cfg,
		settings: make(map[string]string),
	}
}

// Start loads settings and begins periodic refresh.
func (m *Manager) Start(ctx context.Context) {
	logger.Info(ctx, "Settings manager started")

	if err := m.load(ctx); err != nil {
		logger.Error(ctx, "Failed to load initial settings: %v", err)
	}
	// Seed defaults for any missing keys
	m.seedDefaults(ctx)

	interval := defaultRefreshInterval
	if m.cfg != nil && m.cfg.RefreshInterval > 0 {
		interval = m.cfg.RefreshInterval
	}

	logger.Info(ctx, "Settings manager refresh interval: %v", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "Settings manager stopped")
			return
		case <-ticker.C:
			if err := m.load(ctx); err != nil {
				logger.Error(ctx, "Failed to refresh settings: %v", err)
			}
		}
	}
}

// Refresh triggers an immediate reload from the repository.
func (m *Manager) Refresh(ctx context.Context) {
	if err := m.load(ctx); err != nil {
		logger.Error(ctx, "Failed to refresh settings: %v", err)
	}
}

// GetCORSConfig returns the current CORS configuration.
func (m *Manager) GetCORSConfig() *utils.CORSConfig {
	m.mu.RLock()
	raw, ok := m.settings[KeyCORSConfig]
	m.mu.RUnlock()

	if !ok || raw == "" {
		return utils.DefaultCORSConfig()
	}

	var cfg utils.CORSConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return utils.DefaultCORSConfig()
	}
	return &cfg
}

// GetAll returns all settings as a map of key -> parsed JSON value.
func (m *Manager) GetAll() map[string]json.RawMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]json.RawMessage, len(m.settings))
	for k, v := range m.settings {
		result[k] = json.RawMessage(v)
	}
	return result
}

func (m *Manager) load(ctx context.Context) error {
	settings, err := m.repo.GetAppSettings(ctx)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range settings {
		m.settings[s.Key] = s.Value
	}
	return nil
}

func (m *Manager) seedDefaults(ctx context.Context) {
	m.mu.RLock()
	_, hasCORS := m.settings[KeyCORSConfig]
	m.mu.RUnlock()

	if !hasCORS {
		defaultCORS := utils.DefaultCORSConfig()
		raw, _ := json.Marshal(defaultCORS)
		setting := &models.AppSetting{Key: KeyCORSConfig, Value: string(raw)}
		if err := m.repo.UpsertAppSetting(ctx, setting); err != nil {
			logger.Error(ctx, "Failed to seed default CORS config: %v", err)
			return
		}
		m.mu.Lock()
		m.settings[KeyCORSConfig] = string(raw)
		m.mu.Unlock()
		logger.Info(ctx, "Seeded default CORS config")
	}
}
