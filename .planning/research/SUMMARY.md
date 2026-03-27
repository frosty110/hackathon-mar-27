# Project Research Summary

**Project:** Pricing Radar — Competitive Pricing Intelligence Agent
**Domain:** AI/SaaS competitive pricing intelligence (scraping + LLM extraction + vector search + dashboard)
**Researched:** 2026-03-27
**Confidence:** HIGH

## Executive Summary

Pricing Radar is a competitive pricing intelligence agent for AI startup product teams. It solves a gap that every surveyed competitor (Tierly, Prisync, Visualping, Priceva) fails to address: normalizing heterogeneous AI pricing models (per-seat, per-token, credits, flat tiers) to a common unit with transparent assumptions, then grounding strategic recommendations in the team's own positioning documents. The architecture is a two-process split: a Go backend that owns all business logic, and a Streamlit frontend that is display-only. The Go backend exposes a Connect-Go API (Protobuf contract, HTTP/JSON transport); Streamlit calls it with `requests` — no protobuf deps on the Python side. The Go stack uses net/http + goquery for concurrent scraping, pgx for Ghost Postgres, errgroup for goroutine coordination, and Connect-Go for RPC handlers. This architecture maps directly to sponsor requirements and has well-documented patterns for every integration.

The key differentiators are (1) cross-model normalization with explicit assumption transparency, (2) strategy-grounded recommendations via RAG over uploaded positioning docs, and (3) "who prices like us?" clustering via Aerospike HNSW. None of these exist in any surveyed tool. The hardest design problem is normalization: it is an epistemic task (surfacing assumptions about usage patterns), not just a unit conversion. If assumptions are hidden, the demo comparison table is factually untrustworthy and will fail under judge scrutiny.

The primary risks are: JS-rendered pricing pages that return empty HTML to net/http (mitigate by pre-validating all 8 target URLs before coding starts, with mock fallback files ready); demo latency from a live scan during the demo (mitigate by pre-running the scan and displaying cached DB results); and scope creep enabled by AI-accelerated development (mitigate with a hard feature freeze at hour 5). The architectural build order in ARCHITECTURE.md is decisive — follow it to avoid discovering at hour 4 that the DB schema doesn't support change detection.

## Key Findings

### Recommended Stack

The stack is built around a Go backend with a Python Streamlit frontend. The two processes communicate over HTTP/JSON via Connect-Go (proto contract). The Go backend owns all business logic; Streamlit is swappable display. TrueFoundry is accessed via the OpenAI-compatible endpoint (set `base_url` + `api_key`; no extra SDK). Aerospike Vector Search 4.x uses a unified client API. The proto file is the single source of truth for the API shape — defined first, Go code generated via buf.

**Core technologies:**
- Go (backend): Runtime — all business logic, scraping, extraction, normalization, detection, embedding, RAG
- Connect-Go: RPC framework — serves HTTP/JSON natively so Streamlit calls it with plain `requests`
- net/http + goquery: Concurrent scraping — goroutines via errgroup; goquery strips boilerplate HTML (replaces httpx + bs4/lxml)
- pgx v5: Ghost Postgres driver — idiomatic Go async driver; no ORM overhead needed (replaces asyncpg)
- errgroup: Goroutine coordination — structured concurrency for concurrent fetches and LLM calls
- TrueFoundry (OpenAI-compatible REST from Go): LLM routing — single HTTP client, two model name strings; cheap for extraction, expensive for strategy
- aerospike-vector-search (Go client): Vector similarity + RAG — sponsor requirement; HNSW index
- Protobuf / buf: API contract — `.proto` defines all RPC shapes; `buf generate` produces Go stubs
- Python 3.11 + Streamlit 1.55.0 (frontend): Dashboard — display-only; calls Go API via `requests`; no protobuf deps required
- uv: Python dependency management for the frontend — fast installs

### Expected Features

**Must have (table stakes):**
- Automated scraping of 5-8 pre-selected AI company pricing pages — no scraping = no product
- LLM extraction to structured data (tier names, prices, features, model type) — all downstream features require this
- Side-by-side competitor comparison table — the core mental model for the target persona (Head of Product)
- Change detection with diff display — makes the product a monitor, not a one-shot lookup
- Historical scan storage in Ghost Postgres — required for change detection and clustering
- Pricing model identification (per-seat, per-token, credits, flat) — AI-specific, absent from e-commerce tools

**Should have (competitive):**
- Cross-model normalization with explicit assumption transparency — the key differentiator; no competitor does this
- Strategy-grounded recommendations via RAG (TrueFoundry expensive model + Aerospike docs) — unique capability
- "Who prices like us?" clustering via Aerospike — vector similarity on pricing profiles; novel, no competitor offers it
- Two-tier LLM model routing — cheap model for extraction, expensive model (Opus/GPT-4) for strategy; satisfies TrueFoundry sponsor story
- Change event narrative — plain-English summary of what changed and why it matters, generated post-diff

