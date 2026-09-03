package middleware

import (
	"net/http"
	"strings"

	"dms/backend/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	contextUserID   = "user_id"
	contextUsername = "username"
	contextRole     = "role"
	roleAdmin = "ADMIN"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			unauthorized(c, "authorization header is required")
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c, "invalid authorization header")
			return
		}

		claims, err := auth.ParseToken(
			parts[1],
			secret,
		)
		if err != nil {
			unauthorized(c, "invalid or expired token")
			return
		}

		c.Set(contextUserID, claims.UserID)
		c.Set(contextUsername, claims.Username)
		c.Set(contextRole, claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(contextRole)

		if !exists {
			unauthorized(c, "authentication required")
			return
		}

		if role != roleAdmin {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"success": false,
					"message": "admin access required",
				},
			)
			return
		}

		c.Next()
	}
}

func unauthorized(
	c *gin.Context,
	message string,
) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{
			"success": false,
			"message": message,
		},
	)
}
