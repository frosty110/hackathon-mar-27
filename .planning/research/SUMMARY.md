# Project Research Summary

**Project:** Pricing Radar — Competitive Pricing Intelligence Agent
**Domain:** AI/SaaS competitive pricing intelligence (scraping + LLM extraction + vector search + dashboard)
**Researched:** 2026-03-27
**Confidence:** HIGH

## Executive Summary

Pricing Radar is a competitive pricing intelligence agent for AI startup product teams. It solves a gap that every surveyed competitor (Tierly, Prisync, Visualping, Priceva) fails to address: normalizing heterogeneous AI pricing models (per-seat, per-token, credits, flat tiers) to a common unit with transparent assumptions, then grounding strategic recommendations in the team's own positioning documents. The recommended approach is a Python async pipeline — httpx for concurrent page fetches, instructor + TrueFoundry for two-tier LLM routing (cheap extraction, expensive strategy), Ghost Postgres for structured scan history, Aerospike for vector similarity clustering and RAG, and Streamlit as the demo surface. This stack maps directly to sponsor requirements and has well-documented patterns for every integration.

The key differentiators are (1) cross-model normalization with explicit assumption transparency, (2) strategy-grounded recommendations via RAG over uploaded positioning docs, and (3) "who prices like us?" clustering via Aerospike HNSW. None of these exist in any surveyed tool. The hardest design problem is normalization: it is an epistemic task (surfacing assumptions about usage patterns), not just a unit conversion. If assumptions are hidden, the demo comparison table is factually untrustworthy and will fail under judge scrutiny.

The primary risks are: JS-rendered pricing pages that return empty HTML to httpx (mitigate by pre-validating all 8 target URLs before coding starts, with mock fallback files ready); demo latency from synchronous Streamlit + async pipeline mismatch (mitigate by pre-running the scan and displaying cached DB results during the demo); and scope creep enabled by AI-accelerated development (mitigate with a hard feature freeze at hour 5). The architectural build order in ARCHITECTURE.md is decisive — follow it to avoid discovering at hour 4 that the DB schema doesn't support change detection.

## Key Findings

### Recommended Stack

The stack is built around Python 3.11 with uv for dependency management. Every core library maps to a sponsor technology or a well-justified best-in-class choice. TrueFoundry is accessed via the OpenAI-compatible endpoint (set `base_url` + `api_key`; no extra SDK). Aerospike Vector Search 4.x uses a unified client API (the pre-4.0 split client pattern is obsolete). instructor 1.14.5 requires pydantic v2 — do not mix with v1 models. Streamlit requires Python 3.10+, which is satisfied by Python 3.11.

**Core technologies:**
- Python 3.11: Runtime — sweet spot for async perf without package-breaking changes in 3.12+
- httpx 0.28.1: Async scraping — native asyncio, HTTP/2, clean API; Playwright explicitly banned
- instructor 1.14.5: Structured LLM extraction — pydantic-native, auto-retry/validation, 3M+ monthly downloads
- pydantic 2.12.5: Schema validation — required by instructor; 5-50x faster than v1
- openai SDK (TrueFoundry base_url): LLM routing — single SDK covers both cheap and expensive model calls
- aerospike-vector-search 4.x: Vector similarity + RAG — sponsor requirement; HNSW index, unified Python client
- asyncpg 0.31.0: Ghost Postgres driver — 5x faster than psycopg3 for async; no ORM overhead needed
- streamlit 1.55.0: Dashboard — sponsor requirement; latest version with improved dataframe/metric controls
- tenacity 9.1.4: Retry logic — wraps httpx requests and LLM calls with backoff
- beautifulsoup4 + lxml: HTML preprocessing — strip boilerplate before LLM to target <8K tokens per page
- uv: Dependency management — 10-100x faster installs; saves 15-30 min in an 8-hour build

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
- Two-tier LLM model routing — Haiku/Flash for extraction, Opus/GPT-4 for strategy; satisfies TrueFoundry sponsor story
- Change event narrative — plain-English summary of what changed and why it matters, generated post-diff

