package middlewares

import (
	"api-gateway/internal/service/authentication"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/customerrors"
)

func (s *MiddlewareService) GetAuthorizationMiddleware(authService authentication.Service, allowedRoles []string, allowedPermissions []string) (Middleware, error) {
	rolesMap := toMap(allowedRoles)
	permissionsMap := toMap(allowedPermissions)

	return func(ctx *gin.Context) error {
		userClaims, err := authService.Authenticate(ctx)
		if err != nil {
			return err
		}
		if !isAllowed(rolesMap, userClaims.Roles) && !isAllowed(permissionsMap, userClaims.Permissions) {
			return customerrors.ERROR_UNAUTHORISED
		}
		return nil
	}, nil
}

func isAllowed(allowed map[string]bool, current []string) bool {
	for _, v := range current {
		if _, found := allowed[v]; found {
			return true
		}
	}
	return false
}

func toMap(arr []string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range arr {
		m[v] = true
	}
	return m
}
