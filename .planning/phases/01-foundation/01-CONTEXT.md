# Phase 1: Foundation - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning
**Mode:** Auto-generated (infrastructure phase — discuss skipped)

<domain>
## Phase Boundary

The project runs, Ghost DB exists with the right schema, and the scraper successfully fetches all 5-8 target pages. Specifically: Proto contract (pricing.proto) is defined, buf generate produces Go code, Connect-Go server boots on :8080. Go server connects to Ghost Postgres via pgx and creates DB/tables autonomously. All target pricing pages return non-empty HTML via Go net/http; failures fall back to cached local HTML. Raw HTML is stripped of nav/footer/scripts/CSS via goquery. Ghost DB schema supports scan_run_id, competitor key, and scraped_at timestamp.

</domain>

<decisions>
## Implementation Decisions

### Claude's Discretion
All implementation choices are at Claude's discretion — pure infrastructure phase. Use ROADMAP phase goal, success criteria, and the architecture defined in .planning/research/ARCHITECTURE.md to guide decisions. Follow the recommended project structure from ARCHITECTURE.md (proto/, cmd/, internal/, gen/, frontend/).

</decisions>

<code_context>
## Existing Code Insights

### Reusable Assets
- No existing code — greenfield project

### Established Patterns
- Architecture defined in .planning/research/ARCHITECTURE.md: Go backend (Connect-Go API, pgx, goquery) + Streamlit frontend (display only)
- Project structure: proto/ → cmd/server/ → internal/{handler,scraper,storage,config}/ → gen/ → frontend/
- Contract-first development: .proto files are the single source of truth

### Integration Points
- Ghost Postgres: autonomous DB creation via Ghost API, pgx driver for queries
- TrueFoundry: OpenAI-compatible REST API (not needed in Phase 1, but config should be ready)
- Target pricing pages: net/http with conservative timeouts, pre-validated URLs

</code_context>

<specifics>
## Specific Ideas

No specific requirements — infrastructure phase. Refer to ROADMAP phase description and success criteria.

</specifics>

<deferred>
## Deferred Ideas

None — infrastructure phase.

</deferred>
