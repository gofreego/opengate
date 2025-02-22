package middlewares

import "github.com/gin-gonic/gin"

type Middleware func(ctx *gin.Context) error
