package main

import (
	"context"
	"encoding/json"
	"log"

	"dms/backend/internal/config"
	"dms/backend/internal/database"
	"dms/backend/internal/handlers"
	"dms/backend/internal/middleware"
	"dms/backend/internal/repositories"
	"dms/backend/internal/services"
	"dms/backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	// Repository
	deviceRepository := repositories.NewDeviceRepository(db)
	userRepository := repositories.NewUserRepository(db)
	reportRepository := repositories.NewReportRepository(db)

	// Service
	deviceService := services.NewDeviceService(
		deviceRepository,
	)
	authService := services.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.JWTExpiration,
	)
	reportService := services.NewReportService(
		reportRepository,
	)
	statusMonitor := services.NewStatusMonitor(
		deviceRepository,
		cfg.HeartbeatTimeout,
		cfg.StatusCheckInterval,
	)

	// Handler
	deviceHandler := handlers.NewDeviceHandler(
		deviceService,
	)
	authHandler := handlers.NewAuthHandler(authService)
	reportHandler := handlers.NewReportHandler(
		reportService,
	)
	healthHandler := handlers.NewHealthHandler(db)

	// Start status monitor in background
	go statusMonitor.Start(context.Background())

	// WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	go func() {
		for event := range deviceService.StatusChanged() {

			message, err := json.Marshal(event)

			if err != nil {
				log.Println(
					"failed to marshal device event:",
					err,
				)

				continue
			}

			hub.Broadcast(message)
		}

	}()

	go func() {
		for event := range statusMonitor.StatusChanged() {

			message, err := json.Marshal(event)

			if err != nil {
				log.Println(
					"failed to marshal device status event:",
					err,
				)

				continue
			}

			hub.Broadcast(message)
		}

	}()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			cfg.CORSOrigin,
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET("/health", healthHandler.Check)

	api := router.Group("/api/v1")
	api.GET("/ws", websocket.Handler(hub))

	// Public authentication
	api.POST("/auth/login", authHandler.Login)
	api.POST("/devices/:device_id/heartbeat", deviceHandler.Heartbeat)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	protected.Use(middleware.AdminOnly())

	devices := protected.Group("/devices")
	{
		devices.POST("", deviceHandler.Create)
		devices.GET("", deviceHandler.FindAll)
		devices.GET("/:device_id", deviceHandler.FindByDeviceID)
		devices.PUT("/:device_id", deviceHandler.Update)
		devices.DELETE("/:device_id", deviceHandler.Delete)
	}

	reports := protected.Group("/reports")
	{
		reports.GET("/summary", reportHandler.GetSummary)
		reports.GET("/devices/export", reportHandler.ExportDevices)
	}

	log.Printf(
		"DMS backend running on :%s",
		cfg.AppPort,
	)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
