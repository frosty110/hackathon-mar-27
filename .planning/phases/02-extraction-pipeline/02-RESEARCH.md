# Phase 2: Extraction Pipeline - Research

**Researched:** 2026-03-27
**Domain:** LLM extraction + normalization via TrueFoundry OpenAI-compatible API in Go
**Confidence:** HIGH (TrueFoundry endpoint confirmed from .env, Go patterns standard)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**LLM Extraction**
- Use the extraction schema from the design doc: company, pricing_model (per-seat|per-token|per-minute|credits|hybrid|other), tiers[] (name, base_price, unit, features, limits), normalized_monthly_cost, normalization_assumptions, confidence (0.0-1.0), pricing_public, last_scanned
- Use Claude Haiku via TrueFoundry OpenAI-compatible endpoint as the cheap extraction model
- System prompt with JSON schema definition + 1 few-shot example per pricing model type (per-seat example, per-token example). Balance reliability vs token cost.
- TrueFoundry model routing: set two model name env vars in config.go — TRUEFOUNDRY_CHEAP_MODEL for extraction, TRUEFOUNDRY_EXPENSIVE_MODEL for strategy (Phase 4)

**Normalization**
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

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| EXTR-03 | LLM extracts structured pricing data (company, tiers, prices, features, limits, pricing model type) into typed Go structs | TrueFoundry endpoint + Go JSON pattern documented below |
| EXTR-04 | Extraction handles "Contact Sales" tiers by flagging as opaque with confidence score | Schema includes `pricing_public` bool + confidence float; "contact_sales" tier pattern documented |
| NORM-01 | Heterogeneous pricing models (per-seat, per-token, credits, hybrid) normalized to common monthly cost | Normalization algorithm for all 6 model types documented below |
| NORM-02 | Normalization uses explicit reference workload ("50-person team, 1M tokens/month") | Reference workload is a locked decision; normalizer takes it as a constant |
| NORM-03 | Normalization assumptions displayed inline with every normalized figure | `normalization_assumptions` string field carries this; stored in DB alongside cost |
| NORM-04 | Pricing models that cannot be meaningfully normalized display raw pricing with flag | `CannotNormalize` return path documented; "Cannot normalize: [reason]" string pattern |
| SPNS-02 | TrueFoundry routes cheap model for extraction and expensive model for strategy | TRUEFOUNDRY_CHEAP_MODEL already in config.go; extractor uses it; two-model routing verified |
</phase_requirements>

---

## Summary

Phase 2 builds `internal/extractor/llm.go` and `internal/normalizer/normalize.go`, then wires them into `internal/handler/pricing.go` between the existing strip and store steps. The extractor POSTs stripped HTML to the TrueFoundry OpenAI-compatible endpoint using raw `net/http`, receives a JSON response, and unmarshals it into a typed `RawPricingData` struct. The normalizer applies a 50-seat / 1M-token reference workload to produce a `NormalizedPricing` struct with an explicit assumptions string. Storage needs one new table (`extracted_pricing`) with columns for all extracted + normalized fields. The proto contract needs new response messages so the Streamlit frontend can eventually display this data.

**Primary recommendation:** Use `response_format: {"type": "json_object"}` in the TrueFoundry API call (JSON mode) with a detailed system prompt + 2 few-shot examples. This is simpler than function calling, avoids schema JSON complexity, and gpt-4o-mini supports it reliably. Parse `choices[0].message.content` with `json.Unmarshal` into the typed struct.

**Critical discovery:** The `.env` already has TrueFoundry credentials configured with real values:
- Base URL: `https://llm-gateway.truefoundry.com/api/inference/openai`
- Cheap model: `openai-main/gpt-4o-mini`
- Expensive model: `openai-main/gpt-4o`
- Auth: `Authorization: Bearer <TRUEFOUNDRY_API_KEY>`

The extractor and normalizer directories (`internal/extractor/`, `internal/normalizer/`) are empty — both packages must be created from scratch.

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `net/http` (stdlib) | Go 1.26 | POST to TrueFoundry `/v1/chat/completions` | Already used for scraping; zero new deps |
| `encoding/json` (stdlib) | Go 1.26 | Marshal request body, unmarshal LLM response + extracted JSON | Zero dep; sufficient for the schema |
| `pgx/v5` | v5.9.1 (in go.mod) | Extend `pricing_snapshots` table with new columns | Already in project |
| `slog` (stdlib) | Go 1.26 | Log extraction confidence, fallback events | Already used in handler |

