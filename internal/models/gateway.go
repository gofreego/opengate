package models

import (
	"time"
)

// ServiceRoute represents a service route configuration
type ServiceRoute struct {
	Name           string          `json:"name" yaml:"Name"`
	PathPrefix     string          `json:"pathPrefix" yaml:"PathPrefix"`
	TargetURL      string          `json:"targetURL" yaml:"TargetURL"`
	StripPrefix    bool            `json:"stripPrefix" yaml:"StripPrefix"`
	Authentication *Authentication `json:"authentication" yaml:"Authentication"`
	Middleware     []string        `json:"middleware" yaml:"Middleware"`
	Timeout        time.Duration   `json:"timeout" yaml:"Timeout"`
	UpdatedAt      int64           `json:"-" yaml:"-"` // Unix timestamp of last update
}

type Authentication struct {
	Required bool `json:"required" yaml:"Required"`
	// if required is true, then Excepted path and methods does not require authentication
	// if required is false, then Excepted path and methods require authentication
	Except []struct {
		Path    string   `json:"path" yaml:"Path"`
		Methods []string `json:"methods" yaml:"Methods"`
	} `json:"except" yaml:"Except"`
}

func (auth *Authentication) IsAuthenticationRequired(path, method string) bool {
	if auth == nil {
		return false
	}
	if auth.Required {
		return !auth.isExcepted(path, method)
	}
	return auth.isExcepted(path, method)
}

func (auth *Authentication) isExcepted(path, method string) bool {
	for _, except := range auth.Except {
		if except.Path == path {
			if len(except.Methods) == 0 {
				return true
			}
			for _, m := range except.Methods {
				if m == method {
					return true
				}
			}
		}
	}
	return false
}
