package controllers

import (
	"bookstore/database"
	"bookstore/models"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// controllers/books.go
func UploadBook(c *gin.Context) {
	description := c.PostForm("description")
	tags := c.PostFormArray("tags")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	cleanedTags := []string{}
	for _, tag := range tags {
		if tag != "" {
			cleanedTags = append(cleanedTags, tag)
		}
	}
	title := c.PostForm("title")
	author := c.PostForm("author")
	// Handle Cover Image File
	coverFileHeader, err := c.FormFile("coverImageFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cover image file is required"})
		return
	}
	coverFile, err := coverFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open cover image"})
		return
	}
	defer coverFile.Close()
	coverImageData, err := io.ReadAll(coverFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read cover image"})
		return
	}

	// Handle PDF File (similar logic)
	pdfFileHeader, err := c.FormFile("pdfFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file is required"})
		return
	}
	pdfFile, err := pdfFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open PDF file"})
		return
	}
	defer pdfFile.Close()
	pdfData, err := io.ReadAll(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF file"})
		return
	}

	book := models.Book{
		Title:       title,
		Author:      author,
		CoverImage:  coverImageData,
		PDFData:     pdfData,
		Tags:        cleanedTags,
		Description: description,
	}

	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book uploaded successfully via multipart", "book_id": book.ID})
}
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)

	var responseBooks []map[string]interface{}
	for _, book := range books {
		responseBooks = append(responseBooks, map[string]interface{}{
			"id":          book.ID,
			"title":       book.Title,
			"author":      book.Author,
			"likes":       book.Likes,
			"tags":        book.Tags,
			"description": book.Description,
		})
	}

	c.JSON(http.StatusOK, responseBooks)
}
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
	type CommentResponse struct {
		ID     uint   `json:"id"`
		UserID uint   `json:"userId"`
		Text   string `json:"text"`
	}
	bookID := c.Param("id")
	var book models.Book
	var comments []models.Comment

	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := database.DB.Where("book_id = ?", book.ID).Find(&comments).Error; err != nil {
		comments = []models.Comment{}
	}
	responseComments := make([]CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responseComments = append(responseComments, CommentResponse{
			ID:     comment.ID,
			UserID: comment.UserID,
			Text:   comment.Content,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"title":    book.Title,
		"author":   book.Author,
		"likes":    book.Likes,
		"comments": responseComments,
	})
}

func GetBookCover(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book

	if err := database.DB.Select("cover_image").First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if len(book.CoverImage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cover image not found for this book"})
		return
	}

	contentType := http.DetectContentType(book.CoverImage)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, book.CoverImage)
}

func GetBookPDF(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book

	if err := database.DB.Select("pdf_data", "title").First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if len(book.PDFData) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found for this book"})
		return
	}

	contentType := "application/pdf"
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename=\""+book.Title+".pdf\"")
	c.Data(http.StatusOK, contentType, book.PDFData)
}

func GetFeaturedBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Order("likes desc").Limit(3).Find(&books)

	var responseBooks []map[string]interface{}
	for _, book := range books {
		responseBooks = append(responseBooks, map[string]interface{}{
			"id":          book.ID,
			"title":       book.Title,
			"author":      book.Author,
			"likes":       book.Likes,
			"tags":        book.Tags,
			"description": book.Description,
		})
	}

	c.JSON(http.StatusOK, responseBooks)
}
