package handlers

type CreateDeviceRequest struct {
	DeviceID     string  `json:"device_id" binding:"required"`
	DeviceName   string  `json:"device_name" binding:"required"`
	SerialNumber string  `json:"serial_number" binding:"required"`
	OSVersion    *string `json:"os_version"`
	IPAddress    *string `json:"ip_address"`
	Location     *string `json:"location"`
}

type UpdateDeviceRequest struct {
	DeviceName   string  `json:"device_name" binding:"required"`
	SerialNumber string  `json:"serial_number" binding:"required"`
	OSVersion    *string `json:"os_version"`
	IPAddress    *string `json:"ip_address"`
	Location     *string `json:"location"`
}
