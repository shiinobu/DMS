package services

import "dms/backend/internal/models"

type DeviceStatusChangedEvent struct {
	Type	 string              `json:"type"`
	DeviceID string              `json:"device_id"`
	Status   models.DeviceStatus `json:"status"`
	Device   models.Device       `json:"device"`
}