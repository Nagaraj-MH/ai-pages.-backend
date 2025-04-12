package controllers

import (
	"bookstore/database"
	"bookstore/models"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload a new book
func UploadBook(c *gin.Context) {
	var bookInput struct {
		Title      string `json:"title"`
		Author     string `json:"author"`
		CoverImage string `json:"cover_image"`
		PDFData    string `json:"pdf_data"`
		Likes      int    `json:"likes"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&bookInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	// Decode base64 CoverImage
	coverImageData, err := base64.StdEncoding.DecodeString(bookInput.CoverImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cover image data"})
		return
	}

	// Decode base64 PDFData
	pdfData, err := base64.StdEncoding.DecodeString(bookInput.PDFData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PDF data"})
		return
	}

	// Create book model instance
	book := models.Book{
		Title:      bookInput.Title,
		Author:     bookInput.Author,
		CoverImage: coverImageData,
		PDFData:    pdfData,
		Likes:      bookInput.Likes,
	}

	// Save to database
	database.DB.Create(&book)
	c.JSON(http.StatusCreated, gin.H{"message": "Book uploaded successfully", "book": book})
}

// Fetch all books
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)

	// Convert binary fields to base64 for JSON response
	var responseBooks []map[string]interface{}
	for _, book := range books {
		responseBooks = append(responseBooks, map[string]interface{}{
			"id":         book.ID,
			"title":      book.Title,
			"author":     book.Author,
			"cover_image": base64.StdEncoding.EncodeToString(book.CoverImage),
			"pdf_data":    base64.StdEncoding.EncodeToString(book.PDFData),
			"likes":      book.Likes,
		})
	}

	c.JSON(http.StatusOK, responseBooks)
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

// Get book content
func GetBookContent(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book

	// Fetch book by ID
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Return only book content with base64 encoding
	c.JSON(http.StatusOK, gin.H{
		"title":       book.Title,
		"author":      book.Author,
		"cover_image": base64.StdEncoding.EncodeToString(book.CoverImage),
		"pdf_data":    base64.StdEncoding.EncodeToString(book.PDFData),
	})
}
