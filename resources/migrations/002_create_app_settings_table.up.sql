-- Migration: Create app_settings table
-- Version: 002
-- Description: Creates the app_settings table for storing key-value gateway configurations

CREATE TABLE IF NOT EXISTS app_settings (
    key        TEXT PRIMARY KEY,
    value      TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_app_settings_updated_at
    BEFORE UPDATE ON app_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE app_settings IS 'Stores key-value gateway configuration entries (e.g. CORS policy)';
COMMENT ON COLUMN app_settings.key IS 'Unique setting key (e.g. cors_config)';
COMMENT ON COLUMN app_settings.value IS 'JSON-encoded setting value';
COMMENT ON COLUMN app_settings.updated_at IS 'Timestamp of the last update';
