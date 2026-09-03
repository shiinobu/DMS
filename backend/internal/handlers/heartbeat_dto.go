package handlers

type HeartbeatRequest struct {
	IPAddress *string `json:"ip_address"`
}