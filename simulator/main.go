package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HeartbeatRequest struct {
	IPAddress string `json:"ip_address"`
}

type HeartbeatResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		DeviceID      string    `json:"device_id"`
		Status        string    `json:"status"`
		LastSeen      time.Time `json:"last_seen"`
		StatusChanged bool      `json:"status_changed"`
	} `json:"data"`
}

type DeviceConfig struct {
	DeviceID     string
	DeviceName   string
	SerialNumber string
	OSVersion    string
	IPAddress    string
	Location     string
}

type Config struct {
	APIURL   string
	Username string
	Password string
	Interval time.Duration
	Devices  []DeviceConfig
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"token"`
	} `json:"data"`
}

type CreateDeviceRequest struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device_name"`
	SerialNumber string `json:"serial_number"`
	OSVersion    string `json:"os_version"`
	IPAddress    string `json:"ip_address"`
	Location     string `json:"location"`
}

func loadConfig() Config {
	apiURL := getEnv(
		"DMS_API_URL",
		"http://backend:8080/api/v1",
	)
	username := getEnv(
		"DMS_USERNAME",
		"admin",
	)
	password := getEnv(
		"DMS_PASSWORD",
		"admin123",
	)

	return Config{
		APIURL:   apiURL,
		Username: username,
		Password: password,
		Interval: 10 * time.Second,
		Devices: []DeviceConfig{
			{
				DeviceID:     "DMS-001",
				DeviceName:   "Office PC 001",
				SerialNumber: "SN-DMS-001",
				OSVersion:    "Windows 11 Pro",
				IPAddress:    "192.168.1.101",
				Location:     "Office - Floor 1",
			},
			{
				DeviceID:     "DMS-002",
				DeviceName:   "Office PC 002",
				SerialNumber: "SN-DMS-002",
				OSVersion:    "Windows 11 Pro",
				IPAddress:    "192.168.1.102",
				Location:     "Office - Floor 2",
			},
			{
				DeviceID:     "DMS-003",
				DeviceName:   "Office PC 003",
				SerialNumber: "SN-DMS-003",
				OSVersion:    "Windows 10 Pro",
				IPAddress:    "192.168.1.103",
				Location:     "Office - Floor 3",
			},
			{
				DeviceID:     "DMS-004",
				DeviceName:   "Warehouse PC 001",
				SerialNumber: "SN-DMS-004",
				OSVersion:    "Windows 11 Pro",
				IPAddress:    "192.168.1.104",
				Location:     "Warehouse",
			},
			{
				DeviceID:     "DMS-005",
				DeviceName:   "IT Administrator PC",
				SerialNumber: "SN-DMS-005",
				OSVersion:    "Windows 11 Pro",
				IPAddress:    "192.168.1.105",
				Location:     "IT Room",
			},
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func login(
	client *http.Client,
	cfg Config,
) (string, error) {
	requestBody := LoginRequest{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"%s/auth/login",
		cfg.APIURL,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf(
			"login failed with status: %s",
			resp.Status,
		)
	}

	var result LoginResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf(
			"login failed: %s",
			result.Message,
		)
	}

	return result.Data.AccessToken, nil
}

func deviceExists(
	client *http.Client,
	cfg Config,
	token string,
	device DeviceConfig,
) (bool, error) {
	url := fmt.Sprintf(
		"%s/devices/%s",
		cfg.APIURL,
		device.DeviceID,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil

	case http.StatusNotFound:
		return false, nil

	default:
		return false, fmt.Errorf(
			"check device failed for %s: %s",
			device.DeviceID,
			resp.Status,
		)
	}
}

func registerDevice(
	client *http.Client,
	cfg Config,
	token string,
	device DeviceConfig,
) error {
	exists, err := deviceExists(
		client,
		cfg,
		token,
		device,
	)

	if err != nil {
		return err
	}

	if exists {
		log.Printf(
			"[REGISTER] %s already exists, skipping",
			device.DeviceID,
		)
		return nil
	}

	requestBody := CreateDeviceRequest{
		DeviceID:     device.DeviceID,
		DeviceName:   device.DeviceName,
		SerialNumber: device.SerialNumber,
		OSVersion:    device.OSVersion,
		IPAddress:    device.IPAddress,
		Location:     device.Location,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/devices",
		cfg.APIURL,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf(
			"[REGISTER] %s - %s registered",
			device.DeviceID,
			device.DeviceName,
		)
		return nil
	}

	return fmt.Errorf(
		"registration failed for %s: %s",
		device.DeviceID,
		resp.Status,
	)
}

func registerAllDevices(
	client *http.Client,
	cfg Config,
	token string,
) error {
	for _, device := range cfg.Devices {
		if err := registerDevice(
			client,
			cfg,
			token,
			device,
		); err != nil {
			return err
		}
	}

	return nil
}

func sendHeartbeat(cfg Config, device DeviceConfig) error {
	url := fmt.Sprintf(
		"%s/devices/%s/heartbeat",
		cfg.APIURL,
		device.DeviceID,
	)

	payload := HeartbeatRequest{
		IPAddress: device.IPAddress,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var result map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf(
				"heartbeat failed with HTTP status %d",
				resp.StatusCode,
			)
		}

		return fmt.Errorf(
			"heartbeat failed with HTTP status %d: %v",
			resp.StatusCode,
			result,
		)
	}

	var result HeartbeatResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	log.Printf(
		"heartbeat sent | device=%s | ip=%s | status=%s | changed=%t | last_seen=%s",
		result.Data.DeviceID,
		device.IPAddress,
		result.Data.Status,
		result.Data.StatusChanged,
		result.Data.LastSeen.Format(time.RFC3339),
	)

	return nil
}

func runDevice(cfg Config, device DeviceConfig) {
	log.Printf(
		"starting device simulator | device=%s | ip=%s",
		device.DeviceID,
		device.IPAddress,
	)

	// Send heartbeat immediately.
	if err := sendHeartbeat(cfg, device); err != nil {
		log.Printf(
			"heartbeat error | device=%s | error=%v",
			device.DeviceID,
			err,
		)
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := sendHeartbeat(cfg, device); err != nil {
			log.Printf(
				"heartbeat error | device=%s | error=%v",
				device.DeviceID,
				err,
			)
		}
	}
}

func main() {
	cfg := loadConfig()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Println("DMS Simulator starting...")

	log.Println("Logging in...")

	token, err := login(
		client,
		cfg,
	)
	if err != nil {
		log.Fatalf(
			"login failed: %v",
			err,
		)
	}

	log.Println("Login successful")

	log.Println(
		"Registering devices...",
	)

	if err := registerAllDevices(
		client,
		cfg,
		token,
	); err != nil {
		log.Fatalf(
			"device registration failed: %v",
			err,
		)
	}

	log.Println(
		"All devices registered successfully",
	)

	for _, device := range cfg.Devices {
		go runDevice(
			cfg,
			device,
		)
	}

	log.Println(
		"All device simulators started",
	)

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	signalReceived := <-signalChannel

	log.Println("--------------------------------------")
	log.Println("received signal:", signalReceived)
	log.Println("stopping all device simulators...")
	log.Println("DMS Device Simulator stopped")
}
