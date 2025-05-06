package routes

import (
	"bookstore/controllers"
	"bookstore/middlewares"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.RouterGroup) {
	books := router.Group("/books") // `/api/v1/books`
	{
		books.GET("/", controllers.GetBooks)
		books.GET("/featured", controllers.GetFeaturedBooks)
		books.GET("/:id", controllers.GetBookContent)
		books.POST("/upload", controllers.UploadBook)
		books.GET("/:id/cover", controllers.GetBookCover)
		books.GET("/:id/pdf", controllers.GetBookPDF)
	}
	protected := router.Group("/auth") // `/api/v1/auth`
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/:id/like", controllers.LikeBook)
		protected.POST("/:id/comment", controllers.AddComment)
	}
}