**Defer (v2+):**
- Playwright/browser rendering for JS-heavy pages — 3-4 hour build risk, breaks demo reliability
- User authentication and multi-tenant — doubles scope; single-user demo is sufficient
- Scheduled recurring scans (cron-based) — blocked by single-user session model in hackathon build
- Competitor auto-discovery — Tierly takes 13 minutes; unreliable for a deterministic demo

### Architecture Approach

The system is a two-process architecture: Go backend (all business logic) on `:8080`, Streamlit frontend (display only) on `:8501`. No container orchestration needed for the hackathon. The Go backend exposes four Connect-Go RPC handlers: `RunScan`, `GetComparison`, `GetChanges`, `GetRecommendation`. The Protobuf contract (`.proto` files) is defined first and is the single source of truth — both Go codegen and the Streamlit developer reference it.

The key pattern is Extract-Normalize-Store (Go): `internal/scraper`, `internal/extractor`, and `internal/normalizer` are distinct packages with typed Go struct interfaces. This makes each stage independently testable and lets you retry extraction without re-fetching. A second pattern is Two-Tier LLM Routing: cheap model for templated extraction tasks, expensive model (Opus/GPT-4) only for synthesis and strategy generation — both routed through TrueFoundry's OpenAI-compatible gateway from Go. The Streamlit dashboard is read-only — it calls the Go API with `requests` and renders what is returned; it never touches Ghost DB or Aerospike directly.

**Major components:**
1. Connect-Go API (internal/handler/pricing.go) — thin RPC handlers that wire requests to internal business logic
2. Scan Runner (internal/handler/pricing.go RunScan) — orchestrates the full scrape → extract → normalize → store → detect cycle using errgroup
3. Scraper (internal/scraper/fetcher.go + parser.go) — net/http concurrent fetcher; goquery strips boilerplate HTML
4. LLM Extractor (internal/extractor/llm.go) — TrueFoundry cheap model via OpenAI-compatible REST; structured Go output
5. Normalizer (internal/normalizer/normalize.go + assumptions.go) — maps heterogeneous pricing to common unit; records explicit assumptions
6. Change Detector (internal/detector/changes.go) — diffs current normalized snapshot against previous scan in Ghost DB
7. Embedder (internal/embedder/embed.go) — generates pricing profile vectors; upserts to Aerospike HNSW index
8. Pricing Architect (internal/architect/strategy.go) — RAG retrieval from Aerospike + TrueFoundry expensive model for strategy
9. Ghost Postgres (internal/storage/ghost.go) — pgx client; scan_runs, pricing_snapshots, change_events; autonomous DB creation via Ghost API
10. Aerospike AVS (internal/storage/aerospike.go) — pricing profile embeddings for similarity clustering; strategy docs for RAG retrieval
11. Streamlit Dashboard (frontend/app.py) — comparison table, change alerts, cluster visual, transparency map; calls Go API via `requests`

### Critical Pitfalls

1. **Raw HTML overload kills extraction quality** — 80K+ tokens of nav/footer/scripts drowns the pricing signal. Pre-process with goquery to strip boilerplate before sending to the LLM; target 2K-8K tokens per page. Never send full raw HTML to the LLM under any circumstances.

2. **Hidden normalization assumptions expose methodology flaws in demo Q&A** — Normalization is an epistemic task, not just unit conversion. Every normalized Go struct must carry its normalization assumptions. Display assumption footnotes inline in the dashboard. A bare normalized number without a footnote will fail under judge scrutiny.

3. **JS-rendered pages return empty HTML to net/http silently** — Many AI company pricing pages load tiers via JavaScript. net/http returns 200 OK with an empty pricing div. Pre-validate all 8 target URLs before coding begins: check raw HTML for known price strings. Prepare mock fallback HTML files for 2-3 pages before the clock starts.

4. **Ghost DB schema designed for storage, not for change detection** — If scans are stored without a `scan_run_id`, `competitor` key, and `scraped_at` timestamp, change detection queries are impossible or require a schema rewrite in hour 4. Design the schema for the change detection query from day one.

5. **Demo scan latency** — A live scan button that re-fires on every click wastes LLM budget and freezes the UI. Pre-run the scan; display cached DB results during the demo. Gate the scan button with session state. Add `st.progress()` and `st.status()` — never show a frozen screen.

## Implications for Roadmap

Based on combined research, the architectural build order from ARCHITECTURE.md is the right phase structure. Each phase has a clear deliverable and the dependency chain is hard: each later phase fails without the earlier one working.

