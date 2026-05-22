package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak ditemukan"})
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if userRole == "superadmin" {
			isAllowed = true
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke resource ini"})
			c.Abort()
			return
		}

		c.Next()
	}
}