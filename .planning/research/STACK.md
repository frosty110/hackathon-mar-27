# Stack Research

**Domain:** Competitive pricing intelligence agent (scraping + LLM extraction + vector search + dashboard)
**Researched:** 2026-03-27
**Confidence:** HIGH (sponsor technologies verified; supporting library versions pulled from PyPI)

---

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Python | 3.11+ | Runtime | 3.11 gives significant perf gains over 3.10 for async; 3.12+ breaks some packages. Sweet spot for hackathon. |
| httpx | 0.28.1 | Async HTTP scraping | Native asyncio support, HTTP/2, connection pooling. The project explicitly bans Playwright — httpx is the right async alternative. aiohttp is faster in raw benchmarks but httpx has a cleaner API and HTTP/2 mitigates some bot-detection. |
| instructor | 1.14.5 | Structured LLM extraction from HTML | De facto standard (3M+ monthly downloads, 11k GitHub stars). Pydantic-native, supports retries/validation automatically, works with any provider via `from_provider()`. Avoids hand-rolling JSON parsing from LLM output. |
| pydantic | 2.12.5 | Data validation + schema definition | Required by instructor. V2 is 5-50x faster than V1 due to Rust core. Defines the normalized pricing schema that instructor populates. |
| openai | latest | LLM API client | TrueFoundry exposes an OpenAI-compatible endpoint — set `base_url` to the TrueFoundry gateway URL, `api_key` to TRUEFOUNDRY_API_KEY. No separate SDK needed. |
| aerospike-vector-search | 4.x | Vector similarity search for pricing profiles | Project requirement. Supports HNSW index, Python asyncio client, LangChain integration. Use for "who prices like us?" clustering and strategy doc retrieval. |
| asyncpg | 0.31.0 | Async Postgres driver (Ghost DB) | 5x faster than psycopg3 in benchmarks. Ghost provides managed Postgres — asyncpg is the right driver for async Python. No ORM overhead needed for this scope. |
| streamlit | 1.55.0 | Dashboard UI | Project requirement. Latest version (1.55.0, March 2026) has improved st.dataframe, st.metric, and layout controls. Single-file dashboards are fast to build — critical for 8-hour hackathon. |
| tenacity | 9.1.4 | Retry logic for scraping + LLM calls | Industry standard for Python retry logic. httpx requests and LLM calls both need retry-with-backoff. Decorator API integrates cleanly with instructor's retry system. |

### Supporting Libraries

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| beautifulsoup4 | 4.12.x | HTML pre-processing before LLM | Use to strip boilerplate (nav, footer, scripts) from raw HTML before sending to LLM — dramatically reduces token cost and noise. Pair with lxml parser for speed. |
| lxml | 5.x | Fast HTML parser backend for bs4 | Use as the bs4 parser (`BeautifulSoup(html, "lxml")`). C-backed, significantly faster than html.parser. |
| python-dotenv | 1.x | Environment variable management | Load API keys (TrueFoundry, OpenAI, Aerospike) from .env without hardcoding. |
| pytest-asyncio | 0.24.x | Async test support | All scraping and DB calls are async — need this to test coroutines. Only needed if writing tests during the hackathon. |
| rich | 13.x | Terminal output formatting | Nice-to-have for CLI progress output during demo. Zero-config pretty printing of extraction results in terminal. |

### Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| uv | Dependency management + venv | Replaces pip+venv. 10-100x faster installs. Use `uv init`, `uv add <package>`, `uv run streamlit run app.py`. Saves 15-30 minutes on dependency setup during an 8-hour hackathon. |
| Macroscope | Code review during build | 4th sponsor tool. Integrate into PR flow — create branches, open PRs, let Macroscope review. Counts as load-bearing sponsor use without adding runtime complexity. |

---

## Installation

```bash
# Bootstrap (use uv for speed)
pip install uv
uv init pricing-radar
cd pricing-radar

# Core
uv add httpx instructor pydantic openai tenacity streamlit asyncpg

# HTML parsing
uv add beautifulsoup4 lxml

# Aerospike Vector Search
uv add aerospike-vector-search

# Dev / utilities
uv add python-dotenv rich

# Run dashboard
uv run streamlit run app.py
```

---

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| httpx | aiohttp | If raw throughput across hundreds of concurrent requests is the bottleneck — aiohttp has lower per-request overhead. Irrelevant at 5-8 target pages. |
| httpx | Playwright | If target pages require heavy JS rendering (React SPA, lazy-loaded pricing tables). Explicitly out of scope for this project due to build risk. |
| instructor | LangChain extraction chains | If you need document chunking, multi-step pipelines, or agent orchestration beyond single LLM calls. Overkill here; instructor is lighter and more direct. |
| instructor | raw json_mode + manual parsing | If you want zero dependencies. Don't do this — instructor handles retries, validation failures, and partial output recovery that you'd have to reimplement. |
| asyncpg | psycopg3 | If you need a familiar psycopg2-style API or are using Django/SQLAlchemy. For raw async performance, asyncpg wins by 5x. |
| asyncpg | SQLAlchemy async | If schema migrations and ORM abstractions are needed. Ghost's autonomous DB creation makes raw asyncpg sufficient for this scope. |
| openai SDK (TrueFoundry base_url) | truefoundry[ml] SDK | If you need TrueFoundry model registry, experiment tracking, or deployment features. For LLM routing only, the OpenAI-compatible endpoint is simpler. |
| uv | pip + venv | If the team is unfamiliar with uv — pip works fine. But uv's speed matters during a timed hackathon. |
| text-embedding-3-small (via TrueFoundry) | sentence-transformers (local) | If you need fully offline embeddings. For a demo with network access, API embeddings are simpler — no GPU needed, no model download time. |