**Defer (v2+):**
- Playwright/browser rendering for JS-heavy pages — 3-4 hour build risk, breaks demo reliability
- User authentication and multi-tenant — doubles scope; single-user demo is sufficient
- Scheduled recurring scans (cron-based) — blocked by single-user session model in hackathon build
- Competitor auto-discovery — Tierly takes 13 minutes; unreliable for a deterministic demo

### Architecture Approach

The system is a four-layer pipeline: Presentation (Streamlit dashboard) → Agent Orchestration (Scan Runner, Change Detector, Pricing Architect) → Processing (Scraper, LLM Extractor, Normalizer, Embedder) → Data (Ghost Postgres for structured history, Aerospike for vectors). The key pattern is Extract-Normalize-Store (EtNS): scraper, extractor, and normalizer are distinct modules with typed pydantic inputs/outputs. This makes each stage independently testable and lets you retry extraction without re-fetching. A second pattern is Two-Tier LLM Routing: cheap models (Haiku/Flash) for templated extraction tasks, expensive model (Opus/GPT-4) only for synthesis and strategy generation. The dashboard is read-only — it queries both stores but never calls the agent directly.

**Major components:**
1. Scan Runner (agent/runner.py) — orchestrates the full scrape → extract → normalize → store → detect cycle
2. Scraper (scraper/fetcher.py) — httpx.AsyncClient + asyncio.gather for concurrent page fetches
3. LLM Extractor (extractor/llm_extract.py) — TrueFoundry cheap model call with instructor structured output
4. Normalizer (normalizer/normalize.py) — maps heterogeneous pricing models to common unit; records explicit assumptions
5. Change Detector (agent/change_detector.py) — diffs current normalized snapshot against previous scan in Ghost DB
6. Embedder (embedder/embed.py) — generates pricing profile vectors; upserts to Aerospike HNSW index
7. Pricing Architect (agent/architect.py) — RAG retrieval from Aerospike + TrueFoundry expensive model for strategy
8. Streamlit Dashboard (dashboard/app.py) — comparison table, change alerts, cluster visual, transparency map
9. Ghost Postgres — scan_runs, pricing_snapshots, change_events; versioned by scan_run_id and competitor key
10. Aerospike AVS — pricing profile embeddings for similarity clustering; strategy docs for RAG retrieval

### Critical Pitfalls

1. **Raw HTML overload kills extraction quality** — 80K+ tokens of nav/footer/scripts drowns the pricing signal. Pre-process with BeautifulSoup to strip boilerplate; target 2K-8K tokens per page. Never send full raw HTML to the LLM under any circumstances.

2. **Hidden normalization assumptions expose methodology flaws in demo Q&A** — Normalization is an epistemic task, not just unit conversion. Every normalized value must carry its `normalization_assumptions` dict. Display assumption footnotes inline in the dashboard. A bare normalized number without a footnote will fail under judge scrutiny.

3. **JS-rendered pages return empty HTML to httpx silently** — Many AI company pricing pages load tiers via JavaScript. httpx returns 200 OK with an empty pricing div. Pre-validate all 8 target URLs before coding begins: grep raw HTML for known price strings. Prepare mock fallback HTML files for 2-3 pages before the clock starts.

4. **Ghost DB schema designed for storage, not for change detection** — If scans are stored without a `scan_run_id`, `competitor` key, and `scraped_at` timestamp, change detection queries are impossible or require a schema rewrite in hour 4. Design the schema for the change detection query from day one.

5. **Demo scan latency + Streamlit UI freeze** — Streamlit reruns the full script on every widget interaction. A live scan button without `st.session_state` guards re-fires LLM calls on every rerender. Pre-run the scan; display cached DB results during the demo. Gate the scan button with session state. Add `st.progress()` and `st.status()` — never show a frozen screen.

## Implications for Roadmap

Based on combined research, the architectural build order from ARCHITECTURE.md is the right phase structure. Each phase has a clear deliverable and the dependency chain is hard: each later phase fails without the earlier one working.

### Phase 0: Pre-work and Target Validation
**Rationale:** JS-rendering failures (Pitfall 3) are discovered during coding, not before. Pre-validation eliminates the highest-risk demo failure mode before any code is written.
**Delivers:** Validated list of 5-8 scrapable pricing pages; mock HTML fallbacks for any JS-heavy pages; project scaffolding (uv, .env, config.py, pydantic schemas)
**Addresses:** Pitfall 3 (JS pages), Pitfall 5 (scope creep) — feature list is locked before building starts
**Avoids:** Discovering mid-build that 3 of 8 target pages are JS-rendered

