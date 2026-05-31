// internal/middleware/jwt_middleware.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware menerima config
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract userID dan role dari claims
		userID, ok := claims["user_id"]
		if !ok {
			response.SendErrorResponse(c, http.StatusUnauthorized, "User ID not found in token")
			c.Abort()
			return
		}

		role, ok := claims["role"]
		if !ok {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("userRole", role)

		c.Next()
	}
}

// AuthMiddlewareWithSecret - VERSI LAMA untuk kompatibilitas (jika dibutuhkan)
func AuthMiddlewareWithSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			response.SendErrorResponse(c, http.StatusUnauthorized, "User ID not found in token")
			c.Abort()
			return
		}

		role, ok := claims["role"]
		if !ok {
			response.SendErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("userRole", role)

		c.Next()
	}
}