### No New Dependencies Required
The entire extraction pipeline is implementable with the packages already in `go.mod`. No new `go get` needed.

**Installation:** None — all packages already present.

---

## Architecture Patterns

### Recommended Package Structure

```
internal/
├── extractor/
│   └── llm.go              # TrueFoundry API client + JSON parsing
├── normalizer/
│   ├── normalize.go        # Heterogeneous → common monthly cost
│   └── models.go           # Shared Go structs: RawPricingData, NormalizedPricing, Tier
└── storage/
    └── ghost.go            # Extend with SaveExtractedPricing() method
```

### Shared Structs (internal/normalizer/models.go)

Define structs here rather than in extractor to avoid circular imports. Both extractor and storage import from normalizer.

```go
// Source: design.md lines 99-120 (locked schema)
package normalizer

type Tier struct {
    Name      string            `json:"name"`
    BasePrice float64           `json:"base_price"`
    Unit      string            `json:"unit"`       // "seat/month", "per 1M tokens", etc.
    Features  []string          `json:"features"`
    Limits    map[string]string `json:"limits"`
    IsOpaque  bool              `json:"is_opaque"`  // true when "Contact Sales"
}

type RawPricingData struct {
    Company                  string  `json:"company"`
    URL                      string  `json:"url"`
    PricingModel             string  `json:"pricing_model"` // per-seat|per-token|per-minute|credits|hybrid|other
    Tiers                    []Tier  `json:"tiers"`
    NormalizedMonthlyCost    float64 `json:"normalized_monthly_cost"`
    NormalizationAssumptions string  `json:"normalization_assumptions"`
    Confidence               float64 `json:"confidence"`
    PricingPublic            bool    `json:"pricing_public"`
    CannotNormalize          bool    `json:"cannot_normalize"`
    CannotNormalizeReason    string  `json:"cannot_normalize_reason"`
    LastScanned              string  `json:"last_scanned"` // RFC3339
}
```

### Pattern 1: TrueFoundry HTTP Client (internal/extractor/llm.go)

**What:** Raw `net/http` POST to `{TRUEFOUNDRY_BASE_URL}/v1/chat/completions` with `Authorization: Bearer {TRUEFOUNDRY_API_KEY}`. Use `response_format: {"type": "json_object"}` for reliable JSON output.

**Request struct pattern:**

```go
// internal/extractor/llm.go
type chatRequest struct {
    Model          string        `json:"model"`
    Messages       []chatMessage `json:"messages"`
    ResponseFormat responseFormat `json:"response_format"`
    Temperature    float64       `json:"temperature"`
}

type chatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type responseFormat struct {
    Type string `json:"type"` // "json_object"
}

// Response envelope — only what we need
type chatResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error *struct {
        Message string `json:"message"`
    } `json:"error"`
}
```

**Request execution pattern:**

```go
func (e *Extractor) Extract(ctx context.Context, competitor, strippedHTML string) (*normalizer.RawPricingData, error) {
    payload := chatRequest{
        Model: e.cfg.TrueFoundryCheapModel, // "openai-main/gpt-4o-mini"
        Messages: []chatMessage{
            {Role: "system", Content: systemPrompt},
            {Role: "user", Content: "Extract pricing from this HTML:\n\n" + truncate(strippedHTML, 12000)},
        },
        ResponseFormat: responseFormat{Type: "json_object"},
        Temperature:    0.0, // deterministic for extraction
    }
    body, _ := json.Marshal(payload)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost,
        e.cfg.TrueFoundryBaseURL+"/v1/chat/completions",
        bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+e.cfg.TrueFoundryAPIKey)

    resp, err := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    var chatResp chatResponse
    json.NewDecoder(resp.Body).Decode(&chatResp)

    // The content IS the JSON — unmarshal into typed struct
    var result normalizer.RawPricingData
    if err := json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &result); err != nil {
        return nil, fmt.Errorf("malformed LLM JSON for %s: %w", competitor, err)
    }
    return &result, nil
}
```

**HTML truncation:** Stripped HTML can still be 20-50KB. Truncate to ~12,000 characters (roughly 3,000 tokens) before sending. Pricing data is almost always in the first half of the page.

### Pattern 2: Normalization Algorithm (internal/normalizer/normalize.go)

