# Pricing Radar

An autonomous agent that scrapes competitor pricing pages, extracts structured data via LLM, normalizes across heterogeneous pricing models (per-seat, per-token, credits, etc.), detects changes, and generates strategy-grounded responses using internal positioning docs. Built as a hackathon project targeting a 3-minute live demo.

## Quick Start

```bash
cp .env.example .env
# Fill in DATABASE_URL, TRUEFOUNDRY_API_KEY, etc. in .env
go run ./cmd/server
```

## Stack

- **Go backend** — Connect-Go API server, pgx (Ghost Postgres), goquery (HTML scraping)
- **Streamlit frontend** — thin display layer, calls Go API via HTTP/JSON
- **Ghost** — managed Postgres for scan data storage
- **TrueFoundry** — AI gateway for LLM extraction and strategic analysis
- **Aerospike** — vector DB for pricing profile similarity search

## Development

```bash
# Generate protobuf stubs (after editing proto/pricing/v1/pricing.proto)
buf generate

# Run the server
go run ./cmd/server
```
