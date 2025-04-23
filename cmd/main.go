package main

import (
	"github.com/gin-gonic/gin"
	"olxkz/config"
	"olxkz/routes"
)

func main() {
	r := gin.Default()

	config.ConnectDatabase()
	routes.RegisterAuthRoutes(r)
	routes.RegisterCategoryRoutes(r)
	routes.RegisterProductRoutes(r)
	routes.RegisterFavoriteRoutes(r)

	routes.RegisterUserRoutes(r)

	r.Run(":8080")
}