**Reference workload constants:**

```go
const (
    RefSeats          = 50
    RefTokensPerMonth = 1_000_000
    RefMinutesPerMonth = 10_000 // reasonable voice/video usage
)
```

**Per-model normalization logic:**

| Model Type | Algorithm | Example |
|------------|-----------|---------|
| `per-seat` | `cheapest_qualifying_tier.base_price * RefSeats` | $49/seat × 50 = $2,450/mo |
| `per-token` | `price_per_token * RefTokensPerMonth` | $0.003/1K tokens × 1000 = $3,000/mo |
| `per-minute` | `price_per_minute * RefMinutesPerMonth` | $0.01/min × 10,000 = $100/mo |
| `credits` | `cost_per_credit * credits_needed_for_workload` | Requires credit-to-unit mapping from tiers |
| `hybrid` | `base_fee + (usage_component * ref_workload)` | $200 base + per-token rate |
| `other` / outcome-based | `cannot_normalize = true` | Display raw pricing + flag |

**Tier selection for per-seat:** Find the lowest-priced tier that accommodates 50 seats (some tiers have minimum seat counts). If all tiers have "Contact Sales" pricing, `pricing_public = false`, confidence = 0.1-0.3.

**Confidence scoring rules:**

```
0.85-1.0 : Public pricing, full tier data, prices are numeric
0.5-0.84 : Public pricing but some gaps (missing features, partial limits)
0.1-0.3  : "Contact Sales" for all commercial tiers
0.0      : No pricing data found at all (page fetch failed or pure redirect)
```

For each missing required field, deduct 0.1 from confidence.

### Pattern 3: Storage Extension

**Approach:** Add new columns to `pricing_snapshots` rather than a new table. Keeps storage.SaveSnapshot as the single write path; Phase 3 (change detection) can JOIN or filter in one query.

**Migration SQL (add to AutoMigrate):**

```sql
ALTER TABLE pricing_snapshots
    ADD COLUMN IF NOT EXISTS pricing_model         TEXT,
    ADD COLUMN IF NOT EXISTS tiers_json            JSONB,
    ADD COLUMN IF NOT EXISTS normalized_monthly_cost NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS normalization_assumptions TEXT,
    ADD COLUMN IF NOT EXISTS confidence            NUMERIC(4,3),
    ADD COLUMN IF NOT EXISTS pricing_public        BOOLEAN,
    ADD COLUMN IF NOT EXISTS cannot_normalize      BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS cannot_normalize_reason TEXT,
    ADD COLUMN IF NOT EXISTS extraction_error      TEXT;
```

**Why `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`:** Safe to run on every boot (idempotent). PostgreSQL acquires only a brief ACCESS EXCLUSIVE lock for metadata-only column addition with nullable columns — no table rewrite. Already demonstrated by `CREATE TABLE IF NOT EXISTS` in the existing `AutoMigrate`.

**New storage method signature:**

```go
func (db *GhostDB) SaveExtractedPricing(ctx context.Context, snapshotID int64, data *normalizer.RawPricingData) error
```

Or, alternatively, extend `SaveSnapshot` to accept an optional `*normalizer.RawPricingData` and update the row in a second query. The two-query approach is cleaner: insert the snapshot first (Phase 1 pattern preserved), then update extracted fields once extraction completes.

### Pattern 4: Handler Wiring (internal/handler/pricing.go)

Extraction slots between strip and store. Run concurrent extraction with errgroup (same pattern as scraping):

```go
// After stripping, before SaveSnapshot — concurrent extraction
type extractResult struct {
    competitor string
    data       *normalizer.RawPricingData
    err        error
}
extractResults := make([]extractResult, len(rawPages))
var g errgroup.Group
for i, page := range rawPages {
    i, page := i, page
    stripped := strippedPages[i]
    g.Go(func() error {
        raw, err := h.extractor.Extract(ctx, page.Competitor, stripped)
        if err != nil {
            slog.Warn("extraction failed", "competitor", page.Competitor, "error", err)
            extractResults[i] = extractResult{competitor: page.Competitor, err: err}
            return nil // non-fatal
        }
        normalized := h.normalizer.Normalize(raw)
        extractResults[i] = extractResult{competitor: page.Competitor, data: normalized}
        return nil
    })
}
_ = g.Wait()
```

### Anti-Patterns to Avoid

