package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iammrsea/social-app/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresConfig holds configuration for PostgreSQL connection
type postgresConfig struct {
	uri          string
	maxPoolSize  int32
	minPoolSize  int32
	connIdleTime time.Duration
	timeout      time.Duration
}

// SetupPostgres sets up and connects to PostgreSQL
func SetupPostgreSQL(ctx context.Context, cf *config.PostgresConfig) (*pgxpool.Pool, func() error) {
	pool, err := connect(ctx, cf)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	disconnectPool := func() error {
		return disconnect(pool)
	}

	return pool, disconnectPool
}

// connect establishes a connection to PostgreSQL and returns the connection pool
func connect(ctx context.Context, cf *config.PostgresConfig) (*pgxpool.Pool, error) {
	// Create a context with timeout for the connection
	connectionCtx, cancel := context.WithTimeout(ctx, cf.Timeout)
	defer cancel()

	// Configure connection pool settings
	poolConfig, err := pgxpool.ParseConfig(cf.Uri)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PostgreSQL URI: %w", err)
	}

	poolConfig.MaxConns = cf.MaxPoolSize
	poolConfig.MinConns = cf.MinPoolSize
	poolConfig.MaxConnIdleTime = cf.ConnIdleTime

	log.Println("Connecting to PostgreSQL...")
	// Connect to PostgreSQL
	pool, err := pgxpool.NewWithConfig(connectionCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %w", err)
	}
	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, cf.Timeout)
	defer cancel()

	conn, err := pool.Acquire(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire connection from pool: %w", err)
	}
	defer conn.Release()

	if err := conn.Conn().Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("unable to ping PostgreSQL: %w", err)
	}
	log.Println("✅ Connected to PostgreSQL")
	return pool, nil
}

// disconnect closes the PostgreSQL connection pool
func disconnect(pool *pgxpool.Pool) error {
	if pool == nil {
		return nil
	}
	pool.Close()
	return nil
}
