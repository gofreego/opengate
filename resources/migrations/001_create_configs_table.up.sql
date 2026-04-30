-- Migration: Create configs table
-- Version: 001
-- Description: Creates the configs table for storing service route configurations

CREATE TABLE IF NOT EXISTS configs (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL UNIQUE,
    path_prefix     VARCHAR(512) NOT NULL,
    target_url      VARCHAR(512) NOT NULL,
    strip_prefix    BOOLEAN NOT NULL DEFAULT false,
    authentication  JSONB,
    middleware      JSONB DEFAULT '[]'::jsonb,
    timeout         BIGINT NOT NULL DEFAULT 30000000000, -- Default 30 seconds in nanoseconds
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_configs_name ON configs(name);

-- Create index on path_prefix for route matching
CREATE INDEX IF NOT EXISTS idx_configs_path_prefix ON configs(path_prefix);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_configs_updated_at
    BEFORE UPDATE ON configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comment to table
COMMENT ON TABLE configs IS 'Stores service route configurations for the API gateway';
COMMENT ON COLUMN configs.name IS 'Unique identifier name for the route configuration';
COMMENT ON COLUMN configs.path_prefix IS 'URL path prefix to match for this route';
COMMENT ON COLUMN configs.target_url IS 'Target URL to proxy requests to';
COMMENT ON COLUMN configs.strip_prefix IS 'Whether to strip the path prefix when forwarding';
COMMENT ON COLUMN configs.authentication IS 'JSON object containing authentication settings';
COMMENT ON COLUMN configs.middleware IS 'JSON array of middleware names to apply';
COMMENT ON COLUMN configs.timeout IS 'Request timeout in nanoseconds';
