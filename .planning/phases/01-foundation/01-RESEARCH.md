# Phase 1: Foundation - Research

**Researched:** 2026-03-27
**Domain:** Go + Connect-Go + buf + Ghost Postgres + goquery web scraping
**Confidence:** HIGH

## Summary

Phase 1 establishes the full infrastructure: a Connect-Go server wired to Ghost Postgres via pgx, buf-generated protobuf code, and a concurrent net/http scraper with goquery HTML stripping. All eight target pricing pages return 200 OK with pricing content embedded in server-side HTML — none require JS rendering for the static pricing data. Ghost is already installed and authenticated; the connection string format is standard `postgresql://` with no manual TLS config needed.

The key unknowns from STATE.md blockers are now resolved: (1) Ghost delivers a blank Postgres database — you must CREATE TABLE manually on startup; (2) pgx connects to Ghost's TimescaleDB endpoint with the raw connection string from `ghost_create` — no `?sslmode=` parameter needed; (3) buf uses `go get -tool` to install local codegen plugins into `go.mod`; (4) Macroscope integration is purely GitHub-PR-based — install the GitHub App, create PRs during the build, and Macroscope auto-reviews them.

**Primary recommendation:** Use `go get -tool` to install `protoc-gen-go` and `protoc-gen-connect-go` as Go tool dependencies — this avoids needing a separate `buf.gen.yaml` remote plugin configuration and keeps all deps in `go.mod`.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
All implementation choices are at Claude's discretion — pure infrastructure phase. Use ROADMAP phase goal, success criteria, and the architecture defined in .planning/research/ARCHITECTURE.md to guide decisions. Follow the recommended project structure from ARCHITECTURE.md (proto/, cmd/, internal/, gen/, frontend/).

### Claude's Discretion
All implementation choices at Claude's discretion.

### Deferred Ideas (OUT OF SCOPE)
None — infrastructure phase.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| EXTR-01 | Agent fetches 5-8 pre-selected AI company pricing pages concurrently via Go (net/http + goroutines) | errgroup + net/http pattern verified; all 8 target URLs return 200 OK with content |
| EXTR-02 | Raw HTML pre-processed (strip nav, footer, scripts, CSS) via goquery before LLM extraction | goquery Find().Remove() pattern confirmed; full code example available |
| EXTR-05 | Failed fetches fall back to cached HTML from local file and flag as "cached" in output | Implement with `FallbackHTML` field in `RawPage` struct; local file in `demo-data/` |
| STOR-01 | Ghost Postgres DB created autonomously by the agent at startup | `ghost_create` MCP tool confirmed: `name`, `wait` params; returns `connection_string` + `id`; blank DB — tables must be created by Go code |
| STOR-02 | Each scan run inserts new rows with scan_run_id, competitor, scraped_at timestamp | Schema with `scan_runs` + `pricing_snapshots` tables; pgxpool verified connecting to Ghost |
| STOR-03 | Historical scan data is queryable for change detection comparison | `scan_run_id` FK + `scraped_at` timestamp on `pricing_snapshots` enables this |
| SPNS-01 | Ghost is used for all persistent data storage with autonomous DB creation | Ghost MCP tool `ghost_create` confirmed; Go code calls CREATE TABLE on first startup |
| SPNS-04 | Macroscope is connected for code review during the build process | GitHub App install + create PRs; `@Macroscope-app review` comment triggers manual review |
</phase_requirements>

## Project Constraints (from CLAUDE.md)

