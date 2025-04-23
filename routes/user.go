package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"olxkz/config"
	"olxkz/middleware" // если middleware есть
	"olxkz/models"
)

// Барлық user роуттарын тіркеу
func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware()) // 🔐 Защищенный маршрут
	{
		users.GET("", GetAllUsers)
		users.DELETE("/:id", DeleteUser) // Новый маршрут для удаления пользователя
	}
}

// Барлық қолданушыларды алу
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Запрос к базе данных
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Қолданушылар табылмады"})
		return
	}

	// Просто возвращаем пользователей без очистки пароля
	c.JSON(http.StatusOK, users)
}

// Удаление пользователя
func DeleteUser(c *gin.Context) {
	id := c.Param("id") // Получаем id пользователя из URL

	var user models.User

	// Ищем пользователя по ID
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Удаляем пользователя
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}
