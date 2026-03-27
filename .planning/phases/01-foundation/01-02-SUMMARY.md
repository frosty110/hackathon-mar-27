---
phase: 01-foundation
plan: 02
subsystem: api
tags: [protobuf, connect-go, buf, go, config, http]

requires:
  - phase: 01-01
    provides: Go module with all dependencies installed (connectrpc.com/connect, buf.gen.yaml, go.mod tool section with protoc-gen-go and protoc-gen-connect-go)

provides:
  - pricing.proto with 5 RPCs covering all 4 phases (RunScan, GetComparison, GetChanges, GetRecommendation, GetClusters)
  - buf-generated Go types: RunScanRequest, RunScanResponse, CompetitorResult in gen/pricing/v1/pricing.pb.go
  - buf-generated Connect-Go interface: PricingServiceHandler in gen/pricing/v1/pricingv1connect/pricing.connect.go
  - Config struct loaded from env via godotenv in internal/config/config.go
  - Stub PricingHandler satisfying generated interface in internal/handler/pricing.go
  - Connect-Go server booting on :8080 with HTTP/1.1 + h2c in cmd/server/main.go
  - /healthz endpoint returning 200 ok

affects:
  - 01-03 (storage — imports config.Config, uses DatabaseURL)
  - 01-04 (scraper/handler wire — replaces stub RunScan with real implementation)
  - All later phases that import gen/pricing/v1 message types

tech-stack:
  added:
    - buf CLI (v1.66.1) — proto codegen via buf generate
    - proto/pricing/v1/pricing.proto — API contract source of truth
    - gen/pricing/v1/ — buf-generated Go code (gitignored, regenerated on demand)
  patterns:
    - Contract-first: proto defines API shape before any handler implementation
    - Simple Connect-Go interface: buf generates (ctx, *ReqType) -> (*RespType, error) signatures via `opt: simple`
    - http.Protocols for h2c: uses Go 1.24+ stdlib SetUnencryptedHTTP2(true), NOT golang.org/x/net/http2/h2c
    - godotenv for dev config: loads .env silently, production reads env directly

key-files:
  created:
    - proto/pricing/v1/pricing.proto
    - internal/config/config.go
    - internal/handler/pricing.go
    - cmd/server/main.go
  modified: []

key-decisions:
  - "Generated interface uses simple option (ctx, *ReqType -> *RespType, error) not connect.Request wrappers — buf.gen.yaml had opt: simple set from 01-01"
  - "Handler stub returns empty proto responses directly; connect.NewResponse() wrapper not needed with simple interface"

patterns-established:
  - "Pattern: buf.gen.yaml with local [go, tool, protoc-gen-go] plugin — run buf generate from project root"
  - "Pattern: http.Protocols{}.SetHTTP1(true) + SetUnencryptedHTTP2(true) for h2c without x/net dependency"
  - "Pattern: config.Load() fails fast on missing DATABASE_URL — server exits 1 at startup if misconfigured"

requirements-completed: [EXTR-01]

duration: 4min
completed: 2026-03-27
---

# Phase 01 Plan 02: Proto, buf generate, Connect-Go server, and stub handler

**pricing.proto defines 5-RPC PricingService; buf generate produces typed Go stubs; Connect-Go server boots on :8080 with h2c via stdlib http.Protocols**

## Performance

- **Duration:** ~4 min
- **Started:** 2026-03-27T21:43:33Z
- **Completed:** 2026-03-27T21:47:02Z
- **Tasks:** 2
- **Files modified:** 4 created

## Accomplishments

- proto/pricing/v1/pricing.proto written with all 5 RPCs (RunScan + 4 placeholder RPCs for phases 2-4)
- `buf generate` ran cleanly producing pricing.pb.go and pricingv1connect/pricing.connect.go
- PricingHandler stub satisfies generated PricingServiceHandler interface (simple signature variant)
- config.go reads all required env vars with DATABASE_URL required, PORT defaulting to 8080
- main.go boots Connect-Go server with HTTP/1.1 + h2c using Go 1.26 stdlib http.Protocols
- `go build ./...` exits 0 on first attempt

## Task Commits

1. **Task 1: Write pricing.proto and run buf generate** - `2717eed` (feat)
2. **Task 2: Write config.go, stub handler, and main.go** - `d93cbf1` (feat)

**Plan metadata:** (docs commit below)

## Files Created/Modified

- `proto/pricing/v1/pricing.proto` - API contract with 5 RPCs and all Phase 1-4 message types
- `internal/config/config.go` - Config struct + Load() reading DATABASE_URL, PORT, TrueFoundry, Aerospike env vars
- `internal/handler/pricing.go` - Stub PricingHandler implementing all 5 generated RPC methods
- `cmd/server/main.go` - Connect-Go server entrypoint with h2c, health check, and config loading

## Decisions Made

- The generated interface uses the `simple` option (set in buf.gen.yaml from plan 01-01), which produces `(ctx context.Context, req *pricingv1.RunScanRequest) (*pricingv1.RunScanResponse, error)` signatures rather than `*connect.Request[...]` wrappers. Handler was written to match.
- Stub handler returns empty proto responses directly (no `connect.NewResponse()` wrapper needed with simple interface).

## Deviations from Plan

None - plan executed exactly as written. The only adaptation was matching the generated interface signature (simple option produces raw proto types, not connect.Request wrappers) — this was already documented in the generated code and consistent with what buf.gen.yaml specified.

## Issues Encountered

None - `go build ./...` passed on first attempt. buf generate ran cleanly without errors.

## Known Stubs

- `internal/handler/pricing.go` - All 5 RPC methods return empty proto responses. This is intentional per the plan; real RunScan implementation is wired in plan 01-04. The stub allows the server to compile and boot while later plans add real implementations.

## User Setup Required

None - no external service configuration required by this plan. DATABASE_URL must be set in .env (from Ghost DB creation in plan 01-03).

## Next Phase Readiness

- 01-03 (storage): can import `internal/config` and use `cfg.DatabaseURL` for pgxpool connection
- 01-04 (scraper/handler wire): can replace stub `RunScan` body with real scraper call; generated types are ready
- gen/ is gitignored — downstream agents must run `buf generate` before importing gen/ packages

---
*Phase: 01-foundation*
*Completed: 2026-03-27*
