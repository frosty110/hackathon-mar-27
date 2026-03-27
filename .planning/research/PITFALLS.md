# Pitfalls Research

**Domain:** Competitive pricing intelligence agent — LLM-based web scraping, normalization, hackathon demo
**Researched:** 2026-03-27
**Confidence:** HIGH (extraction/scraping), MEDIUM (normalization), HIGH (hackathon patterns)

---

## Critical Pitfalls

### Pitfall 1: Raw HTML Overload Kills LLM Extraction Quality

**What goes wrong:**
The full HTML of a modern SaaS pricing page averages 80K+ tokens — navigation, cookie banners, JS bundles, CSS, footers, sidebar ads, nested div soup. Sending this to an LLM causes the model to extract the wrong "Price" field (e.g., a promotional price in the footer instead of the plan price), miss tiered pricing rows buried in noise, and return inconsistent structure across pages. Token cost balloons 10-40x with no quality improvement.

**Why it happens:**
It feels safe to send everything — "the LLM will figure it out." But relevant pricing evidence becomes diluted when surrounded by 90% irrelevant markup, which directly increases hallucination risk. Teams discover this only after seeing extraction errors on page 3 or 4.

**How to avoid:**
Pre-process HTML in Go before sending to the LLM. Use goquery to strip `<script>`, `<style>`, `<nav>`, `<footer>`, `<head>`, cookie banners, and social widgets. Keep: main content div, pricing tables, tier cards, feature lists. Target 2K–8K tokens per page, not 80K. In `internal/scraper/parser.go`, select `.pricing`, `main`, `article`, or similar semantic containers before the LLM call in `internal/extractor/llm.go`.

**Warning signs:**
- Extraction prompt returns prices from "Related products" section
- Same page gives different prices on consecutive runs
- Token count per page exceeds 15K
- LLM returns `null` for clearly visible prices

**Phase to address:** Extraction pipeline (Phase 1 / scraping core)

---

### Pitfall 2: Pricing Model Heterogeneity Makes Normalization a Lies-to-Judges Problem

**What goes wrong:**
Comparing "$20/month/seat" to "$0.003/1K tokens" to "500 credits = $49" is not a math problem — it's an assumption problem. If normalization assumptions are wrong or hidden, the comparison table shown in the demo will be factually wrong, and judges asking "how did you get this number?" will expose it immediately. Silent incorrect normalization is worse than no normalization.

**Why it happens:**
Teams treat normalization as a purely technical task (unit conversion) but it's actually an epistemic task — you need explicit assumptions about usage patterns (e.g., "assume 1M tokens/month/user" or "assume 10 seats"). Without explicit assumptions, any per-seat vs. per-token comparison is undefined.

**How to avoid:**
Every normalized value in Go must carry its assumption set. Use a `NormalizedPricing` struct in `internal/normalizer/normalize.go`:
```go
type NormalizedPricing struct {
    RawValue                  string             `json:"raw_value"`
    RawModel                  string             `json:"raw_model"`
    NormalizedMonthlyUSD      float64            `json:"normalized_monthly_usd"`
    NormalizationAssumptions  map[string]any     `json:"normalization_assumptions"`
    Confidence                string             `json:"confidence"`
}
```
Record assumptions explicitly in `internal/normalizer/assumptions.go`. Surface assumptions in the UI transparency map. Never display a bare normalized number without the assumption footnote.

**Warning signs:**
- Normalization function takes only `price` and `unit` as parameters (not assumption set)
- Comparison table shows a single number per competitor with no footnote
- "per credit" entries are silently skipped or defaulted to zero
- No `Confidence` or `Assumptions` field in the Go struct

**Phase to address:** Normalization engine (Phase 2)

---

### Pitfall 3: JS-Rendered Pricing Pages Break net/http Silently

**What goes wrong:**
Many SaaS pricing pages (especially AI companies) load pricing tiers via JavaScript after the initial HTML response. Go's `net/http` fetches the initial HTML shell — which may contain only a loading spinner or empty `<div id="pricing"></div>`. The page appears to scrape successfully (200 OK, HTML returned), but the LLM extracts nothing or returns placeholder text. This fails silently — no error, just bad data.

**Why it happens:**
`net/http` is an HTTP client, not a browser. It does not execute JavaScript. Teams discover this when they manually visit the page and see pricing, but the scraped HTML contains none of it.

