# Architecture Research

**Domain:** Competitive pricing intelligence agent (scrape → extract → normalize → detect → recommend)
**Researched:** 2026-03-27
**Confidence:** HIGH (core pipeline patterns), MEDIUM (Aerospike AVS specifics), HIGH (TrueFoundry routing)

## System Overview

Two-process architecture: Go backend (all business logic) communicates with Streamlit frontend (display only) via Protobuf/Connect-Go over HTTP/JSON.

```
┌─────────────────────────────────────────────────────────────────────┐
│                   Streamlit Frontend (Python)                        │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  Dashboard: comparison table │ change alerts │ cluster visual │   │
│  │  Calls Go API via requests (HTTP/JSON)                       │   │
│  └──────────────────────────────┬───────────────────────────────┘   │
└─────────────────────────────────┼───────────────────────────────────┘
                                  │ HTTP/JSON (Connect-Go)
                                  │ Contract: .proto files
┌─────────────────────────────────┼───────────────────────────────────┐
│                      Go Backend (Connect-Go API)                     │
│                                 │                                    │
│  ┌──────────────────────────────▼───────────────────────────────┐   │
│  │  Connect-Go RPC Handlers                                      │   │
│  │  RunScan │ GetComparison │ GetChanges │ GetRecommendation    │   │
│  └──────┬──────────┬──────────┬──────────────┬──────────────────┘   │
│         │          │          │              │                        │
│  ┌──────▼───────┐  │  ┌──────▼───────┐  ┌──▼──────────────────┐    │
│  │  Scan Runner │  │  │  Change      │  │  Pricing Architect  │    │
│  │  (goroutines)│  │  │  Detector    │  │  (RAG + strategy)   │    │
│  └──────┬───────┘  │  └──────────────┘  └─────────┬───────────┘    │
│         │          │                               │                 │
│  ┌──────▼───────┐  │                     ┌─────────▼───────────┐    │
│  │  Scraper     │  │                     │  Embedder           │    │
│  │  (net/http   │  │                     │  (pricing profile   │    │
│  │  + goquery)  │  │                     │   → vector)         │    │
│  └──────┬───────┘  │                     └─────────┬───────────┘    │
│         │          │                               │                 │
│  ┌──────▼───────┐  │                               │                 │
│  │  LLM Extract │──┘                               │                 │
│  │  (TrueFoundry│                                  │                 │
│  │   cheap)     │                                  │                 │
│  └──────┬───────┘                                  │                 │
│         │  Normalizer                              │                 │
│  ┌──────▼───────┐                                  │                 │
│  │  Normalize   │                                  │                 │
│  └──────────────┘                                  │                 │
│                                                     │                 │
│  ┌──────────────────────┐    ┌─────────────────────▼──────────────┐ │
│  │  Ghost (Postgres)    │    │  Aerospike Vector DB               │ │
│  │  via pgx             │    │  - pricing profile embeddings      │ │
│  │  - scan_runs         │    │  - strategy docs (RAG)             │ │
│  │  - pricing_snapshots │    │  - similarity clusters             │ │
│  │  - change_events     │    │                                    │ │
│  └──────────────────────┘    └────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Key Architectural Decisions

1. **Go backend owns all logic.** Streamlit never touches Ghost DB or Aerospike directly. This makes the Go service the product; Streamlit is swappable.
2. **Protobuf is the contract.** `.proto` files are the single source of truth for the API. Connect-Go serves HTTP/JSON natively, so Streamlit calls it with `requests` — no protobuf deps on the Python side.
3. **Two processes, one machine.** Go on `:8080`, Streamlit on `:8501`. No container orchestration needed for the hackathon.

### Component Responsibilities

| Component | Language | Responsibility | Implementation |
|-----------|----------|----------------|----------------|
| Connect-Go API | Go | RPC handlers, request routing | Connect-Go handlers on net/http |
| Scan Runner | Go | Trigger scrape cycle, coordinate goroutines | errgroup for concurrent fetches |
| Scraper | Go | Fetch pricing pages concurrently | net/http + goquery for HTML stripping |
| LLM Extractor | Go | Parse pricing HTML into structured data | TrueFoundry cheap model via OpenAI-compatible API |
| Normalizer | Go | Map heterogeneous pricing to common unit | Pure Go with explicit assumption recording |
| Change Detector | Go | Diff current vs previous scan | Hash + field-level diff against Ghost DB |
| Embedder | Go | Embed pricing profiles as vectors | TrueFoundry embedding model → Aerospike upsert |
| Pricing Architect | Go | RAG retrieval + strategic recommendation | TrueFoundry expensive model + Aerospike docs |
| Ghost (Postgres) | — | Persistent store for scan data, changes | Autonomous DB creation via Ghost API, pgx driver |
| Aerospike | — | Vector store for embeddings + strategy docs | AVS HNSW index |
| Streamlit Dashboard | Python | Display tables, charts, alerts | Calls Go API via requests, renders with Streamlit |

## Recommended Project Structure

```
pricing-radar/
├── proto/
│   └── pricing/
│       └── v1/
│           └── pricing.proto          # API contract (single source of truth)
├── cmd/
│   └── server/
│       └── main.go                    # Go entrypoint — starts Connect-Go server
├── internal/
│   ├── handler/
│   │   └── pricing.go                 # Connect-Go RPC handler implementations
│   ├── scraper/
│   │   ├── fetcher.go                 # net/http concurrent fetcher
│   │   ├── parser.go                  # goquery HTML stripping
│   │   └── targets.go                 # pre-selected page configs
│   ├── extractor/
│   │   └── llm.go                     # TrueFoundry cheap model extraction
│   ├── normalizer/
│   │   ├── normalize.go               # heterogeneous → common unit
│   │   └── assumptions.go             # explicit assumption recording
│   ├── detector/
│   │   └── changes.go                 # diff engine for scan comparison
│   ├── architect/
│   │   └── strategy.go                # RAG retrieval + expensive model call
│   ├── embedder/
│   │   └── embed.go                   # pricing profile → vector → Aerospike
│   ├── storage/
│   │   ├── ghost.go                   # pgx client for Ghost Postgres
│   │   └── aerospike.go               # Aerospike vector client
│   └── config/
│       └── config.go                  # env vars, model names, target URLs
├── gen/
│   └── pricing/
│       └── v1/                        # buf-generated Go code (do not edit)
├── buf.gen.yaml                       # buf codegen config
├── buf.yaml                           # buf module config
├── go.mod
├── go.sum
├── frontend/
│   ├── app.py                         # Streamlit dashboard
│   ├── pyproject.toml                 # uv project
│   └── .python-version
├── demo-data/                         # seeded pricing data + mock pages
└── .env                               # API keys (not committed)
```

### Structure Rationale

- **proto/**: Contract lives at the repo root — both Go codegen and Python developers reference it
- **cmd/server/**: Standard Go project layout — `cmd/` for entrypoints
- **internal/**: All business logic is internal (not importable by external packages)
- **gen/**: Generated protobuf code kept separate from hand-written code
- **frontend/**: Streamlit is a self-contained Python project with its own deps
- **handler/**: Thin layer that wires Connect-Go RPCs to internal business logic

## Architectural Patterns

### Pattern 1: Protobuf Contract-First Development

**What:** Define the `.proto` file first. Generate Go server stubs. Streamlit calls JSON endpoints. The proto is the single source of truth for the API shape.
**When to use:** Any multi-language system where you want type safety at the boundary.
**Trade-offs:** Proto definition upfront, but eliminates "did you change the API?" bugs between Go and Python. Connect-Go's JSON support means Python doesn't need protobuf deps.

### Pattern 2: Extract-Normalize-Store Pipeline (Go)

**What:** Separate scrape, extract, normalize, store into distinct packages with typed interfaces.
**When to use:** When source data is heterogeneous and you need auditable transformation steps.
**Trade-offs:** More packages than a monolith, but each is independently testable.

```go
// internal/handler/pricing.go
func (h *Handler) RunScan(ctx context.Context, req *connect.Request[v1.RunScanRequest]) (*connect.Response[v1.RunScanResponse], error) {
    rawPages := h.scraper.FetchAll(ctx, h.targets)      // concurrent goroutines
    extracted := h.extractor.ExtractAll(ctx, rawPages)    // TrueFoundry cheap model
    normalized := h.normalizer.NormalizeAll(extracted)     // common unit + assumptions
    scanRun := h.storage.StoreScan(ctx, normalized)        // Ghost DB
    changes := h.detector.DetectChanges(ctx, scanRun)      // diff vs previous
    h.embedder.UpsertAll(ctx, normalized)                  // Aerospike vectors
    return connect.NewResponse(&v1.RunScanResponse{...}), nil
}
```

### Pattern 3: Two-Tier LLM Routing (Go)

**What:** Route extraction to cheap model, strategy to expensive model, both via TrueFoundry's OpenAI-compatible endpoint.
**Trade-offs:** Single HTTP client, two model strings. TrueFoundry gateway handles provider dispatch.

### Pattern 4: Dual-Store Architecture

**What:** Postgres (Ghost) for relational data (scan history, changes). Aerospike for vectors (embeddings, strategy docs).
**Trade-offs:** Two clients, but clean boundary: Postgres = "what happened when", Aerospike = "what is similar".

## Data Flow

### Scan Trigger Flow

```
Streamlit: POST /pricing.v1.PricingService/RunScan (JSON)
    ↓
