package middleware

import (
	"net/http"

	"github.com/entertrans/backend-bogor.git/pkg/response"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware middleware untuk check role user
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			response.SendErrorResponse(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid user role")
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		response.SendErrorResponse(c, http.StatusForbidden, "Access denied: insufficient permissions")
		c.Abort()
	}
}
