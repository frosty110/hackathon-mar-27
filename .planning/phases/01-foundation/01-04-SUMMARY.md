---
phase: 01-foundation
plan: "04"
subsystem: scraper
tags: [scraper, goquery, errgroup, html-stripping, cached-fallback, runscan]
dependency_graph:
  requires: [01-02, 01-03]
  provides: [scraper-package, runscan-live, cached-html-fallbacks]
  affects: [internal/handler/pricing.go, internal/scraper/, demo-data/cached/]
tech_stack:
  added: [golang.org/x/sync/errgroup, github.com/PuerkitoBio/goquery]
  patterns: [concurrent-fan-out, cache-fallback, html-stripping]
key_files:
  created:
    - internal/scraper/targets.go
    - internal/scraper/fetcher.go
    - internal/scraper/parser.go
    - demo-data/cached/coderabbit.html
    - demo-data/cached/codacy.html
    - demo-data/cached/deepsource.html
    - demo-data/cached/sourcery.html
    - demo-data/cached/qodo.html
    - demo-data/cached/snyk.html
    - demo-data/cached/greptile.html
    - demo-data/cached/codeant.html
  modified:
    - internal/handler/pricing.go
decisions:
  - "Sourcery AI named 'Sourcery AI' to match target-pages.md (plan used 'Sourcery' — corrected)"
  - "Handler uses plain proto types not connect.Request wrappers (matches generated simple option interface)"
  - "All 8 pages fetched live successfully — no JS-only pages encountered; from_cache=false for all"
metrics:
  duration: "~6 minutes"
  completed: "2026-03-27T21:58:50Z"
  tasks_completed: 2
  files_created: 11
  files_modified: 1
---

# Phase 01 Plan 04: Scraper Implementation Summary

Concurrent scraper with goquery HTML stripping wired into RunScan handler, backed by 8 cached HTML fallback files and Ghost DB storage.

## What Was Built

### Task 1: Scraper Package (commit: 1b49fad)

Three new files in `internal/scraper/`:

- `targets.go`: `DefaultTargets()` returns 8 `Target` structs with Name, URL, FallbackPath for all competitors from `demo-data/target-pages.md`
- `fetcher.go`: `FetchAll(ctx, targets)` uses `errgroup` for concurrent fan-out. Each goroutine handles errors internally and returns nil (never cancels the group). Falls back to `os.ReadFile(t.FallbackPath)` on any fetch failure. Sets `User-Agent: Mozilla/5.0 (compatible; PricingRadar/1.0)` on all requests.
- `parser.go`: `StripBoilerplate(rawHTML)` uses goquery to remove nav, footer, script, style, header, noscript, and `[aria-hidden='true']` elements, then returns `doc.Find("body").Html()`

8 cached HTML fallback files in `demo-data/cached/` fetched via curl with matching User-Agent. All files are non-empty (46KB–563KB each).

### Task 2: Handler Wiring (commit: fcfaac6)

`internal/handler/pricing.go` updated:
- Added `targets []scraper.Target` field to `PricingHandler`
- `NewPricingHandler` initializes targets from `scraper.DefaultTargets()`
- `RunScan` now calls `NewScanRun` → `FetchAll` → `StripBoilerplate` → `SaveSnapshot` → `FinishScanRun`
- Returns `RunScanResponse` with 8 `CompetitorResult` entries including `raw_html_stripped`, `from_cache`, `scraped_at`

## End-to-End Test Results

RunScan called via curl after server startup:

```
scan_run_id: b01e17d5-8c33-4661-af75-416496225c10
results: 8 entries
competitors: CodeRabbit, Codacy, DeepSource, Sourcery AI, Qodo, Snyk, Greptile, CodeAnt AI
from_cache: false for all (all 8 pages fetched live)
scraped_at: 2026-03-27T14:58:32-07:00 (all entries)
raw_html_stripped: non-empty HTML for all entries
response size: ~869KB
```

### Cache Fallback Status

All 8 pages returned live HTML (from_cache=false). No pages fell back to cache during this scan. Cached files serve as resilience fallbacks for network failures.

### Cached File Sizes

| Competitor | File | Size |
|-----------|------|------|
| Codacy | codacy.html | 379KB |
| CodeAnt AI | codeant.html | 563KB |
| CodeRabbit | coderabbit.html | 167KB |
| DeepSource | deepsource.html | 243KB |
| Greptile | greptile.html | 214KB |
| Qodo | qodo.html | 120KB |
| Snyk | snyk.html | 253KB |
| Sourcery AI | sourcery.html | 46KB |

## Deviations from Plan

### Auto-fixed Issues

None. Plan executed with one minor naming adjustment:

**1. Sourcery AI naming** (non-code deviation)
- Plan's `DefaultTargets()` example used `"Sourcery"` but `demo-data/target-pages.md` lists `"Sourcery AI"`
- Used `"Sourcery AI"` to match the authoritative source document

**2. Handler signature adaptation**
- Plan showed handler with `connect.Request[...]` wrappers
- Generated code (`pricingv1connect.PricingServiceHandler`) uses simple proto type signatures (from the `simple` buf.gen.yaml option set in 01-01)
- Adapted handler to match the actual generated interface (plain `*pricingv1.RunScanRequest` not `*connect.Request[...]`)

## Known Stubs

The following methods remain stubs returning empty responses (intentional — future plans):
- `GetComparison` — Phase 2
- `GetChanges` — Phase 2
- `GetRecommendation` — Phase 3
- `GetClusters` — Phase 3

These do not block the plan's goal (RunScan returning real stripped HTML) which is fully achieved.

## Self-Check: PASSED
