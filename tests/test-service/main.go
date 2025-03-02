package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/logger"
)

// start a server using gin which will accept any request and return a 200 response code and log everything in the console

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// log path, method, request body, form, query, header on separate lines
func handleRequest(ctx *gin.Context) {
	logger.Debug(ctx, "handleRequest")
	// Log path
	logger.Debug(ctx, "Path: %s", ctx.Request.URL.Path)

	// Log method
	logger.Debug(ctx, "Method: %s", ctx.Request.Method)

	// Log request body
	body, err := ctx.GetRawData()
	if err == nil {
		logger.Debug(ctx, "Body: %s", string(body))
	} else {
		logger.Error(ctx, "Body: error reading body")
	}

	// Log form
	logger.Debug(ctx, "Form: %s", ctx.Request.PostForm.Encode())

	// Log query
	logger.Debug(ctx, "Query: %s", ctx.Request.URL.RawQuery)

	// Log header
	for key, values := range ctx.Request.Header {
		for _, value := range values {
			logger.Debug(ctx, "Header: %s: %s", key, value)
		}
	}

	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

var (
	port = flag.String("port", "8080", "port to listen on")
)

func main() {
	flag.Parse()
	r := gin.Default()
	r.Any("/*path", handleRequest)
	logger.Info(context.Background(), "Starting server on http://localhost:%s", *port)
	err := r.Run(fmt.Sprintf(":%s", *port))
	if err != nil {
		logger.Error(context.Background(), "Error: %s", err)
	}

}
