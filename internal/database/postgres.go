package database

import (
	"context"
	"fmt"
	"time"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	config "github.com/swaindhruti/pharmastock-backend/internal/config"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func (db *PostgresDB) Health(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
