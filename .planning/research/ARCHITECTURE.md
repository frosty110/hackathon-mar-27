# Architecture Research

**Domain:** Competitive pricing intelligence agent (scrape → extract → normalize → detect → recommend)
**Researched:** 2026-03-27
**Confidence:** HIGH (core pipeline patterns), MEDIUM (Aerospike AVS specifics), HIGH (TrueFoundry routing)

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Presentation Layer                           │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  Streamlit Dashboard                                          │   │
│  │  (comparison table │ change alerts │ cluster visual │ map)   │   │
│  └──────────────────────────────┬───────────────────────────────┘   │
└─────────────────────────────────┼───────────────────────────────────┘
                                  │ reads
┌─────────────────────────────────┼───────────────────────────────────┐
│                        Agent Orchestration Layer                    │
│                                 │                                   │
│  ┌──────────────┐  ┌────────────┴───────┐  ┌──────────────────┐    │
│  │  Scan Runner │  │  Change Detector   │  │  Pricing Architect│    │
│  │  (scheduler) │  │  (diff engine)     │  │  (strategy LLM)   │    │
│  └──────┬───────┘  └────────────────────┘  └─────────┬────────┘    │
│         │                                             │             │
└─────────┼─────────────────────────────────────────────┼────────────┘
          │ triggers                                     │ retrieves
