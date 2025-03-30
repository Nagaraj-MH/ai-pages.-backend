package tests

import (
	"bookstore/controllers"
	"bookstore/database"
	"bookstore/models"
	"bytes"
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
	}
	comments := router.Group("/comments")
	{
		comments.POST("/", controllers.AddComment)
		comments.GET("/:bookID", controllers.GetComments)
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

	book := `{"title": "Deep Learning", "author": "Ian Goodfellow", "cover_image": "cover.jpg", "content": "Deep Learning explained..."}`
	req, _ := http.NewRequest("POST", "/books/upload", bytes.NewBuffer([]byte(book)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

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


func TestGetComments(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/comments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
