package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GhostDB is the pgx storage client for Ghost Postgres (TimescaleDB).
type GhostDB struct {
	pool *pgxpool.Pool
}

// NewGhostDB creates a new connection pool using the Ghost DATABASE_URL.
// Does NOT create tables — call AutoMigrate separately on startup.
func NewGhostDB(ctx context.Context, databaseURL string) (*GhostDB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	// Verify the connection is live.
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ghost db ping: %w", err)
	}
	return &GhostDB{pool: pool}, nil
}

// AutoMigrate creates all application tables if they don't already exist.
// Safe to call on every server boot — uses CREATE TABLE IF NOT EXISTS.
// MUST be called before any other DB operation.
func (db *GhostDB) AutoMigrate(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS scan_runs (
			id          TEXT PRIMARY KEY,
			started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			finished_at TIMESTAMPTZ
		);

		CREATE TABLE IF NOT EXISTS pricing_snapshots (
			id           SERIAL PRIMARY KEY,
			scan_run_id  TEXT NOT NULL REFERENCES scan_runs(id),
			competitor   TEXT NOT NULL,
			raw_html     TEXT,
			from_cache   BOOLEAN NOT NULL DEFAULT FALSE,
			scraped_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("AutoMigrate: %w", err)
	}
	return nil
}

// NewScanRun inserts a new row in scan_runs and returns the scan_run_id (UUID v4).
func (db *GhostDB) NewScanRun(ctx context.Context) (string, error) {
	// Generate UUID in Postgres to avoid a Go UUID dep.
	var id string
	err := db.pool.QueryRow(ctx,
		`INSERT INTO scan_runs (id, started_at) VALUES (gen_random_uuid()::text, NOW()) RETURNING id`,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("NewScanRun: %w", err)
	}
	return id, nil
}

// SaveSnapshot inserts one row into pricing_snapshots for a given scan run.
func (db *GhostDB) SaveSnapshot(ctx context.Context, scanRunID, competitor, rawHTML string, fromCache bool) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO pricing_snapshots (scan_run_id, competitor, raw_html, from_cache, scraped_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		scanRunID, competitor, rawHTML, fromCache,
	)
	if err != nil {
		return fmt.Errorf("SaveSnapshot %s: %w", competitor, err)
	}
	return nil
}

// FinishScanRun sets finished_at on the scan run row.
func (db *GhostDB) FinishScanRun(ctx context.Context, scanRunID string) error {
	_, err := db.pool.Exec(ctx,
		`UPDATE scan_runs SET finished_at = NOW() WHERE id = $1`,
		scanRunID,
	)
	if err != nil {
		return fmt.Errorf("FinishScanRun: %w", err)
	}
	return nil
}
