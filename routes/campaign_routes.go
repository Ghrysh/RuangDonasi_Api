package routes

import (
	"donation-api/controllers"
	"donation-api/middlewares"

	"github.com/gin-gonic/gin"
)

func CampaignRoutes(r *gin.Engine) {
	r.GET("/api/campaigns", controllers.GetCampaigns)

	adminAccess := r.Group("/api/campaigns")
	adminAccess.Use(middlewares.AuthMiddleware())
	adminAccess.Use(middlewares.RoleMiddleware("admin", "superadmin"))
	{
		adminAccess.POST("", controllers.CreateCampaign)

		adminAccess.PUT("/:id", controllers.UpdateCampaign)
		
		adminAccess.DELETE("/:id", controllers.DeleteCampaign) 
	}
}