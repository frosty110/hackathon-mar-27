# Final Idea Comparison — Deep Agents Hackathon

Deep Agents Hackathon, RSAC 2026. March 27, 2026.
13 ideas evaluated. 5 advanced to deep research (landscape analysis, competitor mapping, premise challenges, alternative approaches). Scores below reflect v2 (sharpened) versions.

Formula: (Autonomy x3) + (Idea x3) + (Tech x2) + (Tools x2) + (Demo x2) + (Build x1). Max 65.

---

## Final Rankings (v2)

| Rank | # | Idea | Auto | Idea | Tech | Tools | Demo | Build | v1 | **v2** | Delta |
|------|---|------|------|------|------|-------|------|-------|----|----|-------|
| **1** | 5 | **Pricing Radar** | 5 | 5 | 4 | 5 | 5 | 4 | 59 | **62** | +3 |
| **2** | 13 | Deal Memory | 5 | 4 | 4 | 5 | 5 | 3 | 58 | **60** | +2 |
| **2** | 11b | Churn Rescue | 5 | 4 | 4 | 5 | 5 | 3 | 58 | **60** | +2 |
| **4** | 12b | Revenue Anomaly | 4 | 4 | 4 | 5 | 5 | 3 | 55 | **58** | +3 |
| **5** | 7 | Pricing Change Comms | 4 | 4 | 3.5 | 3.5 | 4.5 | 4 | 52 | **55** | +3 |

---

## #1: Pricing Radar (62/65)

**Pitch:** An autonomous agent that continuously monitors competitor pricing pages, normalizes across pricing models, detects changes in real-time, and delivers strategic analysis with recommended responses.

**Reframe (v2):** "Pricing Radar" — continuous surveillance, not point-in-time analysis. Differentiates from Tierly (closest competitor, which does on-demand reports).

### Why It Wins

- **No direct competitor** in SaaS pricing monitoring. E-commerce scrapers (Prisync, Priceva) can't handle tiered SaaS pricing. Website monitors (Visualping, Distill) detect changes but don't understand pricing. Tierly does analysis but not continuous monitoring.
- **Highest autonomy.** Fully end-to-end: scrape, extract, normalize, compare, detect changes, analyze, advise. Zero human intervention.
- **The meta-moment.** "This scan cost $0.04 and took 38 seconds. A human analyst takes 2 days and $800." Thematically perfect for a pricing intelligence product.
- **Buildability improved** (v1: 3, v2: 4). Dropping Playwright and using httpx with pre-selected pages removes the #1 build risk.

### Sponsor Stack

| Sponsor | Role | Essential? |
|---------|------|-----------|
| **Ghost** | Agent creates per-scan Postgres DB to store extracted pricing data, fork for historical comparison | Yes — the agent managing its own database is a core autonomy proof point |
| **TrueFoundry** | Routes extraction tasks (cheap model for HTML parsing, expensive model for strategic analysis) | Yes — model routing is architecturally meaningful |
| **Aerospike** | Vector search for "who prices like us?" semantic similarity across competitor profiles | Yes — enables the comparison/clustering feature |
| **Overmind** | Monitors agent's own LLM costs. Powers the meta-moment: "this scan cost $0.04" | Strong — thematically perfect but product works without it |

### LinkedIn Demand Validation

**Confirmed real demand:**
- Carmen Insignares Newell (scored 95): Manually analyzed 125 AI startup pricing pages. DMs open for the spreadsheet — the output itself has demand.
- Koshima Satija (scored 92): Analyzed 200+ Voice AI companies over 3 months. Found 12 distinct pricing models, only 45% have public pricing.
- Wisam/Stripe (scored 98): Companies with 100%+ growth change pricing 3x more frequently. Monitoring window is weekly, not quarterly.
- Glenn Turner (scored 94): Figma's pricing change wiped $4B+ in market value. One change reshapes the landscape.

**Validated risks:**
1. **55% hide pricing behind "contact sales."** Flip the gap: transparency map IS intelligence. "6 of 10 competitors are transparent. 4 are hiding."
2. **Normalization requires assumptions.** Make them explicit: "Normalized to: 50-person team, 1M tokens/month."
3. **Public pricing is the least interesting signal.** Scope honestly — this monitors the public signal. Speed + normalization + change detection is the value.

### Demo Script (3 min)

- **0:00-0:30** — "Carmen Insignares Newell spent weeks analyzing 125 AI pricing pages. This agent does it in 38 seconds."
- **0:30-1:00** — Show the agent scanning 10 real AI competitor pages. Normalized comparison table populates live.
- **1:00-1:30** — Transparency map: "7 public, 3 hidden behind contact sales. The hiding is itself intelligence."
- **1:30-2:15** — CHANGE DETECTION: a pre-staged competitor page changes price. Agent detects within 30 seconds, delivers strategic analysis: "Competitor X raised Pro tier 60%. You're now the cheapest option in your tier. Recommendation: hold pricing, emphasize value gap in marketing."
- **2:15-2:45** — "Pricing Architect" mode: agent recommends your optimal response to the change.
- **2:45-3:00** — Meta-moment: "Total cost of this analysis: $0.04. Time: 38 seconds. A human analyst: $800 and 2 days."

