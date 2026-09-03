package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv              string
	AppPort             string
	DatabaseURL         string
	JWTSecret           string
	JWTExpiration       time.Duration
	HeartbeatTimeout    time.Duration
	StatusCheckInterval time.Duration
	CORSOrigin          string
}

func Load() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	jwtExpiration, err := parseDurationEnv(
		"JWT_EXPIRATION",
		"24h",
	)
	if err != nil {
		return nil, err
	}

	heartbeatTimeout, err := parseDurationEnv(
		"HEARTBEAT_TIMEOUT",
		"30s",
	)
	if err != nil {
		return nil, err
	}

	statusCheckInterval, err := parseDurationEnv(
		"STATUS_CHECK_INTERVAL",
		"5s",
	)
	if err != nil {
		return nil, err
	}

	if jwtExpiration <= 0 {
		return nil, errors.New(
			"JWT_EXPIRATION must be greater than 0",
		)
	}

	if heartbeatTimeout <= 0 {
		return nil, errors.New(
			"HEARTBEAT_TIMEOUT must be greater than 0",
		)
	}

	if statusCheckInterval <= 0 {
		return nil, errors.New(
			"STATUS_CHECK_INTERVAL must be greater than 0",
		)
	}

	return &Config{
		AppEnv:              getEnv("APP_ENV", "development"),
		AppPort:             getEnv("APP_PORT", "8080"),
		DatabaseURL:         databaseURL,
		JWTSecret:           jwtSecret,
		JWTExpiration:       jwtExpiration,
		HeartbeatTimeout:    heartbeatTimeout,
		StatusCheckInterval: statusCheckInterval,
		CORSOrigin: getEnv(
			"CORS_ORIGIN",
			"http://localhost:3000",
		),
	}, nil
}

func loadEnvFile() error {
	err := godotenv.Load()

	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return fmt.Errorf("failed to load .env: %w", err)
}

func parseDurationEnv(
	key string,
	fallback string,
) (time.Duration, error) {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid %s value %q: %w",
			key,
			value,
			err,
		)
	}

	return duration, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
