package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"olxkz/config"
	"olxkz/models"
	"strings"
)

var jwtKey = []byte("secret")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token отсутствует"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный token"})
			c.Abort()
			return
		}

		// Извлекаем username
		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		// Ищем user в БД
		var user models.User
		if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			c.Abort()
			return
		}

		// Кладем в контекст и username, и userID
		c.Set("username", user.Username)
		c.Set("userID", user.ID)

		c.Next()
	}
}
