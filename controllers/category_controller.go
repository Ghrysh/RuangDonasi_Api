package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name: input.Name,
		Slug: input.Slug,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}