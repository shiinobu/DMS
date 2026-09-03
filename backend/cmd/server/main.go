package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Fatalf(
			"Failed to Load Config: %v",
			err,
		)
	}

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf(
			"Failed to Connect Database: %v",
			err,
		)
	}
	defer db.Close()

	deviceRepository := repositories.NewDeviceRepository(db)
	userRepository := repositories.NewUserRepository(db)
	reportRepository := repositories.NewReportRepository(db)

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

	deviceHandler := handlers.NewDeviceHandler(
		deviceService,
	)
	authHandler := handlers.NewAuthHandler(authService)
	reportHandler := handlers.NewReportHandler(
		reportService,
	)
	healthHandler := handlers.NewHealthHandler(db)

	hub := websocket.NewHub()
	go hub.Run()
	go broadcastEvents(hub, deviceService.StatusChanged())
	go broadcastEvents(hub, statusMonitor.StatusChanged())

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
	api.POST("/auth/login", authHandler.Login)
	api.POST("/devices/:device_id/heartbeat", deviceHandler.Heartbeat)

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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go statusMonitor.Start(ctx)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf(
			"DMS Backend Running on :%s",
			cfg.AppPort,
		)

		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(
				"HTTP server failed: %v",
				err,
			)
		}

	case <-ctx.Done():
		log.Println(
			"Shutting Down DMS Backend",
		)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf(
			"Failed to Shutdown HTTP Server: %v",
			err,
		)
	}
}

func broadcastEvents(
	hub *websocket.Hub,
	events <-chan services.DeviceStatusChangedEvent,
) {
	for event := range events {
		message, err := json.Marshal(event)
		if err != nil {
			log.Printf(
				"Failed to Marshal Device Status Event: %v",
				err,
			)
			continue
		}

		hub.Broadcast(message)
	}
}
