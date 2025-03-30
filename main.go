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
	os.Setenv("DATABASE_URL", "postgres://avnadmin:AVNS_h838AM4ym48oGQAHCBi@pg-3ce0a6e1-superusxr-b8ec.g.aivencloud.com:20039/defaultdb?sslmode=require")
	log.Println(os.Getenv("DATABASE_URL"))
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
