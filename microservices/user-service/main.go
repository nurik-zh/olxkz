package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

func main() {
	r := gin.Default()

	// Middleware: логирует Request ID, метод, статус и длительность
	r.Use(func(c *gin.Context) {
		requestID := uuid.New().String()
		start := time.Now()

		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		duration := time.Since(start)
		log.Printf("[%s] [Requested: %s] %s %s - %d - Duration: %.3fms\n",
			time.Now().Format(time.RFC3339),
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			float64(duration.Microseconds())/1000.0,
		)

	})

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"id":   id,
			"name": "TestUser",
		})
	})

	r.Run(":8081")
}
