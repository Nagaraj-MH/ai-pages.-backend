package main

import (
	"bookstore/database"
	"bookstore/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change this for production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight response for 12 hours
	}))
	v1 := router.Group("/api/v1")
	{
		routes.AuthRoutes(v1)
		routes.BookRoutes(v1)
		routes.CommentRoutes(v1)
	}
	router.POST("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Login successful"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(router.Run(":" + port))
}
