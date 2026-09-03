package models

import "time"

type ReportSummary struct {
	TotalDevices   int64          `json:"total_devices"`
	OnlineDevices  int64          `json:"online_devices"`
	OfflineDevices int64          `json:"offline_devices"`
	LastOnline     *LastOnline    `json:"last_online,omitempty"`
	LastOffline    *LastOffline   `json:"last_offline,omitempty"`
}

type LastOnline struct {
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	LastOnlineAt time.Time `json:"last_online_at"`
}

type LastOffline struct {
	DeviceID      string    `json:"device_id"`
	DeviceName    string    `json:"device_name"`
	LastOfflineAt time.Time `json:"last_offline_at"`
}