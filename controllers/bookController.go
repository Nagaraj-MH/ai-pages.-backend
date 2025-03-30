package controllers

import (
	"bookstore/database"
	"bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload a new book
func UploadBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	database.DB.Create(&book)
	c.JSON(http.StatusCreated, gin.H{"message": "Book uploaded successfully", "book": book})
}

// Fetch all books
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

// Like a book
func LikeBook(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book

	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	book.Likes++
	database.DB.Save(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Book liked", "likes": book.Likes})
}

func GetBookContent(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book

	// Fetch book by ID
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Return only book content
	c.JSON(http.StatusOK, gin.H{
		"title":   book.Title,
		"author":  book.Author,
		"PDFData": book.PDFData,
	})
}