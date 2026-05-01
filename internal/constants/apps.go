package constants

const (
	HTTP_SERVER = "HTTP_SERVER"
	GRPC_SERVER = "GRPC_SERVER"
	JWT_CLAIMS  = "jwt_claims"

	HEADER_USER_ID       = "X-USER-ID"
	HEADER_USER_UUID     = "X-USER-UUID"
	HEADER_PROFILE_IDS   = "X-PROFILE-IDS"
	HEADER_PERMISSIONS   = "X-PERMISSIONS"
	HEADER_AUTHORIZATION = "Authorization"

	COOKIE_AUTHORIZATION = "authorization"
)

// Route permissions
const (
	PERMISSION_ROUTES_READ  = "routes:read"
	PERMISSION_ROUTES_WRITE = "routes:write"
)
