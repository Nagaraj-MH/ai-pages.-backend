package routes

import (
	"bookstore/controllers"
	"bookstore/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	auth := router.Group("/user") // `/api/v1/user`
	{
		auth.GET("/getprofile/:username", controllers.GetUserImage)
	}
	protected := router.Group("/user") // `/api/v1/user`
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/getme", controllers.GetMe)
		protected.POST("/upload-profile", controllers.UploadUserProfile)
	}
}
