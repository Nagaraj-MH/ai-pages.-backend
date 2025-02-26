package routes

import (
	"bookstore/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	comment := router.Group("/comments")
	{
		comment.POST("/", controllers.AddComment)
		comment.GET("/:bookID", controllers.GetComments)
	}
}
