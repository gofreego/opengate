package models

import "time"

// Config represents a route configuration stored in the database
type Config struct {
	ID             int64           `json:"id"`
	Name           string          `json:"name"`
	PathPrefix     string          `json:"pathPrefix"`
	TargetURL      string          `json:"targetURL"`
	StripPrefix    bool            `json:"stripPrefix"`
	Authentication *Authentication `json:"authentication"`
	Middleware     []string        `json:"middleware"`
	Timeout        time.Duration   `json:"timeout"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// ToServiceRoute converts a Config to a ServiceRoute
func (c *Config) ToServiceRoute() *ServiceRoute {
	return &ServiceRoute{
		Name:           c.Name,
		PathPrefix:     c.PathPrefix,
		TargetURL:      c.TargetURL,
		StripPrefix:    c.StripPrefix,
		Authentication: c.Authentication,
		Middleware:     c.Middleware,
		Timeout:        c.Timeout,
		UpdatedAt:      c.UpdatedAt.UnixMilli(),
	}
}

// ConfigFilter represents filter options for listing configs
type ConfigFilter struct {
	Search string
	Limit  int
	Offset int
}
