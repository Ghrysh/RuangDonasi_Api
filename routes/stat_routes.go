package routes

import (
	"donation-api/controllers"
	"donation-api/middlewares"

	"github.com/gin-gonic/gin"
)

func StatRoutes(r *gin.Engine) {
	api := r.Group("/api/stats")

	api.Use(middlewares.AuthMiddleware())
	api.Use(middlewares.RoleMiddleware("admin", "superadmin"))
	{
		api.GET("/", controllers.GetDashboardStats)
	}
}