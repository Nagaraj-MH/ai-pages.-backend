package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"bookstore/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	fmt.Println("Database connected successfully")

	db.AutoMigrate(&models.Book{}, &models.Comment{})

	DB = db
}
