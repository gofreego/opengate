package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	sqlutils "github.com/gofreego/goutils/databases/connections/sql"
	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/opengate/internal/models"
)

// Repository implements the service.Repository interface using PostgreSQL
type Repository struct {
	connManager sqlutils.DBManager
}

// NewRepository creates a new PostgreSQL repository instance
func NewRepository(ctx context.Context, cfg *sqlutils.Config) (*Repository, error) {
	connManager, err := sqlutils.NewDBManager(cfg)
	if err != nil {
		return nil, err
	}
	repo := &Repository{connManager: connManager}
	if err := repo.initAppSettingsTable(ctx); err != nil {
		return nil, fmt.Errorf("failed to init app_settings table: %w", err)
	}
	return repo, nil
}

func (r *Repository) initAppSettingsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS app_settings (
			key        TEXT PRIMARY KEY,
			value      TEXT NOT NULL DEFAULT '',
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err := r.connManager.Primary().ExecContext(ctx, query)
	return err
}

// Ping checks the database connection
func (r *Repository) Ping(ctx context.Context) error {
	return r.connManager.Primary().Ping()
}

// GetRoutes retrieves all routes for the routing manager
func (r *Repository) GetRoutes(ctx context.Context) ([]*models.ServiceRoute, error) {
	query := `
		SELECT id, name, path_prefix, target_url, strip_prefix, 
		       authentication, middleware, timeout, created_at, updated_at
		FROM configs
		ORDER BY name
	`

	rows, err := r.connManager.Primary().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}
	defer rows.Close()

	var routes []*models.ServiceRoute
	for rows.Next() {
		config, err := r.scanConfig(rows)
		if err != nil {
			logger.Error(ctx, "Failed to scan config row: %v", err)
			continue
		}
		routes = append(routes, config.ToServiceRoute())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating config rows: %w", err)
	}

	return routes, nil
}

// CreateConfig creates a new config in the database
func (r *Repository) CreateConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	authJSON, err := json.Marshal(config.Authentication)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authentication: %w", err)
	}

	middlewareJSON, err := json.Marshal(config.Middleware)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware: %w", err)
	}

	query := `
		INSERT INTO configs (name, path_prefix, target_url, strip_prefix, authentication, middleware, timeout)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	var id int64
	var createdAt, updatedAt time.Time

	err = r.connManager.Primary().QueryRowContext(ctx, query,
		config.Name,
		config.PathPrefix,
		config.TargetURL,
		config.StripPrefix,
		authJSON,
		middlewareJSON,
		config.Timeout,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		// Check for unique constraint violation
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("config with name '%s' already exists", config.Name)
		}
		return nil, fmt.Errorf("failed to insert config: %w", err)
	}

	config.ID = id
	config.CreatedAt = createdAt
	config.UpdatedAt = updatedAt

	return config, nil
}

// GetConfigByID retrieves a config by its ID
func (r *Repository) GetConfigByID(ctx context.Context, id int64) (*models.Config, error) {
	query := `
		SELECT id, name, path_prefix, target_url, strip_prefix, 
		       authentication, middleware, timeout, created_at, updated_at
		FROM configs
		WHERE id = $1
	`

	row := r.connManager.Primary().QueryRowContext(ctx, query, id)
	config, err := r.scanConfigRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("config with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}

// ListConfigs retrieves configs with pagination and optional search
func (r *Repository) ListConfigs(ctx context.Context, filter *models.ConfigFilter) ([]*models.Config, int, error) {
	// Build the query
	baseQuery := `FROM configs`
	var args []interface{}
	argIndex := 1

	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" WHERE name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.connManager.Primary().QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count configs: %w", err)
	}

	// Get paginated results
	selectQuery := fmt.Sprintf(`
		SELECT id, name, path_prefix, target_url, strip_prefix, 
		       authentication, middleware, timeout, created_at, updated_at
		%s
		ORDER BY name
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.connManager.Primary().QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.Config
	for rows.Next() {
		config, err := r.scanConfig(rows)
		if err != nil {
			logger.Error(ctx, "Failed to scan config row: %v", err)
			continue
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating config rows: %w", err)
	}

	return configs, total, nil
}

