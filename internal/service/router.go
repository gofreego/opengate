package service

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/logger"
	"github.com/gofreego/openauth/pkg/jwtutils"
	"github.com/gofreego/opengate/internal/constants"
	"github.com/gofreego/opengate/internal/models"
)

func (s *Service) RouteRequest(ctx *gin.Context) {
	// Get the route for this request
	route := s.routeManager.GetRouteByRequest(ctx.Request)
	if route == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No route found for this request"})
		return
	}

	// Check authentication if required
	if route.Authentication.IsAuthenticationRequired(ctx.Request.URL.Path, ctx.Request.Method) {
		if err := s.authManager.Authenticate(ctx); err != nil {
			logger.Warn(ctx, "Authentication failed for route: %s, error: %v", route.Name, err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
	}

	// Use ProxyPass to handle the request
	s.proxyPass(ctx, route)
}

// proxyPass handles the proxying of requests to the target service
func (s *Service) proxyPass(ctx *gin.Context, route *models.ServiceRoute) {
	// Parse target URL
	targetURL, err := url.Parse(route.TargetURL)
	if err != nil {
		logger.Error(ctx, "Failed to parse target URL: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
		return
	}

	// Create reverse proxy with httputil
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Configure proxy settings
	s.configureProxy(ctx, proxy, route)

	// Handle path modification if needed
	if route.StripPrefix {
		originalPath := ctx.Request.URL.Path
		if strings.HasPrefix(originalPath, route.PathPrefix) {
			ctx.Request.URL.Path = strings.TrimPrefix(originalPath, route.PathPrefix)
			if !strings.HasPrefix(ctx.Request.URL.Path, "/") && ctx.Request.URL.Path != "" {
				ctx.Request.URL.Path = "/" + ctx.Request.URL.Path
			}
		}
	}

	// Execute the proxy
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (s *Service) configureProxy(ctx *gin.Context, proxy *httputil.ReverseProxy, route *models.ServiceRoute) {
	// Set timeout by configuring transport
	timeout := route.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}

	// Configure transport with timeout
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: timeout,
		IdleConnTimeout:       timeout,
	}

	// Configure error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error(r.Context(), "Proxy error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error": "Service unavailable"}`))
	}

	// Configure director for request modification
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Clear user headers to prevent spoofing
		req.Header.Del(constants.HEADER_AUTHORIZATION)
		req.Header.Del(constants.HEADER_USER_ID)
		req.Header.Del(constants.HEADER_USER_UUID)
		req.Header.Del(constants.HEADER_PROFILE_IDS)
		req.Header.Del(constants.HEADER_PERMISSIONS)

		// Add forwarding headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", getClientIP(req))
		req.Header.Set("X-Forwarded-Proto", getScheme(req))

		// Add user headers from JWT claims if authentication was required
		if claims, exists := ctx.Get(constants.JWT_CLAIMS); exists {
			if jwtClaims, ok := claims.(*jwtutils.JWTClaims); ok {
				if jwtClaims.UserID != 0 {
					req.Header.Set(constants.HEADER_USER_ID, fmt.Sprintf("%d", jwtClaims.UserID))
				}
				if jwtClaims.UserUUID != "" {
					req.Header.Set(constants.HEADER_USER_UUID, jwtClaims.UserUUID)
				}
				if len(jwtClaims.Profiles) > 0 {
					profileIDs := make([]string, len(jwtClaims.Profiles))
					for i, p := range jwtClaims.Profiles {
						profileIDs[i] = fmt.Sprintf("%d", p.Id)
					}
					req.Header.Set(constants.HEADER_PROFILE_IDS, strings.Join(profileIDs, ","))
				}
				if len(jwtClaims.Permissions) > 0 {
					req.Header.Set(constants.HEADER_PERMISSIONS, strings.Join(jwtClaims.Permissions, ","))
				}
			}
		}
	}
}

// getClientIP extracts the real client IP from request
func getClientIP(req *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		if idx := strings.Index(xff, ","); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(req.RemoteAddr, ":"); idx > 0 {
		return req.RemoteAddr[:idx]
	}
	return req.RemoteAddr
}

// getScheme determines the request scheme
func getScheme(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