### Phase 0: Pre-work and Target Validation
**Rationale:** JS-rendering failures (Pitfall 3) are discovered during coding, not before. Pre-validation eliminates the highest-risk demo failure mode before any code is written.
**Delivers:** Validated list of 5-8 scrapable pricing pages; mock HTML fallbacks for any JS-heavy pages; project scaffolding (proto definition, buf.gen.yaml, go.mod, .env, internal/config)
**Addresses:** Pitfall 3 (JS pages), Pitfall 5 (scope creep) — feature list is locked before building starts
**Avoids:** Discovering mid-build that 3 of 8 target pages are JS-rendered

### Phase 1: Proto Contract + Go Server Bootstrap
**Rationale:** The proto file is the contract between Go and Streamlit. It must exist before any handler or dashboard code is written. Getting the Go server booting with env vars loaded is the foundation for every subsequent phase.
**Delivers:** proto/pricing/v1/pricing.proto; buf-generated Go stubs; cmd/server/main.go boots and listens on :8080; internal/config loads env vars
**Uses:** Connect-Go, buf CLI, net/http
**Avoids:** API drift between Go and Streamlit (Anti-Pattern 2 in ARCHITECTURE.md)

### Phase 2: Data Foundation (Ghost DB + Scraper)
**Rationale:** Ghost DB schema must exist and be correct before any data flows through it. Change detection requires the right schema from scan 1. The scraper must be verified working before LLM extraction is added on top.
**Delivers:** Ghost Postgres schema (scan_runs, pricing_snapshots, change_events) via pgx; net/http async fetcher returning raw HTML; verified URL fetches for all targets; goquery stripping boilerplate
**Uses:** pgx v5, Ghost Postgres, net/http, goquery, errgroup
**Avoids:** Pitfall 4 (schema misfit for change detection)

### Phase 3: Extraction and Normalization Pipeline
**Rationale:** Extraction and normalization are tightly coupled and share Go struct types, but must stay distinct packages. This is the integration milestone — once the RunScan handler works end-to-end with a single competitor, the rest of the pipeline is incremental.
**Delivers:** LLM extraction (TrueFoundry cheap model) producing structured Go structs; normalization to common unit with explicit assumption recording; RunScan handler wiring the full scrape → extract → normalize → store cycle
**Uses:** TrueFoundry cheap model (OpenAI-compatible REST from Go), errgroup for concurrent LLM calls
**Implements:** internal/extractor, internal/normalizer, internal/handler RunScan (partial)
**Avoids:** Pitfall 1 (raw HTML overload), Pitfall 2 (hidden normalization assumptions)

### Phase 4: Change Detection and Aerospike Integration
**Rationale:** Change detection is the feature that makes Pricing Radar a monitor rather than a one-shot tool. Aerospike integration (embeddings + similarity) can proceed in parallel once normalized data exists in Postgres.
**Delivers:** Change Detector diffing current vs. previous scan; change event storage in Ghost DB; pricing profile embeddings upserted to Aerospike HNSW; "who prices like us?" similarity query working
**Uses:** pgx v5 (diff queries against Ghost DB), Aerospike Go client, TrueFoundry embedding model
**Implements:** internal/detector, internal/embedder, internal/storage/aerospike

### Phase 5: Pricing Architect (RAG + Strategy)
**Rationale:** This phase requires working Aerospike (for doc retrieval), working change detection (to trigger recommendations), and the expensive TrueFoundry model route. It is the most visible differentiator for judges.
**Delivers:** RAG retrieval of strategy docs from Aerospike; strategy prompt construction with competitor context + change summary; TrueFoundry expensive model call; strategic recommendation returned via Connect-Go RPC
**Uses:** TrueFoundry expensive model (Opus/GPT-4), Aerospike AVS (RAG mode)
**Implements:** internal/architect, two-tier LLM routing

### Phase 6: Streamlit Dashboard and Demo Hardening
**Rationale:** Dashboard is built last because it only reads from the Go API, which must already exist. Demo hardening (pre-seeded data, session state guards, progress indicators) is non-negotiable — the demo is the product for the judges.
**Delivers:** Streamlit dashboard (frontend/app.py) with comparison table, change alerts, normalization transparency map, cluster visual — all via `requests` calls to Go API; pre-seeded scan data in DB; demo script timed under 3 minutes
**Uses:** Streamlit 1.55.0, st.cache_data, st.session_state, st.progress, requests
**Avoids:** Pitfall 5 (demo scan latency), UX pitfalls (raw diffs, unlabeled clusters)

### Phase Ordering Rationale

