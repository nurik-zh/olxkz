package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"log"
	"time"
)

func GetUserByID(id string) (string, error) {
	client := resty.New()

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Printf("[Resty] Requesting: %s %s\n", req.Method, req.URL)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		fmt.Printf("[Resty] Response Code: %d\n", resp.StatusCode())
		return nil
	})

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("http://localhost:8081/user/%s", id))

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func main() {
	r := gin.Default()

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

	// Эндпоинт, вызывающий user-service
	r.GET("/product/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		user, err := GetUserByID(userId)
		if err != nil {
			c.JSON(500, gin.H{"error": "User fetch failed"})
			return
		}

		c.JSON(200, gin.H{
			"product": "Sample Product",
			"user":    user,
		})
	})

	r.Run(":8082")
}