**How to avoid:**
Pre-validate all 5-8 target pages before the hackathon clock starts. For each page: fetch with `net/http` in `internal/scraper/fetcher.go`, scan the raw HTML for a known price string (e.g., "$20" or "per seat"). If not present, that page needs a mock or is disqualified. Build the mock fallback HTML files in `demo-data/` for 2-3 pages before you start coding — not as an afterthought. Configure mock fallback URLs in `internal/scraper/targets.go`.

**Warning signs:**
- `net/http` response body contains `<div id="pricing"></div>` or similar empty containers
- LLM returns `{"plans": []}` for a page that visibly has pricing
- Response HTML size is under 10KB for a complex pricing page
- Known price strings don't appear when scanning the raw response body

**Phase to address:** Target page validation (pre-work before Phase 1)

---

### Pitfall 4: Demo Crashes Because the Scan Takes Too Long Live

**What goes wrong:**
A 3-minute demo budget with 5-8 pages scraped + LLM extraction + normalization + vector upsert + strategy generation means the "scan" button press cannot take 90 seconds. If the Go backend processes pages sequentially without goroutines, each HTTP fetch + LLM call adds 10-15 seconds — easily 60-90 seconds total for 8 pages. The Streamlit UI shows a frozen spinner for the duration.

**Why it happens:**
The pipeline is built to work correctly first, then concurrency is bolted on. Without `errgroup` in the Scan Runner, `fetcher.go` falls back to sequential fetches. Separately, if Streamlit calls the Go API synchronously with no progress indication, judges see a blank screen.

**How to avoid:**
Use `errgroup` for concurrent goroutines in `internal/scraper/fetcher.go` — all 5-8 pages fetched in parallel. Cap `net/http` timeouts at 10 seconds per request via `http.Client{Timeout: 10 * time.Second}`. Pre-run the scan before the demo and cache results in Ghost DB — the demo shows _replay_ of a scan result, with the option to re-trigger for the live change detection story. In Streamlit, use `st.session_state` to track `scan_in_progress` and show `st.progress()` bars while polling the Go API for scan status.

**Warning signs:**
- `fetcher.go` loops over URLs with no goroutines
- `http.Client` created with no `Timeout` field set
- Streamlit calls the Go API and blocks the entire UI thread with no status polling
- Demo script requires waiting for a fresh scan result during the 3 minutes

**Phase to address:** Scan Runner goroutines (Phase 1) and Streamlit UI (Phase 4) and demo preparation (final phase)

---

### Pitfall 5: Hackathon Scope Creep via "Nice to Have" Features

**What goes wrong:**
With AI acceleration, a team can ship features 3-5x faster than expected — so they keep adding them. The result by hour 6 is a product with 8 features at 60% polish each instead of 4 features at 95% polish. The winning demo has fewer features that all work visibly and confidently, not more features that partially work.

**Why it happens:**
Momentum feels like progress. Each additional feature ("let's add historical trend charts", "let's add a competitor similarity score badge") feels low-effort mid-build. The cumulative effect is an unstable demo with untested edge cases.

**How to avoid:**
Define a hard "feature freeze" at hour 5 (3 hours before end). After freeze: no new features, only polish, testing, and demo script rehearsal. Assign one team member as "feature sheriff" whose only job after hour 3 is to say no.

**Warning signs:**
- "This will only take 20 minutes" said more than twice
- Dashboard has more than 4 distinct views or panels
- Any feature added after hour 4 hasn't been tested with real data
- No demo script written by hour 6

**Phase to address:** All phases — enforce at planning time, not retrospectively

---

### Pitfall 6: Ghost Postgres Schema Designed for Hackathon, Not for Change Detection

**What goes wrong:**
Change detection requires comparing _current_ scan vs. _previous_ scan for the same competitor. If the DB schema stores scans as flat rows without a versioned scan_run_id and a competitor key, the change detection query in `internal/detector/changes.go` becomes fragile or impossible. Teams discover this when writing the diff logic in hour 4 and realize the schema can't support it without a rewrite.

**Why it happens:**
Schemas get designed for "storing the data" not "querying the data." The temptation is to store the LLM output JSON blob as a single column and figure out querying later.