Go handler: creates scan_run record in Ghost Postgres
    ↓
fetcher.go: errgroup → goroutines fetch 5-8 URLs concurrently
    ↓
parser.go: goquery strips boilerplate HTML
    ↓
llm.go: concurrent TrueFoundry calls (cheap model) → structured pricing data
    ↓
normalize.go: heterogeneous → common unit + assumptions
    ↓
ghost.go: upsert pricing_snapshots for this scan_run
    ↓
embed.go: embed normalized pricing → vector, upsert to Aerospike
    ↓
changes.go: diff current snapshots against previous scan_run
    ↓
ghost.go: write change_event records
    ↓
Response JSON → Streamlit renders updated dashboard
```

### Strategic Recommendation Flow

```
Streamlit: POST /pricing.v1.PricingService/GetRecommendation (JSON)
    ↓
Go handler: query Aerospike for strategy docs (semantic search)
    ↓
strategy.go: build prompt = strategy docs + competitor context + change summary
    ↓
TrueFoundry gateway: route to expensive model (Opus/GPT-4)
    ↓
Response JSON → Streamlit displays recommendation panel
```

### Similarity Clustering Flow

```
Streamlit: POST /pricing.v1.PricingService/GetClusters (JSON)
    ↓
Go handler: Aerospike ANN query for each vendor embedding
    ↓