- **Sending raw (un-stripped) HTML to the LLM:** Token cost balloons 3-5x; nav/footer noise confuses extraction. Always use `scraper.StripBoilerplate` output.
- **Not truncating stripped HTML:** Even stripped, some pages are 30KB+. Truncate to 12,000 chars. Pricing data loads first; tail is usually FAQs.
- **Silently dropping "Contact Sales" tiers:** Required by EXTR-04. Set `is_opaque: true` on the tier, `pricing_public: false` on the company, confidence 0.1-0.3.
- **Parsing JSON with regex or string scanning:** Use `json.Unmarshal` into the typed struct. If it fails, record the error and continue — don't crash the scan.
- **Using a new table instead of extending `pricing_snapshots`:** Adds a JOIN for Phase 3 change detection for no benefit. Single table, nullable columns for extraction fields.
- **Using function calling instead of json_object mode:** Function calling requires a JSON schema definition object; json_object mode is simpler, and gpt-4o-mini handles it well with a clear system prompt.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Reliable JSON output from LLM | Custom parser / regex | `response_format: {"type": "json_object"}` + `json.Unmarshal` | JSON mode guarantees valid JSON envelope; struct unmarshaling catches schema drift |
| OpenAI-compatible HTTP client | Custom retry logic, SDK | Raw `net/http` + `encoding/json` | TrueFoundry is plain REST; existing project pattern; zero new deps |
| DB schema migration | golang-migrate, goose | `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` in AutoMigrate | Already using `CREATE TABLE IF NOT EXISTS` idiom; overkill for 4-table schema |
| HTML token counting | tiktoken Go port | Hard truncation at 12,000 chars | Accurate token counting adds a dep; truncation is safe and fast |
| Credit model normalization | Complex credit exchange | Per-tier credit-to-unit annotation in the prompt | LLM extracts the conversion; normalizer applies arithmetic |

---

## Prompt Engineering

### System Prompt Architecture

The system prompt must do three things: define the output schema, define the reference workload, and give 2 few-shot examples (per-seat and per-token). The LLM never sees the Go structs — it sees a JSON schema in the prompt.

**Recommended system prompt structure:**

```
You are a pricing data extraction specialist. Extract structured pricing data from HTML pricing pages.

Always respond with ONLY valid JSON matching this schema — no markdown, no explanation:
{
  "company": string,
  "url": string,
  "pricing_model": "per-seat" | "per-token" | "per-minute" | "credits" | "hybrid" | "other",
  "tiers": [
    {
      "name": string,
      "base_price": number | null,
      "unit": string,
      "features": [string],
      "limits": {key: value},
      "is_opaque": boolean
    }
  ],
  "normalized_monthly_cost": number | null,
  "normalization_assumptions": string,
  "confidence": number (0.0-1.0),
  "pricing_public": boolean,
  "cannot_normalize": boolean,
  "cannot_normalize_reason": string | null,
  "last_scanned": string (ISO8601)
}

Reference workload for normalization: 50-person team, 1M tokens/month.

Rules:
- Set is_opaque=true and base_price=null for "Contact Sales" tiers
- Set pricing_public=false if ALL commercial tiers are "Contact Sales"
- Confidence: 0.85+ for full public pricing, 0.1-0.3 for Contact Sales only
- Set cannot_normalize=true for pure outcome-based pricing (e.g., "pay per result")
- Always include the normalization_assumptions string even when cannot_normalize=true

--- EXAMPLE 1 (per-seat) ---
Input: [HTML showing Pro tier $49/seat/month, Enterprise "Contact Sales"]
Output: {
  "company": "Acme AI",
  "pricing_model": "per-seat",
  "tiers": [
    {"name": "Pro", "base_price": 49.00, "unit": "seat/month", "is_opaque": false, ...},
    {"name": "Enterprise", "base_price": null, "unit": "seat/month", "is_opaque": true, ...}
  ],
  "normalized_monthly_cost": 2450.00,
  "normalization_assumptions": "50 seats at $49/seat/month. Enterprise tier excluded (Contact Sales).",
  "confidence": 0.75,
  "pricing_public": true,
  ...
}

--- EXAMPLE 2 (per-token) ---
Input: [HTML showing $3.00 per 1M input tokens, $15.00 per 1M output tokens]
Output: {
  "company": "LLM Corp",
  "pricing_model": "per-token",
  "tiers": [
    {"name": "Input", "base_price": 3.00, "unit": "per 1M tokens", "is_opaque": false, ...},
    {"name": "Output", "base_price": 15.00, "unit": "per 1M tokens", "is_opaque": false, ...}
  ],
  "normalized_monthly_cost": 3000.00,
  "normalization_assumptions": "1M tokens/month at $3.00/1M input tokens. Output cost excluded (assumes read-heavy workload).",
  "confidence": 0.90,
  "pricing_public": true,
  ...
}
```

