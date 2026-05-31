package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware middleware untuk check role user
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from JWT token
		// Asumsi: role disimpan di claim JWT dengan key "role"
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role tidak ditemukan"})
			c.Abort()
			return
		}
		
		role, ok := roleInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role tidak valid"})
			c.Abort()
			return
		}
		
		// Check if user has required role
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses ditolak"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}