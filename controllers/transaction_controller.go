package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"donation-api/utils"
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type DirectDonationInput struct {
	CampaignID uint    `json:"campaign_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,min=10000"`
}

func CreateDirectTransaction(c *gin.Context) {
	var input DirectDonationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	var campaign models.Campaign
	if err := config.DB.First(&campaign, input.CampaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}

	myStaticQRIS := "00020101021126570011ID.DANA.WWW011893600915302458384202090245838420303UMI51440014ID.CO.QRIS.WWW0215ID10265232339420303UMI5204504553033605802ID5912Ruang Donasi6012Kab. Bandung61054011163049F30"

	var uniqueCode int
	var totalAmount float64

	for {
		rand.Seed(time.Now().UnixNano())
		uniqueCode = rand.Intn(900) + 100
		totalAmount = input.Amount + float64(uniqueCode)

		var existingTransaction models.Transaction
		err := config.DB.Where("total_amount = ? AND status = ?", totalAmount, "pending").First(&existingTransaction).Error
		if err != nil {
			break
		}
	}

	dynamicQRText := utils.GenerateDynamicQRIS(myStaticQRIS, totalAmount)

	var png []byte
	png, err := qrcode.Encode(dynamicQRText, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat gambar QRIS"})
		return
	}

	base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	transaction := models.Transaction{
		UserID:      userID.(uint),
		CampaignID:  input.CampaignID,
		Amount:      input.Amount,
		UniqueCode:  uniqueCode,
		TotalAmount: totalAmount,
		Status:      "pending",
		PaymentURL:  base64QR,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil dibuat, silakan scan QRIS untuk membayar",
		"data": gin.H{
			"transaction_id":   transaction.ID,
			"nominal_asli":     transaction.Amount,
			"kode_unik":        transaction.UniqueCode,
			"total_pembayaran": transaction.TotalAmount,
			"status":           transaction.Status,
			"qris_image":       base64QR,
		},
	})
}

func GetTransactionStatus(c *gin.Context) {
	transactionID := c.Param("id")

	expireLimit := time.Now().Add(-5 * time.Minute) 
	config.DB.Model(&models.Transaction{}).
		Where("status = ? AND created_at < ?", "pending", expireLimit).
		Update("status", "gagal")

	var transaction models.Transaction
	if err := config.DB.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status transaksi berhasil diambil",
		"data": gin.H{
			"status": transaction.Status,
		},
	})
}

func GetAllTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	expireLimit := time.Now().Add(-5 * time.Minute) 
	
	config.DB.Model(&models.Transaction{}).
		Where("status = ? AND created_at < ?", "pending", expireLimit).
		Update("status", "gagal")

	var currentUser models.User
	if err := config.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var transactions []models.Transaction
	query := config.DB.Preload("User").Preload("Campaign")

	if currentUser.Role == "user" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	var responseData []gin.H
	for _, t := range transactions {
		responseData = append(responseData, gin.H{
			"id":      t.ID,
			"user":    t.User.Name,
			"program": t.Campaign.Title,
			"amount":  t.TotalAmount,
			"date":    t.CreatedAt.Format("02 Jan 2006"),
			"status":  t.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data transaksi berhasil diambil",
		"data":    responseData,
	})
}