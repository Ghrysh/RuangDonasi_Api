package main

import (
	"donation-api/config"
	"donation-api/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true 
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "ngrok-skip-browser-warning", "Accept", "X-Pinggy-No-Screen"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/uploads", "./uploads")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server Donasi Berjalan Lancar!"})
	})

	routes.AuthRoutes(router)
	routes.CampaignRoutes(router)
	routes.TransactionRoutes(router)
	routes.StatRoutes(router)
	routes.CountryRoutes(router)
	routes.UserRoutes(router)
	routes.ArticleRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	router.Run(":" + port)
}