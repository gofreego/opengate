package dao

type JWTKeySource string

const (

	// JWT Key Types to fetch the signing key from different sources
	JWT_KEY_SOURCE_ENV            JWTKeySource = "env"            // fetch from an environment variable
	JWT_KEY_SOURCE_API            JWTKeySource = "api"            // fetch from an API
	JWT_KEY_SOURCE_SECRET_MANAGER JWTKeySource = "secret_manager" // fetch from a secret manager
	//this method is not recommended for production use
	JWT_KEY_SOURCE_CONFIG JWTKeySource = "config" // fetch from a configuration file
	// Environment variable names
	SIGNING_KEY_ENV_NAME = "JWT_SIGNING_KEY"
)

type JWTKeySourceConfig struct {
	KeySource        JWTKeySource `yaml:"key_source" bson:"keySource"`
	SigningKey       string       `yaml:"signing_key" bson:"signingKey"`
	API              JWTAPIConfig `yaml:"api" bson:"api"`
	SecretManagerKey string       `yaml:"secret_manager" bson:"secretManager"`
}

type JWTAPIConfig struct {
	URL      string `yaml:"url" bson:"url"`
	Method   string `yaml:"method" bson:"method"`
	BodyPath string `yaml:"body_path" bson:"bodyPath"` // JSON path to the JWT key
}
