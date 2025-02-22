package service

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ReverseProxy forwards requests to the appropriate backend service
func ReverseProxy(c *gin.Context, target string) {
	targetURL, err := url.Parse(target)
	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	c.Request.Host = targetURL.Host
	proxy.ServeHTTP(c.Writer, c.Request)
}
