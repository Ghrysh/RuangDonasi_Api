package routes

import (
	"donation-api/controllers"
	"donation-api/middlewares"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.Engine) {
	api := r.Group("/api")

	userAccess := api.Group("/transactions").Use(middlewares.AuthMiddleware())
	{
		userAccess.POST("/", controllers.CreateDirectTransaction)
		userAccess.GET("", controllers.GetAllTransactions)
		userAccess.GET("/:id", controllers.GetTransactionStatus)
	}

	webhookAccess := api.Group("/webhooks")
	{
		webhookAccess.POST("/mutation", controllers.MutationWebhook)
	}
}