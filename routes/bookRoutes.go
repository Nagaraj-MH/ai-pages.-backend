package routes

import (
	"bookstore/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	book := router.Group("/books")
	{
		book.GET("/:id", controllers.GetBookContent)
		book.GET("/", controllers.GetBooks)
		book.POST("/upload", controllers.UploadBook)
		book.POST("/:id/like", controllers.LikeBook)
	}
}
