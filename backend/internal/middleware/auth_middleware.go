package middleware

import (
	"net/http"
	"strings"

	"dms/backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "authorization header is required",
			})

			c.Abort()
			return
		}

		parts := strings.SplitN(
			authHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid authorization header",
			})

			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ParseToken(
			tokenString,
			secret,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid or expired token",
			})

			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "role information is missing",
			})

			c.Abort()
			return
		}

		role, ok := roleValue.(string)

		if !ok || role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "admin access required",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
