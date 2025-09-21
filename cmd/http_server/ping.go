package http_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/logger"
)

func (s *HTTPServer) ping(ctx *gin.Context) {

	logger.Debug(ctx, "Ping request received")
	res, err := s.service.Ping(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
