# Feature Research

**Domain:** Competitive pricing intelligence for AI/SaaS products
**Researched:** 2026-03-27
**Confidence:** MEDIUM — existing tools (Tierly, Prisync, Priceva, Visualping) verified via official sites and reviews; market gaps inferred from multiple converging sources

---

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Automated pricing page scraping | Manual tracking is the pain being replaced; if you don't automate it, the product has no reason to exist | MEDIUM | httpx + LLM extraction; JS-heavy pages excluded from MVP |
| Structured data extraction (tiers, prices, features) | Raw HTML is useless; users need tier names, prices, features per tier in a table | MEDIUM | LLM extraction handles heterogeneous page layouts better than regex |
| Side-by-side competitor comparison table | Core mental model for Head of Product; spreadsheet replacement must look like a spreadsheet | LOW | Streamlit table; columns = competitors, rows = pricing dimensions |
| Change detection with alerts | Pricing pages change silently; without detection, tool is a one-shot not a monitor | MEDIUM | Diff against stored previous scan; highlight what changed |
| Historical scan storage | Can't detect changes without history; also needed for trend analysis | LOW | Ghost Postgres stores each scan with timestamp |
| Pricing model identification | AI pricing has 5+ model types; users need to know if competitor switched from seats to tokens | LOW | LLM classifies model type as part of extraction |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valued.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Cross-model normalization to common unit | Prisync, Priceva, and Visualping are ecommerce-focused and don't handle token/credit/seat normalization; Tierly scores tiers but doesn't normalize to $/unit | HIGH | Hardest problem in the domain; requires explicit assumptions ("100 MAUs at average usage"); the assumptions being surfaced is itself differentiating |
| Explicit assumption transparency | When normalizing credits to $/MAU, the assumptions are surfaced so the user can challenge them | LOW | Display normalization assumptions inline with the table; Tierly does scoring but not transparent cost math |
| "Who prices like us?" clustering | Vector similarity on pricing profiles surfaces structural positioning insights ("you price like Anthropic, not OpenAI") | HIGH | Aerospike vector search; novel capability not present in any surveyed tool |
| Strategy-grounded recommendations | Tierly gives tactical scoring; no tool grounds recommendations in your own positioning docs and strategy | HIGH | RAG over uploaded strategy docs; TrueFoundry routes to high-capability model for this step only |
| Model routing (cheap extraction, expensive strategy) | Cost optimization is meaningful for a tool that runs on-demand; LLM costs matter | MEDIUM | TrueFoundry model routing; Haiku/Flash for scraping, Opus/GPT-4 for strategy |
| Normalization assumption export | Users can copy the assumptions into their own analysis; supports the spreadsheet workflow they're replacing | LOW | Simple export of assumption table alongside comparison |
| Change event narrative | When a competitor changes pricing, the tool writes a plain-English summary of what changed and why it matters | LOW | LLM summarization step after diff detection |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Arbitrary URL support (any competitor URL) | Seems like flexibility | JS-rendered pages fail with httpx; Playwright adds 3-4 hours of build risk in an 8-hour hackathon; demo breaks on live JS-heavy page | Pre-selected pages that are verified to work with httpx; localhost mock fallback for demo stability |
| Real-time WebSocket price feeds | Feels more impressive in demo | Adds infrastructure complexity (WebSocket server, connection management) with no UX benefit over polling for a batch-scan tool | Polling on-demand scan; "last scanned" timestamp displayed |
| User authentication and multi-tenant | Seems like a complete product | Doubles build scope; demo is single-user by nature; auth complexity increases demo failure risk | Single-user; document that auth is a post-hackathon concern |
| Automated repricing (push changes to your own product) | Closes the loop between intelligence and action | Requires product API integration, approval workflows, rollback logic; scope is 8x larger | Strategic recommendation text only; human executes pricing change |
| Email alert subscriptions | Looks like a production feature | Requires email infrastructure (SMTP, templates, deliverability); demo can't show email receipt in 3 minutes | In-app change badge + Slack webhook as a stretch goal |
| Dynamic pricing rules engine | Natural extension of competitor data | B2B SaaS pricing changes quarterly, not hourly; rules engine solves e-commerce, not the AI startup persona | Human-in-the-loop with strategic recommendations |
| Competitor discovery (auto-find competitors) | Reduces setup friction | Tierly's auto-discovery takes 13 minutes; unverified competitors yield scraping failures; demo needs deterministic results | Manual competitor list pre-loaded at demo start |

