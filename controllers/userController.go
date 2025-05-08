package controllers

import (
	"bookstore/constants"
	"bookstore/database"
	"bookstore/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	email, exists := c.Get(string(constants.ContextUserEmailKey))
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user with email exists"})
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "email": user.Email, "name": user.Name})

}
func UploadUserProfile(c *gin.Context) {
	var user models.User
	email, exists := c.Get(string(constants.ContextUserEmailKey))
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user with email exists"})
		return
	}

	pictureHeader, err := c.FormFile("profilePicture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coverFile, err := pictureHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	defer coverFile.Close()
	profileImageData, err := io.ReadAll(coverFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
		return
	}

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	user.ProfilePicture = profileImageData
	database.DB.Save(&user)
	c.JSON(http.StatusAccepted, gin.H{"message": "succuess"})
}
func GetUserImage(c *gin.Context) {
	username := c.Param("username")
	print(username)
	var user models.User
	if err := database.DB.Select("profile_picture").Where("username", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(user.ProfilePicture) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile image not found for this User"})
		return
	}

	contentType := http.DetectContentType(user.ProfilePicture)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, user.ProfilePicture)
}
