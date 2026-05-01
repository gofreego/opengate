package utils

import (
	"context"
	"strings"

	"github.com/gofreego/opengate/internal/constants"
	"google.golang.org/grpc/metadata"
)

// GetPermissionsFromContext extracts permissions from gRPC metadata headers
func GetPermissionsFromContext(ctx context.Context) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	// Get permissions from header (lowercase key in grpc-gateway)
	perms := md.Get(strings.ToLower(constants.HEADER_PERMISSIONS))
	if len(perms) == 0 {
		return nil
	}

	// Permissions are comma-separated
	return strings.Split(perms[0], ",")
}

// HasPermission checks if a specific permission exists in the context
func HasPermission(ctx context.Context, permission string) bool {
	perms := GetPermissionsFromContext(ctx)
	for _, p := range perms {
		if strings.TrimSpace(p) == permission {
			return true
		}
	}
	return false
}

// HasRoutesReadPermission checks if user has routes read permission
func HasRoutesReadPermission(ctx context.Context) bool {
	return HasPermission(ctx, constants.PERMISSION_ROUTES_READ)
}

// HasRoutesWritePermission checks if user has routes write permission
func HasRoutesWritePermission(ctx context.Context) bool {
	return HasPermission(ctx, constants.PERMISSION_ROUTES_WRITE)
}
