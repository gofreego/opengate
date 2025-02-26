package authentication

import (
	"api-gateway/internal/models/dao"
	"api-gateway/internal/service/authentication/jwt"
	"context"

	"github.com/gin-gonic/gin"
)

type Service interface {
	// authenticate the user and add the user details to the context like user id, roles, permissions etc and returns userClaims.
	Authenticate(ctx *gin.Context) (*dao.UserClaims, error)
}

type AuthType string

const (
	AUTH_TYPE_JWT   AuthType = "jwt"
	AUTH_TYPE_OAUTH AuthType = "oauth" // OAuth 2.0
)

type Config struct {
	Type AuthType   `yaml:"type" bson:"type"`
	JWT  jwt.Config `yaml:"jwt" bson:"jwt"`
}

func NewService(ctx context.Context, cfg *Config) Service {
	switch cfg.Type {
	case AUTH_TYPE_JWT:
		return jwt.NewJWTAuthenticationService(ctx, &cfg.JWT)
	case AUTH_TYPE_OAUTH:
		panic("OAuth Authentication Service not implemented")
	default:
		panic("invalid authentication type, valid values are jwt, oauth")
	}
}