---

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| Playwright / Selenium | Explicit project constraint. Also: 30-60 minute setup risk, Chromium download, async complexity. Zero reward for this demo. | httpx + LLM to handle JS-rendered content by extracting what httpx can reach |
| LangChain | 10x more abstraction than needed. Debugging LangChain chains during a hackathon is a time sink. Simple instructor + openai SDK is faster to write and debug. | instructor + openai SDK directly |
| SQLAlchemy ORM | Adds migration complexity and model boilerplate. Ghost creates the DB; you just need to run SQL. | asyncpg raw queries or psycopg3 |
| pandas for storage | Pricing data is small — 5-8 companies, tens of rows. Pandas adds import time and complexity. | Python dicts + pydantic models until you need DataFrame display in Streamlit (use polars or native st.dataframe) |
| requests (sync) | Blocking I/O; scraping 5-8 pages sequentially takes 3-5x longer than async. Streamlit can tolerate it but it makes the demo feel slow. | httpx AsyncClient with asyncio.gather |
| json.loads(llm_response) | LLMs hallucinate JSON structure. This fails silently or throws unrecoverable errors. | instructor — validates, retries, and coerces automatically |

---

## Stack Patterns by Variant

**For the extraction pipeline (scrape + LLM parse):**
- Use `httpx.AsyncClient` with `asyncio.gather` to fetch all pricing pages concurrently
- Strip HTML with bs4/lxml before sending to LLM (target: < 2000 tokens of relevant content)
- Use instructor with Haiku/Flash (via TrueFoundry) for extraction — cheap, fast, good enough for structured pricing data

**For the strategy analysis (Pricing Architect):**
- Use instructor with Opus/GPT-4 (via TrueFoundry) — expensive model, called only once per scan
- Retrieve relevant internal positioning docs from Aerospike vector search before prompting
- Pass normalized pricing comparison + retrieved docs as context

**For the vector similarity ("who prices like us?"):**
- Generate embeddings of normalized pricing profiles using text-embedding-3-small via TrueFoundry
- Upsert into Aerospike AVS index with HNSW
- Query with current product's pricing profile to find nearest neighbors

**For the Streamlit dashboard:**
- Use `st.cache_data` to avoid re-running scans on every widget interaction
- Run the scan pipeline in a background thread (`threading.Thread`) or use `asyncio.run()` inside the Streamlit button callback — Streamlit runs synchronously, wrap async calls appropriately
- Use `st.rerun()` to refresh after scan completes

---

## Version Compatibility

| Package | Compatible With | Notes |
|---------|----------------|-------|
| instructor 1.14.5 | pydantic >=2.0 | instructor v1.x requires pydantic v2. Do not mix with pydantic v1 models. |
| httpx 0.28.1 | Python >=3.8 | Stable release. 1.0.dev3 exists but is a pre-release — avoid for hackathon. |
| asyncpg 0.31.0 | Python >=3.8, PostgreSQL >=9.6 | Ghost provides Postgres 15+ — no compatibility issues. |
| streamlit 1.55.0 | Python >=3.10 | Requires Python 3.10+. If using Python 3.11 (recommended), fully compatible. |
| aerospike-vector-search 4.x | Python >=3.9 | 4.0.0 merged admin and standard clients into single object — use the unified client API, not the pre-4.0 split pattern. |
| tenacity 9.1.4 | Python >=3.10 | Recent versions require Python 3.10+. Compatible with Python 3.11. |

---

## TrueFoundry Integration Pattern

TrueFoundry exposes an OpenAI-compatible endpoint. No separate SDK needed for LLM routing:

```python
import os
from openai import AsyncOpenAI

# Cheap model for extraction (Haiku / Flash)
extraction_client = AsyncOpenAI(
    api_key=os.environ["TRUEFOUNDRY_API_KEY"],
    base_url=os.environ["TRUEFOUNDRY_GATEWAY_URL"],
)

# Expensive model for strategy (Opus / GPT-4)
strategy_client = AsyncOpenAI(
    api_key=os.environ["TRUEFOUNDRY_API_KEY"],
    base_url=os.environ["TRUEFOUNDRY_GATEWAY_URL"],
)

# instructor wraps these directly
import instructor
extraction_instructor = instructor.from_openai(extraction_client)
strategy_instructor = instructor.from_openai(strategy_client)
```

Route by passing different `model=` strings to TrueFoundry — the gateway handles provider dispatch.

---

## Sources

- PyPI / streamlit — version 1.55.0 verified March 2026
- PyPI / instructor — version 1.14.5 verified March 2026
- PyPI / httpx — version 0.28.1 verified (stable); 1.0.dev3 is pre-release
- PyPI / pydantic — version 2.12.5 verified
- PyPI / asyncpg — version 0.31.0 verified
- PyPI / tenacity — version 9.1.4 verified
- [Aerospike AVS Python client docs](https://aerospike.com/docs/vector/develop/python) — 4.0.0 unified client release confirmed; exact latest patch unverified (PyPI blocked JS)
- [Instructor docs](https://python.useinstructor.com/) — HIGH confidence, 3M+ monthly downloads, actively maintained
- [TrueFoundry AI Gateway](https://www.truefoundry.com/ai-gateway) — OpenAI-compatible endpoint confirmed; MEDIUM confidence on exact gateway URL format (environment-specific)
- [Aerospike AVS 4.0.0 release notes](https://aerospike.com/docs/vector/release-notes/python/python-4.0.0-release-notes) — unified client API confirmed
- [httpx async scraping](https://brightdata.com/blog/web-data/web-scraping-with-httpx) — HTTP/2 + async patterns MEDIUM confidence (external blog, verified against httpx docs)

---

*Stack research for: Pricing Radar — competitive pricing intelligence agent*
*Researched: 2026-03-27*
