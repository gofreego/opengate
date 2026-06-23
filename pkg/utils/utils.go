package utils

import (
	"net/http"
	"strconv"
)

// CORSConfig holds the CORS policy read from the settings store.
type CORSConfig struct {
	Enabled        bool   `json:"enabled"`
	AllowedOrigins string `json:"allowedOrigins"`
	AllowedMethods string `json:"allowedMethods"`
	AllowedHeaders string `json:"allowedHeaders"`
	MaxAge         int    `json:"maxAge"`
}

// DefaultCORSConfig returns a permissive default that allows all origins.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:        true,
		AllowedOrigins: "*",
		AllowedMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowedHeaders: "Accept, Authorization, Content-Type, X-CSRF-Token, X-User-Id, X-User-Perms",
		MaxAge:         3600,
	}
}

// CorsMiddleware adds CORS headers using a dynamic config provider so the
// policy can be changed at runtime without restarting the server.
func CorsMiddleware(next http.Handler, getConfig func() *CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := getConfig()
		if cfg == nil {
			cfg = DefaultCORSConfig()
		}

		if cfg.Enabled {
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", cfg.AllowedOrigins)
				w.Header().Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
