package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"olxkz/config"
	"olxkz/middleware" // егер middleware бар болса
	"olxkz/models"
)

// Барлық user роуттарын тіркеу
func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware()) // 🔐 Қорғалған маршрут
	{
		users.GET("", GetAllUsers)
	}
}

// Барлық қолданушыларды алу
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// DB сұраныс
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Қолданушылар табылмады"})
		return
	}

	// Парольді шығарып тастаймыз (қауіпсіздік үшін)
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}
