package dao

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofreego/goutils/customerrors"
	"gopkg.in/yaml.v3"
)

type RouteAuthConfig struct {
	AuthType string   `yaml:"auth_type" bson:"authType"`
	Required bool     `yaml:"required" bson:"required"`
	Roles    []string `yaml:"roles" bson:"roles,omitempty"`
}

// MatchConfig holds matching rules for a route
type MatchConfig struct {
	Host    string   `yaml:"host" bson:"host,omitempty"`
	Prefix  string   `yaml:"prefix" bson:"prefix,omitempty"`
	Regex   string   `yaml:"regex" bson:"regex,omitempty"`
	Methods []string `yaml:"methods" bson:"methods,omitempty"`
}

var validMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
	http.MethodHead:    true,
	http.MethodConnect: true,
	http.MethodTrace:   true,
}

func (m *MatchConfig) Validate() error {
	if m.Host == "" && m.Prefix == "" && m.Regex == "" {
		return customerrors.BAD_REQUEST_ERROR("at least one of host, prefix or regex is required")
	}
	if len(m.Methods) >= 0 {
		for _, method := range m.Methods {
			if _, ok := validMethods[method]; !ok {
				return customerrors.BAD_REQUEST_ERROR("invalid method: %s", method)
			}
		}
	}
	return nil
}

// RouteConfig defines a single API Gateway route
type RouteConfig struct {
	ID          string          `yaml:"id" bson:"_id"`
	Name        string          `yaml:"name" bson:"name"`
	Description string          `yaml:"description" bson:"description"`
	Match       MatchConfig     `yaml:"match" bson:"match"`
	Auth        RouteAuthConfig `yaml:"auth" bson:"auth"`
	Target      string          `yaml:"target" bson:"target"`
}

func (r *RouteConfig) Validate() error {
	if r.ID == "" {
		return customerrors.BAD_REQUEST_ERROR("ID is required")
	}

	if r.Name == "" {
		return customerrors.BAD_REQUEST_ERROR("name is required")
	}

	// check if target is a valid URL
	_, err := url.Parse(r.Target)
	if err != nil {
		return customerrors.BAD_REQUEST_ERROR("invalid target URL")
	}
	err = r.Match.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *RouteConfig) String(isJson bool) string {
	if isJson {
		bytes, _ := json.Marshal(r)
		return string(bytes)
	}
	bytes, _ := yaml.Marshal(r)
	return string(bytes)
}