---

## Feature Dependencies

```
Scraping Pipeline
    └──requires──> Structured Extraction (LLM)
                       └──requires──> Normalization Engine
                                          └──requires──> Historical Storage (Ghost DB)
                                                             └──requires──> Change Detection
                                                                                └──enables──> Change Narrative (LLM)

Normalization Engine ──requires──> Assumption Transparency (output surface)

Strategy Documents (uploaded)
    └──enables──> Strategy-Grounded Recommendations (RAG)

Historical Storage (Ghost DB)
    └──enables──> Pricing Profile Embeddings
                      └──requires──> Aerospike Vector DB
                                         └──enables──> "Who prices like us?" Clustering

TrueFoundry Model Routing
    └──enhances──> Scraping Pipeline (cheap model) + Strategy Recommendations (expensive model)

Change Detection ──conflicts-with──> First-run scan (no history exists yet; must handle gracefully)
```

### Dependency Notes

- **Normalization requires Structured Extraction:** You cannot normalize prices to a common unit until tiers, prices, and model types have been extracted as structured data. Raw text is insufficient.
- **Change Detection requires Historical Storage:** The diff is between current scan and previous scan. Ghost DB must have at least one prior scan stored. First run shows "baseline established" not "change detected."
- **Clustering requires Embeddings requires Storage:** Aerospike vector search operates on pricing profile vectors. These are generated from normalized comparison data after extraction and storage. Cannot cluster on first run.
- **Strategy Recommendations enhance (not require) Normalization:** Recommendations are richer when strategy docs are provided, but the comparison table stands alone without them. RAG is an enhancement layer.
- **Assumption Transparency conflicts with opaque normalization:** If normalization is done without surfacing assumptions, users distrust the numbers. Transparency is non-optional for this persona (Head of Product is analytical; they will probe the math).

---

## MVP Definition

### Launch With (v1 — Hackathon demo, 8 hours)

Minimum viable product — what's needed to validate the concept and win the demo.

- [ ] Scraping pipeline for 5-8 pre-selected AI company pricing pages — core value delivery; no scraping = no product
- [ ] LLM extraction to structured data (tier names, prices, features, model type) — required before anything else works
- [ ] Normalization to common unit with explicit assumptions surfaced — the key differentiator missing from all competitors; must be present to make the demo claim
- [ ] Ghost DB storage of scan results — required for change detection and clustering; also satisfies sponsor integration
- [ ] Change detection with diff display — makes the demo story work (catch Macroscope transition live or via mock)
- [ ] Strategy-grounded recommendations via TrueFoundry model routing — satisfies both autonomy judging criteria and TrueFoundry sponsor integration
- [ ] "Who prices like us?" clustering via Aerospike — satisfies Aerospike sponsor integration; genuinely novel feature
- [ ] Streamlit dashboard with comparison table, change alerts, normalization assumptions, cluster visual — demo surface; judges evaluate on presentation

### Add After Validation (v1.x)

Features to add once core is working and post-hackathon.

- [ ] Scheduled recurring scans (cron-based) — adds continuous monitoring; blocked by the single-user session model in hackathon build
- [ ] Slack/webhook notifications on change detection — closes the alerting loop; needs email/Slack infrastructure
- [ ] Normalization assumption export to CSV — users will want to take the comparison into their own models
- [ ] More competitor pages beyond pre-selected set — needs robust JS-rendering handling (Playwright) for arbitrary URLs

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] Playwright/browser rendering for JS-heavy pages — significant reliability increase; 3-4 hours of build work that breaks demo reliability
- [ ] User authentication and multi-tenant — required for any commercial product; not needed for single-user validation
- [ ] Competitor auto-discovery — Tierly takes 13 minutes; valuable but unreliable during live demo
- [ ] Pricing elasticity analysis (internal) — different product; requires your own pricing experiment data not competitor data

