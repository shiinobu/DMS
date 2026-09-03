package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(
	db *pgxpool.Pool,
) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Check(
	c *gin.Context,
) {
	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"success": false,
				"message": "Database is Unavailable",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "DMS Service is Running",
		},
	)
}