**Why 2 few-shot examples vs more:** Token cost on gpt-4o-mini is $0.15/1M input tokens. Each example adds ~200 tokens. 2 examples cover the most common pricing models (seat + token-based AI tools). More examples push total prompt toward 2,000+ tokens per call × 8 competitors = expensive for a hackathon.

### JSON Mode vs Function Calling Decision

**Use JSON mode** (`response_format: {"type": "json_object"}`).

- Function calling requires a `tools` array with a full JSON Schema definition — more code, more prompt tokens, more failure modes.
- JSON mode guarantees the outer envelope is valid JSON. Schema adherence comes from the system prompt + examples.
- gpt-4o-mini handles JSON mode reliably when the system prompt includes the schema.
- Simpler error path: if `json.Unmarshal` fails, the content string is logged for debugging.

---

## Common Pitfalls

### Pitfall 1: LLM Returns JSON Wrapped in Markdown Code Fences

**What goes wrong:** Response is ` ```json\n{...}\n``` ` instead of bare `{...}`.
**Why it happens:** Models trained on markdown-heavy data default to fenced code blocks even in JSON mode. gpt-4o-mini is better than older models but not immune.
**How to avoid:** Set `response_format: {"type": "json_object"}` — this forces bare JSON output in compliant implementations. Add `"Return ONLY raw JSON, no markdown"` to the system prompt as a belt-and-suspenders measure.
**Warning signs:** `json.Unmarshal` error containing "invalid character '`'"; log `content[:100]` on error.
**Recovery:** If unmarshal fails, try stripping ```json ... ``` wrapper with `strings.TrimPrefix` / `strings.TrimSuffix` before returning error.

### Pitfall 2: Confidence Field Outside 0.0-1.0 Range

**What goes wrong:** LLM returns `confidence: 85` (treating it as a percentage) instead of `0.85`.
**Why it happens:** Model interprets "0.0-1.0 float" ambiguously.
**How to avoid:** Explicitly state in prompt: "confidence as a decimal between 0.0 and 1.0, NOT a percentage."
**Recovery:** Post-process: if `data.Confidence > 1.0`, divide by 100.

### Pitfall 3: HTML Too Long Causes Context Window Errors

**What goes wrong:** HTTP 400 from TrueFoundry with "context_length_exceeded" error.
**Why it happens:** Some pricing pages (Snyk, enterprise tools) have very long HTML even after stripping.
**How to avoid:** Hard-truncate stripped HTML at 12,000 characters before building the LLM request. gpt-4o-mini context window is 128K tokens but the prompt itself adds ~800 tokens; 12K chars ≈ 3,000 tokens is safe.
**Warning signs:** HTTP 400 response from TrueFoundry with error JSON in body; always check `chatResp.Error != nil`.

### Pitfall 4: Normalization Divide-by-Zero on Missing Prices

**What goes wrong:** Normalizer panics because `base_price` is null/0 for all tiers.
**Why it happens:** "Contact Sales" only pages — all tiers have `is_opaque: true` and `base_price: 0`.
**How to avoid:** Check `pricing_public` before normalization. If false, set `cannot_normalize = true`, reason = "All pricing is Contact Sales".
**Warning signs:** `normalized_monthly_cost` of 0.0 combined with `pricing_public: false`.

### Pitfall 5: Per-Token Units Are Not Standardized

**What goes wrong:** One page says "$3.00 per 1M tokens", another says "$0.003 per 1K tokens", another says "$0.000003 per token". Normalization produces wildly different results.
**Why it happens:** LLM extracts the unit literally from the page.
**How to avoid:** Prompt explicitly: "Normalize token pricing to 'per 1M tokens' as the unit." Then the normalizer can assume all token prices are per-1M.
**Warning signs:** `unit` field containing "per token" or "per 1K tokens" instead of "per 1M tokens".

