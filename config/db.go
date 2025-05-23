package config

import (
	"fmt"
	"olxkz/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := getDSN("olx")
	connect(dsn)
}

func ConnectTestDatabase() {
	host := "localhost"
	port := "5433"
	user := "postgres"
	password := "postgres"
	dbname := "olxkz_test"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty",
		host, user, password, dbname, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	DB = database
	DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{})
}

func getDSN(dbname string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)
}

func connect(dsn string) {
	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Попытка подключения #%d неудачна: %s\n", i+1, err.Error())
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("Ошибка подключения к базе данных: " + err.Error())
	}

}
