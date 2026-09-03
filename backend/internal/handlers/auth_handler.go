package handlers

import (
	"errors"
	"net/http"

	"dms/backend/internal/repositories"
	"dms/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(
	service *services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	token, user, err := h.service.Login(
		c.Request.Context(),
		req.Username,
		req.Password,
	)

	if err != nil {

		if errors.Is(err, services.ErrInvalidCredentials) ||
			errors.Is(err, repositories.ErrUserNotFound) {

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid username or password",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to login",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
	})
}