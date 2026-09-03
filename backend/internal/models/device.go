package models

import "time"

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "ONLINE"
	DeviceStatusOffline DeviceStatus = "OFFLINE"
)

type Device struct {
	ID            int64        `json:"id"`
	DeviceID      string       `json:"device_id" gorm:"uniqueIndex;not null" validate:"required"`
	DeviceName    string       `json:"device_name" validate:"required"`
	SerialNumber  string       `json:"serial_number" validate:"required"`
	OSVersion     *string      `json:"os_version,omitempty"`
	IPAddress     *string      `json:"ip_address,omitempty"`
	Location      *string      `json:"location,omitempty"`
	Status        DeviceStatus `json:"status"`
	LastSeen      *time.Time   `json:"last_seen,omitempty"`
	LastOnlineAt  *time.Time   `json:"last_online_at,omitempty"`
	LastOfflineAt *time.Time   `json:"last_offline_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
