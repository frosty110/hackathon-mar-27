# Pricing Radar

## What This Is

An autonomous agent that scrapes competitor pricing pages, extracts structured data via LLM, normalizes across heterogeneous pricing models (per-seat, per-token, credits, etc.), detects changes, and generates strategy-grounded responses using internal positioning docs. Built as a hackathon project with a 3-minute live demo targeting judges evaluating autonomy, idea quality, technical implementation, tool use, and presentation.

## Core Value

Continuous competitive pricing intelligence that turns a 2-day manual spreadsheet into a 38-second automated scan with strategic recommendations.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Extraction pipeline scrapes 5-8 pre-selected AI company pricing pages via httpx + LLM
- [ ] Normalization engine converts heterogeneous pricing models to a common comparison unit with explicit assumptions
- [ ] Change detection compares current scan against previous scan stored in Ghost DB
- [ ] Pricing Architect generates strategic responses grounded in internal strategy docs from Aerospike
- [ ] "Who prices like us?" vector similarity clustering via Aerospike
- [ ] Streamlit dashboard displays comparison table, transparency map, change alerts, and cluster visual
- [ ] Ghost Postgres DB stores all scan data with autonomous DB creation
- [ ] TrueFoundry routes cheap model for extraction, expensive model for strategic analysis
- [ ] Macroscope integrated for code review during build (4th sponsor)
- [ ] Full demo runs under 3 minutes without breaking

### Out of Scope

- Playwright/browser rendering — build risk too high, httpx + LLM only
- Arbitrary URL support — MVP limited to pre-selected pages that work with httpx
- Real-time WebSocket updates — polling-based is sufficient for demo
- Mobile app — web-only via Streamlit
- User authentication — single-user demo tool
- Overmind integration — cut in favor of TrueFoundry for architectural depth

## Context

**Hackathon constraints:** 8-hour build, 3-minute demo, 3+ sponsor tools required. Judged on autonomy, idea, technical implementation, tool use, presentation.

**Sponsor stack:**
| Sponsor | Role |
|---------|------|
| Ghost | Postgres DB — stores scan data, autonomous DB creation, fork for historical snapshots |
| TrueFoundry | Model routing — cheap model (Haiku/Flash) for extraction, expensive model (Opus/GPT-4) for strategy |
| Aerospike | Vector DB — pricing profile embeddings for similarity search + strategy doc retrieval |
| Macroscope | Code review — reviews PRs during the build process |

**Demo story:** Macroscope (a sponsor) is transitioning to usage-based pricing. Pricing Radar scans their competitors, normalizes pricing, and if we catch the change live, the Pricing Architect generates a strategic response. Fallback: seed DB with current pricing and mock the "after" page with usage-based model.

**Target pages:** 5-8 AI company pricing pages including Macroscope competitors. Pages must work with httpx (no heavy JS rendering).

**Demand evidence:** Carmen Insignares Newell analyzed 125 AI startup pricing pages manually. Koshima Satija spent 3 months on 200+ Voice AI companies. Pricing consultants do this 15x/year. Companies with 100%+ growth change pricing 3x more frequently.

## Constraints

- **Timeline**: 8 hours total build time (hackathon)
- **Demo**: 3-minute live demo, must not break
- **Sponsors**: Must use 3+ sponsor tools in load-bearing way (Ghost, TrueFoundry, Aerospike + Macroscope for code review)
- **Tech stack**: Python + Streamlit for dashboard, httpx for scraping (no Playwright)
- **Reliability**: Pre-selected pages + localhost mock fallback for demo stability

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Approach B: Ghost + Aerospike + TrueFoundry | Architecturally strongest sponsor story. TrueFoundry model routing is meaningful, not cosmetic | — Pending |
| httpx only, no Playwright | Build risk too high for headless browser in 8 hours | — Pending |
| Pre-selected pages for demo | Eliminates external dependency risk during live demo | — Pending |
| Macroscope as 4th sponsor (code review) | Fits naturally into dev workflow without forcing an unnatural integration | — Pending |
| Macroscope competitors as target pages | Creates compelling narrative — analyzing a sponsor's competitive landscape | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-27 after initialization*
