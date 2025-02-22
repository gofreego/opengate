package dao

// MatchConfig holds matching rules for a route
type MatchConfig struct {
	Host    string   `yaml:"host,omitempty" bson:"host,omitempty"`
	Prefix  string   `yaml:"prefix,omitempty" bson:"prefix,omitempty"`
	Regex   string   `yaml:"regex,omitempty" bson:"regex,omitempty"`
	Methods []string `yaml:"methods,omitempty" bson:"methods,omitempty"`
}

// RouteConfig defines a single API Gateway route
type RouteConfig struct {
	ID          string      `yaml:"id" bson:"_id"`
	Name        string      `yaml:"name" bson:"name"`
	Description string      `yaml:"description" bson:"description"`
	Match       MatchConfig `yaml:"match" bson:"match"`
	Target      string      `yaml:"target" bson:"target"`
}
