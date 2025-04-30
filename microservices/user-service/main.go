package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"olxkz/config"
	"olxkz/models"
)

func main() {
	config.ConnectDatabase()
	r := gin.Default()

	// Middleware
	r.Use(func(c *gin.Context) {
		requestID := uuid.New().String()
		start := time.Now()
		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
		duration := time.Since(start)
		log.Printf("[%s] [RequestID: %s] %s %s - %d - Duration: %.3fms\n",
			time.Now().Format(time.RFC3339),
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			float64(duration.Microseconds())/1000.0,
		)
	})

	// Get user by ID
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(200, user)
	})

	r.Run(":8081") // User Service is running on port 8081
}
