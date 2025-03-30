package routes

import (
	"bookstore/controllers"
	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.RouterGroup) { 
	comments := router.Group("/comments") // `/api/v1/comments`
	{
		comments.POST("/", controllers.AddComment)
		comments.GET("/:bookID", controllers.GetComments)
	}
}
