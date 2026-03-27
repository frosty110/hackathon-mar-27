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
Pre-process HTML before sending to LLM. Strip: `<script>`, `<style>`, `<nav>`, `<footer>`, `<head>`, cookie banners, social widgets. Keep: main content div, pricing tables, tier cards, feature lists. Target 2K–8K tokens per page, not 80K. Use BeautifulSoup to extract `.pricing`, `main`, `article`, or similar semantic containers before LLM call.

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
Every normalized value must carry its assumption set. Use a `NormalizedPrice` schema:
```
{
  "raw_value": "...",
  "raw_model": "per_token",
  "normalized_monthly_usd": 42.00,
  "normalization_assumptions": {"tokens_per_month": 1000000},
  "confidence": "medium"
}
```
Surface assumptions in the UI transparency map. Never display a bare normalized number without the assumption footnote.

**Warning signs:**
- Normalization function takes only `price` and `unit` as inputs (not assumptions)
- Comparison table shows a single number per competitor with no footnote
- "per credit" entries are silently skipped or defaulted to zero
- No `confidence` or `assumption` field in the data schema

**Phase to address:** Normalization engine (Phase 2)

---

### Pitfall 3: JS-Rendered Pricing Pages Break httpx Silently

**What goes wrong:**
Many SaaS pricing pages (especially AI companies) load pricing tiers via JavaScript after the initial HTML response. httpx fetches the initial HTML shell — which may contain only a loading spinner or empty `<div id="pricing"></div>`. The page appears to scrape successfully (200 OK, HTML returned), but the LLM extracts nothing or returns placeholder text. This fails silently — no error, just bad data.

**Why it happens:**
httpx is an HTTP client, not a browser. It does not execute JavaScript. Teams discover this when they manually visit the page and see pricing, but the scraped HTML contains none of it.

**How to avoid:**
Pre-validate all 5-8 target pages before the hackathon clock starts. For each page: fetch with httpx, grep the raw HTML for a known price string (e.g., "$20" or "per seat"). If not present, that page needs a mock or is disqualified. Build the mock fallback HTML files for 2-3 pages before you start coding — not as an afterthought.

**Warning signs:**
- httpx response body contains `<div id="pricing"></div>` or similar empty containers
- LLM returns `{"plans": []}` for a page that visibly has pricing
- Response HTML size is under 10KB for a complex pricing page
- Known price strings don't appear in `grep` of raw response

**Phase to address:** Target page validation (pre-work before Phase 1)

---

### Pitfall 4: Demo Crashes Because the Scan Takes Too Long Live

**What goes wrong:**
A 3-minute demo budget with 5-8 pages scraped + LLM extraction + normalization + vector upsert + strategy generation means the "scan" button press cannot take 90 seconds. If the scan blocks the UI thread or runs synchronously with no progress indicators, judges see a frozen screen. If it runs in a background thread with Streamlit, session state race conditions cause silent failures or duplicate LLM calls.

**Why it happens:**
Teams build the pipeline to work correctly, then bolt on the UI last. Streamlit reruns the entire script on every interaction — if a scan is triggered by a button click without proper session state guards, it re-fires on every rerender. Async scraping + Streamlit's threading model is a known incompatibility.

**How to avoid:**
Use `st.session_state` to track `scan_in_progress = True` and gate the button. Pre-run the scan before the demo and cache results in the database — the demo shows _replay_ of a scan result, with the option to re-trigger for the live change detection story. Show `st.progress()` bars and `st.status()` containers so the UI is never frozen. Cap httpx timeouts at 10 seconds per page.

**Warning signs:**
- Scan triggered by `if st.button("Scan"):` without `st.session_state` guard
- No timeout on httpx requests
- No loading indicator during LLM extraction
- Demo script requires waiting for a fresh scan result during the 3 minutes

**Phase to address:** Streamlit UI (Phase 4) and demo preparation (final phase)

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
Change detection requires comparing _current_ scan vs. _previous_ scan for the same competitor. If the DB schema stores scans as flat rows without a versioned scan_id and a competitor key, the change detection query becomes fragile or impossible. Teams discover this when writing the change detection logic in hour 4 and realize the schema can't support it without a rewrite.

**Why it happens:**
Schemas get designed for "storing the data" not "querying the data." The temptation is to store the LLM output JSON blob as a single column and figure out querying later.

