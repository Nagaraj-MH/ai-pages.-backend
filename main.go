package main

import (
	"bookstore/database"
	"bookstore/routes"
	"fmt"
	"log"
	"os"


	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(os.Getenv("DATABASE_URL"))
	database.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(router.Run(":" + port))
}
