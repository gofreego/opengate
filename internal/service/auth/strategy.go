package auth

import "github.com/gin-gonic/gin"

type Strategy interface {
	Authenticate(ctx *gin.Context) error
}
