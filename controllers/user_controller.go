package controllers

import (
	"donation-api/config"
	"donation-api/models"
	"donation-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateProfileInput struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"name":        user.Name,
			"email":       user.Email,
			"phone":       user.Phone,
			"role":        user.Role,
			"is_google":   user.GoogleID != nil,
		},
	})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isGoogleAccount := user.GoogleID != nil

	if input.Name != "" {
		user.Name = input.Name
	}

	if isGoogleAccount {
		if (input.Email != "" && input.Email != user.Email) || input.NewPassword != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Akun yang terintegrasi dengan Google tidak dapat mengubah Email atau Password"})
			return
		}
	} else {

		if input.Email != "" && input.Email != user.Email {
			var existing models.User
			if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan oleh pengguna lain"})
				return
			}
			user.Email = input.Email
		}

		if input.NewPassword != "" {
			if input.OldPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama wajib diisi untuk keamanan"})
				return
			}
			if input.NewPassword != input.ConfirmNewPassword {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Password baru dan konfirmasi tidak cocok"})
				return
			}
			if match := utils.CheckPasswordHash(input.OldPassword, *user.Password); !match {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
				return
			}
			hashedPassword, _ := utils.HashPassword(input.NewPassword)
			user.Password = &hashedPassword
		}
	}

	if input.Phone != "" {
		isSamePhone := false
		if user.Phone != nil && *user.Phone == input.Phone {
			isSamePhone = true
		}

		if !isSamePhone {
			var existing models.User
			if err := config.DB.Where("phone = ?", input.Phone).First(&existing).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor telepon sudah digunakan oleh pengguna lain"})
				return
			}
			user.Phone = &input.Phone
		}
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembaruan profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui",
		"data": gin.H{
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
		},
	})
}

type UpdateRoleInput struct {
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	var usersResponse []gin.H
	for _, u := range users {

		var successCount int64
		config.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND status = ?", u.ID, "success").
			Count(&successCount)

		usersResponse = append(usersResponse, gin.H{
			"id":                u.ID,
			"name":              u.Name,
			"email":             u.Email,
			"phone":             u.Phone,
			"role":              u.Role,
			"is_active":         u.IsActive, 
			"transaction_count": successCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data users berhasil diambil",
		"data":    usersResponse,
	})
}

func UpdateUserRole(c *gin.Context) {
	targetUserID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"})
		return
	}

	if input.Role != "" {
		user.Role = input.Role
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user berhasil diperbarui",
		"data": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"role":      user.Role,
			"is_active": user.IsActive,
		},
	})
}

func DeleteUser(c *gin.Context) {
	targetUserID := c.Param("id")

	if err := config.DB.Delete(&models.User{}, targetUserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus permanen"})
}