- Phase 0 must precede all others: JS page failures discovered mid-build cost 1-2 hours; discovered in pre-work they cost 30 minutes.
- Phase 1 establishes the contract before any handler or dashboard code — both Go and Streamlit know the API shape.
- Phases 2-3 are strict dependencies: you cannot detect changes without stored data; you cannot store data without a schema; you cannot extract without a scraper.
- Phase 4 is independent once Phase 3 delivers normalized data to diff against and embed.
- Phase 5 depends on Phase 4 for Aerospike docs and on Phase 3 for change events.
- Phase 6 is always last: building UI before the Go API exists forces building against mocks, which adds reconciliation time.
- The feature freeze should be enforced at the end of Phase 5 — no new features in Phase 6, only polish and demo prep.

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 4 (Aerospike integration):** MEDIUM confidence on exact Aerospike Go client API details for AVS 4.x vector operations (exact method signatures for batch upsert and ANN query need verification against live docs)
- **Phase 5 (TrueFoundry model routing):** MEDIUM confidence on exact gateway URL format and available model name strings (environment-specific configuration; verify from TrueFoundry dashboard before coding)

Phases with standard patterns (skip research-phase):
- **Phase 1 (Proto + server bootstrap):** Connect-Go and buf patterns are well-documented with official getting-started guides
- **Phase 2 (DB + Scraper):** pgx v5 + Ghost Postgres patterns are well-documented; net/http + goquery async patterns are established
- **Phase 3 (Extraction + Normalization):** Go struct-based pipeline with OpenAI-compatible REST client is a standard pattern
- **Phase 6 (Streamlit):** st.session_state, st.cache_data, and `requests` call patterns are well-documented

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Connect-Go, pgx v5, goquery, errgroup are well-documented Go libraries; TrueFoundry OpenAI-compatible pattern confirmed; Aerospike Go client 4.x confirmed |
| Features | MEDIUM | Competitor tools (Tierly, Prisync, Visualping, Priceva) verified from official sites; market gaps inferred from converging sources but not user-tested |
| Architecture | HIGH | Two-process split (Go + Streamlit), proto-first development, EtNS pipeline, two-tier LLM routing are well-documented patterns; Aerospike-specific HNSW API details are MEDIUM |
| Pitfalls | HIGH | HTML token overload quantified from published research (80K+ avg); JS rendering failure pattern is well-documented; hackathon patterns from practitioner accounts |

**Overall confidence:** HIGH

### Gaps to Address

- **TrueFoundry gateway URL format:** Environment-specific. Confirm the exact `TRUEFOUNDRY_GATEWAY_URL` value and available model name strings (cheap extraction model, expensive strategy model) from the TrueFoundry dashboard before the coding clock starts. This is the single highest-risk unknown.
- **Aerospike Go client exact method signatures:** The 4.x unified client API is confirmed but exact Go method calls for vector upsert and ANN query (e.g., parameter names, index configuration) should be verified against live docs or a quick test before Phase 4 begins.
- **Ghost DB autonomous creation behavior:** Research notes Ghost creates the DB autonomously, but the schema must still be explicitly created via pgx. Verify whether Ghost's autonomous creation pre-creates any tables or delivers a blank Postgres instance, to avoid a false assumption in Phase 2 setup.
- **Pre-validated target page list:** The 5-8 specific competitor pricing pages (e.g., OpenAI, Anthropic, Cohere, Macroscope) need to be pre-validated with net/http before the hackathon clock starts. This is Phase 0 and is blocking for everything else.

## Sources

### Primary (HIGH confidence)
- [Connect-Go docs](https://connectrpc.com/docs/go/getting-started/) — RPC framework, HTTP/JSON transport
- [buf CLI](https://buf.build/docs/) — proto codegen
- [pgx v5](https://github.com/jackc/pgx) — Go Postgres driver
- [goquery](https://github.com/PuerkitoBio/goquery) — HTML parsing in Go
- [Aerospike Go client](https://github.com/aerospike/aerospike-client-go) — vector store client
- [DRIPPER paper](https://openreview.net/pdf/e2b774a7481c9ccba439fa31dd837e9e32088b81.pdf) — 80K avg token count for raw HTML

### Secondary (MEDIUM confidence)
- [TrueFoundry AI Gateway](https://www.truefoundry.com/ai-gateway) — OpenAI-compatible endpoint confirmed; gateway URL format is environment-specific
- [Tierly](https://tierly.app/), [Prisync](https://prisync.com/), [Priceva](https://priceva.com/), [Visualping](https://visualping.io/) — competitor feature gaps verified from official sites
- [Aerospike Vector Search Developer Guide](https://aerospike.com/blog/aerospike-vector-search-guide/) — HNSW index patterns

### Tertiary (LOW confidence)
- [Things I Learned by Participating in GenAI Hackathons](https://towardsdatascience.com/things-i-learnt-by-participating-in-genai-hackathons-over-the-past-6-months/) — feature creep and demo polish patterns; practitioner account, unverified methodology

---
*Research completed: 2026-03-27 (updated: Go backend + Streamlit frontend architecture)*
*Ready for roadmap: yes*
