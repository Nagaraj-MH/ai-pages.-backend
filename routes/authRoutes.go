package routes

import (
	"bookstore/controllers"
	"bookstore/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth") // `/api/v1/auth`
	{
		auth.GET("/getprofile/:id", controllers.GetUserImage)
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
		auth.GET("/check-username", controllers.CheckUsername)
	}
	protected := router.Group("/auth") // `/api/v1/auth`
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/getme", controllers.GetMe)
		protected.POST("/upload-profile", controllers.UploadUserProfile)
	}
}
