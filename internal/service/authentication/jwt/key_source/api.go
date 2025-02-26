package keysources

type JWTAPIConfig struct {
	URL      string `yaml:"url" bson:"url"`
	Method   string `yaml:"method" bson:"method"`
	BodyPath string `yaml:"body_path" bson:"bodyPath"` // JSON path to the JWT key
}
