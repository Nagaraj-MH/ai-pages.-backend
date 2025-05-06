package controllers

import (
	"bookstore/constants"
	"bookstore/database"
	"bookstore/models"
	"bookstore/utils"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"

	"gorm.io/gorm"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user.Username = strings.ToLower(user.Username)
	var existingUser models.User
	err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

// Login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Debug password comparison
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Reset Password
func ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ?", request.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if token has expired
	if time.Now().After(user.TokenExpires) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Reset token has expired"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Update password and remove token
	user.Password = hashedPassword
	user.ResetToken = ""
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// Forgot Password - Generate Reset Token
func ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate a secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}
	user.ResetToken = hex.EncodeToString(tokenBytes)
	user.TokenExpires = time.Now().Add(1 * time.Hour) // Token expires in 1 hour

	// Save token to database
	database.DB.Save(&user)

	//Send reset email with the token (Using email)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to email", "token": user.ResetToken})
}
func CheckUsername(c *gin.Context) {
	username := c.Query("username")
	username = strings.ToLower(username)
	var existingUser models.User
	err := database.DB.Where("username = ?", username).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"available": false})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"available": true})
}

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
	username := c.Param("id")
	var user models.User
	if err := database.DB.Select("profile_picture").Where("username", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if len(user.ProfilePicture) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cover image not found for this book"})
		return
	}

	contentType := http.DetectContentType(user.ProfilePicture)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, user.ProfilePicture)
}
