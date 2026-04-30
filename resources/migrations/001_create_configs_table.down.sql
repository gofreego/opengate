-- Migration: Drop configs table
-- Version: 001
-- Description: Drops the configs table and associated objects

-- Drop trigger first
DROP TRIGGER IF EXISTS update_configs_updated_at ON configs;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (will be dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_configs_path_prefix;
DROP INDEX IF EXISTS idx_configs_name;

-- Drop table
DROP TABLE IF EXISTS configs;