Response JSON → Streamlit renders cluster visual
```

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| Ghost (Postgres) | Ghost client → autonomous DB creation on startup; pgx for queries | Create tables on first run; scan_run is top-level entity |
| TrueFoundry | OpenAI-compatible REST API from Go; two model name strings | Cheap for extraction, expensive for strategy |
| Aerospike | Go client for K/V + AVS REST/gRPC for vector operations | HNSW index; batch upsert after each scan |
| Target pricing pages | net/http with conservative timeouts (10s); no auth, no JS | Pre-validate all URLs; localhost mock fallback |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Streamlit → Go | HTTP/JSON via Connect-Go (proto contract) | Streamlit uses `requests`; no protobuf deps |
| Scraper → Extractor | Go structs: `RawPage{URL, HTML, FetchedAt}` | Scraper does not parse; extractor does not fetch |
| Extractor → Normalizer | Go struct: `RawPricingData` | LLM output parsed into typed struct |
| Normalizer → Storage | Go struct: `NormalizedPricing` | Storage writes to Ghost + triggers embedder |
| Change Detector → Architect | Go struct: `ChangeEvent` | Architect invoked only when changes exist |

## Suggested Build Order

```
1. proto/pricing/v1/pricing.proto + buf generate (contract first)
        ↓
2. cmd/server/main.go + internal/config (Go server boots, env vars load)
        ↓
