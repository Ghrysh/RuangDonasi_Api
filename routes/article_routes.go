package routes

import (
	"donation-api/controllers"
	"donation-api/middlewares"

	"github.com/gin-gonic/gin"
)

func ArticleRoutes(r *gin.Engine) {
	r.GET("/api/articles", controllers.GetArticles)

	adminAccess := r.Group("/api/articles")
	adminAccess.Use(middlewares.AuthMiddleware())
	adminAccess.Use(middlewares.RoleMiddleware("admin", "superadmin"))
	{
		adminAccess.POST("", controllers.CreateArticle)
		
		adminAccess.PUT("/:id", controllers.UpdateArticle)
		adminAccess.DELETE("/:id", controllers.DeleteArticle)
	}
}