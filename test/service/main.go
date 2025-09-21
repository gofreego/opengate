package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	router := gin.Default()

	// Test service routes
	testService := router.Group("/testservice")

	// Ping endpoint
	testService.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "Test service is alive",
			"service":  "testservice",
			"endpoint": "ping",
			"timestamp": gin.H{
				"unix": 1726656000,
				"iso":  "2024-09-18T12:00:00Z",
			},
		})
	})

	testService.POST("/ping", func(c *gin.Context) {
		var request struct {
			Message string `json:"message"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":           "success",
			"message":          "Ping received",
			"received_message": request.Message,
			"service":          "testservice",
			"endpoint":         "ping",
		})
	})

	// Authorise endpoint
	testService.GET("/authorise", func(c *gin.Context) {
		token := c.GetHeader("authorization")

		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status":     "success",
				"authorized": false,
				"message":    "No token provided",
				"service":    "testservice",
				"endpoint":   "authorise",
			})
			return
		}

		// Simple mock authorization logic
		authorized := token == "valid_token" || token == "test_token"

		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"authorized":  authorized,
			"token":       token,
			"message":     "Authorization check completed",
			"service":     "testservice",
			"endpoint":    "authorise",
			"permissions": []string{"read", "write"},
		})
	})

	// Start server on port 8081
	router.Run(":8081")
}
