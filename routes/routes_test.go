package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"olxkz/config"
	"olxkz/models"
	"testing"
)

// Функция для настройки тестового роутера
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", GetAllUsers)
	router.POST("/users", Register)
	router.GET("/categories", GetCategories)
	router.POST("/categories", CreateCategory)
	router.PUT("/categories/:id", UpdateCategory)
	router.DELETE("/categories/:id", DeleteCategory)
	router.GET("/products", GetProducts)
	router.POST("/products", CreateProduct)
	router.PUT("/products/:id", UpdateProduct)
	router.DELETE("/products/:id", DeleteProduct)
	return router
}

// Инициализация тестовой базы данных
func initTestDB() {
	// Здесь нужно настроить отдельную тестовую базу данных
	config.ConnectDatabase()

	// Очистка данных перед тестами
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM categories")
	config.DB.Exec("DELETE FROM products")
}

// Тест создания пользователя
func TestCreateUser(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	user := models.User{Username: "Test User", Email: "test@example.com"}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, createdUser.Email)
}

// Тест получения пользователей
func TestGetUsers(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Добавляем пользователя вручную для теста
	config.DB.Create(&models.User{Username: "Jane", Email: "jane@example.com"})

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)
	assert.Equal(t, "Jane", users[0].Username)
}

// Тест создания категории
func TestCreateCategory(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	category := models.Category{Name: "Test Category"}
	body, _ := json.Marshal(category)

	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCategory models.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	assert.NoError(t, err)
	assert.Equal(t, category.Name, createdCategory.Name)
}

// Тест получения категорий
func TestGetCategories(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Добавляем категорию вручную для теста
	config.DB.Create(&models.Category{Name: "Electronics"})

	req, _ := http.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var categories []models.Category
	err := json.Unmarshal(w.Body.Bytes(), &categories)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(categories), 1)
	assert.Equal(t, "Electronics", categories[0].Name)
}

// Тест обновления категории
func TestUpdateCategory(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Создаем категорию для обновления
	category := models.Category{Name: "Old Category"}
	config.DB.Create(&category)

	updatedCategory := models.Category{Name: "Updated Category"}
	body, _ := json.Marshal(updatedCategory)

	// Используем fmt.Sprintf для преобразования ID в строку
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%d", category.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.Category
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, updatedCategory.Name, updated.Name)
}

// Тест удаления категории
func TestDeleteCategory(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Создаем категорию для удаления
	category := models.Category{Name: "To Be Deleted"}
	result := config.DB.Create(&category)
	assert.NoError(t, result.Error) // Убедитесь, что категория успешно создана

	// Используем fmt.Sprintf для преобразования ID в строку
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/categories/%d", category.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что ответ имеет статус 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что категория была удалена
	var deletedCategory models.Category
	err := config.DB.First(&deletedCategory, category.ID).Error
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected error to be gorm.ErrRecordNotFound")
}

// Тест создания продукта
func TestCreateProduct(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	product := models.Product{Name: "Test Product", CategoryID: 1}
	body, _ := json.Marshal(product)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct models.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(t, err)
	assert.Equal(t, product.Name, createdProduct.Name)
}

// Тест получения продуктов
func TestGetProducts(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Добавляем продукт вручную для теста
	config.DB.Create(&models.Product{Name: "Product A", CategoryID: 1})

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var products []models.Product
	err := json.Unmarshal(w.Body.Bytes(), &products)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(products), 1)
	assert.Equal(t, "Product A", products[0].Name)
}

// Тест обновления продукта
func TestUpdateProduct(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Создаем продукт для обновления
	product := models.Product{Name: "Old Product", CategoryID: 1}
	config.DB.Create(&product)

	updatedProduct := models.Product{Name: "Updated Product", CategoryID: 1}
	body, _ := json.Marshal(updatedProduct)

	// Используем fmt.Sprintf для преобразования ID в строку
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", product.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.Product
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, updatedProduct.Name, updated.Name)
}

// Тест удаления продукта
func TestDeleteProduct(t *testing.T) {
	initTestDB()
	router := setupTestRouter()

	// Создаем продукт для удаления
	product := models.Product{Name: "To Be Deleted", CategoryID: 1}
	config.DB.Create(&product)

	// Используем fmt.Sprintf для преобразования ID в строку
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%d", product.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что продукт был удален
	var deletedProduct models.Product
	err := config.DB.First(&deletedProduct, product.ID).Error
	assert.Error(t, err) // Ожидаем ошибку, так как продукт должен быть удален
}