---

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Scraping pipeline (httpx + LLM extraction) | HIGH | MEDIUM | P1 |
| Normalization with explicit assumptions | HIGH | HIGH | P1 |
| Historical storage (Ghost DB) | HIGH | LOW | P1 |
| Change detection + diff display | HIGH | MEDIUM | P1 |
| Strategy recommendations (TrueFoundry + RAG) | HIGH | MEDIUM | P1 |
| Pricing similarity clustering (Aerospike) | MEDIUM | HIGH | P1 (sponsor) |
| Streamlit dashboard | HIGH | LOW | P1 |
| Scheduled recurring scans | MEDIUM | MEDIUM | P2 |
| Slack/webhook alerts | MEDIUM | MEDIUM | P2 |
| Assumption export to CSV | MEDIUM | LOW | P2 |
| Playwright JS rendering | MEDIUM | HIGH | P3 |
| Competitor auto-discovery | LOW | HIGH | P3 |
| User authentication | LOW | HIGH | P3 |

**Priority key:**
- P1: Must have for launch (hackathon demo)
- P2: Should have, add post-hackathon validation
- P3: Nice to have, future consideration

---

## Competitor Feature Analysis

| Feature | Tierly | Prisync | Visualping | Priceva | Pricing Radar (this project) |
|---------|--------|---------|------------|---------|------------------------------|
| Target market | SaaS founders | E-commerce retailers | General web monitoring | E-commerce/retail | AI startup product teams |
| Pricing model normalization (seats/tokens/credits) | No — scores tiers, no cost math | No — SKU prices only | No — visual diff only | No — retail prices only | YES — explicit $/unit with assumptions |
| Cross-model comparison | Partial (tier scoring) | No | No | No | YES — normalize to common unit |
| Change detection | Reanalysis only (manual) | YES | YES | YES | YES |
| LLM-powered extraction | YES (tier features) | No | Partial (AI summarize) | No | YES (structured extraction) |
| Historical scan storage | No | YES | YES | YES | YES (Ghost DB) |
| Strategic recommendations grounded in your docs | No | No | No | No | YES (RAG over strategy docs) |
| Vector similarity clustering | No | No | No | No | YES (Aerospike) |
| Model routing (cost optimization) | No | No | No | No | YES (TrueFoundry) |
| Assumption transparency | No | No | No | No | YES (surfaces normalization math) |
| Handles "Contact Sales" tiers | Partial | No | No | No | YES (flags as opaque, estimates range) |
| Price: free/accessible | Pay-per-analysis | $99+/mo | Free tier available | Custom | Hackathon demo (free) |

### Key Gap Summary

All surveyed tools fail in at least one of these ways for the AI startup persona:

1. **E-commerce bias:** Prisync and Priceva built for SKU prices (one number per product), not heterogeneous SaaS tiers with feature gates, usage limits, and multiple billing dimensions.
2. **No normalization math:** Tierly scores tiers but never answers "which competitor is actually cheaper at 500 MAUs?" The cost comparison is absent.
3. **Monitoring without meaning:** Visualping detects page changes but produces a screenshot diff, not a structured "they added a new tier at $X with these features."
4. **No strategy grounding:** Zero tools connect competitor pricing data to your own positioning strategy. They describe the landscape but don't say what you should do about it.
5. **Manual recurrence:** Tierly requires re-running analysis manually; no continuous monitoring. Visualping monitors but doesn't interpret.

---

## Sources

- [Tierly homepage — features and tier analysis](https://tierly.app/)
- [Tierly: Best Pricing Intelligence Tools for SaaS 2026](https://tierly.app/blog/best-pricing-intelligence-tools)
- [Tierly: Competitive Intelligence Report Guide for SaaS 2026](https://tierly.app/blog/competitive-intelligence-report-guide)
- [Prisync: AI competitor price tracking features](https://prisync.com/)
- [Priceva: Price intelligence software](https://priceva.com/price-intelligence)
- [Visualping: Top Competitor Price Tracking Tools 2026](https://visualping.io/blog/top-tools-competitor-price-tracking)
- [Visualping: Competitor Pricing Monitoring](https://visualping.io/blog/competitor-pricing-change-alerts)
- [Brightdata: Best Price Intelligence Tools 2026](https://brightdata.com/blog/web-data/best-price-intelligence-tools)
- [Getmonetizely: The 2026 Guide to SaaS, AI, and Agentic Pricing Models](https://www.getmonetizely.com/blogs/the-2026-guide-to-saas-ai-and-agentic-pricing-models)
- [Data-Mania: AI Pricing Models Explained: Usage, Seats, Credits, and Outcome-Based Options](https://www.data-mania.com/blog/ai-pricing-models-explained-usage-seats-credits-outcome-based-options/)

---
*Feature research for: competitive pricing intelligence — AI/SaaS focus*
*Researched: 2026-03-27*
