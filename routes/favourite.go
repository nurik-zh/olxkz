package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"olxkz/config"
	"olxkz/middleware"
	models2 "olxkz/models"
	"strconv"
)

func RegisterFavoriteRoutes(r *gin.Engine) {
	favorites := r.Group("/favorites")
	favorites.Use(middleware.AuthMiddleware())
	{
		favorites.GET("", GetFavorites)
		favorites.POST("/:productId", AddToFavorites)
		favorites.DELETE("/:productId", RemoveFromFavorites)
	}
}

func GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")

	var favorites []models2.Favorite
	err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить избранное"})
		return
	}

	// Возвращаем только данные продукта
	var products []models2.Product
	for _, fav := range favorites {
		products = append(products, fav.Product)
	}

	c.JSON(http.StatusOK, products)
}

func AddToFavorites(c *gin.Context) {
	userID := c.GetUint("userID") // Предположим, мы сохраняем userID в middleware
	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукта"})
		return
	}

	fav := models2.Favorite{
		UserID:    userID,
		ProductID: uint(productId),
	}
	config.DB.Create(&fav)

	c.JSON(http.StatusCreated, gin.H{"message": "Добавлено в избранное"})
}

func RemoveFromFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукта"})
		return
	}

	config.DB.Where("user_id = ? AND product_id = ?", userID, productId).Delete(&models2.Favorite{})
	c.JSON(http.StatusOK, gin.H{"message": "Удалено из избранного"})
}