**How to avoid:**
Design schema with change detection in mind from the start, applied via `internal/storage/ghost.go` on first run:
```sql
CREATE TABLE scan_runs (
  id UUID PRIMARY KEY,
  triggered_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE pricing_snapshots (
  id SERIAL PRIMARY KEY,
  scan_run_id UUID NOT NULL REFERENCES scan_runs(id),
  competitor TEXT NOT NULL,
  scraped_at TIMESTAMPTZ NOT NULL,
  raw_html TEXT,
  extracted_json JSONB,
  normalized_json JSONB
);

CREATE TABLE change_events (
  id SERIAL PRIMARY KEY,
  scan_run_id UUID NOT NULL REFERENCES scan_runs(id),
  competitor TEXT NOT NULL,
  field TEXT NOT NULL,
  old_value TEXT,
  new_value TEXT,
  detected_at TIMESTAMPTZ NOT NULL
);
```
Change detection = `SELECT * FROM pricing_snapshots WHERE competitor = $1 ORDER BY scraped_at DESC LIMIT 2` via pgx. Seed with at least two scan runs before demo.

**Warning signs:**
- Schema has no `competitor` identifier column (only URLs)
- Snapshots stored as a single blob with no timestamp
- No `scan_run_id` to group pages scraped together
- Change detection query written without looking at actual stored data

**Phase to address:** Ghost DB setup (Phase 1 foundation)

---

### Pitfall 7: Go LLM Extraction Returns Unvalidated JSON Strings

**What goes wrong:**
TrueFoundry's OpenAI-compatible API returns a JSON string in the `content` field. If `internal/extractor/llm.go` does a naïve `json.Unmarshal` on the raw response without validating the extracted schema, malformed or incomplete LLM outputs propagate silently through the pipeline. Downstream, `normalizer.go` panics on nil pointer dereference or produces zero-value floats that look like real prices.

**Why it happens:**
The Go extractor parses the HTTP response correctly (valid JSON envelope) but doesn't validate the domain schema inside `content`. "It unmarshalled" and "it's valid pricing data" are two different checks.

**How to avoid:**
After unmarshalling the LLM response into the `RawPricingData` struct, validate required fields before returning. Add a `1` retry with `2s` backoff on extraction failure. Log the raw `content` string on validation failure for post-mortem debugging.

**Warning signs:**
- `extractor/llm.go` has no field validation after `json.Unmarshal`
- `NormalizedPricing.NormalizedMonthlyUSD` is 0.0 for multiple competitors in a scan
- No retry logic in `llm.go`
- Panic in `normalizer.go` with nil pointer on `RawPricingData.Plans`

**Phase to address:** Extraction pipeline (Phase 1)

---

## Technical Debt Patterns

Shortcuts that seem reasonable in a hackathon context.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Send full raw HTML to LLM | No goquery preprocessing code | 10-40x token cost, hallucination-prone extraction | Never — even for demo |
| Hardcode normalization assumptions | No config system needed | Wrong numbers shown as facts | Only if assumptions are displayed inline in UI |
| No retry on LLM extraction failure | Simpler Go pipeline | Demo fails on transient TrueFoundry API error | Never — add 1 retry with 2s backoff |
| Skip mock fallback pages | Saves 30 min | Demo breaks if any target page fails during live demo | Never — mock at least 2 pages in demo-data/ |
| Store full scan in Streamlit session state | Simple to implement | Rerenders wipe scan results | Never — always persist to Ghost DB |
| Run Go scan synchronously (no errgroup) | Easy to code | Sequential fetches take 60-90s for 8 pages | Never — use errgroup from the start |

---

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Aerospike Vector Search | Creating the vector client on every request — has startup overhead | Create once in `storage/aerospike.go` at server startup, reuse via dependency injection |
| Aerospike Vector Search | Forgetting to include vector field in results (excluded by default) | Explicitly request vector field when needed in ANN query |
| Aerospike Vector Search | Not closing client on shutdown | Use `defer client.Close()` in `cmd/server/main.go` |
| TrueFoundry model routing | Using expensive model (Opus/GPT-4) for both extraction and strategy | Route: cheap model (Haiku/Flash) for `extractor/llm.go`, expensive only for `architect/strategy.go` |
| Ghost Postgres | Relying on autonomous DB creation without verifying schema before first scan | Create and verify schema in `storage/ghost.go` init, call it from `cmd/server/main.go` on startup |
| net/http scraper | No timeout set on `http.Client` — hangs indefinitely on unresponsive pages | Set `Timeout: 10 * time.Second` on the `http.Client` in `scraper/fetcher.go` |
| net/http scraper | No User-Agent header — some pages return 403 | Set a browser-like User-Agent string on each request |
| Connect-Go + Streamlit | Streamlit calls Connect-Go endpoint with wrong Content-Type | Connect-Go JSON mode requires `Content-Type: application/json`; use `requests.post(..., json=payload)` |
| pgx | Opening a new connection per request instead of using a pool | Use `pgxpool.New` in `storage/ghost.go`; inject pool into handlers |
| Go errgroup | Ignoring errors returned from goroutines | Always check `g.Wait()` error; propagate fetch failures to the handler |

