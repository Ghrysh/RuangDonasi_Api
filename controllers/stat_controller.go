package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var totalDonations float64
	var totalCampaigns int64
	var totalTransactions int64

	config.DB.Model(&models.Transaction{}).
		Where("status = ?", "success").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalDonations)

	config.DB.Model(&models.Campaign{}).Count(&totalCampaigns)

	config.DB.Model(&models.Transaction{}).
		Where("status = ?", "success").
		Count(&totalTransactions)

	c.JSON(http.StatusOK, gin.H{
		"message": "Statistik berhasil diambil",
		"data": gin.H{
			"total_saldo_donasi": totalDonations,
			"total_campaign":     totalCampaigns,
			"total_donatur":      totalTransactions,
		},
	})
}