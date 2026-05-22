package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Excerpt  string `json:"excerpt"`
	Category string `json:"category"`
	ReadTime int    `json:"read_time"`
	Color    string `json:"color"`
	Accent   string `json:"accent"`
	Image    string `json:"image"`
}

func CreateArticle(c *gin.Context) {
	var input ArticleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readTime := input.ReadTime
	if readTime == 0 {
		readTime = (len(input.Content) / 200) + 1
	}

	article := models.Article{
		Title:    input.Title,
		Content:  input.Content,
		Excerpt:  input.Excerpt,
		Category: input.Category,
		ReadTime: readTime,
		Color:    input.Color,
		Accent:   input.Accent,
		Image:    input.Image,
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat artikel baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Artikel berhasil dibuat dengan format teks Base64!",
		"data":    article,
	})
}

func GetArticles(c *gin.Context) {
	var articles []models.Article

	if err := config.DB.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": articles,
	})
}

func UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")

	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	var input ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		article.Title = input.Title
	}
	if input.Content != "" {
		article.Content = input.Content
	}
	if input.Excerpt != "" {
		article.Excerpt = input.Excerpt
	}
	if input.Category != "" {
		article.Category = input.Category
	}
	if input.ReadTime != 0 {
		article.ReadTime = input.ReadTime
	}
	if input.Color != "" {
		article.Color = input.Color
	}
	if input.Accent != "" {
		article.Accent = input.Accent
	}
	if input.Image != "" {
		article.Image = input.Image
	}

	if err := config.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil diperbarui!",
		"data":    article,
	})
}

func DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")

	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil dihapus",
	})
}