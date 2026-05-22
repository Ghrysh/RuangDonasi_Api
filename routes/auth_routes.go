package routes

import (
	"donation-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/google", controllers.GoogleLogin)

		authGroup.POST("/forgot-password", controllers.ForgotPassword)
		authGroup.POST("/reset-password", controllers.ResetPassword)
	}
}