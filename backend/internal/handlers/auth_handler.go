package handlers

import (
	"errors"
	"net/http"

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
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "invalid request",
			},
		)
		return
	}

	result, err := h.service.Login(
		c.Request.Context(),
		request.Username,
		request.Password,
	)
	if err != nil {
		if errors.Is(
			err,
			services.ErrInvalidCredentials,
		) {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"success": false,
					"message": "invalid username or password",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "failed to login",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "login successful",
			"data": gin.H{
				"token": result.Token,
				"user": gin.H{
					"id":       result.User.ID,
					"username": result.User.Username,
					"role":     result.User.Role,
				},
			},
		},
	)
}