┌─────────┼─────────────────────────────────────────────┼────────────┐
│                        Processing Layer                             │
│         │                                             │             │
│  ┌──────▼───────┐  ┌────────────────────┐  ┌─────────▼────────┐    │
│  │  Scraper     │  │  Normalizer        │  │  Embedder        │    │
│  │  (httpx      │  │  (heterogeneous    │  │  (pricing profile│    │
│  │   async)     │  │   model mapping)   │  │   → vector)      │    │
│  └──────┬───────┘  └──────────┬─────────┘  └─────────┬────────┘    │
│         │ raw HTML            │ structured            │ embedding   │
│  ┌──────▼───────┐             │                       │             │
│  │  LLM         │ ────────────┘                       │             │
│  │  Extractor   │                                     │             │
│  │  (TrueFoundry│                                     │             │
│  │   cheap model│                                     │             │
│  └──────────────┘                                     │             │
└───────────────────────────────────────────────────────┼────────────┘
                                                        │
┌───────────────────────────────────────────────────────┼────────────┐
│                        Data Layer                                   │
│                                                        │            │
│  ┌──────────────────────┐         ┌────────────────────▼─────────┐  │
│  │  Ghost (Postgres)    │         │  Aerospike Vector DB         │  │
│  │  - scan_runs         │         │  - pricing profile embeddings│  │
│  │  - pricing_snapshots │         │  - strategy docs (RAG)       │  │
│  │  - change_events     │         │  - similarity clusters       │  │
│  └──────────────────────┘         └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Implementation |
|-----------|----------------|----------------|
| Scan Runner | Trigger scrape cycle, coordinate async fetches, persist scan run record | Python asyncio entry point |
| Scraper | Async-fetch 5-8 pricing pages concurrently, return raw HTML/text | httpx.AsyncClient + asyncio.gather |
| LLM Extractor | Parse unstructured pricing HTML into structured JSON (plan names, prices, units, features) | TrueFoundry cheap model (Haiku/Flash) with structured output schema |
| Normalizer | Map heterogeneous pricing models (per-seat, per-token, credits, tiers) to a common comparison unit; record explicit assumptions | Pure Python with pydantic models |
| Change Detector | Diff current normalized snapshot against previous scan stored in Postgres; emit change events | Hash comparison + field-level diff against Ghost DB |
| Embedder | Embed normalized pricing profiles as vectors for similarity search | Embedding model call, store in Aerospike |
| Pricing Architect | RAG: retrieve strategy docs from Aerospike, synthesize competitive response | TrueFoundry expensive model (Opus/GPT-4) |
| Streamlit Dashboard | Display comparison table, change alerts, cluster visual, transparency map | Streamlit app querying Postgres + Aerospike |
| Ghost (Postgres) | Persistent store for all scan data, change history, strategy docs metadata | Autonomous DB creation via Ghost API |
| Aerospike | Vector store for pricing embeddings (similarity clustering) and strategy doc retrieval (RAG) | AVS HNSW index |
| TrueFoundry Gateway | Route extraction calls to cheap model, strategy calls to expensive model | LLM router with cost-aware dispatch |

## Recommended Project Structure

```
pricing_radar/
├── agent/
│   ├── runner.py           # orchestrates full scan cycle
│   ├── change_detector.py  # diffs current vs previous scan
│   └── architect.py        # Pricing Architect — RAG + strategy generation
├── scraper/
│   ├── fetcher.py          # httpx async fetcher, concurrent page loads
│   └── targets.py          # pre-selected page configs (URL, hints, fallback)
├── extractor/
│   ├── llm_extract.py      # TrueFoundry cheap model call, structured output
│   └── schema.py           # pydantic models for raw extracted data
├── normalizer/
│   ├── normalize.py        # maps pricing models to common unit
│   ├── models.py           # per-seat, per-token, credits, flat definitions
│   └── assumptions.py      # records normalization assumptions per vendor
├── storage/
│   ├── ghost_db.py         # Ghost Postgres client — scan runs, snapshots
│   └── aerospike_db.py     # Aerospike client — embeddings, RAG docs
├── embedder/
│   └── embed.py            # generate + upsert pricing profile vectors
├── dashboard/
│   └── app.py              # Streamlit dashboard entry point
├── config.py               # env vars, model names, target URLs
└── main.py                 # CLI entry: run scan, run dashboard
```

### Structure Rationale

- **agent/**: Orchestration lives separately from data processing — the runner wires components, the detector and architect are distinct concerns
- **scraper/**: Isolated fetch layer makes it easy to swap httpx for Playwright if needed post-hackathon, and keeps target configs out of business logic
- **extractor/**: Extraction from LLM is distinct from normalization — the LLM extracts what the page says, the normalizer interprets it
- **normalizer/**: Normalization is the hardest domain logic; isolating it enables unit testing and makes the assumption-recording transparent
- **storage/**: One module per backing store; dashboard reads from both without needing to know how data got there

## Architectural Patterns

### Pattern 1: Extract-Normalize-Store Pipeline (EtNS)

**What:** Separate the three concerns of getting data (scrape), understanding data (LLM extract), and making data comparable (normalize) into distinct, synchronous stages that pass typed models forward.
**When to use:** Any time source data is heterogeneous and you need auditable transformation steps — pricing pages are the canonical case.
**Trade-offs:** More modules than a monolithic scrape-and-store, but each stage is independently testable. In an 8-hour build, the overhead is worth it because debugging is vastly easier.

**Example:**
```python
# runner.py
async def run_scan(targets: list[Target]) -> ScanRun:
    raw_pages = await fetch_all(targets)           # scraper layer
    extracted = [llm_extract(p) for p in raw_pages] # extractor layer (cheap model)
    normalized = [normalize(e) for e in extracted]  # normalizer layer
    scan_run = await store_scan(normalized)          # storage layer
    changes = detect_changes(scan_run)               # change detector
    return scan_run, changes
```

### Pattern 2: Two-Tier LLM Routing (Extraction vs. Strategy)

**What:** Route structurally simple extraction tasks (parse a pricing page into JSON) to a fast, cheap model. Route reasoning-heavy tasks (generate competitive response grounded in strategy docs) to a frontier model. Gate routing at the TrueFoundry gateway level.
**When to use:** Any pipeline where some LLM calls are templated/structured and others require synthesis and judgment.
**Trade-offs:** Adds TrueFoundry as a dependency but reduces cost by ~80% and improves throughput for the extraction-heavy path.

**Example:**
```python
# extractor/llm_extract.py  — cheap model
def llm_extract(html: str) -> RawPricingData:
    return truefoundry_client.complete(
        model="haiku-3-5",  # routed by TrueFoundry to cheap tier
        messages=[{"role": "user", "content": EXTRACTION_PROMPT + html}],
        response_format=RawPricingData,
    )

# agent/architect.py  — expensive model
def generate_strategy(context: StrategyContext) -> StrategicResponse:
    return truefoundry_client.complete(
        model="claude-opus-4",  # routed by TrueFoundry to expensive tier
        messages=[{"role": "user", "content": build_strategy_prompt(context)}],
    )
```

### Pattern 3: Dual-Store Architecture (Relational + Vector)

**What:** Use Postgres (Ghost) for structured relational data (scan history, change events, pricing tables) and Aerospike for vector data (pricing profile embeddings, strategy doc retrieval). Each store serves what it's best at — no forcing vectors into Postgres or scan history into a vector DB.
**When to use:** When you need both time-series change detection (relational, indexed by timestamp/vendor) and similarity clustering or semantic retrieval (vector).
**Trade-offs:** Two clients to manage, but the boundary is clean: Postgres = "what happened when", Aerospike = "what is similar / what is relevant".

## Data Flow

### Scan Trigger Flow

```
User clicks "Run Scan" (Streamlit)
    ↓
runner.py: creates scan_run record in Ghost Postgres
    ↓
fetcher.py: asyncio.gather → httpx fetches 5-8 URLs concurrently
    ↓
llm_extract.py: parallel LLM calls via TrueFoundry (cheap model)
    → returns structured RawPricingData per vendor
    ↓
normalize.py: RawPricingData → NormalizedPricing (common unit + assumptions)
    ↓
ghost_db.py: upsert pricing_snapshots for this scan_run
    ↓
embed.py: embed NormalizedPricing → vector, upsert to Aerospike
    ↓
change_detector.py: diff current snapshots against previous scan_run
    → emits ChangeEvent records if diffs found
    ↓
ghost_db.py: write ChangeEvent records
    ↓
Dashboard: Streamlit re-reads Postgres + Aerospike, renders updated state
```

### Strategic Recommendation Flow (Pricing Architect)

```
ChangeEvent detected (e.g. Macroscope switches to usage-based)
    ↓
architect.py: query Aerospike for strategy docs (RAG retrieval)
    → semantic search: "usage-based pricing response"
    ↓
architect.py: build prompt = strategy docs + competitor context + change summary
    ↓
TrueFoundry gateway: route to expensive model (Opus/GPT-4)
    ↓
StrategicResponse returned, stored in Ghost Postgres
    ↓
Dashboard: display recommendation in change alert panel
```

### Similarity Clustering Flow

```
After each scan_run:
    NormalizedPricing for all vendors → embedder → Aerospike upsert

Dashboard "cluster" view:
    Aerospike ANN query: k-nearest neighbors for each vendor embedding
    → cluster groups rendered in Streamlit as visual
```

### Key Data Flows

1. **Extraction path:** HTML → TrueFoundry cheap model → structured JSON (synchronous within async gather)
2. **Change detection path:** current snapshot hash vs. Ghost DB stored hash → diff if mismatch
3. **RAG path:** change event → Aerospike semantic search → strategy context → TrueFoundry expensive model → recommendation
4. **Similarity path:** normalized pricing vector → Aerospike HNSW → nearest neighbors → "who prices like us" cluster

## Scaling Considerations

This is a hackathon tool; scale is not a concern. The notes below are informational only.

| Scale | Architecture Adjustments |
|-------|--------------------------|
| Demo (5-8 pages) | Single process, everything in-process, polling on button click |
| 50-100 pages | Add task queue (Celery/ARQ), separate scrape workers from API |
| 1000+ pages | Distributed scraping workers, rate-limit management, incremental Aerospike indexing |

### Scaling Priorities (post-hackathon)

1. **First bottleneck:** LLM extraction latency — mitigate with concurrent calls and caching extracted HTML that hasn't changed
2. **Second bottleneck:** Aerospike index rebuild on large vector upserts — use batch upserts and async indexing already built into AVS

## Anti-Patterns

### Anti-Pattern 1: Monolithic Scrape-Extract-Store Function

**What people do:** Write one `scrape_and_store(url)` function that fetches, parses, normalizes, and writes in a single blob.
**Why it's wrong:** When the LLM extraction fails (and it will), you can't retry just the extraction. When normalization logic changes, you reparse everything. Debugging is a nightmare in a 3-minute demo.
**Do this instead:** Separate scraper, extractor, and normalizer as distinct modules with typed inputs/outputs. Cache raw HTML so you can re-run extraction without re-fetching.

### Anti-Pattern 2: Single LLM for All Tasks

**What people do:** Send every call — both "parse this pricing table" and "generate a strategic response" — to the same expensive model.
**Why it's wrong:** Extraction is a structured, deterministic task that Haiku/Flash handles at 95% quality for 5% of the cost. Wasting GPT-4 on HTML parsing also hides the TrueFoundry routing value.
**Do this instead:** Use TrueFoundry's routing to dispatch extraction to cheap models and reserve the expensive model for synthesis tasks. This is also the sponsor story.

### Anti-Pattern 3: Storing Raw HTML in Postgres

**What people do:** Persist raw HTML in the database as the primary data artifact.
**Why it's wrong:** Blobs in Postgres slow queries, obscure the actual data model, and create re-processing complexity. The structured normalized pricing is what you query, diff, and display.
**Do this instead:** Store normalized pricing snapshots in Postgres. Cache raw HTML on disk or in memory only if you need LLM extraction retry without re-fetch.

### Anti-Pattern 4: Blocking I/O for Scraping

**What people do:** Sequential `requests.get()` calls per URL in a loop.
**Why it's wrong:** Fetching 8 pages sequentially at 2-3 seconds each = 16-24 seconds before any extraction begins. Demo-killing latency.
**Do this instead:** `httpx.AsyncClient` with `asyncio.gather` to fetch all pages concurrently. With 8 targets, this collapses fetch time to the slowest single page (~2-3s vs. 24s).

### Anti-Pattern 5: Hardcoding the Normalization Unit

**What people do:** Pick one unit ("per month per user") and silently coerce everything to it.
**Why it's wrong:** Pricing normalization involves real assumptions (e.g., "Anthropic token pricing assumed 1M tokens/month per user"). Hidden assumptions erode trust and fail in demo Q&A.
**Do this instead:** Record assumptions explicitly in a `normalization_assumptions` field per snapshot. Display them in the dashboard's "transparency map." This is a differentiator, not overhead.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| Ghost (Postgres) | Ghost client → autonomous DB creation on startup; psycopg2 / asyncpg for queries | Create tables on first run if not exist; use scan_run as the top-level entity |
| TrueFoundry | OpenAI-compatible REST API with model name routing; single client instance | Two model names: one for cheap (extraction), one for expensive (strategy) |
| Aerospike | Aerospike Python client + AVS for vector operations; separate namespace for vectors vs. docs | HNSW index; batch upsert after each scan for cluster freshness |
| Target pricing pages | httpx.AsyncClient with conservative timeouts (10s); no auth, no JS | Pre-validate all URLs before demo; keep localhost mock ready as fallback |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Scraper → Extractor | Typed dict / dataclass: `{url, html, fetched_at}` | Scraper does not parse; extractor does not fetch |
| Extractor → Normalizer | Pydantic `RawPricingData` model | LLM output validated against schema before passing forward |
| Normalizer → Storage | Pydantic `NormalizedPricing` model | Storage writes both Postgres snapshot and triggers embedder |
| Agent → Dashboard | Postgres + Aerospike reads (no shared in-process state) | Dashboard is read-only; it never calls the agent directly |
| Change Detector → Architect | `ChangeEvent` dataclass passed explicitly | Architect is invoked only when change events are present; not on every scan |

## Suggested Build Order

Build in dependency order — each layer only requires the layer below it to be complete.

```
1. config.py + schema.py (pydantic models)
        ↓
2. ghost_db.py (DB init, table creation, basic writes)
        ↓
3. fetcher.py + targets.py (async fetch, verify all URLs work)
        ↓
4. llm_extract.py (TrueFoundry extraction call, structured output)
        ↓
5. normalize.py (normalization logic + assumption recording)
        ↓
6. runner.py (wire scrape → extract → normalize → store in one cycle)
        ↓
7. change_detector.py (diff against previous scan_run)
        ↓
8. aerospike_db.py + embed.py (vector upsert, similarity query)
        ↓
9. architect.py (RAG retrieval + TrueFoundry expensive model call)
        ↓
10. app.py (Streamlit dashboard — reads from both stores)
        ↓
11. main.py + demo script (seed DB, trigger scan, run demo flow)
```

**Why this order:**
- Steps 1-2 establish data models and persistence before anything else touches them
- Step 3 validates scraping works before adding LLM extraction on top
- Steps 4-5 can be tested with a single target URL before wiring all 8
- Step 6 (runner) is the integration milestone — once this works end-to-end, the rest is incremental
- Steps 7-9 are independent enhancements that can be demoed as progressive capability
- Dashboard (step 10) is last because it only displays; it can be built against seed data

## Sources

- [How to Build Effective Competitive Pricing Intelligence Systems](https://www.getmonetizely.com/articles/how-to-build-effective-competitive-pricing-intelligence-systems-a-complete-guide)
- [Understanding Price Intelligence in 2025 — Impact Analytics](https://www.impactanalytics.co/blog/price-intelligence)
- [Aerospike Vector Search Developer Guide](https://aerospike.com/blog/aerospike-vector-search-guide/)
- [Aerospike RAG vector database use case](https://aerospike.com/solutions/use-cases/rag-vector-database/)
- [TrueFoundry: What is LLM Router?](https://www.truefoundry.com/blog/what-is-llm-router)
- [TrueFoundry: LLM Load Balancing](https://www.truefoundry.com/blog/llm-load-balancing)
- [Web Scraping with Python HTTPX — Scrapfly](https://scrapfly.io/blog/posts/web-scraping-with-python-httpx)
- [Async Web Scraping Concurrency Patterns](https://use-apify.com/blog/async-web-scraping-concurrency-patterns)
- [Agentic RAG — What is Agentic RAG? Weaviate](https://weaviate.io/blog/what-is-agentic-rag)
- [RAG Architecture + LLM Agent](https://www.k2view.com/blog/rag-architecture-llm-agent/)

---
*Architecture research for: competitive pricing intelligence agent*
*Researched: 2026-03-27*