| Directive | Detail |
|-----------|--------|
| Tech stack | Go backend (Connect-Go, pgx, goquery) + Streamlit frontend (display only). No Playwright. |
| No Playwright/Selenium | Explicit constraint. Use net/http + goquery + LLM extraction. |
| No LangChain | Use raw net/http for LLM calls. |
| No ORM | Use pgx with raw SQL. Schema is 3-4 tables. |
| No gin/echo/fiber | Connect-Go already provides the HTTP server. |
| No gRPC on Python side | Connect-Go serves JSON; call with `requests`. |
| Protobuf contract-first | `.proto` files are single source of truth. Generate before writing handlers. |
| Streamlit is display-only | No direct DB or Aerospike access from Python. All data flows through Go API. |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| connectrpc.com/connect | v1.19.1 | Connect-Go HTTP/JSON RPC server | Serves HTTP/JSON natively — Streamlit calls with plain `requests` |
| github.com/jackc/pgx/v5 | v5.9.1 | Postgres driver for Ghost DB | Best Go Postgres driver; native types, connection pool |
| github.com/PuerkitoBio/goquery | v1.12.0 | HTML stripping before LLM extraction | jQuery-like CSS selector API; Remove() for nav/footer/script/style |
| google.golang.org/protobuf | v1.36.11 | Protobuf runtime | Required by Connect-Go |
| golang.org/x/sync | v0.20.0 | errgroup for concurrent scraping | Structured goroutine lifecycle; propagates first error |
| github.com/joho/godotenv | v1.5.1 | Load .env in development | Load API keys from .env; production reads env directly |
| buf CLI | v1.66.1 | Protobuf linting and codegen | `buf generate` from .proto to Go stubs |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| connectrpc.com/connect/cmd/protoc-gen-connect-go | latest | buf codegen plugin for Connect-Go handlers | Required by `buf generate`; installed via `go get -tool` |
| google.golang.org/protobuf/cmd/protoc-gen-go | latest | buf codegen plugin for Go structs | Required by `buf generate`; installed via `go get -tool` |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| buf remote plugins | buf local plugins via `go get -tool` | Remote plugins require no local install but add network latency. Local plugins installed via `go get -tool` are faster and offline-capable. For hackathon: use local plugins via `go get -tool`. |
| pgxpool | pgx.Connect | pgxpool is concurrency-safe for multi-goroutine handler use; pgx.Connect is single-connection |
| errgroup | sync.WaitGroup | errgroup propagates first error automatically; simpler for concurrent scraping |

**Installation:**
```bash
# Install buf CLI
brew install bufbuild/buf/buf

# Initialize Go module
go mod init github.com/yourname/pricing-radar

# Core runtime deps
go get connectrpc.com/connect
go get github.com/jackc/pgx/v5
go get github.com/PuerkitoBio/goquery
go get golang.org/x/sync
go get github.com/joho/godotenv

# Codegen tools (installed into go.mod as tool deps)
go get -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest
go get -tool connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

**Version verification:** Verified against Go module proxy on 2026-03-27.

## Architecture Patterns

### Recommended Project Structure
```
pricing-radar/
├── proto/
│   └── pricing/
│       └── v1/
│           └── pricing.proto         # API contract (single source of truth)
├── cmd/
│   └── server/
│       └── main.go                   # Go entrypoint — starts Connect-Go server
├── internal/
│   ├── handler/
│   │   └── pricing.go                # Connect-Go RPC handler implementations
│   ├── scraper/
│   │   ├── fetcher.go                # net/http concurrent fetcher
│   │   ├── parser.go                 # goquery HTML stripping
│   │   └── targets.go                # pre-selected page configs (URL, fallback path)
│   ├── storage/
│   │   └── ghost.go                  # pgx client, autonomous DB+table creation
│   └── config/
│       └── config.go                 # env vars, target URLs, Ghost DB config
├── gen/
│   └── pricing/
│       └── v1/                       # buf-generated Go code (do not edit)
├── demo-data/
│   ├── cached/                       # local HTML fallback files
│   └── target-pages.md               # URL manifest (already exists)
├── buf.gen.yaml                      # buf codegen config
├── buf.yaml                          # buf module config
├── go.mod
├── go.sum
└── .env                              # API keys and Ghost DB connection string
```

### Pattern 1: buf config files (contract-first setup)

**What:** Define proto first, generate Go code, implement handler interface.

`buf.yaml`:
```yaml
version: v2
modules:
  - path: proto
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

`buf.gen.yaml`:
```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/yourname/pricing-radar/gen
plugins:
  - local: [go, tool, protoc-gen-go]
    out: gen
    opt: paths=source_relative
  - local: [go, tool, protoc-gen-connect-go]
    out: gen
    opt:
      - paths=source_relative
      - simple
inputs:
  - directory: proto
```

Generate: `buf generate`

Generated output:
```
gen/pricing/v1/pricing.pb.go
gen/pricing/v1/pricingv1connect/pricing.connect.go
```

### Pattern 2: Minimal pricing.proto skeleton

```protobuf
syntax = "proto3";
package pricing.v1;

option go_package = "github.com/yourname/pricing-radar/gen/pricing/v1;pricingv1";

message RunScanRequest {}

message RunScanResponse {
  string scan_run_id = 1;
  repeated CompetitorResult results = 2;
}

message CompetitorResult {
  string competitor = 1;
  string raw_html_stripped = 2;
  bool from_cache = 3;
  string scraped_at = 4;
}

service PricingService {
  rpc RunScan(RunScanRequest) returns (RunScanResponse) {}
}
```

