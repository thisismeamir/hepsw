package client

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index/cache"
	"github.com/thisismeamir/hepsw/internal/index/queries"
	"github.com/thisismeamir/hepsw/internal/utils"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Client is the main interface to the HepSW package index
type Client struct {
	IndexConfig *configuration.IndexConfig
	db          *sql.DB
	queries     *queries.Queries
	cache       *cache.Cache
}

var tableIDColumns = map[string]string{
	"packages":     "id",
	"versions":     "id",
	"dependencies": "id",
}

type TableSchema struct {
	Name      string
	CreateSQL string
	IDColumn  string
}

// Sync synchronizes the local index with the remote Turso database.
// Only the local database is updated — remote is always the source of truth.
// Sync is incremental: only rows created/updated after config.LastSync are pulled.
// Call this on-demand; it is not timer-driven.
func (c *Client) Sync(config *configuration.Configuration) error {
	localDB, localErr := OpenLocalDatabase()
	if localErr != nil {
		return fmt.Errorf("cannot open local database: %w\n  hint: run with --force to recreate it", localErr)
	}
	defer localDB.Close()

	remoteDB, remoteErr := OpenRemoteDatabase(config)
	if remoteErr != nil {
		return fmt.Errorf("cannot connect to remote database: %w\n  hint: check your internet connection or auth token", remoteErr)
	}
	defer remoteDB.Close()

	tables, err := fetchRemoteTables(remoteDB)
	if err != nil {
		return fmt.Errorf("failed to read remote schema: %w", err)
	}

	// Initialize LastSeenIDs map if needed
	if config.IndexConfig.LastSeenIDs == nil {
		config.IndexConfig.LastSeenIDs = make(map[string]int64)
	}

	// Sync each table
	for _, table := range tables {
		lastID := config.IndexConfig.LastSeenIDs[table.Name]
		newLastID, err := syncTable(localDB, remoteDB, table, lastID)
		if err != nil {
			return fmt.Errorf("failed to sync table %q: %w", table.Name, err)
		}
		if newLastID > lastID {
			config.IndexConfig.LastSeenIDs[table.Name] = newLastID
		}
	}

	if err := config.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: sync succeeded but could not save progress: %v\n", err)
	}

	return nil
}

// fetchRemoteTables discovers all tables from the remote database
func fetchRemoteTables(remoteDB *sql.DB) ([]TableSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := remoteDB.QueryContext(ctx, `
        SELECT name, sql
        FROM sqlite_master
        WHERE type = 'table'
          AND name NOT LIKE 'sqlite_%'
          AND name NOT LIKE 'libsql_%'
          AND name NOT LIKE '_litestream_%'
        ORDER BY name
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to query remote schema: %w", err)
	}
	defer rows.Close()

	var tables []TableSchema
	for rows.Next() {
		var t TableSchema
		if err := rows.Scan(&t.Name, &t.CreateSQL); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}
		t.IDColumn = tableIDColumns[t.Name]
		tables = append(tables, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating remote tables: %w", err)
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("remote database has no tables")
	}

	return tables, nil
}

// syncTable syncs a single table from remote to local
func syncTable(localDB, remoteDB *sql.DB, table TableSchema, lastSeenID int64) (int64, error) {
	if err := ensureLocalTable(localDB, table); err != nil {
		return lastSeenID, err
	}

	rows, maxID, err := fetchRemoteRows(remoteDB, table, lastSeenID)
	if err != nil {
		return lastSeenID, err
	}
	if len(rows) == 0 {
		return lastSeenID, nil
	}

	if err := upsertRows(localDB, table, rows); err != nil {
		return lastSeenID, err
	}

	return maxID, nil
}

// ensureLocalTable creates the table locally if needed
func ensureLocalTable(localDB *sql.DB, table TableSchema) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if table.IDColumn == "" {
		// Full wipe: drop and recreate
		_, err := localDB.ExecContext(ctx, fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, table.Name))
		if err != nil {
			return fmt.Errorf("failed to drop table %q: %w", table.Name, err)
		}
	}

	createSQL := ensureIfNotExists(table.CreateSQL)
	if _, err := localDB.ExecContext(ctx, createSQL); err != nil {
		return fmt.Errorf("failed to create table %q: %w", table.Name, err)
	}

	return nil
}

// ensureIfNotExists adds IF NOT EXISTS to CREATE TABLE
func ensureIfNotExists(createSQL string) string {
	upper := strings.ToUpper(createSQL)
	if strings.Contains(upper, "IF NOT EXISTS") {
		return createSQL
	}
	idx := strings.Index(upper, "CREATE TABLE")
	if idx == -1 {
		return createSQL
	}
	insertAt := idx + len("CREATE TABLE")
	return createSQL[:insertAt] + " IF NOT EXISTS" + createSQL[insertAt:]
}

// fetchRemoteRows fetches rows from remote table
func fetchRemoteRows(remoteDB *sql.DB, table TableSchema, lastSeenID int64) ([]map[string]any, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var query string
	if table.IDColumn != "" {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE "%s" > %d ORDER BY "%s" ASC`,
			table.Name, table.IDColumn, lastSeenID, table.IDColumn,
		)
	} else {
		query = fmt.Sprintf(`SELECT * FROM "%s"`, table.Name)
	}

	sqlRows, err := remoteDB.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query table %q: %w", table.Name, err)
	}
	defer sqlRows.Close()

	cols, err := sqlRows.Columns()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get columns: %w", err)
	}

	idColIdx := -1
	if table.IDColumn != "" {
		for i, c := range cols {
			if c == table.IDColumn {
				idColIdx = i
				break
			}
		}
	}

	var (
		result []map[string]any
		maxID  int64
	)

	for sqlRows.Next() {
		scanPtrs := make([]any, len(cols))
		scanVals := make([]any, len(cols))
		for i := range scanVals {
			scanPtrs[i] = &scanVals[i]
		}

		if err := sqlRows.Scan(scanPtrs...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]any, len(cols))
		for i, col := range cols {
			row[col] = scanVals[i]
		}

		if idColIdx >= 0 {
			if id, ok := toInt64(scanVals[idColIdx]); ok && id > maxID {
				maxID = id
			}
		}

		result = append(result, row)
	}

	if err := sqlRows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return result, maxID, nil
}

