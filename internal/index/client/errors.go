package client

import "errors"

var (
	// Configuration errors
	ErrMissingDatabaseURL = errors.New("database URL is required")
	ErrMissingAuthToken   = errors.New("auth token is required")

	// Query errors
	ErrPackageNotFound = errors.New("package not found")
	ErrVersionNotFound = errors.New("version not found")
	ErrInvalidQuery    = errors.New("invalid query parameters")
	ErrDatabaseError   = errors.New("database error")

	// Connection errors
	ErrConnectionFailed   = errors.New("failed to connect to database")
	ErrTimeout            = errors.New("operation timed out")
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")
)
