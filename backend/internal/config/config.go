package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	AppPort  string
	DatabaseURL string

	JWTSecret     string
	JWTExpiration time.Duration

	HeartbeatTimeout    time.Duration
	StatusCheckInterval time.Duration

	CORSOrigin string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiration, err := time.ParseDuration(
		getEnv("JWT_EXPIRATION", "24h"),
	)
	if err != nil {
		return nil, err
	}

	heartbeatTimeout, err := time.ParseDuration(
		getEnv("HEARTBEAT_TIMEOUT", "30s"),
	)
	if err != nil {
		return nil, err
	}

	statusCheckInterval, err := time.ParseDuration(
		getEnv("STATUS_CHECK_INTERVAL", "5s"),
	)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("APP_PORT", "8080"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: jwtExpiration,

		HeartbeatTimeout:    heartbeatTimeout,
		StatusCheckInterval: statusCheckInterval,

		CORSOrigin: getEnv(
			"CORS_ORIGIN",
			"http://localhost:3000",
		),
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}