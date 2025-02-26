package middlewares

import (
	"api-gateway/internal/service/authentication"

	"github.com/gin-gonic/gin"
)

func (s *MiddlewareService) GetAuthenticationMiddleware(authService authentication.Service) (Middleware, error) {
	return func(ctx *gin.Context) error {
		_, err := authService.Authenticate(ctx)
		if err != nil {
			return err
		}
		return nil
	}, nil
}