---

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Sequential net/http page fetches (no errgroup) | Scan takes 60-90s for 8 pages | Use `errgroup` + goroutines in `scraper/fetcher.go` | At 4+ pages |
| Full HTML sent to LLM | $0.50+ per scan run, 30s latency per page | Strip to main content with goquery first, target <8K tokens | At 3+ pages |
| Aerospike client created per request | Slow upserts, connection overhead | Create client once at startup in `storage/aerospike.go` | At 5+ upserts |
| Streamlit polls Go API on every widget interaction | Duplicate scan triggers, race conditions | Cache scan status in `st.session_state`; poll only during active scan | Immediately |
| pgx pool not used | Connection exhaustion under concurrent goroutines | Use `pgxpool` — required when multiple goroutines write to Ghost DB simultaneously | At 3+ concurrent writes |

---

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| Hardcoding API keys (TrueFoundry, Aerospike) in Go source | Keys visible in demo screen share or GitHub | Load from `.env` via `internal/config/config.go` using `os.Getenv`; never commit `.env` |
| Logging raw scraped HTML in Go | Competitor HTML in logs if demo screen shared | Log only metadata (URL, status code, token count) in `scraper/fetcher.go`, not raw content |
| No rate limiting on net/http scraper | IP ban from target pages mid-demo | Add `time.Sleep(1 * time.Second)` between requests in the fetcher loop |

---

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Showing raw normalized numbers with no assumptions | Judges ask "how did you calculate this?" — exposes methodology flaw | Show inline assumption footnotes: "Assumes 1M tokens/month" (surfaced from `NormalizationAssumptions` in Go response) |
| Change detection shown as raw JSON diff | Incomprehensible to non-technical judges | Show human-readable diff: "OpenAI raised Pro plan from $20 to $25/month" (format in Streamlit from `ChangeEvent` struct) |
| Cluster visualization without labels | "Prices like us" graph has no context | Label each node with competitor name + model type |
| Scan button triggers full re-scrape in demo | 60s blank screen mid-demo | Pre-run scan, show results from Ghost DB via Go API, offer "refresh" only as optional |
| Dashboard shows loading spinners for all 8 pages simultaneously | Overwhelming, appears broken | Show sequential progress with per-competitor status updates polled from Go API |

---

## "Looks Done But Isn't" Checklist

- [ ] **Extraction pipeline:** Go extractor returns valid JSON — verify it validates against the `RawPricingData` struct fields, not just that `json.Unmarshal` succeeds
- [ ] **Normalization:** Numbers in comparison table look correct — verify they match manually calculated values for at least 2 competitors
- [ ] **Change detection:** Alert fires for changed pricing — verify it correctly ignores re-runs of identical data (no false positives)
- [ ] **Aerospike similarity:** "Prices like us" cluster returns results — verify it returns meaningful groupings, not random ordering
- [ ] **TrueFoundry routing:** Code routes cheap model for extraction — verify expensive model is NOT called for extraction (check TrueFoundry gateway logs/costs)
- [ ] **Mock fallback:** Mock pages exist in `demo-data/` — verify mock HTML actually contains extractable pricing (not placeholder text)
- [ ] **Go API boots cleanly:** `cmd/server/main.go` starts, connects to Ghost DB, initializes Aerospike client, and listens on `:8080` without errors
- [ ] **Streamlit connects to Go API:** `frontend/app.py` posts to `localhost:8080` and receives valid JSON responses
- [ ] **Demo script:** Full demo runs — time it with a stopwatch, must be under 3 minutes end-to-end

