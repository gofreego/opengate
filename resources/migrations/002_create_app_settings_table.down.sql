-- Migration: Drop app_settings table
-- Version: 002
-- Description: Drops the app_settings table and associated objects

DROP TRIGGER IF EXISTS update_app_settings_updated_at ON app_settings;

DROP TABLE IF EXISTS app_settings;
