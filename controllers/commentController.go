package controllers

import (
	"bookstore/database"
	"bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add a comment to a book
func AddComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment data"})
		return
	}

	database.DB.Create(&comment)
	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully", "comment": comment})
}

// Get comments for a book
func GetComments(c *gin.Context) {
	bookID := c.Param("bookID")
	var comments []models.Comment
	database.DB.Where("book_id = ?", bookID).Find(&comments)
	c.JSON(http.StatusOK, comments)
}
