package main

import (
	// "bookstore/database"
	"bookstore/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// log.Println(os.Getenv("DATABASE_URL"))
	// database.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.BookRoutes(router)
	routes.CommentRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(router.Run(":" + port))
}