**How to avoid:**
Design schema with change detection in mind from the start:
```sql
CREATE TABLE scans (
  id SERIAL PRIMARY KEY,
  scan_run_id UUID NOT NULL,  -- groups all pages in one scan
  competitor TEXT NOT NULL,
  scraped_at TIMESTAMPTZ NOT NULL,
  raw_html TEXT,
  extracted_json JSONB,
  normalized_json JSONB
);
```
Change detection = `SELECT * FROM scans WHERE competitor = $1 ORDER BY scraped_at DESC LIMIT 2`. Seed with at least two scan runs before demo.

**Warning signs:**
- Schema has no `competitor` identifier column (only URLs)
- Scans stored as a single blob with no timestamp
- No scan_run_id to group pages scraped together
- Change detection query written without looking at actual stored data

**Phase to address:** Ghost DB setup (Phase 1 foundation)

---

## Technical Debt Patterns

Shortcuts that seem reasonable in a hackathon context.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Send full raw HTML to LLM | No preprocessing code | 10-40x token cost, hallucination-prone extraction | Never — even for demo |
| Hardcode normalization assumptions | No config system needed | Wrong numbers shown as facts | Only if assumptions are displayed inline in UI |
| No retry on LLM extraction failure | Simpler pipeline | Demo fails on transient API error | Never — add 1 retry with 2s backoff |
| Skip mock fallback pages | Saves 30 min | Demo breaks if any target page fails during live demo | Never — mock at least 2 pages |
| Store full scan in session state | Simple to implement | Streamlit rerenders wipe scan results | Never — always persist to DB |
| Run scan synchronously in Streamlit | Easy to code | UI freezes, judges see blank screen | Never — use pre-run + cached display |

---

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Aerospike Vector Search | Creating Index object on every request — has overhead per creation | Create once, store reference, reuse across queries |
| Aerospike Vector Search | Forgetting to include vector field in results (excluded by default with Index object) | Explicitly pass `include_fields=["vector"]` when needed |
| Aerospike Vector Search | Not closing client on app shutdown | Use context manager or explicit `.close()` call |
| TrueFoundry model routing | Using expensive model (Opus/GPT-4) for both extraction and strategy | Route: cheap model (Haiku/Flash) for extraction, expensive only for strategy generation |
| Ghost Postgres | Relying on autonomous DB creation without verifying schema before first scan | Create and verify schema in setup phase, don't assume auto-creation |
| httpx | No timeout set — hangs indefinitely on unresponsive pages | Set `timeout=10.0` on every request |
| httpx | No User-Agent header — some pages return 403 | Set a browser-like User-Agent string |
| Streamlit + LLM | Triggering LLM call on every rerender via `if st.button(...)` without state guard | Guard with `if "scan_done" not in st.session_state` |

---

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Sequential httpx page fetches | Scan takes 60+ seconds for 8 pages | Use `asyncio` + `httpx.AsyncClient` for concurrent fetches | At 4+ pages |
| Full HTML sent to LLM | $0.50+ per scan run, 30s latency per page | Strip to main content first, target <8K tokens | At 3+ pages |
| Aerospike index created in loop | Slow upserts, resource leak | Create index once at startup | At 5+ upserts |
| Streamlit re-runs LLM call on widget interaction | Duplicate API calls, race conditions | Cache results in `st.session_state` or DB before rendering | Immediately |

---

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| Hardcoding API keys (TrueFoundry, Aerospike) in source code | Keys visible in demo screen share or GitHub | Use `.env` + `python-dotenv`, never commit keys |
| Logging raw scraped HTML | Competitor HTML in logs if demo screen shared | Log only metadata (URL, status, token count), not raw content |
| No rate limiting on httpx | IP ban from target pages mid-demo | Add `asyncio.sleep(1)` between requests, respect robots.txt minimally |

---

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Showing raw normalized numbers with no assumptions | Judges ask "how did you calculate this?" — exposes methodology flaw | Show inline assumption footnotes: "Assumes 1M tokens/month" |
| Change detection shown as raw JSON diff | Incomprehensible to non-technical judges | Show human-readable diff: "OpenAI raised Pro plan from $20 to $25/month" |
| Cluster visualization without labels | "Prices like us" graph has no context | Label each node with competitor name + model type |
| Scan button triggers full re-scrape in demo | 60s blank screen mid-demo | Pre-run scan, show results from DB, offer "refresh" only as optional |
| Dashboard shows loading spinners for all 8 pages simultaneously | Overwhelming, appears broken | Show sequential progress with per-competitor status updates |

---

## "Looks Done But Isn't" Checklist

