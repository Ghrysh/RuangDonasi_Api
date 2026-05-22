package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type MacroDroidPayload struct {
	Amount float64 `json:"amount"` 
	Text   string  `json:"text"` // Untuk menyimpan isi notifikasi asli ("Berhasil menerima Rp 50.123...")
}

func MutationWebhook(c *gin.Context) {
	secret := c.GetHeader("X-Callback-Token")
	if secret != os.Getenv("MUTATION_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses Webhook Ditolak"})
		return
	}

	var mutasi MacroDroidPayload
	if err := c.ShouldBindJSON(&mutasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak sesuai"})
		return
	}

	var transaction models.Transaction

	err := config.DB.Where("total_amount = ? AND status = ?", mutasi.Amount, "pending").First(&transaction).Error
	
	if err == nil {
		transaction.Status = "success"
		config.DB.Save(&transaction)

		var campaign models.Campaign
		if err := config.DB.First(&campaign, transaction.CampaignID).Error; err == nil {
			campaign.CurrentAmount += transaction.Amount 
			campaign.DonatorCount += 1
			config.DB.Save(&campaign)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Donasi Berhasil Diproses", "trx_id": transaction.ID})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada tagihan yang cocok dengan nominal tersebut"})
}