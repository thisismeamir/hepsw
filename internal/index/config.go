package index

import "time"

type Config struct {
	DatabaseURL string
	AuthToken   string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	CacheTTL    time.Duration
	EnableCache bool
}