3. internal/storage/ghost.go (Ghost DB init, table creation, basic writes via pgx)
        ↓
4. internal/scraper/ (fetcher + parser — verify all target URLs work)
        ↓
5. internal/extractor/llm.go (TrueFoundry extraction, structured output)
        ↓
6. internal/normalizer/ (normalization logic + assumption recording)
        ↓
7. internal/handler/pricing.go RunScan (wire scrape → extract → normalize → store)
        ↓
8. internal/detector/changes.go (diff against previous scan_run)
        ↓
9. internal/storage/aerospike.go + internal/embedder/ (vector upsert, similarity query)
        ↓
10. internal/architect/strategy.go (RAG retrieval + expensive model call)
        ↓
11. frontend/app.py (Streamlit dashboard — calls Go API with requests)
        ↓
12. Demo hardening (seed data, mock pages, timing)
```

**Why this order:**
- Step 1 establishes the contract before any code — both Go and Streamlit know the API shape
- Steps 2-4 get the foundation running (server boots, DB exists, scraper works)
- Steps 5-7 build the extraction pipeline end-to-end
- Steps 8-10 add intelligence features independently
- Step 11 (dashboard) is last because it only displays — can be built against mock API responses
- Step 12 is non-negotiable demo polish

## Scaling Considerations

Hackathon scope — not a concern. Notes for posterity:

| Scale | Architecture Adjustments |
|-------|--------------------------|
| Demo (5-8 pages) | Single Go process, goroutines sufficient |
| 50-100 pages | Add worker pool with rate limiting |
| 1000+ pages | Distributed workers, message queue, incremental Aerospike indexing |

## Anti-Patterns

### Anti-Pattern 1: Streamlit Directly Accessing Ghost/Aerospike

**What people do:** Have Streamlit connect to Postgres and Aerospike directly.
**Why it's wrong:** Defeats the purpose of the Go backend. Creates two data access paths, makes the proto contract meaningless.
**Do this instead:** All data flows through the Go API. Streamlit is display-only.

### Anti-Pattern 2: Skipping the Proto Definition

**What people do:** Define the API ad-hoc in Go handlers and match it in Python by convention.
**Why it's wrong:** Inevitable drift between Go and Python. Bugs surface at demo time.
**Do this instead:** Define proto first, generate Go code, Streamlit calls the documented JSON endpoints.

### Anti-Pattern 3: Fat Streamlit with Business Logic

**What people do:** Put normalization or change detection logic in Python because "it's easier."
**Why it's wrong:** Splits business logic across two languages. Debugging crosses process boundaries.
**Do this instead:** All logic in Go. Streamlit renders what the API returns.

### Anti-Pattern 4: Single LLM for All Tasks

**What people do:** Send extraction and strategy tasks to the same expensive model.
**Why it's wrong:** Wastes cost, hides the TrueFoundry routing value (sponsor story).
**Do this instead:** Cheap model for extraction, expensive model for strategy.

---

## Sources

- [Connect-Go docs](https://connectrpc.com/docs/go/getting-started/)
- [buf CLI](https://buf.build/docs/)
- [pgx v5](https://github.com/jackc/pgx)
- [goquery](https://github.com/PuerkitoBio/goquery)
- [Aerospike Go client](https://github.com/aerospike/aerospike-client-go)
- [TrueFoundry AI Gateway](https://www.truefoundry.com/ai-gateway)

---
*Architecture research for: Pricing Radar — Go backend + Protobuf/Connect-Go + Streamlit frontend*
*Researched: 2026-03-27 (updated from Python-only to Go+Streamlit split)*
