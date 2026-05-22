package routes

import (
	"donation-api/controllers"

	"github.com/gin-gonic/gin"
)

func CountryRoutes(r *gin.Engine) {
	api := r.Group("/api/countries")
	{
		api.GET("/codes", controllers.GetCountryCodes)
	}
}