### Phase 1: Data Foundation (DB + Scraper)
**Rationale:** Ghost DB schema must exist and be correct before any data flows through it. Change detection requires the right schema from scan 1. Scraper must be verified working before LLM extraction is added on top.
**Delivers:** Ghost Postgres schema (scan_runs, pricing_snapshots, change_events); httpx async fetcher returning raw HTML; verified URL fetches for all targets
**Uses:** asyncpg, Ghost Postgres, httpx, tenacity
**Avoids:** Pitfall 6 (schema misfit for change detection)

### Phase 2: Extraction and Normalization Pipeline
**Rationale:** Extraction and normalization are tightly coupled and share pydantic schemas, but must stay distinct modules. This is the integration milestone — once runner.py works end-to-end with a single competitor, the rest of the pipeline is incremental.
**Delivers:** LLM extraction (TrueFoundry cheap model + instructor) producing structured RawPricingData; normalization to common unit with explicit assumption recording; runner.py wiring the full scrape → extract → normalize → store cycle
**Uses:** instructor, TrueFoundry cheap model, pydantic, beautifulsoup4/lxml, tenacity
**Implements:** LLM Extractor, Normalizer, Scan Runner (partial)
**Avoids:** Pitfall 1 (raw HTML overload), Pitfall 2 (hidden normalization assumptions)

### Phase 3: Change Detection and Aerospike Integration
**Rationale:** Change detection is the feature that makes Pricing Radar a monitor rather than a one-shot tool. Aerospike integration (embeddings + similarity) can proceed in parallel once normalized data exists in Postgres.
**Delivers:** Change Detector diffing current vs. previous scan; change event storage in Ghost DB; pricing profile embeddings upserted to Aerospike HNSW; "who prices like us?" similarity query working
**Uses:** Aerospike Vector Search 4.x, TrueFoundry embedding model (text-embedding-3-small)
**Implements:** Change Detector, Embedder, Aerospike client

### Phase 4: Pricing Architect (RAG + Strategy)
**Rationale:** This phase requires working Aerospike (for doc retrieval), working change detection (to trigger recommendations), and the expensive TrueFoundry model route. It is the most visible differentiator for judges.
**Delivers:** RAG retrieval of strategy docs from Aerospike; strategy prompt construction with competitor context + change summary; TrueFoundry expensive model call; StrategicResponse stored in Ghost DB
**Uses:** TrueFoundry expensive model (Opus/GPT-4), Aerospike AVS (RAG mode), instructor
**Implements:** Pricing Architect, two-tier LLM routing

### Phase 5: Streamlit Dashboard and Demo Hardening
**Rationale:** Dashboard is built last because it only reads from stores that must already exist. Demo hardening (pre-seeded data, session state guards, progress indicators) is non-negotiable — the demo is the product for the judges.
**Delivers:** Streamlit dashboard with comparison table, change alerts, normalization transparency map, cluster visual; pre-seeded scan data in DB; demo script timed under 3 minutes
**Uses:** Streamlit 1.55.0, st.cache_data, st.session_state, st.progress
**Avoids:** Pitfall 4 (demo scan latency), UX pitfalls from PITFALLS.md (raw diffs, unlabeled clusters)

### Phase Ordering Rationale

- Phase 0 must precede all others: JS page failures discovered mid-build cost 1-2 hours; discovered in pre-work they cost 30 minutes.
- Phases 1-2 are strict dependencies: you cannot detect changes without stored data; you cannot store data without a schema; you cannot extract without a scraper.
- Phase 3 is independent once Phase 2 delivers normalized data to diff against and embed.
- Phase 4 depends on Phase 3 for Aerospike docs and on Phase 2 for change events.
- Phase 5 is always last: building UI before data exists forces building against mocks, which adds reconciliation time.
- The feature freeze (Pitfall 5) should be enforced at the end of Phase 4 — no new features in Phase 5, only polish and demo prep.

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 3 (Aerospike integration):** MEDIUM confidence on exact AVS 4.x API details (PyPI was blocked JS during research; unified client API confirmed from release notes but exact method signatures need verification against live docs)
- **Phase 4 (TrueFoundry model routing):** MEDIUM confidence on exact gateway URL format and available model name strings (environment-specific configuration; verify from TrueFoundry dashboard before coding)

