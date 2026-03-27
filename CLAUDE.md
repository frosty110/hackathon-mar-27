<!-- GSD:project-start source:PROJECT.md -->
## Project

**Pricing Radar**

An autonomous agent that scrapes competitor pricing pages, extracts structured data via LLM, normalizes across heterogeneous pricing models (per-seat, per-token, credits, etc.), detects changes, and generates strategy-grounded responses using internal positioning docs. Built as a hackathon project with a 3-minute live demo targeting judges evaluating autonomy, idea quality, technical implementation, tool use, and presentation.

**Core Value:** Continuous competitive pricing intelligence that turns a 2-day manual spreadsheet into a 38-second automated scan with strategic recommendations.

### Constraints

- **Timeline**: 8 hours total build time (hackathon)
- **Demo**: 3-minute live demo, must not break
- **Sponsors**: Must use 3+ sponsor tools in load-bearing way (Ghost, TrueFoundry, Aerospike + Macroscope for code review)
- **Tech stack**: Go backend (Connect-Go API, pgx, goquery) + Streamlit frontend (display only). Protobuf contract between them. No Playwright.
- **Reliability**: Pre-selected pages + localhost mock fallback for demo stability
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Architecture: Go Backend + Streamlit Frontend
## Recommended Stack
### Go Backend
| Technology | Purpose | Why Recommended |
|------------|---------|-----------------|
| Go 1.22+ | Backend runtime | Goroutines for concurrent scraping, strong typing, single binary. User preference. |
| Connect-Go (connectrpc.com/connect) | API framework | Protobuf-first RPC that speaks HTTP/JSON natively. Streamlit calls it with plain `requests` — no protobuf deps on the Python side. |
| buf | Protobuf toolchain | Linting, codegen, breaking change detection for `.proto` files. Replaces raw `protoc`. |
| net/http | HTTP scraping | Standard library HTTP client with conservative timeouts for fetching pricing pages. No third-party HTTP client needed. |
| goquery (github.com/PuerkitoBio/goquery) | HTML parsing | jQuery-like HTML selection for stripping boilerplate (nav, footer, scripts) before LLM extraction. |
| golang.org/x/sync/errgroup | Concurrent scraping | Structured goroutine lifecycle management. Propagates first error, cancels context on failure. stdlib-adjacent, no extra dep beyond sync. |
| pgx (github.com/jackc/pgx/v5) | Postgres driver (Ghost DB) | Best Go Postgres driver. Connection pooling, native types, prepared statements. |
| net/http (raw REST calls) | LLM API client | TrueFoundry exposes an OpenAI-compatible REST endpoint. Call `/v1/chat/completions` directly via `net/http` — no SDK needed. Decode JSON response into Go structs. |
| aerospike-client-go (github.com/aerospike/aerospike-client-go/v7) | Aerospike K/V client | Official Go client for Aerospike K/V operations. For vector search, use the Aerospike AVS REST/gRPC API. |
| slog | Structured logging | stdlib in Go 1.21+. Zero-dependency structured logging. |
| godotenv (github.com/joho/godotenv) | Env var management | Load API keys from `.env` in development. Production reads env directly. |
### Protobuf Contract Layer
| Technology | Purpose | Why Recommended |
|------------|---------|-----------------|
| Protocol Buffers v3 | API contract definition | Single source of truth for all request/response types between Go and Streamlit. |
| buf CLI | Proto management | `buf lint`, `buf generate`, `buf breaking` — modern protobuf toolchain. |
| Connect-Go plugin | Go server codegen | Generates typed Go handlers from `.proto` files. |
| No Python protobuf deps needed | — | Connect-Go serves HTTP/JSON. Streamlit calls endpoints with `requests` and gets plain JSON back. |
### Streamlit Frontend (Python)
| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Python 3.11+ | — | Streamlit runtime | Required by Streamlit. |
| streamlit | 1.55.0 | Dashboard UI | Fast to build, built-in charting, sufficient for hackathon demo. |
| requests | latest | HTTP client to call Go API | Simple, synchronous — fine for Streamlit's execution model. |
| uv | latest | Dependency management | Fast installs for the thin Python layer. |
### Development Tools
| Tool | Purpose | Notes |
|------|---------|-------|
| buf | Protobuf linting + codegen | `buf generate` produces Go server stubs. |
| uv | Python dependency management | For the Streamlit frontend only. |
| Macroscope | Code review during build | 4th sponsor tool. PR review workflow. |
| air or gow | Go live reload (optional) | Hot-reload during development. Nice-to-have. |
## Installation
# Go backend setup
# Core Go deps
# Install buf CLI (protobuf toolchain)
# macOS:
# Streamlit frontend setup
# Proto codegen
# Run
## Protobuf Contract Pattern
# frontend/app.py — calls Go API with plain requests
## Alternatives Considered
| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| Connect-Go | gRPC + grpcio (Python) | If you need full protobuf typing on both sides. Adds grpcio dep to Python — overkill for a thin Streamlit layer. |
| Connect-Go | OpenAPI + codegen | If you prefer REST-first design. Weaker type guarantees than protobuf. |
| Connect-Go | plain net/http + JSON | If you want zero framework overhead. Loses the proto contract — the whole point of this architecture. |
| goquery | colly | If you need crawling (following links, pagination). goquery is lighter for targeted page parsing. |
| pgx | database/sql + pq | If you want stdlib interface. pgx is faster and has better Postgres type support. |
| net/http (raw REST) | openai-go SDK | The openai-go SDK works too, but TrueFoundry's API is plain REST — raw `net/http` has zero extra deps and keeps the LLM client explicit. |
## What NOT to Use
| Avoid | Why | Use Instead |
|-------|-----|-------------|
| Playwright / Selenium | Explicit project constraint. Even harder to set up in Go. | net/http + goquery + LLM extraction |
| LangChain (Go or Python) | Overkill abstraction. No mature Go port. | Direct OpenAI-compatible API calls via net/http |
| instructor / pydantic | Python-only structured output libraries. No role in the Go backend. | Parse LLM JSON responses into Go structs manually |
| gRPC on Python side | Adds protobuf compilation + grpcio deps to the thin Streamlit layer | Connect-Go serves JSON — call with requests |
| ORM (GORM, sqlc) | Adds codegen/migration complexity for a small schema | pgx with raw SQL — schema is 3-4 tables |
| gin / echo / fiber | Connect-Go already provides the HTTP server | Use Connect-Go's built-in net/http handler |
## Stack Patterns by Variant
- Use goroutines + `errgroup` to fetch all pricing pages concurrently
- Strip HTML with goquery (remove nav, footer, scripts) before sending to LLM
- Call TrueFoundry cheap model (Haiku/Flash) via raw `net/http` POST to the OpenAI-compatible `/v1/chat/completions` endpoint
- Decode the JSON response into Go structs (no pydantic, no instructor — plain `encoding/json`)
- Call TrueFoundry expensive model (Opus/GPT-4) via the same `net/http` pattern, just swap the model name string
- Retrieve relevant strategy docs from Aerospike vector search before prompting
- Pass normalized pricing comparison + retrieved docs as context
- Generate embeddings via TrueFoundry embedding model (same OpenAI-compatible API, `/v1/embeddings`)
- Upsert into Aerospike AVS index with HNSW via the AVS REST/gRPC API
- Query with current product's pricing profile to find nearest neighbors
- Streamlit calls Go API endpoints with `requests`
- Use `st.cache_data` to avoid re-calling the API on every widget interaction
- All data comes from the Go API — Streamlit never touches Ghost DB or Aerospike directly
- No business logic in Python: normalization, change detection, and clustering all live in Go
## TrueFoundry Integration Pattern (Go)
## Sources
- [Connect-Go docs](https://connectrpc.com/docs/go/getting-started/) — HTTP/JSON + gRPC on same port
- [buf CLI](https://buf.build/docs/) — modern protobuf toolchain
- [pgx v5 docs](https://github.com/jackc/pgx) — best Go Postgres driver
- [goquery](https://github.com/PuerkitoBio/goquery) — jQuery-like HTML parsing for Go
- [Aerospike Go client](https://github.com/aerospike/aerospike-client-go) — official client
- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) — structured goroutine concurrency
- [TrueFoundry AI Gateway](https://www.truefoundry.com/ai-gateway) — OpenAI-compatible REST endpoint
- PyPI / streamlit — version 1.55.0 verified March 2026
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
