package main

import (
	"bookstore/database"
	"bookstore/routes"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		routes.AuthRoutes(v1)
		routes.BookRoutes(v1)
		routes.CommentRoutes(v1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(router.Run(":" + port))
}
