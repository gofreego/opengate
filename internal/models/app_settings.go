package models

import "time"

// AppSetting stores a key-value configuration entry
type AppSetting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"` // JSON-encoded value
	UpdatedAt time.Time `json:"updatedAt"`
}