// UpdateConfig updates an existing config
func (r *Repository) UpdateConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	authJSON, err := json.Marshal(config.Authentication)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authentication: %w", err)
	}

	middlewareJSON, err := json.Marshal(config.Middleware)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware: %w", err)
	}

	query := `
		UPDATE configs
		SET name = $1, path_prefix = $2, target_url = $3, strip_prefix = $4,
		    authentication = $5, middleware = $6, timeout = $7
		WHERE id = $8
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err = r.connManager.Primary().QueryRowContext(ctx, query,
		config.Name,
		config.PathPrefix,
		config.TargetURL,
		config.StripPrefix,
		authJSON,
		middlewareJSON,
		config.Timeout,
		config.ID,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("config with id %d not found", config.ID)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("config with name '%s' already exists", config.Name)
		}
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	config.CreatedAt = createdAt
	config.UpdatedAt = updatedAt

	return config, nil
}

// DeleteConfig deletes a config by ID
func (r *Repository) DeleteConfig(ctx context.Context, id int64) error {
	query := `DELETE FROM configs WHERE id = $1`

	result, err := r.connManager.Primary().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("config with id %d not found", id)
	}

	return nil
}

// scanConfig scans a row from sql.Rows into a Config struct
func (r *Repository) scanConfig(rows *sql.Rows) (*models.Config, error) {
	var config models.Config
	var authJSON, middlewareJSON []byte
	var timeout int64

	err := rows.Scan(
		&config.ID,
		&config.Name,
		&config.PathPrefix,
		&config.TargetURL,
		&config.StripPrefix,
		&authJSON,
		&middlewareJSON,
		&timeout,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	config.Timeout = time.Duration(timeout)

	if len(authJSON) > 0 {
		if err := json.Unmarshal(authJSON, &config.Authentication); err != nil {
			return nil, fmt.Errorf("failed to unmarshal authentication: %w", err)
		}
	}

	if len(middlewareJSON) > 0 {
		if err := json.Unmarshal(middlewareJSON, &config.Middleware); err != nil {
			return nil, fmt.Errorf("failed to unmarshal middleware: %w", err)
		}
	}

	return &config, nil
}

// scanConfigRow scans a single row from sql.Row into a Config struct
func (r *Repository) scanConfigRow(row *sql.Row) (*models.Config, error) {
	var config models.Config
	var authJSON, middlewareJSON []byte
	var timeout int64

	err := row.Scan(
		&config.ID,
		&config.Name,
		&config.PathPrefix,
		&config.TargetURL,
		&config.StripPrefix,
		&authJSON,
		&middlewareJSON,
		&timeout,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	config.Timeout = time.Duration(timeout)

	if len(authJSON) > 0 {
		if err := json.Unmarshal(authJSON, &config.Authentication); err != nil {
			return nil, fmt.Errorf("failed to unmarshal authentication: %w", err)
		}
	}

	if len(middlewareJSON) > 0 {
		if err := json.Unmarshal(middlewareJSON, &config.Middleware); err != nil {
			return nil, fmt.Errorf("failed to unmarshal middleware: %w", err)
		}
	}

	return &config, nil
}

// GetAppSettings retrieves all app settings
func (r *Repository) GetAppSettings(ctx context.Context) ([]*models.AppSetting, error) {
	query := `SELECT key, value, updated_at FROM app_settings ORDER BY key`
	rows, err := r.connManager.Primary().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query app_settings: %w", err)
	}
	defer rows.Close()

	var settings []*models.AppSetting
	for rows.Next() {
		var s models.AppSetting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan app_setting: %w", err)
		}
		settings = append(settings, &s)
	}
	return settings, rows.Err()
}

// UpsertAppSetting inserts or updates a single app setting
func (r *Repository) UpsertAppSetting(ctx context.Context, setting *models.AppSetting) error {
	query := `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) DO UPDATE
			SET value = EXCLUDED.value, updated_at = NOW()
	`
	_, err := r.connManager.Primary().ExecContext(ctx, query, setting.Key, setting.Value)
	if err != nil {
		return fmt.Errorf("failed to upsert app_setting: %w", err)
	}
	return nil
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation
func isUniqueViolation(err error) bool {
	// PostgreSQL unique violation error code is 23505
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "unique constraint") || strings.Contains(errMsg, "duplicate key")
}
