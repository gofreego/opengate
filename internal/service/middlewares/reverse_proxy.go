package middlewares

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/customerrors"
)

// ReverseProxy forwards requests to the appropriate backend service
func (s *MiddlewareService) getReverseProxyMiddleware(target string) (Middleware, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, customerrors.New(http.StatusInternalServerError, "invalid target url")
	}

	return func(ctx *gin.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		ctx.Request.Host = targetURL.Host
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	}, nil
}
