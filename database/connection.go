package database

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"bookstore/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := "postgres://avnadmin:AVNS_h838AM4ym48oGQAHCBi@pg-3ce0a6e1-superusxr-b8ec.g.aivencloud.com:20039/ai-pages?sslmode=require"
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	fmt.Println("Database connected successfully")

	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Comment{})

	DB = db
}
