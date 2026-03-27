# Target Pricing Pages for Demo

## Primary (8 pages, all httpx-fetchable)

| # | Company | URL | Model | Key Price |
|---|---------|-----|-------|-----------|
| 1 | CodeRabbit | https://www.coderabbit.ai/pricing | Per-seat | Free; Pro $24/seat/mo (annual), $30 monthly |
| 2 | Codacy | https://www.codacy.com/pricing | Per-committer | Free; Team $18/dev/mo (annual), $21 monthly |
| 3 | DeepSource | https://deepsource.com/pricing | Per-seat + AI credits | Free (OSS); Team $24/user/mo + $120/yr AI credit |
| 4 | Sourcery AI | https://www.sourcery.ai/pricing | Per-seat | Free (OSS); Pro $12/seat/mo; Team $24/seat/mo |
| 5 | Qodo | https://www.qodo.ai/pricing/ | Per-seat + credits | Free (75 credits); Teams $30/user/mo |
| 6 | Snyk | https://snyk.io/plans/ | Per-contributing dev | Free; Team $25/dev/mo (min 5, max 10) |
| 7 | Greptile | https://www.greptile.com/pricing | Per-seat + usage | Cloud $30/seat/mo (50 reviews incl, $1/extra) |
| 8 | CodeAnt AI | https://www.codeant.ai/pricing | Per-seat | Premium $24/user/mo |

## "Contact Sales" / Hidden Pricing (for transparency map)

| Company | URL | Notes |
|---------|-----|-------|
| SonarQube/SonarCloud | https://www.sonarsource.com/plans-and-pricing/ | JS-rendered, not httpx-fetchable. Server editions = contact sales. |
| GitHub Copilot Code Review | N/A | Bundled with Copilot tiers ($10-39/user/mo), no standalone page |

## Macroscope (the "client")

- **Price:** $30/dev/month, 5-seat minimum
- **Positioning:** Top of range, tied with Qodo and Greptile
- **Differentiation:** AST + LLM (structural understanding, not just pattern matching)

## Mock Page for Change Detection Demo

A localhost-served HTML page mimicking one competitor's pricing. During demo, a teammate changes the price to trigger the change detection alert and Pricing Architect response.

Suggest: mock CodeRabbit's page, change Pro from $24 to $36/seat. This triggers a Pricing Architect response: "CodeRabbit raised Pro 50%. At $36/seat they now match our $30 floor but without AST analysis. Recommendation: hold price, publish a 'bugs caught' comparison showing depth advantage."
