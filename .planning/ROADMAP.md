# Roadmap: Pricing Radar

## Overview

Four phases from zero to demo-ready competitive pricing intelligence. Phase 1 lays the data foundation (scaffolding, scraper, Ghost DB schema). Phase 2 builds the extraction and normalization pipeline that all downstream features depend on. Phase 3 adds the intelligence layer: change detection and Aerospike vector similarity. Phase 4 wires in the Pricing Architect and builds the Streamlit dashboard with demo hardening. The demo is the product — every phase is aimed at a clean 3-minute live run.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation** - Project scaffolding, Ghost DB schema, httpx scraper verified against all target pages
- [ ] **Phase 2: Extraction Pipeline** - LLM extraction and cross-model normalization with transparent assumptions
- [ ] **Phase 3: Intelligence Layer** - Change detection and Aerospike vector similarity clustering
- [ ] **Phase 4: Architect & Dashboard** - Pricing Architect RAG, Streamlit dashboard, demo hardening

## Phase Details

### Phase 1: Foundation
**Goal**: The project runs, Ghost DB exists with the right schema, and the scraper successfully fetches all 5-8 target pages
**Depends on**: Nothing (first phase)
**Requirements**: EXTR-01, EXTR-02, EXTR-05, STOR-01, STOR-02, STOR-03, SPNS-01, SPNS-04
**Success Criteria** (what must be TRUE):
  1. Proto contract (`pricing.proto`) is defined, `buf generate` produces Go code, and Connect-Go server boots on `:8080`
  2. Go server connects to Ghost Postgres via pgx and creates the DB and tables autonomously without error
  3. All 5-8 target pricing pages return non-empty HTML via Go net/http; pages that fail fall back to cached local HTML and are flagged as "cached"
  4. Raw HTML is stripped of nav, footer, scripts, and CSS via goquery before any further processing
  5. Ghost DB schema supports scan_run_id, competitor key, and scraped_at timestamp so change detection queries will work
**Plans**: TBD

### Phase 2: Extraction Pipeline
**Goal**: A full scrape-to-store cycle works end-to-end for a single competitor, producing structured normalized pricing data with explicit assumption footnotes
**Depends on**: Phase 1
**Requirements**: EXTR-03, EXTR-04, NORM-01, NORM-02, NORM-03, NORM-04, SPNS-02
**Success Criteria** (what must be TRUE):
  1. TrueFoundry cheap model (via OpenAI-compatible endpoint) extracts structured pricing into typed Go structs: company, tiers, prices, features, limits, pricing model type
  2. "Contact Sales" tiers are flagged as opaque with a confidence score rather than silently dropped
  3. All heterogeneous pricing models (per-seat, per-token, credits, hybrid) are normalized to a common monthly cost using an explicit reference workload
  4. Every normalized figure in the output carries its normalization_assumptions dict; pricing models that cannot be normalized display raw pricing with a flag
**Plans**: TBD

### Phase 3: Intelligence Layer
**Goal**: The system detects pricing changes between scans and clusters competitors by pricing similarity
**Depends on**: Phase 2
**Requirements**: CHNG-01, CHNG-02, CHNG-03, CHNG-04, CLST-01, CLST-02, CLST-03, SPNS-03
**Success Criteria** (what must be TRUE):
  1. Running two scans in sequence produces a plain-English change narrative for any competitor whose pricing changed (price magnitude, direction, tier additions/removals, model type shifts)
  2. The first scan for a new competitor shows "baseline established" rather than false change alerts
  3. "Who prices like us?" query returns competitors ranked by pricing similarity using Aerospike HNSW vector search
  4. When a competitor changes pricing, its cluster position shifts are visible in the similarity rankings
**Plans**: TBD

### Phase 4: Architect & Dashboard
**Goal**: Judges can watch a 3-minute demo that shows a live scan, a detected pricing change, a strategy recommendation, and a dashboard with all visualizations
**Depends on**: Phase 3
**Requirements**: ARCH-01, ARCH-02, ARCH-03, ARCH-04, DASH-01, DASH-02, DASH-03, DASH-04, DASH-05, DASH-06, DEMO-01, DEMO-02, DEMO-03
**Success Criteria** (what must be TRUE):
  1. When a pricing change is detected, the Pricing Architect retrieves relevant strategy docs from Aerospike and generates a structured recommendation via TrueFoundry expensive model (Opus/GPT-4)
  2. Streamlit dashboard shows: normalized comparison table, transparency map (public vs "Contact Sales" counts), change alert banners with magnitude and direction, cluster visual, and Pricing Architect recommendation panel
  3. Dashboard reads exclusively from Ghost DB — no live LLM calls happen during page render
  4. The full scan-to-dashboard pipeline completes in under 60 seconds
  5. The demo change detection moment works reliably via the localhost mock fallback (seeded "before" and "after" pricing pages)
**Plans**: TBD
**UI hint**: yes

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 0/TBD | Not started | - |
| 2. Extraction Pipeline | 0/TBD | Not started | - |
| 3. Intelligence Layer | 0/TBD | Not started | - |
| 4. Architect & Dashboard | 0/TBD | Not started | - |
