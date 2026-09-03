package handlers

import (
	"encoding/csv"
	"net/http"
	"time"

	"dms/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(
	service services.ReportService,
) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

func valueOrEmpty(
	value *string,
) string {
	if value == nil {
		return ""
	}

	return *value
}

func formatTime(
	value *time.Time,
	location *time.Location,
) string {
	if value == nil {
		return ""
	}

	return value.In(location).Format("2006-01-02 15:04:05")
}

func (h *ReportHandler) GetSummary(
	c *gin.Context,
) {
	result, err := h.service.GetSummary(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "failed to retrieve report summary",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "report summary retrieved successfully",
			"data":    result,
		},
	)
}

func (h *ReportHandler) ExportDevices(
	c *gin.Context,
) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "failed to load timezone",
			},
		)
		return
	}

	devices, err := h.service.GetDevicesForExport(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "failed to export devices",
			},
		)
		return
	}

	filename := "dms-devices-" +
		time.Now().Format("2006-01-02") +
		".csv"

	c.Header(
		"Content-Type",
		"text/csv; charset=utf-8",
	)

	c.Header(
		"Content-Disposition",
		`attachment; filename="`+
			filename+
			`"`,
	)

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	err = writer.Write([]string{
		"Device ID",
		"Device Name",
		"Serial Number",
		"OS Version",
		"IP Address",
		"Location",
		"Status",
		"Last Seen",
		"Last Online",
		"Last Offline",
		"Created At",
		"Updated At",
	})

	if err != nil {
		return
	}

	for _, device := range devices {
		row := []string{
			device.DeviceID,
			device.DeviceName,
			device.SerialNumber,
			valueOrEmpty(device.OSVersion),
			valueOrEmpty(device.IPAddress),
			valueOrEmpty(device.Location),
			string(device.Status),
			formatTime(device.LastSeen, location),
			formatTime(device.LastOnlineAt, location),
			formatTime(device.LastOfflineAt, location),
			formatTime(&device.CreatedAt, location),
			formatTime(&device.UpdatedAt, location),
		}

		if err := writer.Write(row); err != nil {
			return
		}
	}
}
