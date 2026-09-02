package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"dms/internal/config"
	"dms/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "DMS backend is running",
		})
	})

	log.Printf(
		"DMS backend running on :%s",
		cfg.AppPort,
	)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}