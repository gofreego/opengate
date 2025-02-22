package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/response"
)

func (s *Service) Route(ctx *gin.Context) {
	id, err := s.match.GetMatchID(ctx)
	if err != nil {
		response.WriteError(ctx, err)
		return
	}

	err = s.middlewares.ExecuteMiddleware(ctx, id)
	if err != nil {
		response.WriteError(ctx, err)
		return
	}
}