- [ ] **Extraction pipeline:** LLM returns valid JSON — verify it validates against schema, not just that it's parseable JSON
- [ ] **Normalization:** Numbers in comparison table look correct — verify they match manually calculated values for at least 2 competitors
- [ ] **Change detection:** Alert fires for changed pricing — verify it correctly ignores re-runs of identical data (no false positives)
- [ ] **Aerospike similarity:** "Prices like us" cluster returns results — verify it returns meaningful groupings, not random ordering
- [ ] **TrueFoundry routing:** Code routes cheap model for extraction — verify expensive model is NOT called for extraction (check logs/costs)
- [ ] **Mock fallback:** Mock pages exist — verify mock HTML actually contains extractable pricing (not placeholder text)
- [ ] **Demo script:** Full demo runs — time it with a stopwatch, must be under 3 minutes end-to-end

---

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| JS-rendered page fails during demo | LOW (if prepared) | Switch to pre-seeded mock HTML fallback; narrate as "we mock pages that require browser rendering" |
| LLM extraction returns bad JSON | LOW | Add `json.loads()` with try/except + retry; fallback to last good extraction from DB |
| Normalization shown as wrong numbers | HIGH | If caught before demo: fix assumptions. If caught during demo: pivot to showing raw values + methodology |
| Aerospike vector query returns empty | MEDIUM | Pre-populate with at least 5 competitor embeddings in setup; demo cannot work with empty index |
| Streamlit UI freezes during scan | MEDIUM | Pre-run all scans before demo; demo shows cached DB results, scan is optional live feature |
| Ghost DB schema prevents change detection | HIGH | If caught in first 2 hours: redesign. If caught in hour 5+: hardcode a mock diff for demo |

---

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Raw HTML overload | Phase 1: Extraction pipeline | Token count logged per page; extraction accuracy spot-checked on 3 pages |
| Normalization assumption hiding | Phase 2: Normalization engine | UI shows assumption footnotes; schema has `normalization_assumptions` field |
| JS-rendered page detection | Pre-work: Target page validation | All 8 pages grep-validated for price strings before coding starts |
| Demo scan latency | Phase 4: Streamlit UI | Demo dry run timed under 3 minutes with pre-run scan |
| Feature creep | All phases (planning constraint) | Feature freeze declared at hour 5; no new features added after |
| Ghost schema misfit for change detection | Phase 1: DB setup | Change detection query runs successfully on seeded data before Phase 2 |

---

## Sources

- [Web Scraping with AI: How LLMs Are Transforming Data Extraction in 2026](https://use-apify.com/blog/web-scraping-with-ai-llms-2026) — HTML noise, token optimization, extraction best practices
- [DRIPPER: Token-Efficient Main HTML Extraction](https://openreview.net/pdf/e2b774a7481c9ccba439fa31dd837e9e32088b81.pdf) — 80K avg token count for raw HTML, 90% noise figure
- [Structured Data Extraction from Unstructured Content using LLM Schemas](https://simonwillison.net/2025/Feb/28/llm-schemas/) — JSON schema, hallucination mitigation, validation
- [Web Scraping for Pricing Intelligence](https://www.zyte.com/learn/web-scraping-for-pricing-intelligence/) — change detection, data quality, staleness pitfalls
- [Scraping SaaS Websites for Pricing Intelligence](https://www.just3things.com/scraping-saas-websites-for-pricing-intelligence-a-comprehensive-guide-to-competitive-analysis/) — JS rendering limitation with static scrapers
- [Per-Seat vs Per-Token vs Per-Output: Financial Tradeoffs](https://www.afternoon.co/blog/seat-token-output) — heterogeneous model comparison complexity
- [Aerospike Vector Search Python Client](https://aerospike-vector-search-python-client.readthedocs.io/en/latest/client.html) — Index object overhead, include_fields gotcha
- [Maintaining State When Working with LLMs in Streamlit](https://discuss.streamlit.io/t/maintaining-state-when-working-with-llms/46631) — session state race conditions
- [Things I Learned by Participating in GenAI Hackathons](https://towardsdatascience.com/things-i-learnt-by-participating-in-genai-hackathons-over-the-past-6-months/) — feature creep, scope, demo polish patterns
- [Hallucination in Long-Context LLMs](https://arxiv.org/html/2601.02023) — diluted evidence in large contexts increases hallucination risk

---
*Pitfalls research for: Pricing Radar — competitive pricing intelligence agent*
*Researched: 2026-03-27*