### Pitfall 6: Extractor and Normalizer in Same Package

**What goes wrong:** Handler imports extractor, extractor imports normalizer structs — circular import if normalizer also imports extractor.
**Why it happens:** Struct definitions placed in the wrong package.
**How to avoid:** Shared structs (`RawPricingData`, `Tier`) live in `internal/normalizer/models.go`. Extractor imports from normalizer (one direction only). Storage imports from normalizer. Handler imports all three.

---

## Code Examples

### Complete Extractor Call Pattern

```go
// internal/extractor/llm.go
// Source: design.md + .env confirmed TrueFoundry config

type Extractor struct {
    cfg    *config.Config
    client *http.Client
}

func New(cfg *config.Config) *Extractor {
    return &Extractor{
        cfg:    cfg,
        client: &http.Client{Timeout: 30 * time.Second},
    }
}

func truncate(s string, maxChars int) string {
    if len(s) <= maxChars {
        return s
    }
    return s[:maxChars]
}
```

### Normalization Dispatch

```go
// internal/normalizer/normalize.go
func Normalize(raw *RawPricingData) *RawPricingData {
    if !raw.PricingPublic {
        raw.CannotNormalize = true
        raw.CannotNormalizeReason = "All pricing is Contact Sales"
        return raw
    }
    switch raw.PricingModel {
    case "per-seat":
        normalizeSeat(raw)
    case "per-token":
        normalizeToken(raw)
    case "per-minute":
        normalizeMinute(raw)
    case "credits":
        normalizeCredits(raw)
    case "hybrid":
        normalizeHybrid(raw)
    default:
        raw.CannotNormalize = true
        raw.CannotNormalizeReason = fmt.Sprintf("Cannot normalize: %s pricing model", raw.PricingModel)
    }
    return raw
}

func normalizeSeat(raw *RawPricingData) {
    // Find cheapest non-opaque tier
    for _, t := range raw.Tiers {
        if !t.IsOpaque && t.BasePrice > 0 {
            raw.NormalizedMonthlyCost = t.BasePrice * RefSeats
            raw.NormalizationAssumptions = fmt.Sprintf(
                "%d seats at $%.2f/seat/month (%s tier).", RefSeats, t.BasePrice, t.Name)
            return
        }
    }
    raw.CannotNormalize = true
    raw.CannotNormalizeReason = "No public seat pricing found"
}
```

### Ghost DB Schema Extension

```go
// In AutoMigrate() — append to the existing exec block
_, err = db.pool.Exec(ctx, `
    ALTER TABLE pricing_snapshots
        ADD COLUMN IF NOT EXISTS pricing_model              TEXT,
        ADD COLUMN IF NOT EXISTS tiers_json                 JSONB,
        ADD COLUMN IF NOT EXISTS normalized_monthly_cost    NUMERIC(12,2),
        ADD COLUMN IF NOT EXISTS normalization_assumptions  TEXT,
        ADD COLUMN IF NOT EXISTS confidence                 NUMERIC(4,3),
        ADD COLUMN IF NOT EXISTS pricing_public             BOOLEAN,
        ADD COLUMN IF NOT EXISTS cannot_normalize           BOOLEAN DEFAULT FALSE,
        ADD COLUMN IF NOT EXISTS cannot_normalize_reason    TEXT,
        ADD COLUMN IF NOT EXISTS extraction_error           TEXT;
`)
```

### SaveExtractedPricing Update Method

