package validator

import (
	keysources "api-gateway/internal/service/authentication/jwt/key_source"
	"api-gateway/internal/service/authentication/jwt/validator/hs256"
	"context"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	UpdateSecretKey(secretKey string)
	Validate(token string) (jwt.MapClaims, error)
}

type Config struct {
	ValidatorType string `yaml:"validator_type" bson:"validatorType"`
	// KeyRefreshInterval to refresh the key from the source in seconds
	KeyRefreshInterval int               `yaml:"refresh_interval" bson:"refreshInterval"`
	KeySource          keysources.Config `yaml:"key_source" bson:"keySource"`
}

func NewValidator(ctx context.Context, cfg *Config) Validator {
	keySource := keysources.NewKeySource(&cfg.KeySource)
	secretKey, err := keySource.GetSecretKey()
	if err != nil {
		logger.Panic(ctx, "failed to get secret key")
	}
	var validator Validator
	switch cfg.ValidatorType {
	case "HS256":
		validator = hs256.NewHS256Validator(secretKey)
	default:
		logger.Panic(ctx, "invalid jwt validator type, using default HS256")
	}
	go refreshSecretKey(ctx, cfg.KeyRefreshInterval, keySource, validator)
	return validator
}

// refreshSecretKey refreshes the secret key from the source at the given interval using timer
func refreshSecretKey(ctx context.Context, refreshInterval int, keySource keysources.KeySource, validator Validator) {
	logger.Debug(ctx, "starting secret key refresh timer")
	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
	for range ticker.C {
		logger.Debug(ctx, "refreshing secret key")
		secretKey, err := keySource.GetSecretKey()
		if err != nil {
			logger.Error(ctx, "failed to get secret key from source")
			continue
		}
		validator.UpdateSecretKey(secretKey)
	}
}
