---
phase: 01-foundation
plan: 03
subsystem: database
tags: [ghost, timescaledb, postgres, pgx, pgxpool, storage, migration]

# Dependency graph
requires:
  - phase: 01-02
    provides: "Go server with Connect-Go handler and config struct including DatabaseURL field"

provides:
  - "Live Ghost PostgreSQL 18.3 database (id: fzt2e1otin) with scan_runs and pricing_snapshots tables"
  - "internal/storage/ghost.go — GhostDB type with NewGhostDB, AutoMigrate, NewScanRun, SaveSnapshot, FinishScanRun"
  - "AutoMigrate wired into cmd/server/main.go startup — tables created idempotently on every boot"
  - "handler.PricingHandler accepts *storage.GhostDB — ready for scraper wiring in 01-04"

affects:
  - 01-04
  - 02-extraction
  - 03-vector-store

# Tech tracking
tech-stack:
  added:
    - "github.com/jackc/pgx/v5/pgxpool — connection pooling for concurrent handler goroutines"
    - "Ghost TimescaleDB (PostgreSQL 18.3) — live cloud DB via ghost_create MCP tool"
  patterns:
    - "AutoMigrate pattern: CREATE TABLE IF NOT EXISTS for idempotent startup"
    - "pgxpool.New with raw connection string — no sslmode parameter (Ghost handles TLS)"
    - "gen_random_uuid()::text for UUID generation in Postgres (avoids Go UUID dependency)"

key-files:
  created:
    - "internal/storage/ghost.go"
  modified:
    - "cmd/server/main.go"
    - "internal/handler/pricing.go"
    - "go.mod"
    - "go.sum"
    - ".env"

key-decisions:
  - "Ghost DB id is fzt2e1otin — use for ghost_sql verification queries in future plans"
  - "Ghost host: fzt2e1otin.fgocqo9f3c.tsdb.cloud.timescale.com:39043 (password redacted)"
  - "AutoMigrate called before handler registration in main.go to guarantee tables exist at startup"
  - "pgxpool.New used instead of pgx.Connect — concurrency-safe for multi-goroutine handler use"
  - "UUID generated in Postgres via gen_random_uuid() — no external Go UUID library needed"

patterns-established:
  - "Pattern 1: GhostDB storage client — NewGhostDB connects and pings, AutoMigrate creates tables, methods use raw SQL with pgxpool"
  - "Pattern 2: Database startup sequence — NewGhostDB → AutoMigrate → register handlers"
  - "Pattern 3: Error wrapping — fmt.Errorf with %w for all DB errors"

requirements-completed: [STOR-01, STOR-02, STOR-03, SPNS-01]

# Metrics
duration: 10min
completed: 2026-03-27
---

# Phase 1 Plan 3: Ghost DB Storage Client Summary

**Ghost PostgreSQL 18.3 DB provisioned via ghost_create, pgxpool storage client with AutoMigrate creating scan_runs and pricing_snapshots tables on startup**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-03-27T21:42:00Z
- **Completed:** 2026-03-27T21:52:53Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- Ghost DB `pricing-radar` created (id: fzt2e1otin, PostgreSQL 18.3 on TimescaleDB cloud)
- `.env` updated with real `DATABASE_URL` pointing to live Ghost DB (gitignored, not committed)
- `internal/storage/ghost.go` implemented with GhostDB type using pgxpool for connection pooling
- AutoMigrate creates `scan_runs` and `pricing_snapshots` tables with correct schema and FK constraints
- `NewScanRun`, `SaveSnapshot`, `FinishScanRun` implement the full write interface for the scraper handler
- AutoMigrate wired into `cmd/server/main.go` — runs on every startup before handler registration
- `handler.PricingHandler` updated to accept `*storage.GhostDB` — ready for scraper wiring in 01-04
- Tables verified via `ghost sql fzt2e1otin` — both tables confirmed present after AutoMigrate ran

## Ghost DB Details

- **ID:** `fzt2e1otin`
- **Host:** `fzt2e1otin.fgocqo9f3c.tsdb.cloud.timescale.com:39043`
- **DB name:** `tsdb`
- **User:** `tsdbadmin`
- **Engine:** PostgreSQL 18.3 (TimescaleDB cloud)

## Table Schema Created

```sql
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
```

## Task Commits

Each task was committed atomically:

1. **Task 1: Provision Ghost DB via MCP tool and write .env** - No code commit (.env gitignored; DB provisioned via CLI)
2. **Task 2: Implement internal/storage/ghost.go with AutoMigrate, NewScanRun, SaveSnapshot** - `b42e25c` (feat)

## Files Created/Modified
- `internal/storage/ghost.go` - GhostDB type with pgxpool, AutoMigrate, NewScanRun, SaveSnapshot, FinishScanRun
- `cmd/server/main.go` - Added GhostDB initialization + AutoMigrate call before handler registration
- `internal/handler/pricing.go` - Updated NewPricingHandler to accept *storage.GhostDB
- `go.mod` / `go.sum` - Added pgxpool transitive deps (pgpassfile, pgservicefile, puddle)
- `.env` - Updated with real DATABASE_URL (gitignored)

## Decisions Made
- Used `pgxpool.New` (not `pgx.Connect`) — concurrency-safe for concurrent handler goroutines
- UUID generated in Postgres via `gen_random_uuid()::text` — avoids adding a Go UUID dependency
- Connection string used verbatim from ghost_create output — no `?sslmode=` appended (Ghost handles TLS)
- AutoMigrate runs before handler registration to guarantee table existence at startup

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Added missing pgx transitive dependencies to go.sum**
- **Found during:** Task 2 (building internal/storage/ghost.go)
- **Issue:** `go build` failed — missing go.sum entries for pgxpool, pgconn, and puddle packages
- **Fix:** Ran `go get github.com/jackc/pgx/v5/pgxpool@v5.9.1` and `go get github.com/jackc/pgx/v5/pgconn@v5.9.1`
- **Files modified:** `go.mod`, `go.sum`
- **Verification:** `go build ./...` exits 0
- **Committed in:** `b42e25c` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 blocking dependency)
**Impact on plan:** Auto-fix was necessary to resolve missing transitive deps. No scope creep.

## Issues Encountered
- Ghost DB status was `queued` immediately after creation — polled `ghost list` for ~20 seconds until `running` before running ghost_sql verification. Normal behavior.

## User Setup Required
- `.env` already has `DATABASE_URL` set to live Ghost DB. No additional user setup needed.
- Ghost DB is running and will persist across server restarts.

## Next Phase Readiness
- Storage layer complete — `internal/storage/ghost.go` exports all interfaces needed by plan 01-04
- `handler.PricingHandler` already has `db *storage.GhostDB` field — 01-04 can call storage methods directly
- Tables exist in live Ghost DB — scraper can immediately start saving snapshots

---
*Phase: 01-foundation*
*Completed: 2026-03-27*

## Self-Check: PASSED

- FOUND: internal/storage/ghost.go
- FOUND: cmd/server/main.go
- FOUND: internal/handler/pricing.go
- FOUND: .planning/phases/01-foundation/01-03-SUMMARY.md
- FOUND commit b42e25c: feat(01-03): implement GhostDB storage client
- Ghost DB tables verified: scan_runs, pricing_snapshots both confirmed via ghost_sql