```go
// internal/storage/ghost.go
func (db *GhostDB) UpdateExtractedPricing(
    ctx context.Context,
    scanRunID, competitor string,
    data *normalizer.RawPricingData,
) error {
    tiersJSON, _ := json.Marshal(data.Tiers)
    _, err := db.pool.Exec(ctx, `
        UPDATE pricing_snapshots
        SET
            pricing_model             = $1,
            tiers_json                = $2,
            normalized_monthly_cost   = $3,
            normalization_assumptions = $4,
            confidence                = $5,
            pricing_public            = $6,
            cannot_normalize          = $7,
            cannot_normalize_reason   = $8
        WHERE scan_run_id = $9 AND competitor = $10
    `,
        data.PricingModel,
        tiersJSON,
        data.NormalizedMonthlyCost,
        data.NormalizationAssumptions,
        data.Confidence,
        data.PricingPublic,
        data.CannotNormalize,
        data.CannotNormalizeReason,
        scanRunID, competitor,
    )
    return err
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Regex + scraping rules per site | LLM extraction with JSON mode | 2023-2024 | One prompt handles all pricing page layouts |
| Custom SDK per LLM provider | OpenAI-compatible gateway (TrueFoundry) | 2023+ | Swap cheap↔expensive model by changing one string |
| Pydantic/instructor (Python only) | `response_format: json_object` + Go struct | 2024 | No Python-only structured output library needed in Go |
| Separate migration framework | `ALTER TABLE IF NOT EXISTS` idiom | Established | Sufficient for small schemas; no migration file management |

---

## TrueFoundry API — Confirmed Configuration

**Source:** `/Users/blaisealbuquerque/Projects/hackathon-mar-27/.env` (verified)

| Parameter | Value | Confidence |
|-----------|-------|------------|
| Base URL | `https://llm-gateway.truefoundry.com/api/inference/openai` | HIGH — from .env |
| Chat completions endpoint | `{BASE_URL}/v1/chat/completions` | HIGH — standard OpenAI path |
| Auth header | `Authorization: Bearer {TRUEFOUNDRY_API_KEY}` | HIGH — confirmed from TrueFoundry docs |
| Cheap model | `openai-main/gpt-4o-mini` | HIGH — from .env |
| Expensive model | `openai-main/gpt-4o` | HIGH — from .env |
| JSON mode parameter | `response_format: {"type": "json_object"}` | HIGH — OpenAI standard, supported by gpt-4o-mini |
| Temperature for extraction | `0.0` | HIGH — deterministic output desired |

**Important:** The .env has the TrueFoundry API key set. The config.go already loads `TRUEFOUNDRY_CHEAP_MODEL`. No config.go changes needed — all four TrueFoundry env vars are already loaded.

---

## Target Pages Context

The 8 competitors are all **AI code review tools** (from `internal/scraper/targets.go`):

| Competitor | URL | Likely Pricing Model |
|------------|-----|---------------------|
| CodeRabbit | /pricing | per-seat or per-PR (credits/hybrid) |
| Codacy | /pricing | per-seat |
| DeepSource | /pricing | per-seat (has free tier) |
| Sourcery AI | /pricing | per-seat |
| Qodo | /pricing/ | per-seat or hybrid |
| Snyk | /plans/ | per-seat (complex tiers) |
| Greptile | /pricing | likely "Contact Sales" heavy |
| CodeAnt AI | /pricing | likely "Contact Sales" heavy |

Most are per-seat models — the per-seat normalization path will be exercised most. At least 2 are likely to be "Contact Sales" only, exercising EXTR-04.

---

## Proto Contract Update Required

The current `CompetitorResult` proto message only has `raw_html_stripped`, `from_cache`, `scraped_at`. Phase 2 must add extracted pricing fields so Streamlit can eventually display them. The `GetComparison` RPC also needs real response messages.

**Add to `proto/pricing/v1/pricing.proto`:**

```protobuf
message TierProto {
  string  name       = 1;
  double  base_price = 2;
  string  unit       = 3;
  repeated string features = 4;
  bool    is_opaque  = 5;
}

message ExtractedPricing {
  string  competitor                 = 1;
  string  pricing_model              = 2;
  repeated TierProto tiers           = 3;
  double  normalized_monthly_cost    = 4;
  string  normalization_assumptions  = 5;
  double  confidence                 = 6;
  bool    pricing_public             = 7;
  bool    cannot_normalize           = 8;
  string  cannot_normalize_reason    = 9;
  string  last_scanned               = 10;
}

// Update RunScanResponse to include extracted data
message RunScanResponse {
  string scan_run_id = 1;
  repeated CompetitorResult results = 2;         // existing (scrape data)
  repeated ExtractedPricing extracted = 3;       // new (extraction data)
}

// Populate GetComparison
message GetComparisonResponse {
  repeated ExtractedPricing comparisons = 1;
  string scan_run_id = 2;
}
```

After updating the proto, run `buf generate` to regenerate `gen/` files.

---

## Environment Availability

