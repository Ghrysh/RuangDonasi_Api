package routes

import (
	"donation-api/controllers"
	"donation-api/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	api := r.Group("/api")

	userAccess := api.Group("/users").Use(middlewares.AuthMiddleware())
	{
		userAccess.GET("/profile", controllers.GetProfile)
		userAccess.PUT("/profile", controllers.UpdateProfile)

		userAccess.GET("/", controllers.GetAllUsers)
		userAccess.PUT("/:id", controllers.UpdateUserRole)

		userAccess.DELETE("/:id", controllers.DeleteUser)
	}
}