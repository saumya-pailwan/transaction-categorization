package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
)

var Pool *pgxpool.Pool

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	// Bootstrap: Create extension if not exists using a temporary connection
	// We need to do this before creating the pool with AfterConnect,
	// because RegisterTypes fails if the type doesn't exist.
	tempConfig, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL for bootstrap: %w", err)
	}
	tempConn, err := pgx.ConnectConfig(context.Background(), tempConfig)
	if err != nil {
		// If we can't connect, maybe DB isn't up, but let the pool creation handle the error closer to standard flow
		// or just return here.
		return fmt.Errorf("unable to connect for bootstrap: %w", err)
	}
	_, err = tempConn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS vector")
	tempConn.Close(context.Background())
	if err != nil {
		return fmt.Errorf("failed to bootstrap vector extension: %w", err)
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Register pgvector types
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvector.RegisterTypes(ctx, conn)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func Migrate() error {
	ctx := context.Background()

	_, err := Pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}

	_, err = Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			plaid_id TEXT UNIQUE NOT NULL,
			amount NUMERIC NOT NULL,
			date DATE NOT NULL,
			description TEXT NOT NULL,
			merchant_name TEXT,
			category TEXT,
			embedding vector(1536),
			is_manual BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	_, err = Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rules (
			id SERIAL PRIMARY KEY,
			pattern TEXT NOT NULL,
			category TEXT NOT NULL,
			priority INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create rules table: %w", err)
	}

	return nil
}
