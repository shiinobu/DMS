package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(42, "admin", "ADMIN", secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %q, want admin", claims.Username)
	}
	if claims.Role != "ADMIN" {
		t.Errorf("Role = %q, want ADMIN", claims.Role)
	}
	if claims.Issuer != issuer {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, issuer)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, err := GenerateToken(1, "admin", "ADMIN", "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if _, err := ParseToken(token, "wrong-secret"); err == nil {
		t.Fatal("ParseToken() expected an error for the wrong secret")
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	token, err := GenerateToken(1, "admin", "ADMIN", "test-secret", -time.Second)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if _, err := ParseToken(token, "test-secret"); err == nil {
		t.Fatal("ParseToken() expected an error for an expired token")
	}
}