### Pattern 3: Connect-Go server with HTTP/2 (h2c)

```go
// Source: connectrpc.com/docs/go/getting-started/
package main

import (
    "net/http"
    pricingv1connect "github.com/yourname/pricing-radar/gen/pricing/v1/pricingv1connect"
    "github.com/yourname/pricing-radar/internal/handler"
)

func main() {
    h := &handler.PricingHandler{}
    mux := http.NewServeMux()
    path, httpHandler := pricingv1connect.NewPricingServiceHandler(h)
    mux.Handle(path, httpHandler)

    p := new(http.Protocols)
    p.SetHTTP1(true)
    p.SetUnencryptedHTTP2(true)

    s := &http.Server{
        Addr:      ":8080",
        Handler:   mux,
        Protocols: p,
    }
    s.ListenAndServe()
}
```

**Note:** `http.Protocols` with `SetUnencryptedHTTP2(true)` is the Go 1.24+ way to enable h2c. This eliminates the `golang.org/x/net/http2/h2c` dependency that older tutorials reference. Go 1.26 is available — use this approach.

### Pattern 4: Ghost autonomous DB creation (startup)

```go
// Source: ghost mcp get ghost_create (verified 2026-03-27)
// Ghost MCP tool: ghost_create(name, wait=true) returns connection_string + id
// The Go server uses ghost_create output stored in .env at build time.
// Autonomous pattern: check if DB exists (ghost_list), create if not, store connection string.
//
// In practice for this project: the agent runs ghost_create once via MCP tool,
// stores the connection_string in .env as DATABASE_URL, and pgx reads it on startup.
// "Autonomous DB creation" = server creates tables if they don't exist on first boot.

func (s *Storage) AutoMigrate(ctx context.Context) error {
    _, err := s.pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS scan_runs (
            id          TEXT PRIMARY KEY,
            started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            finished_at TIMESTAMPTZ
        );
        CREATE TABLE IF NOT EXISTS pricing_snapshots (
            id          SERIAL PRIMARY KEY,
            scan_run_id TEXT NOT NULL REFERENCES scan_runs(id),
            competitor  TEXT NOT NULL,
            raw_html    TEXT,
            from_cache  BOOLEAN NOT NULL DEFAULT FALSE,
            scraped_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
        );
    `)
    return err
}
```

### Pattern 5: pgx connection to Ghost

```go
// Source: verified by running pgx v5.9.1 against Ghost TimescaleDB 2026-03-27
// Ghost connection string format: postgresql://tsdbadmin:<pass>@<id>.fgocqo9f3c.tsdb.cloud.timescale.com:<port>/tsdb
// pgxpool.New connects with NO sslmode parameter needed — Ghost handles TLS automatically.

pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
```

### Pattern 6: Concurrent scraping with errgroup + fallback

```go
// Source: pkg.go.dev/golang.org/x/sync/errgroup
type RawPage struct {
    Competitor string
    URL        string
    HTML       string
    FromCache  bool
    FetchedAt  time.Time
}

