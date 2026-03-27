# Requirements: Pricing Radar

**Defined:** 2026-03-27
**Core Value:** Continuous competitive pricing intelligence that turns a 2-day manual spreadsheet into a 38-second automated scan with strategic recommendations

## v1 Requirements

Requirements for hackathon demo. Each maps to roadmap phases.

### Extraction

- [ ] **EXTR-01**: Agent fetches 5-8 pre-selected AI company pricing pages in parallel via httpx
- [ ] **EXTR-02**: Raw HTML is pre-processed (strip nav, footer, scripts, CSS) before LLM extraction to reduce token cost
- [ ] **EXTR-03**: LLM extracts structured pricing data (company, tiers, prices, features, limits, pricing model type) into pydantic schema
- [ ] **EXTR-04**: Extraction handles "Contact Sales" tiers by flagging as opaque with confidence score
- [ ] **EXTR-05**: Failed fetches fall back to cached HTML from local file and flag as "cached" in output

### Storage

- [ ] **STOR-01**: Ghost Postgres DB is created autonomously by the agent at startup
- [ ] **STOR-02**: Each scan run inserts new rows with scan_run_id, competitor, and scraped_at timestamp
- [ ] **STOR-03**: Historical scan data is queryable for change detection comparison

### Normalization

- [ ] **NORM-01**: Heterogeneous pricing models (per-seat, per-token, credits, hybrid) are normalized to a common monthly cost
- [ ] **NORM-02**: Normalization uses an explicit reference workload (e.g., "50-person team, 1M tokens/month")
- [ ] **NORM-03**: Normalization assumptions are displayed inline with every normalized figure
- [ ] **NORM-04**: Pricing models that cannot be meaningfully normalized display raw pricing with a flag

### Change Detection

- [ ] **CHNG-01**: Current scan is compared against most recent previous scan for each competitor
- [ ] **CHNG-02**: Detected changes include: price changes (magnitude + direction), new tiers, removed tiers, model type changes
- [ ] **CHNG-03**: Change detection produces a plain-English narrative of what changed
- [ ] **CHNG-04**: First scan for a competitor shows "baseline established" instead of false change alerts

### Pricing Architect

- [ ] **ARCH-01**: Internal strategy documents are stored in Aerospike for retrieval
- [ ] **ARCH-02**: When a pricing change is detected, the Pricing Architect retrieves relevant strategy docs
- [ ] **ARCH-03**: TrueFoundry routes the Pricing Architect to an expensive model (Opus/GPT-4) for strategic analysis
- [ ] **ARCH-04**: Pricing Architect generates a structured recommendation referencing specific strategy goals

### Similarity Clustering

- [ ] **CLST-01**: Normalized pricing profiles are embedded as vectors and stored in Aerospike
- [ ] **CLST-02**: "Who prices like us?" query returns competitors ranked by pricing similarity
- [ ] **CLST-03**: Cluster shifts are visible when a competitor changes pricing

### Dashboard

- [ ] **DASH-01**: Streamlit dashboard displays normalized comparison table with all competitors
- [ ] **DASH-02**: Transparency map shows count of public vs hidden ("Contact Sales") pricing
- [ ] **DASH-03**: Change alerts display with magnitude and direction when changes are detected
- [ ] **DASH-04**: Market positioning cluster visual shows pricing similarity groupings
- [ ] **DASH-05**: Pricing Architect recommendation panel displays strategic response text
- [ ] **DASH-06**: Dashboard reads from Ghost DB (no live LLM calls during page render)

### Sponsor Integration

- [ ] **SPNS-01**: Ghost is used for all persistent data storage with autonomous DB creation
- [ ] **SPNS-02**: TrueFoundry routes cheap model for extraction and expensive model for strategy
- [ ] **SPNS-03**: Aerospike stores pricing embeddings for similarity search and strategy docs for RAG
- [ ] **SPNS-04**: Macroscope is connected for code review during the build process

### Demo

- [ ] **DEMO-01**: Full scan-to-dashboard pipeline completes in under 60 seconds
- [ ] **DEMO-02**: Live demo runs under 3 minutes without breaking
- [ ] **DEMO-03**: Change detection demo moment works reliably (localhost mock fallback available)

## v2 Requirements

Deferred to post-hackathon. Tracked but not in current roadmap.

### Monitoring

- **MNTR-01**: Scheduled recurring scans via cron
- **MNTR-02**: Slack/webhook notifications on change detection
- **MNTR-03**: Email alerts for pricing changes

### Export

- **XPRT-01**: Normalization assumption export to CSV
- **XPRT-02**: Board-ready PDF export of comparison data

### Scaling

- **SCAL-01**: Playwright/browser rendering for JS-heavy pages
- **SCAL-02**: Support for arbitrary competitor URLs beyond pre-selected set
- **SCAL-03**: Competitor auto-discovery

## Out of Scope

| Feature | Reason |
|---------|--------|
| Playwright/browser rendering | Build risk too high for 8-hour hackathon; httpx + pre-selected pages is reliable |
| User authentication | Single-user demo tool; auth doubles scope |
| Real-time WebSocket updates | Polling sufficient; WebSocket adds infrastructure complexity |
| Automated repricing | Requires product API integration; scope is 8x larger |
| Mobile app | Web-only via Streamlit |
| Dynamic pricing rules engine | B2B SaaS changes quarterly; rules engine solves e-commerce |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| EXTR-01 | — | Pending |
| EXTR-02 | — | Pending |
| EXTR-03 | — | Pending |
| EXTR-04 | — | Pending |
| EXTR-05 | — | Pending |
| STOR-01 | — | Pending |
| STOR-02 | — | Pending |
| STOR-03 | — | Pending |
| NORM-01 | — | Pending |
| NORM-02 | — | Pending |
| NORM-03 | — | Pending |
| NORM-04 | — | Pending |
| CHNG-01 | — | Pending |
| CHNG-02 | — | Pending |
| CHNG-03 | — | Pending |
| CHNG-04 | — | Pending |
| ARCH-01 | — | Pending |
| ARCH-02 | — | Pending |
| ARCH-03 | — | Pending |
| ARCH-04 | — | Pending |
| CLST-01 | — | Pending |
| CLST-02 | — | Pending |
| CLST-03 | — | Pending |
| DASH-01 | — | Pending |
| DASH-02 | — | Pending |
| DASH-03 | — | Pending |
| DASH-04 | — | Pending |
| DASH-05 | — | Pending |
| DASH-06 | — | Pending |
| SPNS-01 | — | Pending |
| SPNS-02 | — | Pending |
| SPNS-03 | — | Pending |
| SPNS-04 | — | Pending |
| DEMO-01 | — | Pending |
| DEMO-02 | — | Pending |
| DEMO-03 | — | Pending |

**Coverage:**
- v1 requirements: 36 total
- Mapped to phases: 0
- Unmapped: 36 ⚠️

---
*Requirements defined: 2026-03-27*
*Last updated: 2026-03-27 after initial definition*
