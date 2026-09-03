package database

import (
	"context"
	"time"

	"dms/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns        = 10
	minConns        = 2
	maxConnLifetime = time.Hour
)

func NewPostgresPool(
	cfg *config.Config,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = maxConns
	poolConfig.MinConns = minConns
	poolConfig.MaxConnLifetime = maxConnLifetime

	db, err := pgxpool.NewWithConfig(
		context.Background(),
		poolConfig,
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
