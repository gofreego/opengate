package jwt

import (
	"api-gateway/internal/customerrors"
	"api-gateway/internal/models/dao"
	"api-gateway/internal/service/authentication/jwt/validator"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimConfig struct {
	UserIDPath      string `yaml:"user_id_path" bson:"userIdPath"`
	PermissionsPath string `yaml:"permissions_path" bson:"permissionsPath"`
	RolesPath       string `yaml:"roles_path" bson:"rolesPath"`
}

type Config struct {
	Validator validator.Config `yaml:"validator" bson:"validator"`
	Claims    ClaimConfig      `yaml:"claims" bson:"claims"`
}

type JWTAuthenticationService struct {
	cfg       *Config
	validator validator.Validator
}

func NewJWTAuthenticationService(ctx context.Context, config *Config) *JWTAuthenticationService {
	return &JWTAuthenticationService{
		cfg:       config,
		validator: validator.NewValidator(ctx, &config.Validator),
	}
}

const (
	AUTHORIZATION_HEADER = "Authorization"
	USER_ID_HEADER       = "user_id"
	PERMISSIONS_HEADER   = "permissions"
	ROLES_HEADER         = "roles"
)

func (s *JWTAuthenticationService) Authenticate(ctx *gin.Context) (*dao.UserClaims, error) {
	token := ctx.GetHeader(AUTHORIZATION_HEADER)
	if token == "" {
		return nil, customerrors.ErrNoJWTToken
	}
	claims, err := s.validator.Validate(token)
	if err != nil {
		return nil, err
	}
	userClaims, err := s.getUserClaims(ctx, claims)
	if err != nil {
		return nil, err
	}
	ctx.Set(USER_ID_HEADER, userClaims.UserID)
	ctx.Set(PERMISSIONS_HEADER, strings.Join(userClaims.Permissions, ","))
	ctx.Set(ROLES_HEADER, strings.Join(userClaims.Roles, ","))
	return userClaims, nil
}

func (s *JWTAuthenticationService) getUserClaims(_ context.Context, claims jwt.MapClaims) (*dao.UserClaims, error) {
	var userClaims dao.UserClaims
	userClaims.UserID = claims[s.cfg.Claims.UserIDPath].(string)
	userClaims.Permissions = claims[s.cfg.Claims.PermissionsPath].([]string)
	userClaims.Roles = claims[s.cfg.Claims.RolesPath].([]string)
	return &userClaims, nil
}
