package utils

import "github.com/gin-gonic/gin"

// GetUserIDFromContext mengambil user_id dari context dengan aman
func GetUserIDFromContext(c *gin.Context) (uint64, bool) {
	// Coba dari berbagai kemungkinan key
	keys := []string{"user_id", "userID", "userId"}

	for _, key := range keys {
		if val, exists := c.Get(key); exists {
			switch v := val.(type) {
			case uint:
				return uint64(v), true
			case uint64:
				return v, true
			case int:
				return uint64(v), true
			case int64:
				return uint64(v), true
			case float64:
				return uint64(v), true
			}
		}
	}

	return 0, false
}
