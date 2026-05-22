package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type CampaignFormInput struct {
	CategoryID   uint    `form:"category_id"`
	Title        string  `form:"title"`
	Location     string  `form:"location"`
	Description  string  `form:"description"`
	TargetAmount float64 `form:"target_amount"`
	EndDate      string  `form:"end_date"`
	IsUrgent     bool    `form:"is_urgent"`
	Fundraiser   string  `form:"fundraiser"`
	ImageURL     string  `form:"image_url"`
}

type CampaignFEResponse struct {
	ID         uint    `json:"id"`
	Kategori   string  `json:"kategori"`
	Daerah     string  `json:"daerah"`
	Judul      string  `json:"judul"`
	Terkumpul  float64 `json:"terkumpul"`
	Target     float64 `json:"target"`
	SisaHari   int     `json:"sisaHari"`
	IsUrgent   bool    `json:"isUrgent"`
	Penggalang string  `json:"penggalang"`
	Donatur    int     `json:"donatur"`
	Deskripsi  string  `json:"deskripsi"`
	ImgSeed    string  `json:"imgSeed"`
}

func CreateCampaign(c *gin.Context) {
	var input CampaignFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedEndDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD"})
		return
	}

	var imagePath string
	file, err := c.FormFile("image")
	if err == nil {
		filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)

		if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file gambar lokal"})
			return
		}
		imagePath = "/uploads/" + filename
	} else {
		if input.ImageURL != "" {
			imagePath = input.ImageURL
		} else {
			imagePath = "default"
		}
	}

	isUrgent := input.IsUrgent
	if c.PostForm("is_urgent") == "true" {
		isUrgent = true
	}

	campaign := models.Campaign{
		CategoryID:   input.CategoryID,
		Title:        input.Title,
		Location:     input.Location,
		Description:  input.Description,
		TargetAmount: input.TargetAmount,
		EndDate:      parsedEndDate,
		IsUrgent:     isUrgent,
		Fundraiser:   input.Fundraiser,
		ImageURL:     imagePath,
		Status:       "active",
	}

	if err := config.DB.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat campaign donasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Campaign donasi berhasil dibuat beserta gambar lokal!",
	})
}

// Fungsi Mengambil List Campaign
func GetCampaigns(c *gin.Context) {
	var campaigns []models.Campaign

	if err := config.DB.Preload("Category").Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data campaign"})
		return
	}

	var formattedCampaigns []CampaignFEResponse
	for _, camp := range campaigns {
		sisaHari := int(time.Until(camp.EndDate).Hours() / 24)
		if sisaHari < 0 {
			sisaHari = 0
		}

		formattedCampaigns = append(formattedCampaigns, CampaignFEResponse{
			ID:         camp.ID,
			Kategori:   camp.Category.Name,
			Daerah:     camp.Location,
			Judul:      camp.Title,
			Terkumpul:  camp.CurrentAmount,
			Target:     camp.TargetAmount,
			SisaHari:   sisaHari,
			IsUrgent:   camp.IsUrgent,
			Penggalang: camp.Fundraiser,
			Donatur:    camp.DonatorCount,
			Deskripsi:  camp.Description,
			ImgSeed:    camp.ImageURL, // FE tinggal panggil BASE_URL + ImgSeed jika diawali dengan /uploads/
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": formattedCampaigns,
	})
}

// Fungsi Menghapus Campaign
func DeleteCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	var campaign models.Campaign
	if err := config.DB.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}

	config.DB.Where("campaign_id = ?", campaignID).Delete(&models.Transaction{})

	if err := config.DB.Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Campaign berhasil dihapus",
	})
}

// Fungsi Mengubah Campaign (Bisa Ganti Gambar Juga)
func UpdateCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	var campaign models.Campaign
	if err := config.DB.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}

	var input CampaignFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// === LOGIKA EDIT/REPLACE UPLOAD GAMBAR BARU ===
	file, err := c.FormFile("image")
	if err == nil {
		filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "uploads/"+filename); err == nil {
			campaign.ImageURL = "/uploads/" + filename // Update path gambar baru ke database
		}
	}

	if input.CategoryID != 0 {
		campaign.CategoryID = input.CategoryID
	}
	if input.Title != "" {
		campaign.Title = input.Title
	}
	if input.Location != "" {
		campaign.Location = input.Location
	}
	if input.Description != "" {
		campaign.Description = input.Description
	}
	if input.TargetAmount != 0 {
		campaign.TargetAmount = input.TargetAmount
	}
	if input.Fundraiser != "" {
		campaign.Fundraiser = input.Fundraiser
	}

	if input.EndDate != "" {
		if parsedTime, err := time.Parse("2006-01-02", input.EndDate); err == nil {
			campaign.EndDate = parsedTime
		}
	}

	if c.PostForm("is_urgent") == "true" {
		campaign.IsUrgent = true
	} else if c.PostForm("is_urgent") == "false" {
		campaign.IsUrgent = false
	}

	if err := config.DB.Save(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembaruan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Campaign berhasil diperbarui!",
		"data":    campaign,
	})
}