package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/cache"
	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/openauth/api/openauth_v1"
	"github.com/gofreego/openauth/pkg/clients/openauth"
	"github.com/gofreego/openauth/pkg/jwtutils"
	"github.com/gofreego/opengate/internal/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type OpenAuthStrategy struct {
	client openauth_v1.OpenAuthClient
	conn   *grpc.ClientConn
	cache  cache.Cache
}

func NewOpenAuthStrategy(ctx context.Context, config *openauth.ClientConfig, cache cache.Cache) (Strategy, error) {
	client, conn, err := openauth.NewOpenAuthClientV1(ctx, config)
	if err != nil {
		return nil, err
	}
	return &OpenAuthStrategy{
		client: client,
		conn:   conn,
		cache:  cache,
	}, nil
}

func (s *OpenAuthStrategy) Authenticate(ctx *gin.Context) error {
	reqContext := ctx.Request.Context()
	token := ctx.GetHeader(constants.HEADER_AUTHORIZATION)
	if authenticated, err := s.isAuthenticatedInCache(token); err == nil && authenticated {
		// Even for cached auth, we need claims for headers
		_, claims, err := s.getJWTDetails(token)
		if err != nil {
			logger.Error(ctx, "Failed to get JWT details from cache: %v", err)
			return err
		}
		ctx.Set(constants.JWT_CLAIMS, claims)
		return nil
	}
	expiresAt, claims, err := s.getJWTDetails(token)
	if err != nil {
		logger.Error(ctx, "Failed to get JWT details: %v", err)
		return err
	}

	if expiresAt.Before(time.Now()) {
		return fmt.Errorf("token is expired")
	}

	// Store claims in gin context
	ctx.Set(constants.JWT_CLAIMS, claims)

	// Set additional headers for the authentication request
	authRequest := &openauth_v1.IsAuthenticatedRequest{
		AccessToken: ctx.GetHeader(constants.HEADER_AUTHORIZATION),
	}

	reqContext = metadata.AppendToOutgoingContext(reqContext, constants.HEADER_AUTHORIZATION, ctx.GetHeader(constants.HEADER_AUTHORIZATION))

	_, err = s.client.IsAuthenticated(reqContext, authRequest)
	if err != nil {
		s.setCache(ctx.GetHeader(constants.HEADER_AUTHORIZATION), false, time.Minute)
		logger.Error(ctx, "Authentication error: %v", err)
		return err
	}
	s.setCache(ctx.GetHeader(constants.HEADER_AUTHORIZATION), true, time.Until(*expiresAt))
	return nil
}

func (s *OpenAuthStrategy) isAuthenticatedInCache(authToken string) (bool, error) {
	var isAuthenticated bool
	if s.cache != nil {
		err := s.cache.GetV(context.Background(), authToken, &isAuthenticated)
		if err == nil {
			return isAuthenticated, nil
		}
		if err.Error() != "redis: nil" && err.Error() != "not found" {
			logger.Error(context.Background(), "Cache get error: %v", err)
			return false, err
		}
	}
	return false, fmt.Errorf("not found in cache")
}

func (s *OpenAuthStrategy) setCache(authToken string, isAuthenticated bool, duration time.Duration) {
	if s.cache != nil {
		err := s.cache.SetWithTimeout(context.Background(), authToken, isAuthenticated, duration)
		if err != nil {
			logger.Error(context.Background(), "Cache set error: %v", err)
		}
	}
}

func (s *OpenAuthStrategy) getJWTDetails(token string) (*time.Time, *jwtutils.JWTClaims, error) {
	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	// JWT tokens have 3 parts separated by dots: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid JWT token format")
	}

	// Decode the payload (second part)
	payload := parts[1]

	// Add padding if necessary for base64 decoding
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// Decode base64
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse JSON payload
	var claims jwtutils.JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JWT claims: %w", err)
	}
	return &claims.RegisteredClaims.ExpiresAt.Time, &claims, nil
}

func (s *OpenAuthStrategy) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
