package index

import "database/sql"

type Client struct {
	config  *Config
	db      *sql.DB
	queries *Query
}
