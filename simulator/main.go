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
	DeviceID  string
	IPAddress string
}

type Config struct {
	APIURL   string
	Interval time.Duration
	Devices  []DeviceConfig
}

func loadConfig() Config {
	apiURL := getEnv(
		"DMS_API_URL",
		"http://localhost:8080/api/v1",
	)

	return Config{
		APIURL:   apiURL,
		Interval: 10 * time.Second,
		Devices: []DeviceConfig{
			{
				DeviceID:  "DMS-001",
				IPAddress: "192.168.1.101",
			},
			{
				DeviceID:  "DMS-002",
				IPAddress: "192.168.1.102",
			},
			{
				DeviceID:  "DMS-003",
				IPAddress: "192.168.1.103",
			},
			{
				DeviceID:  "DMS-004",
				IPAddress: "192.168.1.104",
			},
			{
				DeviceID:  "DMS-005",
				IPAddress: "192.168.1.105",
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

	log.Println("======================================")
	log.Println("DMS Device Simulator")
	log.Println("======================================")
	log.Println("API URL:", cfg.APIURL)
	log.Println("Heartbeat interval:", cfg.Interval)
	log.Println("Total devices:", len(cfg.Devices))
	log.Println("--------------------------------------")

	for _, device := range cfg.Devices {
		log.Printf(
			"device=%s | ip=%s",
			device.DeviceID,
			device.IPAddress,
		)
	}

	log.Println("--------------------------------------")

	for _, device := range cfg.Devices {
		go runDevice(cfg, device)
	}

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