// upsertRows writes rows to local database
func upsertRows(localDB *sql.DB, table TableSchema, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := localDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cols := make([]string, 0, len(rows[0]))
	for col := range rows[0] {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	placeholders := make([]string, len(cols))
	quotedCols := make([]string, len(cols))
	for i, col := range cols {
		placeholders[i] = "?"
		quotedCols[i] = fmt.Sprintf(`"%s"`, col)
	}

	stmt := fmt.Sprintf(
		`INSERT OR REPLACE INTO "%s" (%s) VALUES (%s)`,
		table.Name,
		strings.Join(quotedCols, ", "),
		strings.Join(placeholders, ", "),
	)

	prepared, err := tx.PrepareContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer prepared.Close()

	for _, row := range rows {
		vals := make([]any, len(cols))
		for i, col := range cols {
			vals[i] = row[col]
		}
		if _, err := prepared.ExecContext(ctx, vals...); err != nil {
			return fmt.Errorf("failed to upsert row: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// toInt64 converts various types to int64
func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int64:
		return val, true
	case float64:
		return int64(val), true
	case []byte:
		n, err := strconv.ParseInt(string(val), 10, 64)
		return n, err == nil
	case string:
		n, err := strconv.ParseInt(val, 10, 64)
		return n, err == nil
	}
	return 0, false
}

// OpenRemoteDatabase opens a connection to the remote Turso database.
func OpenRemoteDatabase(config *configuration.Configuration) (*sql.DB, error) {
	if err := config.ValidateRemote(); err != nil {
		return nil, err
	}

	// libsql driver expects: libsql://[host]?authToken=[token]
	// DatabaseURL should already be in the form libsql://xxx.turso.io
	connStr := fmt.Sprintf("%s?authToken=%s", config.IndexConfig.DatabaseURL, config.IndexConfig.AuthToken)
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrConnectionFailed, err)
	}

	// Remote Turso — HTTP/WS underneath, not a traditional pool.
	// 1 conn is enough; Turso handles concurrency on their end.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), config.IndexConfig.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %v", utils.ErrConnectionFailed, err)
	}

	return db, nil
}

// OpenLocalDatabase opens (or creates) the local SQLite index database.
func OpenLocalDatabase() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}
	dbDir := filepath.Join(homeDir, ".hepsw")
	dbPath := filepath.Join(dbDir, "index.db")

	// "file:" URI is the portable way to open a local SQLite with the libsql driver
	db, err := sql.Open("sqlite3", "file://"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open local database, try sync --force for a fresh index: %w", err)
	}

	// Local SQLite — single writer, no pool needed
	db.SetMaxOpenConns(1)

	// Verify the file is actually a valid SQLite database (catches corruption)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("local database is corrupt or unreadable, try sync --force: %w", err)
	}

	return db, nil
}

// New creates a new HepSW index client
func New(IndexConfig *configuration.IndexConfig) (*Client, error) {
	// First tries to initialize using the local database at ~/.hepsw/index.db
	db, err := OpenLocalDatabase()
	if err != nil {
		return nil, err
	}
	if IndexConfig.EnableCache {
		c := cache.New(IndexConfig.CacheTTL)
		return &Client{
			IndexConfig: IndexConfig,
			db:          db,
			queries:     queries.New(db),
			cache:       c,
		}, nil
	}

	client := &Client{
		IndexConfig: IndexConfig,
		db:          db,
		queries:     queries.New(db),
		cache:       cache.New(IndexConfig.CacheTTL),
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
		if err == utils.ErrPackageNotFound || err == utils.ErrVersionNotFound {
			return err
		}
	}

	return fmt.Errorf("%w: %v", utils.ErrMaxRetriesExceeded, lastErr)
}
