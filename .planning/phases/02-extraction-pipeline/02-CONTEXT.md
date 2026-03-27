# Phase 2: Extraction Pipeline - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

A full scrape-to-store cycle works end-to-end, producing structured normalized pricing data with explicit assumption footnotes. TrueFoundry cheap model extracts structured pricing into typed Go structs. "Contact Sales" tiers are flagged as opaque. All heterogeneous pricing models normalize to a common monthly cost. Every normalized figure carries its assumptions.

</domain>

<decisions>
## Implementation Decisions

### LLM Extraction
- Use the extraction schema from the design doc: company, pricing_model (per-seat|per-token|per-minute|credits|hybrid|other), tiers[] (name, base_price, unit, features, limits), normalized_monthly_cost, normalization_assumptions, confidence (0.0-1.0), pricing_public, last_scanned
- Use Claude Haiku via TrueFoundry OpenAI-compatible endpoint as the cheap extraction model
- System prompt with JSON schema definition + 1 few-shot example per pricing model type (per-seat example, per-token example). Balance reliability vs token cost.
- TrueFoundry model routing: set two model name env vars in config.go — TRUEFOUNDRY_CHEAP_MODEL for extraction, TRUEFOUNDRY_EXPENSIVE_MODEL for strategy (Phase 4)

### Normalization
- Reference workload: "50-person team, 1M tokens/month" — covers both seat-based and token-based models
- Un-normalizable pricing (e.g., pure outcome-based): display raw pricing with flag "Cannot normalize: [reason]"
- Confidence scoring: 0.0-1.0 float. Public page with full tier data = 0.85+. "Contact Sales" = 0.1-0.3. Missing fields reduce score proportionally.
- "Contact Sales" tiers: flagged as opaque with low confidence, NOT silently dropped

### Claude's Discretion
- Exact prompt engineering for the extraction system prompt
- How to structure the few-shot examples
- Whether to use JSON mode or function calling for structured output
- Go struct field naming and JSON tags
- Error handling for malformed LLM responses

</decisions>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/scraper/` — fetcher.go returns `[]FetchResult{URL, HTML, Cached, Error}`, parser.go has `StripBoilerplate`
- `internal/storage/ghost.go` — `SaveSnapshot(ctx, scanRunID, competitor, rawHTML, strippedHTML)` already stores raw + stripped HTML
- `internal/config/config.go` — loads env vars from .env via godotenv. Add TRUEFOUNDRY_BASE_URL, TRUEFOUNDRY_API_KEY, TRUEFOUNDRY_CHEAP_MODEL, TRUEFOUNDRY_EXPENSIVE_MODEL
- `internal/handler/pricing.go` — RunScan already fetches + strips + stores. Extraction + normalization slot in between strip and store.

### Established Patterns
- Go backend, all business logic in internal/
- Connect-Go handler calls internal packages
- Config via env vars loaded in config.go
- TrueFoundry uses OpenAI-compatible API — raw net/http POST to /v1/chat/completions

### Integration Points
- Extractor reads stripped HTML from scraper output, returns typed Go struct
- Normalizer takes raw extraction output, applies reference workload, returns normalized struct with assumptions
- Storage.SaveSnapshot needs to be extended to store extracted + normalized data (or new table)

</code_context>

<specifics>
## Specific Ideas

- The normalization schema from the design doc (docs/design.md lines 99-120) is the reference implementation
- Worked example from design doc: "Competitor A charges $49/seat/month. Competitor B charges $0.003/token. For a 50-person team using 1M tokens/month: A = $2,450/mo, B = $3,000/mo."

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>
