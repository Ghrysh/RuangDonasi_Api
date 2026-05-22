package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"donation-api/utils"
	"net/http"
	"context"
	"google.golang.org/api/idtoken"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	CountryCode     string `json:"country_code" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type LoginInput struct {
	EmailOrPhone string `json:"email_or_phone" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password dan Konfirmasi Password tidak cocok!"})
		return
	}

	fullPhone := input.CountryCode + input.Phone

	var existingUser models.User
	if err := config.DB.Where("email = ? OR phone = ?", input.Email, fullPhone).First(&existingUser).Error; err == nil {
		if existingUser.Email == input.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
			return
		}
		if existingUser.Phone != nil && *existingUser.Phone == fullPhone {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor telepon sudah digunakan"})
			return
		}
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    &fullPhone,
		Password: &hashedPassword,
		Role:     "user",
	}

	config.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil!",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ? OR phone = ?", input.EmailOrPhone, input.EmailOrPhone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun tidak ditemukan atau password salah"})
		return
	}

	if user.Password == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun ini terdaftar via Google. Silakan login dengan tombol Google."})
		return
	}

	if match := utils.CheckPasswordHash(input.Password, *user.Password); !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun tidak ditemukan atau password salah"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah dinonaktifkan oleh Superadmin. Hubungi bantuan."})
		return
	}	

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"data": gin.H{
			"token": token,
			"role":  user.Role,
			"name":  user.Name,
		},
	})
}

type GoogleLoginInput struct {
    IDToken string `json:"access_token" binding:"required"`
}

func GoogleLogin(c *gin.Context) {
	var input GoogleLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Token dari Google wajib dikirim"})
		return
	}

	googleClientID := "643019908937-vm93a1clqsnfktbvq1sd2trddj5b771h.apps.googleusercontent.com" 
	
	payload, err := idtoken.Validate(context.Background(), input.IDToken, googleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Google tidak valid atau sudah kedaluwarsa"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	googleID := payload.Subject

	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		user = models.User{
			Name:     name,
			Email:    email,
			GoogleID: &googleID,
			Role:     "user",
		}
		config.DB.Create(&user)
	} else {
		if user.GoogleID == nil {
			user.GoogleID = &googleID
			config.DB.Save(&user)
		}
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Google berhasil!",
		"data": gin.H{
			"token": token,
			"role":  user.Role,
			"name":  user.Name,
		},
	})
}