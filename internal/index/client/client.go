package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thisismeamir/hepsw/internal/index/cache"
	"github.com/thisismeamir/hepsw/internal/index/queries"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Client is the main interface to the HepSW package index
type Client struct {
	IndexConfig *IndexConfig
	db          *sql.DB
	queries     *queries.Queries
	cache       *cache.Cache
}

// New creates a new HepSW index client
func New(IndexConfig *IndexConfig) (*Client, error) {
	if err := IndexConfig.Validate(); err != nil {
		return nil, err
	}

	// Create Turso connection string
	// Format: libsql://[url]?authToken=[token]
	connStr := fmt.Sprintf("%s?authToken=%s", IndexConfig.DatabaseURL, IndexConfig.AuthToken)

	// Open database connection
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), IndexConfig.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	client := &Client{
		IndexConfig: IndexConfig,
		db:          db,
		queries:     queries.New(db),
	}

	// Initialize cache if enabled
	if IndexConfig.EnableCache {
		client.cache = cache.New(IndexConfig.CacheTTL)
	}

	return client, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// DB returns the underlying database connection for advanced usage
func (c *Client) DB() *sql.DB {
	return c.db
}

// Queries returns the queries handler
func (c *Client) Queries() *queries.Queries {
	return c.queries
}

// Cache returns the cache handler
func (c *Client) Cache() *cache.Cache {
	return c.cache
}

// Ping checks if the database connection is alive
func (c *Client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// withRetry executes a function with retry logic
func (c *Client) withRetry(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= c.IndexConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.IndexConfig.RetryDelay):
			}
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry certain errors
		if err == ErrPackageNotFound || err == ErrVersionNotFound {
			return err
		}
	}

	return fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, lastErr)
}
