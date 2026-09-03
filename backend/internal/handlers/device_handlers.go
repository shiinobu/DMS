package handlers

import (
	"errors"
	"net/http"

	"dms/backend/internal/models"
	"dms/backend/internal/repositories"
	"dms/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *services.DeviceService
}

func NewDeviceHandler(
	service *services.DeviceService,
) *DeviceHandler {
	return &DeviceHandler{
		service: service,
	}
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var req CreateDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	device := &models.Device{
		DeviceID:     req.DeviceID,
		DeviceName:   req.DeviceName,
		SerialNumber: req.SerialNumber,
		OSVersion:    req.OSVersion,
		IPAddress:    req.IPAddress,
		Location:     req.Location,
	}

	if err := h.service.Create(
		c.Request.Context(),
		device,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create device",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    device,
	})
}

func (h *DeviceHandler) FindAll(c *gin.Context) {
	devices, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get devices",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    devices,
	})
}

func (h *DeviceHandler) FindByDeviceID(c *gin.Context) {
	deviceID := c.Param("device_id")

	device, err := h.service.FindByDeviceID(
		c.Request.Context(),
		deviceID,
	)

	if errors.Is(err, repositories.ErrDeviceNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "device not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get device",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    device,
	})
}

func (h *DeviceHandler) Update(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req UpdateDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	device := &models.Device{
		DeviceName:   req.DeviceName,
		SerialNumber: req.SerialNumber,
		OSVersion:    req.OSVersion,
		IPAddress:    req.IPAddress,
		Location:     req.Location,
	}

	err := h.service.Update(
		c.Request.Context(),
		deviceID,
		device,
	)

	if errors.Is(err, repositories.ErrDeviceNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "device not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update device",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    device,
	})
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	deviceID := c.Param("device_id")

	err := h.service.Delete(
		c.Request.Context(),
		deviceID,
	)

	if errors.Is(err, repositories.ErrDeviceNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "device not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete device",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "device deleted successfully",
	})
}

func (h *DeviceHandler) Heartbeat(c *gin.Context) {

	deviceID := c.Param("device_id")

	var req HeartbeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	device, previousStatus, err := h.service.Heartbeat(
		c.Request.Context(),
		deviceID,
		req.IPAddress,
	)

	if err != nil {

		if errors.Is(err, repositories.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "device not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to process heartbeat",
		})
		return
	}

	statusChanged := previousStatus != device.Status

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "heartbeat received",
		"data": gin.H{
			"device_id":      device.DeviceID,
			"status":         device.Status,
			"last_seen":      device.LastSeen,
			"status_changed": statusChanged,
		},
	})
}