func FetchAll(ctx context.Context, targets []Target) []RawPage {
    results := make([]RawPage, len(targets))
    var g errgroup.Group
    for i, t := range targets {
        i, t := i, t
        g.Go(func() error {
            html, err := fetchWithTimeout(ctx, t.URL, 10*time.Second)
            if err != nil {
                // fallback to cached file
                html, err = os.ReadFile(t.FallbackPath)
                if err != nil {
                    return err
                }
                results[i] = RawPage{Competitor: t.Name, URL: t.URL, HTML: string(html), FromCache: true, FetchedAt: time.Now()}
                return nil
            }
            results[i] = RawPage{Competitor: t.Name, URL: t.URL, HTML: html, FromCache: false, FetchedAt: time.Now()}
            return nil
        })
    }
    g.Wait() // best-effort: log errors, don't abort entire scan
    return results
}
```

**Note on error handling:** For the scraper, best-effort is correct — a failed fetch should fall back to cache, not abort the whole scan. Use a `sync.Mutex` or channel to collect errors separately rather than returning from errgroup (which cancels all goroutines on first error).

### Pattern 7: goquery HTML stripping

```go
// Source: pkg.go.dev/github.com/PuerkitoBio/goquery (verified 2026-03-27)
func StripBoilerplate(rawHTML string) (string, error) {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
    if err != nil {
        return "", err
    }
    // Remove elements that add token cost without pricing signal
    doc.Find("nav").Remove()
    doc.Find("footer").Remove()
    doc.Find("script").Remove()
    doc.Find("style").Remove()
    doc.Find("header").Remove()
    doc.Find("noscript").Remove()
    doc.Find("[aria-hidden='true']").Remove()

    stripped, err := doc.Find("body").Html()
    if err != nil {
        return "", err
    }
    return stripped, nil
}
```

### Anti-Patterns to Avoid

- **h2c via golang.org/x/net/http2/h2c:** Older tutorials use this package. Go 1.24+ has `http.Protocols` built in. Use `p.SetUnencryptedHTTP2(true)` instead.
- **pgxpool without context cancellation:** Always pass context to pool operations; use `context.Background()` for startup, request context for handlers.
- **errgroup for best-effort fan-out:** errgroup cancels all goroutines on first error. For scraping where fallback is acceptable, collect errors manually with a mutex instead of returning them from goroutines.
- **Storing Ghost DB ID in code:** The `ghost_create` output includes both `id` and `connection_string`. Store the full `connection_string` as `DATABASE_URL` in `.env` — the Go server needs the connection string, not the Ghost database ID.
- **Omitting `CREATE TABLE IF NOT EXISTS`:** Ghost delivers a blank Postgres database. The server must create all tables on first boot. Always use `IF NOT EXISTS` so restarts are idempotent.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Concurrent goroutine lifecycle | Manual WaitGroup + error channels | `golang.org/x/sync/errgroup` | errgroup handles context cancellation and first-error propagation |
| HTML CSS selector removal | String parsing / regex | `goquery.Find().Remove()` | Handles malformed HTML, nested selectors, handles DOM correctly |
| Postgres connection pooling | Manual connection management | `pgxpool.New()` | Handles pool sizing, reconnection, prepared statement caching |
| Proto-to-Go codegen | Manual struct mirroring | `buf generate` | buf handles field mapping, JSON marshaling, versioning |

**Key insight:** The biggest time sink in this domain is wrestling with boilerplate setup. buf + `go get -tool` eliminates the "where do I put protoc plugins" problem entirely.

## Common Pitfalls

### Pitfall 1: Ghost DB connection requires no sslmode — but TimescaleDB is not vanilla Postgres
**What goes wrong:** Developer adds `?sslmode=disable` thinking it simplifies things; Ghost rejects it because it requires TLS.
**Why it happens:** Ghost runs on TimescaleDB in Timescale Cloud — TLS is mandatory and handled automatically by the client.
**How to avoid:** Use the `connection_string` from `ghost_create` verbatim. Do not append `?sslmode=disable`.
**Warning signs:** `tls: no supported versions satisfy MinVersion` or `connection refused` with a modified connection string.

### Pitfall 2: Ghost delivers a blank database — no tables pre-created
**What goes wrong:** Server boots, tries to INSERT into `scan_runs`, gets `relation does not exist`.
**Why it happens:** `ghost_create` provisions a blank Postgres instance. Ghost does not auto-create application tables.
**How to avoid:** Call `AutoMigrate()` in `main.go` before registering Connect-Go handlers, using `CREATE TABLE IF NOT EXISTS`.
**Warning signs:** `ERROR: relation "scan_runs" does not exist (SQLSTATE 42P01)`.

### Pitfall 3: buf `local` plugin mode requires `go get -tool` first
**What goes wrong:** `buf generate` fails with `could not find protoc-gen-go`.
**Why it happens:** `buf.gen.yaml` with `local: [go, tool, protoc-gen-go]` means "run via `go tool protoc-gen-go`" — the tool must be registered in `go.mod` first via `go get -tool`.
**How to avoid:** Run `go get -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest` and `go get -tool connectrpc.com/connect/cmd/protoc-gen-connect-go@latest` before `buf generate`.
**Warning signs:** `buf generate` output: `could not find protoc-gen-go in your PATH`.

### Pitfall 4: JS-heavy pricing pages return empty content
**What goes wrong:** Some pages (SonarQube, potentially others) return a Next.js shell with no pricing data in the HTML.
**Why it happens:** Pricing content is rendered client-side by React/Next.js.
**How to avoid:** All 8 target pages in `demo-data/target-pages.md` are pre-validated as having content in the HTML response. SonarQube is explicitly excluded. Pre-test all URLs before writing production fetch code.
**Warning signs:** Fetched HTML > 50KB but contains no `$` or price-related text.

### Pitfall 5: errgroup cancels all goroutines on first fetch error
**What goes wrong:** One URL times out, errgroup cancels context, all other fetches abort.
**Why it happens:** `errgroup` cancels the shared context when any goroutine returns a non-nil error.
**How to avoid:** Handle errors inside the goroutine (fall back to cache, mark `FromCache: true`), return nil from the goroutine. Use errgroup for lifecycle management, not error propagation.
**Warning signs:** Scan returns only 1-2 results instead of 8 when any single URL is slow.

### Pitfall 6: Macroscope requires an actual GitHub repo and PRs
**What goes wrong:** Developer expects a CLI tool or API to trigger reviews; nothing happens.
**Why it happens:** Macroscope is GitHub-App-based only. It reviews PRs automatically, not arbitrary code diffs.
**How to avoid:** Initialize a git repo, push to GitHub, install Macroscope GitHub App at macroscope.com, and create PRs during the build. Comment `@Macroscope-app review` on any PR to trigger manual review.
**Warning signs:** No Macroscope comments appearing on PRs after installation.

## Code Examples

### Ghost MCP tool usage (via MCP in Claude Code)

The agent calls Ghost MCP tools during setup — these are not Go function calls but MCP tool calls:

```
ghost_create(name="pricing-radar", wait=true)
→ { connection_string: "postgresql://tsdbadmin:...@...tsdb.cloud.timescale.com:31755/tsdb", id: "..." }

