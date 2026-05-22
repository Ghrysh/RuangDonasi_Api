package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCountryCodes(c *gin.Context) {
	countries := []map[string]string{
		{"code": "ID", "dial_code": "+62", "name": "Indonesia"},
		{"code": "MY", "dial_code": "+60", "name": "Malaysia"},
		{"code": "SG", "dial_code": "+65", "name": "Singapore"},
		{"code": "TH", "dial_code": "+66", "name": "Thailand"},
		{"code": "PH", "dial_code": "+63", "name": "Philippines"},
		{"code": "VN", "dial_code": "+84", "name": "Vietnam"},
		{"code": "BN", "dial_code": "+673", "name": "Brunei Darussalam"},
		{"code": "JP", "dial_code": "+81", "name": "Japan"},
		{"code": "KR", "dial_code": "+82", "name": "South Korea"},
		{"code": "US", "dial_code": "+1", "name": "United States"},
		{"code": "GB", "dial_code": "+44", "name": "United Kingdom"},
		{"code": "SA", "dial_code": "+966", "name": "Saudi Arabia"},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data kode negara",
		"data":    countries,
	})
}