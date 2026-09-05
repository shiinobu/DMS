package config

import (
	"os"
	"testing"
	"time"
)

func TestParseDurationEnv(t *testing.T) {
	t.Setenv("TEST_DURATION", "30s")

	duration, err := parseDurationEnv("TEST_DURATION", "1m")
	if err != nil {
		t.Fatalf("parseDurationEnv() error = %v", err)
	}

	if duration != 30*time.Second {
		t.Errorf("duration = %v, want 30s", duration)
	}
}

func TestParseDurationEnvUsesFallback(t *testing.T) {
	os.Unsetenv("TEST_DURATION_FALLBACK")

	duration, err := parseDurationEnv("TEST_DURATION_FALLBACK", "5s")
	if err != nil {
		t.Fatalf("parseDurationEnv() error = %v", err)
	}

	if duration != 5*time.Second {
		t.Errorf("duration = %v, want 5s", duration)
	}
}

func TestParseDurationEnvRejectsInvalidValue(t *testing.T) {
	t.Setenv("TEST_DURATION_INVALID", "not-a-duration")

	if _, err := parseDurationEnv("TEST_DURATION_INVALID", "5s"); err == nil {
		t.Fatal("parseDurationEnv() expected an error for an invalid duration")
	}
}
