package routes

import (
	"bookstore/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.RouterGroup) {
	books := router.Group("/books") // `/api/v1/books`
	{
		books.GET("/", controllers.GetBooks)
		books.GET("/:id", controllers.GetBookContent)
		books.POST("/upload", controllers.UploadBook)
		books.POST("/:id/like", controllers.LikeBook)
		books.POST("/:id/comment", controllers.AddComment)
		books.GET("/:id/cover", controllers.GetBookCover) 
		books.GET("/:id/pdf", controllers.GetBookPDF) 
	}
}