Step 2.6: Checked.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| TrueFoundry API key | LLM extraction | Yes (in .env) | — | Use cached HTML + hardcoded extraction for demo |
| TrueFoundry base URL | LLM extraction | Yes (in .env) | — | — |
| Ghost Postgres | Storage extension | Yes (DATABASE_URL in .env) | TimescaleDB | — |
| `net/http` | TrueFoundry HTTP client | Yes (stdlib) | Go 1.26 | — |
| `encoding/json` | JSON parsing | Yes (stdlib) | Go 1.26 | — |
| `pgx/v5` | DB writes | Yes (go.mod v5.9.1) | 5.9.1 | — |
| `buf` CLI | Proto regeneration | Assumed available | — | Manually copy generated types if missing |

**Missing dependencies with no fallback:** None.

**Note on TrueFoundry API key:** The JWT in `.env` has `exp: 1797033599` (≈ year 2026 — appears valid through the hackathon timeline). Verify it's not expired before the build starts.

---

## Open Questions

1. **TrueFoundry JSON mode support for `openai-main/gpt-4o-mini`**
   - What we know: gpt-4o-mini supports `response_format: json_object` on OpenAI directly. TrueFoundry routes to the same model.
   - What's unclear: Whether TrueFoundry gateway passes through `response_format` transparently or strips it.
   - Recommendation: Test with a simple extraction call in Wave 0/Hour 1. Fallback: remove `response_format` and rely on system prompt instruction alone if the gateway rejects it.

2. **Token cost per scan run**
   - What we know: 8 competitors × ~3,000 tokens (stripped HTML) + ~800 tokens (system prompt + few-shot) ≈ 30,000 input tokens. At gpt-4o-mini rates ($0.15/1M), this is ~$0.005 per scan.
   - What's unclear: Actual stripped HTML sizes for these 8 pages.
   - Recommendation: Not a blocking concern at hackathon scale. Log `len(strippedHTML)` for each page to verify truncation is working.

3. **Credit-based pricing normalization for CodeRabbit**
   - What we know: CodeRabbit may use a credit or per-PR model (not pure per-seat).
   - What's unclear: Whether the LLM can correctly map "N free PRs/month, then $X/PR" to a monthly cost at reference workload.
   - Recommendation: Include a third few-shot example in the prompt for credit/per-unit models if CodeRabbit turns out to use this model. The few-shot example costs ~200 tokens.

---

## Validation Architecture

nyquist_validation is set to `false` in `.planning/config.json`. Section skipped per workflow rules.

---

## Sources

### Primary (HIGH confidence)
- `.env` file (verified) — TrueFoundry base URL, model names, API key
- `internal/storage/ghost.go` — existing `AutoMigrate` pattern; `ALTER TABLE IF NOT EXISTS` extension approach
- `internal/handler/pricing.go` — existing RunScan flow; confirmed insertion point for extractor
- `internal/config/config.go` — all 4 TrueFoundry env vars already loaded
- `docs/design.md` lines 99-120 — locked extraction schema
- `internal/scraper/targets.go` — 8 confirmed competitor targets

### Secondary (MEDIUM confidence)
- [TrueFoundry Auth Docs](https://www.truefoundry.com/docs/gateway/authentication) — Bearer token format confirmed
- [TrueFoundry Langfuse Integration](https://langfuse.com/integrations/gateways/truefoundry) — `openai-main/gpt-4o-mini` model name format confirmed
- [OpenAI Structured Outputs Guide](https://platform.openai.com/docs/guides/structured-outputs) — `response_format: json_object` behavior documented
- [PostgreSQL ALTER TABLE docs](https://www.postgresql.org/docs/current/ddl-alter.html) — `ADD COLUMN IF NOT EXISTS` idempotency confirmed

### Tertiary (LOW confidence)
- WebSearch findings on LLM JSON reliability — community patterns, not officially benchmarked
- Token cost estimates — based on OpenAI published rates, TrueFoundry may differ slightly

---

## Metadata

**Confidence breakdown:**
- TrueFoundry endpoint/auth: HIGH — .env has confirmed working values
- Go HTTP + JSON pattern: HIGH — stdlib, well-established
- Normalization algorithm: HIGH — pure arithmetic, no external deps
- Schema extension: HIGH — PostgreSQL idempotent migration pattern
- Prompt engineering: MEDIUM — LLM behavior is probabilistic; few-shot structure is recommended but output consistency needs testing
- Credit model normalization: LOW — edge case for CodeRabbit; may need iteration

**Research date:** 2026-03-27
**Valid until:** 2026-04-27 (TrueFoundry endpoint/model names stable; LLM behavior patterns stable)
