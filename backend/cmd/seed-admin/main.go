package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dms/backend/internal/config"
	"dms/backend/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" {
		username = "admin"
	}

	if password == "" {
		password = "admin123"
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		INSERT INTO users (
			username,
			password_hash,
			role
		)
		VALUES ($1, $2, 'ADMIN')
		ON CONFLICT (username)
		DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			updated_at = NOW()
	`

	_, err = db.Exec(
		context.Background(),
		query,
		username,
		string(passwordHash),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin user seeded successfully")
	fmt.Println("Username:", username)
}