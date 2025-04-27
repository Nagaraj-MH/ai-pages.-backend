package controllers

import (
	"bookstore/database"
	"bookstore/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Add a comment to a book
func AddComment(c *gin.Context) {
	bookID := c.Param("id")
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment data"})
		return
	}
	id, err := strconv.ParseUint(bookID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	comment.BookID = uint(id)

	database.DB.Create(&comment)
	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully", "comment": gin.H{
		"text":       comment.Content,
		"userId":     comment.UserID,
		"created_at": comment.CreatedAt,
	}})
}
