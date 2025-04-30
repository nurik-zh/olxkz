package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=olx port=5432 sslmode=disable"
	connect(dsn)
}

func ConnectTestDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=olxkz_test port=5432 sslmode=disable"
	connect(dsn)
}

func connect(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Ошибка подключения к базе данных: " + err.Error())
	}
}