ghost_sql(id="...", query="SELECT table_name FROM information_schema.tables WHERE table_schema='public';")
→ verify tables exist after AutoMigrate runs
```

Store the `connection_string` in `.env` as `DATABASE_URL`.

### net/http fetch with timeout and User-Agent

```go
func fetchWithTimeout(ctx context.Context, url string, timeout time.Duration) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; PricingRadar/1.0)")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
    }

    body, err := io.ReadAll(resp.Body)
    return string(body), err
}
```

**Note:** Setting a User-Agent header is important — some CDNs (Cloudflare) block requests with Go's default `Go-http-client/1.1` agent.

### Connect-Go handler skeleton

```go
// internal/handler/pricing.go
type PricingHandler struct {
    scraper *scraper.Fetcher
    storage *storage.GhostDB
}

// Implements gen/pricing/v1/pricingv1connect.PricingServiceHandler
func (h *PricingHandler) RunScan(
    ctx context.Context,
    req *connect.Request[pricingv1.RunScanRequest],
) (*connect.Response[pricingv1.RunScanResponse], error) {
    pages := h.scraper.FetchAll(ctx)
    // Phase 1: return raw stripped HTML; extraction in Phase 2
    resp := &pricingv1.RunScanResponse{}
    for _, p := range pages {
        resp.Results = append(resp.Results, &pricingv1.CompetitorResult{
            Competitor:     p.Competitor,
            RawHtmlStripped: p.StrippedHTML,
            FromCache:      p.FromCache,
            ScrapedAt:      p.FetchedAt.Format(time.RFC3339),
        })
    }
    return connect.NewResponse(resp), nil
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| h2c via `golang.org/x/net/http2/h2c` | `http.Protocols.SetUnencryptedHTTP2(true)` | Go 1.24 | Eliminates `x/net` dependency for h2c |
| `protoc` + manual plugin PATH management | `buf generate` with `go get -tool` plugins | 2024 | All codegen deps in `go.mod`; no PATH juggling |
| `pgx.Connect` for single connection | `pgxpool.New` | pgx v5 | Concurrency-safe; required for concurrent handler goroutines |
| Remote buf plugins (BSR) | Local buf plugins via `go get -tool` | buf v2 | Offline-capable; faster for hackathon environment |

**Deprecated/outdated:**
- `golang.org/x/net/http2/h2c`: Still works but superseded by stdlib in Go 1.24+. Don't use.
- buf `version: v1` format in `buf.yaml`: Use `version: v2`. The `v1` format omits `modules:` key.

## Open Questions

1. **Go module path for this project**
   - What we know: Module path must match the `go_package_prefix` in `buf.gen.yaml`
   - What's unclear: Whether this project will be pushed to GitHub (affects module path)
   - Recommendation: Use `github.com/yourname/pricing-radar` as a placeholder; any path works for local hackathon use since there's no external import

2. **Ghost DB persistence across demo restarts**
   - What we know: Ghost databases persist between sessions (status: running); the same `DATABASE_URL` works across server restarts
   - What's unclear: Whether `ghost pause` / `ghost resume` preserves data
   - Recommendation: Keep the database running (don't pause it during the demo); `CREATE TABLE IF NOT EXISTS` handles idempotent restarts

3. **User-Agent blocking on target pages**
   - What we know: All 8 pages return 200 OK with curl and Mozilla UA
   - What's unclear: Whether Go's HTTP client will be blocked by any CDN without a UA header
   - Recommendation: Always set `User-Agent: Mozilla/5.0 (compatible; PricingRadar/1.0)` in the scraper

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Backend runtime | Yes | 1.26.1 | — |
| Ghost CLI | DB provisioning, MCP tools | Yes | v0.4.5 | — |
| Ghost account (authenticated) | ghost_create MCP tool | Yes | Verified via `ghost status` | — |
| buf CLI | Proto codegen | No | — | Install: `brew install bufbuild/buf/buf` |
| git | Version control, Macroscope PRs | Yes | 2.50.1 | — |
| uv | Streamlit frontend (Phase 4) | Yes | 0.10.12 | — |
| Macroscope GitHub App | Code review (SPNS-04) | No | — | Install at macroscope.com; requires GitHub repo |
| protoc-gen-go | buf local plugin | No | — | Install: `go get -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest` |
| protoc-gen-connect-go | buf local plugin | No | — | Install: `go get -tool connectrpc.com/connect/cmd/protoc-gen-connect-go@latest` |
| All 8 target URLs | Scraper (EXTR-01) | Yes | Verified 200 OK | Cached HTML in demo-data/cached/ |

**Missing dependencies with no fallback:**
- buf CLI — required for proto codegen, no substitute. Install first task of Wave 0.

**Missing dependencies with fallback:**
- protoc-gen-go / protoc-gen-connect-go — installed via `go get -tool` as part of project init.
- Macroscope GitHub App — requires GitHub repo creation; plan PR workflow in Wave 0.

## Sources

### Primary (HIGH confidence)
- Ghost CLI `ghost mcp list` + `ghost mcp get` — verified tool names, params, output schema for `ghost_create`, `ghost_connect`, `ghost_sql`, `ghost_schema` (2026-03-27)
- Ghost live test — created `pricing-radar-test` DB, confirmed connection string format `postgresql://tsdbadmin:...@<id>.fgocqo9f3c.tsdb.cloud.timescale.com:31755/tsdb`
- pgx v5 live test — `pgxpool.New` with raw Ghost connection string succeeds; PostgreSQL 18.3 confirmed
- goquery v1.12.0 — `Find().Remove()` pattern verified from pkg.go.dev
- Connect-Go v1.19.1 — `connectrpc.com/docs/go/getting-started/` — main.go pattern, buf.gen.yaml
- Go 1.26.1 — `http.Protocols.SetUnencryptedHTTP2(true)` confirmed available
- `go get -tool` — confirmed in `go help get` output for Go 1.26
- Target page validation — all 8 URLs return 200 OK with pricing content (curl tested 2026-03-27)

### Secondary (MEDIUM confidence)
- buf CLI quickstart at `buf.build/docs/cli/quickstart/` — `buf config init`, `buf.gen.yaml` format with `inputs:` key
- buf installation — `brew install bufbuild/buf/buf` confirmed from buf.build/docs/cli/installation/
- Macroscope setup docs at `docs.macroscope.com/setup-instructions` — GitHub App install, PR-based workflow

### Tertiary (LOW confidence)
- Macroscope `@Macroscope-app review` manual trigger comment — from WebSearch, not confirmed via official docs

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against Go module proxy 2026-03-27
- Architecture: HIGH — follows ARCHITECTURE.md, all patterns verified via live testing
- Ghost API: HIGH — inspected directly via `ghost mcp list` + `ghost mcp get` + live DB creation
- Target URL availability: HIGH — all 8 URLs tested with curl 2026-03-27
- Macroscope integration: MEDIUM — GitHub App install confirmed; `@Macroscope-app review` trigger is LOW

**Research date:** 2026-03-27
**Valid until:** 2026-04-27 (stable ecosystem; target URLs may change)