Phases with standard patterns (skip research-phase):
- **Phase 1 (DB + Scraper):** asyncpg + Ghost Postgres patterns are well-documented; httpx async patterns are established
- **Phase 2 (Extraction + Normalization):** instructor + pydantic patterns are extensively documented with official examples
- **Phase 5 (Streamlit):** st.session_state, st.cache_data, and threading patterns for Streamlit + async are well-documented

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | All package versions verified from PyPI March 2026; TrueFoundry OpenAI-compatible pattern confirmed; Aerospike 4.0.0 unified client confirmed from release notes |
| Features | MEDIUM | Competitor tools (Tierly, Prisync, Visualping, Priceva) verified from official sites; market gaps inferred from converging sources but not user-tested |
| Architecture | HIGH | Core pipeline patterns (EtNS, two-tier LLM routing, dual-store) are well-documented; Aerospike-specific HNSW API details are MEDIUM |
| Pitfalls | HIGH | HTML token overload quantified from published research (80K+ avg); JS rendering failure pattern is well-documented; hackathon patterns from practitioner accounts |

**Overall confidence:** HIGH

### Gaps to Address

- **TrueFoundry gateway URL format:** Environment-specific. Confirm the exact `TRUEFOUNDRY_GATEWAY_URL` value and available model name strings (Haiku, Flash, Opus equivalent) from the TrueFoundry dashboard before the coding clock starts. This is the single highest-risk unknown.
- **Aerospike AVS exact method signatures:** The 4.0.0 unified client API is confirmed but exact Python method calls (e.g., `client.vector_search()` parameter names, `include_fields` behavior) should be verified against live docs or a quick test upsert/query before Phase 3 begins.
- **Ghost DB autonomous creation behavior:** Research notes Ghost creates the DB autonomously, but the schema must still be explicitly created. Verify whether Ghost's autonomous creation pre-creates any tables or delivers a blank Postgres instance, to avoid a false assumption in Phase 1 setup.
- **Pre-validated target page list:** The 5-8 specific competitor pricing pages (e.g., OpenAI, Anthropic, Cohere, Macroscope) need to be pre-validated with httpx before the hackathon clock starts. This is Phase 0 and is blocking for everything else.

## Sources

### Primary (HIGH confidence)
- PyPI / streamlit 1.55.0, instructor 1.14.5, httpx 0.28.1, pydantic 2.12.5, asyncpg 0.31.0, tenacity 9.1.4 — all versions verified March 2026
- [Aerospike AVS 4.0.0 release notes](https://aerospike.com/docs/vector/release-notes/python/python-4.0.0-release-notes) — unified client API
- [Instructor docs](https://python.useinstructor.com/) — structured extraction patterns
- [DRIPPER paper](https://openreview.net/pdf/e2b774a7481c9ccba439fa31dd837e9e32088b81.pdf) — 80K avg token count for raw HTML

### Secondary (MEDIUM confidence)
- [TrueFoundry AI Gateway](https://www.truefoundry.com/ai-gateway) — OpenAI-compatible endpoint confirmed; gateway URL format is environment-specific
- [Tierly](https://tierly.app/), [Prisync](https://prisync.com/), [Priceva](https://priceva.com/), [Visualping](https://visualping.io/) — competitor feature gaps verified from official sites
- [Aerospike Vector Search Developer Guide](https://aerospike.com/blog/aerospike-vector-search-guide/) — HNSW index patterns
- [Web Scraping with Python HTTPX — Scrapfly](https://scrapfly.io/blog/posts/web-scraping-with-python-httpx) — async patterns

### Tertiary (LOW confidence)
- [Things I Learned by Participating in GenAI Hackathons](https://towardsdatascience.com/things-i-learnt-by-participating-in-genai-hackathons-over-the-past-6-months/) — feature creep and demo polish patterns; practitioner account, unverified methodology

---
*Research completed: 2026-03-27*
*Ready for roadmap: yes*
