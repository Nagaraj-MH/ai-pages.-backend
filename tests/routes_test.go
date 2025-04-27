package tests

import (
	"bookstore/controllers"
	"bookstore/database"
	"bookstore/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func setupTestDB() {
	database.ConnectDB() 
	database.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Comment{})
}


func setupRouter() *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
	}
	books := router.Group("/books")
	{
		books.POST("/upload", controllers.UploadBook)
		books.GET("/", controllers.GetBooks)
		books.POST("/:id/like", controllers.LikeBook)
		books.POST("/:id/comment", controllers.AddComment)
	}

	return router
}


func TestSignup(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	user := `{"name": "John Doe", "email": "johndoe@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte(user)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}


func TestLogin(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	
	user := `{"name": "John Doe", "email": "johndoe@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte(user)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	
	loginData := `{"email": "johndoe@example.com", "password": "password123"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(loginData)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "token") 
}


func TestUploadBook(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Simulated binary data for CoverImage and PDFData
	coverImageData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG file signature
	pdfData := []byte("%PDF-1.4 Sample PDF content")

	// Encode binary data as base64 strings for JSON
	encodedCoverImage := base64.StdEncoding.EncodeToString(coverImageData)
	encodedPDFData := base64.StdEncoding.EncodeToString(pdfData)

	// Create book struct
	book := map[string]interface{}{
		"title":       "Deep Learning",
		"author":      "Ian Goodfellow",
		"cover_image": encodedCoverImage,
		"pdf_data":    encodedPDFData,
		"likes":       0,
	}

	// Convert struct to JSON
	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book JSON: %v", err)
	}

	// Create request
	req, _ := http.NewRequest("POST", "/books/upload", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert HTTP status
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBooks(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/books/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}


func TestLikeBook(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	
	book := `{"title": "Deep Learning", "author": "Ian Goodfellow", "cover_image": "cover.jpg", "content": "Deep Learning explained..."}`
	req, _ := http.NewRequest("POST", "/books/upload", bytes.NewBuffer([]byte(book)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	
	req, _ = http.NewRequest("POST", "/books/1/like", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}


func TestAddComment(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	comment := `{"book_id": 1, "user_id": 1, "content": "Great book!"}`
	req, _ := http.NewRequest("POST", "/comments/", bytes.NewBuffer([]byte(comment)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}


