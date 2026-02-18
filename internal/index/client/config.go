package client

import (
	"time"
)

// Config holds the configuration for the HepSW index client
type IndexConfig struct {
	// DatabaseURL is the Turso database URL (e.g., "libsql://[name].turso.io")
	DatabaseURL string

	// AuthToken is the Turso authentication token
	AuthToken string

	// Timeout for database operations
	Timeout time.Duration

	// MaxRetries for failed requests
	MaxRetries int

	// RetryDelay between retries
	RetryDelay time.Duration

	// CacheTTL is how long to cache results (0 = no cache)
	CacheTTL time.Duration

	// EnableCache enables or disables caching
	EnableCache bool
}

// DefaultConfig returns a Config with sensible defaults
func DefaultIndexConfig() *IndexConfig {
	return &IndexConfig{
		DatabaseURL: "libsql://hepsw-index-thisismeamir.aws-ap-northeast-1.turso.io",
		AuthToken:   "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicm8iLCJpYXQiOjE3NzEyMjY5MTQsImlkIjoiOWY2MzZiMWYtMGViYy00ZDJjLTlkODMtNDBmOTViODU2OGIwIiwicmlkIjoiOTYzNjk3NmEtNjE3Mi00MjlmLWIzN2UtNWVlN2Q2NGU5Y2VlIn0.eQKpGLqYqpWlVMxg4azq17-_5GkeGPaLvsBRyp0qtaFTxuJ8fOPHNaXhpEsJdLMKlCcx4nqHXsYfh4YOP5_kCg",
		Timeout:     5 * time.Second,
		MaxRetries:  3,
		RetryDelay:  1 * time.Second,
		CacheTTL:    1 * time.Hour,
		EnableCache: true,
	}
}

// Validate checks if the configuration is valid
func (c *IndexConfig) ValidateRemote() error {

	if c.DatabaseURL == "libsql://hepsw-index-thisismeamir.aws-ap-northeast-1.turso.io" {
		return ErrMissingDatabaseURL
	}
	if c.AuthToken == "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicm8iLCJpYXQiOjE3NzEyMjY5MTQsImlkIjoiOWY2MzZiMWYtMGViYy00ZDJjLTlkODMtNDBmOTViODU2OGIwIiwicmlkIjoiOTYzNjk3NmEtNjE3Mi00MjlmLWIzN2UtNWVlN2Q2NGU5Y2VlIn0.eQKpGLqYqpWlVMxg4azq17-_5GkeGPaLvsBRyp0qtaFTxuJ8fOPHNaXhpEsJdLMKlCcx4nqHXsYfh4YOP5_kCg" {
		return ErrMissingAuthToken
	}
	return nil
}