### Honest Weaknesses

- No phone call demo moment (unlike #13, #11b). The "whoa" is visual (real-time change detection), not physical (phone ringing).
- Scraping is inherently fragile. Pre-selected pages + httpx mitigates but doesn't eliminate.
- The normalization layer is where credibility lives or dies. If the LLM makes bad assumptions, the comparison is misleading.

### Build Plan (8 hours)

| Hour | Task | Cut point |
|------|------|-----------|
| 1-2 | Core extraction pipeline: httpx + LLM extraction on 3-5 pre-selected pages | Must work by hour 2 or pivot |
| 2-3 | Normalization logic + Ghost DB storage | Skip Ghost, use SQLite if blocked |
| 3-4 | Comparison dashboard + TrueFoundry model routing | Dashboard can be minimal |
| 4-5 | Change detection + Aerospike vector similarity | Aerospike can be cut if behind |
| 5-6 | Strategic analysis / "Pricing Architect" feature | Nice-to-have |
| 6-7 | Overmind integration + meta-cost tracking | Nice-to-have |
| 7-8 | Demo polish, pre-stage the change detection moment | Non-negotiable |

---

## #2 (tie): Deal Memory (60/65)

**Pitch:** "Your company's deal memory, delivered by phone." Agent pulls from Gong, Stripe, and GitHub, finds the historical deal most similar to your current one, and calls the rep with a story about what happened last time.

**Reframe (v2):** From "deal intelligence" to "institutional deal memory." Dodges Rox AI ($1.2B), Gong Mission Andromeda, and Salesforce Agentforce comparisons. Makes vector search (Aerospike) the hero.

### Why It's Strong

- **Best Airbyte story.** All 3 available connectors (Gong, Stripe, GitHub) are naturally load-bearing. This is Airbyte's pitch deck come to life.
- **Phone call with narrative storytelling.** Not "3 risk signals detected" but "Let me tell you about the deal that looked just like this one. It died on Day 22 when the champion disengaged. You're on Day 18."
- **5/5 tool use.** Airbyte (3 connectors), Aerospike (vector similarity), Bland AI (phone call), TrueFoundry (model routing), Ghost (per-deal DB).

### Sponsor Stack

| Sponsor | Role | Essential? |
|---------|------|-----------|
| **Airbyte** | Gong (call sentiment), Stripe (billing history), GitHub (tech stack) | Yes — the triple-connector is the product |
| **Aerospike** | Vector search for similar historical deals | Yes — "deal memory" is the differentiator |
| **Bland AI** | Calls rep with narrative coaching | Yes — the demo moment |
| **TrueFoundry** | Routes analysis tasks across models | Yes — architecturally meaningful |
| **Ghost** | Per-deal analysis DB | Defensible but not essential |

### Honest Weaknesses

- Gong requires enterprise API access. Data will be synthetic. Test connector in hour 1.
- Rox AI ($1.2B), Gong, and Salesforce Agentforce are in this space. Need the "institutional memory" reframe to differentiate.
- Sales domain may not resonate with RSAC-adjacent judges.
- Six sponsors risks looking like stacking. Ghost and Overmind are the ones to cut if pressed.

### Landscape Competitors

- **Rox AI** ($1.2B, March 2026): Autonomous sales agents. Doesn't call reps, doesn't do deal similarity.
- **Gong Mission Andromeda**: AI coaching on their own data. Building the coaching layer this idea proposes.
- **Salesforce Agentforce** ($540M ARR): Agentic CRM is mainstream. But it's a platform, not a focused tool.
- **Gap:** Nobody calls the rep with a narrative story. That remains the genuine differentiator.

---

## #2 (tie): Churn Rescue (60/65)

**Pitch:** "There are 50 tools that predict churn. Zero that do anything about it. We built the one that picks up the phone."

### Why It's Strong

- **Same Airbyte triple-connector strength** as #13. Stripe (billing signals), Gong (call sentiment), GitHub (usage decline).
- **Sharp competitive positioning.** Every existing tool (Gainsight, ChurnZero, Agency.inc) stops at predict-and-alert. None make the call.
- **Phone call is diagnosis-driven** (v2 improvement): "Your API broke after our v3.2 release, here's the fix" — not a generic retention pitch.

### Sponsor Stack

| Sponsor | Role | Essential? |
|---------|------|-----------|
| **Airbyte** | Stripe (billing), Gong (sentiment), GitHub (usage) | Yes — multi-signal churn detection |
| **TrueFoundry** | Multi-model reasoning for churn scoring | Yes — the reasoning layer |
| **Bland AI** | Calls at-risk customers with personalized retention | Yes — the action that differentiates |
| **Ghost** | Per-analysis-run tracking DB | Defensible |
| **Aerospike** | Historical churn pattern matching | Cut first if behind |

### Honest Weaknesses

- Churn prediction is not novel. The innovation is autonomous action, not detection.
- Same Gong access risk as #13.
- "Would you really let AI call customers?" Production concern. For demo, theatrical value is high regardless.
- Seven integrations in 8 hours. Cut Aerospike, Auth0, Overmind first.

### Landscape Competitors

- **Gainsight, ChurnZero, Totango**: Predict and alert. None take action.
- **Agency.inc (Kai)**: Closest competitor. Agentic but limits to digital channels (emails, proposals). Doesn't call.
- **ChurnZero**: Recently launched 12 "agentic AI" agents. Monitor and draft emails. Don't call.
- **Gap:** Autonomous voice outreach via Bland AI is a genuine competitive gap.

---

## #4: Revenue Anomaly (58/65)

**Pitch:** "Revenue dropped 15%. This agent finds out why in 30 seconds by cross-correlating Stripe payments, GitHub deploys, and Gong call sentiment — then calls your CFO with the diagnosis."

### Why It's Interesting

- **Cross-silo investigation** is the technical showcase. Connecting a revenue drop to a specific PR is something no existing tool does.
- **Visual causal timeline** (v2 improvement) is potentially a stronger "whoa" than the phone call — animated connecting lines showing the causal chain materializing.
- **Same strong Airbyte triple-connector** as #13 and #11b.

### Honest Weaknesses

- Architecturally near-identical to #13 and #11b (same connectors, same phone call, same Aerospike). Narrower decision space.
- **Tellius Agent Mode** (Oct 2025) already does autonomous revenue anomaly investigation with ranked root causes.
- **Anodot** does real-time anomaly detection with cross-metric correlation.
- The LLM cross-correlation producing a correct causal hypothesis from non-planted data is unproven.

---

## #5: Pricing Change Comms (55/65)

**Pitch:** "Figma lost $4B in market value from one pricing change. Not because the price was wrong — because the communication was wrong. This agent gets the communication right."

**Reframe (v2):** "Figma Time Machine" — open the demo showing what Figma should have done.

### Why It's Interesting

- **Figma framing** is instantly relatable. Every judge knows the story.
- **Bland AI phone call** to high-value at-risk customers is still a strong demo moment.
- **Timely:** AI pricing shifts (seat-based to usage-based) are accelerating, meaning pricing changes happen more frequently.

### Honest Weaknesses

- Dropped to only 2 live sponsor integrations (Bland AI, Aerospike) to reduce build risk. Weaker tool use score.
- Narrow applicability — pricing changes are infrequent events.
- No direct customer validation. Figma case is compelling but secondhand.
- No existing landscape competitors, which could mean the market doesn't exist.

---

## Cross-Cutting Patterns

### Architecture Families

**Family A: Airbyte Triple-Connector + Bland AI Phone Call**
Ideas #13, #11b, #12b share nearly identical architecture. Same 3 Airbyte connectors (Gong, Stripe, GitHub), same Aerospike vector search, same Bland AI call, same TrueFoundry routing. If you pick one, you're implicitly rejecting the other two. Pick the domain that excites you most.

**Family B: Web Scraping + LLM Extraction**
Idea #5 is architecturally distinct. No Airbyte, no Bland AI. Scraping + extraction + normalization + change detection. Different risk profile (scraping fragility vs API access).

**Family C: Segmentation + Multi-Channel Outreach**
Idea #7 is its own thing. Segmentation logic + personalized messaging + phone calls.

### Shared Risks

- **Gong API access** affects #13, #11b, #12b. Must verify in hour 1 or pre-seed data.
- **Buildability** is universally tight (3-4/5). Nobody has margin for debugging. The winner manages build risk best.
- **Synthetic data** weakens the autonomy narrative for any idea. Pre-seed realistic data during development.

### Demo Differentiators

| Idea | Demo Moment | Sensory Channel |
|------|-------------|-----------------|
| #5 Pricing Radar | Real-time change detection + strategic analysis in 30 sec | Visual (screen) |
| #13 Deal Memory | Phone call telling a story about a similar deal that died | Audio (phone rings) |
| #11b Churn Rescue | Phone call with diagnosis-driven retention offer | Audio (phone rings) |
| #12b Revenue Anomaly | Visual causal timeline + phone briefing | Visual + Audio |
| #7 Pricing Comms | "Figma Time Machine" + phone call to VIP customer | Audio (phone rings) |

---

## Decision Matrix

| If you want... | Build | Score | Risk |
|----------------|-------|-------|------|
| Highest score, no direct competitors | **#5 Pricing Radar** | 62 | Scraping fragility |
| Best sponsor story (Airbyte showcase) | **#13 Deal Memory** or **#11b Churn Rescue** | 60 | Gong API access |
| Phone ringing on the table | **#13, #11b, #12b, or #7** | 55-60 | Build complexity |
| Most technically interesting | **#12b Revenue Anomaly** (cross-correlation) | 58 | Tellius exists |
| Best one-line pitch | **#11b Churn Rescue** | 60 | Churn prediction isn't novel |
| Safest build (fewest integration risks) | **#5 Pricing Radar** (after v2 de-risking) | 62 | Extraction prompt quality |
