package main

import (
	"fmt"
	"log"
	"bookstore/database"
	"bookstore/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)

	port := "8080"
	fmt.Println("Server running on port", port)
	log.Fatal(router.Run(":" + port))
}