---

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| JS-rendered page fails during demo | LOW (if prepared) | Switch to pre-seeded mock HTML in `demo-data/`; narrate as "we mock pages that require browser rendering" |
| LLM extraction returns bad JSON | LOW | Retry logic in `extractor/llm.go` handles transient failures; fallback to last good snapshot from Ghost DB |
| Normalization shown as wrong numbers | HIGH | If caught before demo: fix assumptions in `normalizer/assumptions.go`. If caught during demo: pivot to showing raw values + methodology |
| Aerospike vector query returns empty | MEDIUM | Pre-populate with at least 5 competitor embeddings in setup via `embedder/embed.go`; demo cannot work with empty index |
| Go API scan takes too long live | MEDIUM | Pre-run all scans before demo; Streamlit shows cached Ghost DB results via `GetComparison`, scan is optional live feature |
| Ghost DB schema prevents change detection | HIGH | If caught in first 2 hours: redesign schema in `storage/ghost.go`. If caught in hour 5+: hardcode a mock `ChangeEvent` slice for demo |
| Connect-Go endpoint returns error to Streamlit | MEDIUM | Check Go server logs first; most likely a schema mismatch or missing env var — fix in `config/config.go` |

---

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Raw HTML overload | Phase 1: Extraction pipeline | Token count logged per page in `extractor/llm.go`; extraction accuracy spot-checked on 3 pages |
| Normalization assumption hiding | Phase 2: Normalization engine | UI shows assumption footnotes; `NormalizedPricing` struct has `NormalizationAssumptions` field |
| JS-rendered page detection | Pre-work: Target page validation | All 8 pages grep-validated for price strings before coding starts; mocks in `demo-data/` |
| Demo scan latency | Phase 1: errgroup concurrency | 8-page scan completes in under 15s with goroutines; demo dry run timed under 3 minutes |
| Feature creep | All phases (planning constraint) | Feature freeze declared at hour 5; no new features added after |
| Ghost schema misfit for change detection | Phase 1: DB setup | Change detection query in `detector/changes.go` runs successfully on seeded data before Phase 2 |
| LLM extraction validation gap | Phase 1: Extraction pipeline | `extractor/llm.go` validates required struct fields; retry tested with mocked failure |

---

## Sources

- [Web Scraping with AI: How LLMs Are Transforming Data Extraction in 2026](https://use-apify.com/blog/web-scraping-with-ai-llms-2026) — HTML noise, token optimization, extraction best practices
- [DRIPPER: Token-Efficient Main HTML Extraction](https://openreview.net/pdf/e2b774a7481c9ccba439fa31dd837e9e32088b81.pdf) — 80K avg token count for raw HTML, 90% noise figure
- [Structured Data Extraction from Unstructured Content using LLM Schemas](https://simonwillison.net/2025/Feb/28/llm-schemas/) — JSON schema, hallucination mitigation, validation
- [Web Scraping for Pricing Intelligence](https://www.zyte.com/learn/web-scraping-for-pricing-intelligence/) — change detection, data quality, staleness pitfalls
- [Scraping SaaS Websites for Pricing Intelligence](https://www.just3things.com/scraping-saas-websites-for-pricing-intelligence-a-comprehensive-guide-to-competitive-analysis/) — JS rendering limitation with static scrapers
- [Per-Seat vs Per-Token vs Per-Output: Financial Tradeoffs](https://www.afternoon.co/blog/seat-token-output) — heterogeneous model comparison complexity
- [goquery](https://github.com/PuerkitoBio/goquery) — Go HTML parsing and selection for scraper preprocessing
- [pgx v5](https://github.com/jackc/pgx) — Go Postgres driver; pool usage, query patterns
- [Aerospike Go client](https://github.com/aerospike/aerospike-client-go) — client lifecycle, vector search gotchas
- [Connect-Go getting started](https://connectrpc.com/docs/go/getting-started/) — JSON mode, Content-Type requirements
- [Maintaining State When Working with LLMs in Streamlit](https://discuss.streamlit.io/t/maintaining-state-when-working-with-llms/46631) — session state race conditions
- [Things I Learned by Participating in GenAI Hackathons](https://towardsdatascience.com/things-i-learnt-by-participating-in-genai-hackathons-over-the-past-6-months/) — feature creep, scope, demo polish patterns
- [Hallucination in Long-Context LLMs](https://arxiv.org/html/2601.02023) — diluted evidence in large contexts increases hallucination risk

---
*Pitfalls research for: Pricing Radar — competitive pricing intelligence agent*
*Researched: 2026-03-27 (updated to reflect Go backend + Streamlit frontend architecture)*
