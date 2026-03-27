# Idea #5: Competitive Pricing Intelligence Agent

## One-Line Pitch

An autonomous agent that continuously scrapes, normalizes, and monitors competitor pricing pages across your market vertical, alerting you the moment someone changes their pricing before your sales team finds out from a prospect.

## The Problem

**Who has this problem:** Product managers, heads of pricing, and founders at AI/SaaS companies who need to know what competitors charge. Carmen Insignares Newell (scored 95 in research) manually analyzed 125 AI startup pricing pages. Koshima Satija (scored 92) found that 55% of voice AI companies hide pricing behind "contact sales." Every pricing decision at every SaaS company starts with someone opening 15 browser tabs and copying numbers into a spreadsheet.

**Current workaround:** Manual spreadsheet research. An analyst visits each competitor's pricing page, records the tiers, prices, and packaging, then tries to normalize across wildly different models (per-seat vs per-token vs per-minute vs per-API-call vs credits). This takes days. By the time it's done, the data is already stale. Nobody monitors for changes, so pricing shifts are discovered weeks or months late, often when a prospect says "your competitor just dropped their price."

**Why it sucks:** The data is always outdated. The normalization is always inconsistent (how do you compare a per-seat model to a per-token model?). Nobody re-checks after the initial research. And every company does this from scratch; there's no aggregated, continuously updated source.

## The Autonomous Agent Loop

This is the core of the product and the #1 judging criterion.

**Trigger:** The agent runs on a configurable schedule (e.g., every 6 hours) or can be triggered on demand. For the demo, we trigger it live.

**Step 1 - Discovery & Scraping:**
The agent takes a list of competitor URLs (pricing pages) as input. It autonomously fetches each page, handling JavaScript-rendered content. For pages behind "contact sales" walls, it flags them as opaque and records that fact (this itself is a useful signal).

**Step 2 - Extraction & Normalization:**
The LLM parses each pricing page and extracts structured data: plan names, prices, billing frequency, unit of measurement (seats, tokens, minutes, API calls, credits), feature lists per tier, and any usage limits. It normalizes everything into a common schema. This is where the agent reasoning matters most: it has to decide how to map "100 credits/month" vs "$0.006/token" vs "$25/seat/month" into comparable terms.

**Step 3 - Change Detection:**
The agent compares the newly extracted data against the previous snapshot stored in the database. It identifies: new plans added, plans removed, price increases/decreases, packaging changes (features moved between tiers), and new "contact sales" gates.

**Step 4 - Analysis & Alert:**
When changes are detected, the agent generates an analysis: what changed, the magnitude, likely strategic intent (e.g., "Competitor X dropped their starter tier price 20%, likely to capture SMB market share"), and recommended response options. This gets stored and can trigger notifications.

**Step 5 - Persistence & Memory:**
All snapshots, diffs, and analyses are stored with timestamps, creating a historical record of competitor pricing evolution over time. The agent uses previous analyses as context for new ones ("this is the third price decrease from Competitor X in 6 months, suggesting a land-and-expand strategy shift").

**Key autonomous decisions the agent makes:**
- How to normalize heterogeneous pricing models into comparable units
- Whether a page change is a real pricing change vs a cosmetic update
- What the strategic significance of a change is
- Whether a change warrants an urgent alert vs a routine update

## Sponsor Stack (3+ required)

### 1. Ghost (Postgres via MCP) -- ESSENTIAL
**What it does:** The agent uses Ghost to autonomously create and manage its own Postgres database for storing pricing snapshots, competitor profiles, change history, and normalized comparison tables. The agent creates tables, inserts scraped data, queries for diffs, and forks the DB for "what-if" analysis.

**Why it's load-bearing:** The entire change-detection loop depends on comparing current scrape vs historical data. Without a persistent, queryable database, you have no change detection, no historical tracking, no normalized comparison tables. The agent literally cannot function without storage, and Ghost lets the agent manage its own database schema without human DBA intervention, which reinforces the autonomy story.

**Essential or cosmetic:** Essential. Remove it and you have a one-shot scraper with no memory, no change detection, and no history. The product collapses.

### 2. TrueFoundry (Multi-LLM Gateway) -- ESSENTIAL
**What it does:** Routes the extraction and analysis prompts across multiple LLMs. Pricing page parsing is a structured extraction task (good for a fast, cheap model). Strategic analysis is a reasoning task (needs a stronger model). TrueFoundry's gateway lets the agent pick the right model for each subtask and provides observability on what each step costs.

**Why it's load-bearing:** The normalization task is the hardest part of this product. Different pricing pages require different levels of reasoning to parse. A gateway that routes simple extraction to a cheap model and complex normalization to a strong model is not just an optimization; it's the difference between the agent working reliably and hallucinating pricing data. Also, since this product is about pricing intelligence, using TrueFoundry to show what the agent itself costs to run is a delicious meta-demo moment.

**Essential or cosmetic:** Essential. You could hard-code a single model, but you'd lose the adaptive reasoning that makes extraction reliable across wildly different pricing page formats. And you'd lose the cost observability, which is thematically on-point.

### 3. Aerospike (Sub-ms Database + Vector Search) -- STRONG VALUE-ADD
**What it does:** Stores vector embeddings of pricing page content and competitor descriptions, enabling semantic similarity search. When you ask "who in my space prices most like Competitor X?" or "find competitors with a free tier," the agent queries Aerospike's vector index. Also serves as a fast cache layer for the most recent snapshot of each competitor, enabling sub-millisecond reads during the comparison step.

**Why it's load-bearing:** The normalization problem is fundamentally a similarity problem. When the agent encounters a new pricing model it hasn't seen, vector search against its memory of all previous pricing structures helps it find the closest analog. This makes extraction more reliable over time. Without it, the agent treats every pricing page as novel, missing patterns.

**Essential or cosmetic:** Honest assessment: strong value-add, not strictly essential. You could do the comparison work purely with LLM calls and Ghost for storage. But Aerospike adds real capability (semantic search across the pricing landscape) and speed (sub-ms reads for the comparison loop). It's genuinely useful, not just bolted on.

### 4. Overmind (LLM Optimization/Observability) -- NICE-TO-HAVE BUT THEMATIC
**What it does:** Drop-in SDK that monitors every LLM call the agent makes, showing cost, latency, and token usage per step. Provides optimization suggestions. For a pricing intelligence product, knowing "this agent costs $0.03 per competitor per scan cycle" is a meaningful data point.

**Why it's load-bearing:** Provides the observability layer that lets you show judges exactly what happens inside the agent. During the demo, you can show the Overmind dashboard alongside the agent output: "here's the agent reasoning, here's what each step cost, here's where it chose a cheaper model." It's the x-ray into the agent's economics.

**Essential or cosmetic:** Cosmetic in the strictest sense (the agent works without it), but thematically perfect for a pricing-focused hackathon. A pricing intelligence agent that can't tell you its own cost would be ironic.

## The "Whoa" Demo Moment

The "whoa" comes at the 1:30 mark. You've set up the agent to monitor 5-6 real competitor pricing pages. The agent has already scraped them once (pre-loaded baseline). Then, live on stage, you trigger a re-scan.

But before the re-scan, you've prepared one twist: one of the "competitor" pages is a staging page you control, and you've changed the pricing since the last scan. You change it live on stage in a browser tab judges can see.

The agent scrapes all pages, and within 30 seconds, it flags the change: "ALERT: Competitor Y increased Pro tier from $49/mo to $79/mo (+61%). Enterprise tier unchanged. Free tier removed. Strategic assessment: likely margin compression forcing upmarket move. Recommend: consider introducing a new mid-tier to capture customers they're abandoning."

The judge sees: real competitor data, normalized across different pricing models, a detected change with magnitude and context, and a strategic recommendation, all generated autonomously.

Second "whoa": pull up the Ghost database and show the historical table. "Here's every pricing change we've tracked across this market over the last [simulated] 90 days, with trend lines." A pricing analyst's dream.

## 3-Minute Demo Script

**0:00-0:30 -- Setup (The Problem)**
"Raise your hand if your company has ever been surprised by a competitor's pricing change." [Pause.] "Carmen Insignares Newell analyzed 125 AI startup pricing pages by hand for her research. Every company does this. Nobody has a system for it. We built one."

**0:30-0:50 -- Show the Baseline**
Show the Ghost database with the current pricing snapshot: a table of 5-6 real competitors, their tiers, prices, and packaging, all normalized into a common format. "This is what the agent already knows. It scraped these pages this morning and normalized per-seat, per-token, and per-minute pricing into comparable monthly costs."

**0:50-1:10 -- The Live Change**
Open a browser tab showing the staging "competitor" pricing page. Change the price visibly. "I just changed this competitor's pricing. In the real world, this happens and you don't find out for weeks."

**1:10-1:50 -- The Agent Runs**
Trigger the agent. Show it hitting each URL, extracting data via TrueFoundry's LLM gateway (show the model routing: "extraction using GPT-4o-mini, analysis using Claude"), storing results in Ghost, comparing against the baseline in Aerospike. The change is detected. The alert appears with the diff, magnitude, and strategic analysis.

**1:50-2:20 -- The Intelligence Layer**
Show the competitive landscape view: a normalized comparison table across all competitors. Show vector similarity: "Which competitor prices most like us?" Show the historical trend: "Competitor X has decreased pricing 3 times in 6 months." Show the Overmind dashboard: "This entire scan cycle cost $0.04 and took 45 seconds."

**2:20-2:50 -- The Autonomous Loop**
"This runs every 6 hours without human intervention. It creates its own database tables via Ghost. It picks the right model for each task via TrueFoundry. It builds memory of the competitive landscape in Aerospike. And Overmind shows us exactly what it costs."

**2:50-3:00 -- The Landing**
"Every company does competitive pricing research. Nobody does it continuously. This agent does. It never forgets to check. It never gets the normalization wrong twice. And it costs four cents per scan."

## Technical Architecture

```
                    ┌─────────────┐
                    │   Trigger   │
                    │ (cron/manual│
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Orchestrator│  (Python agent)
                    │   Agent     │
                    └──┬───┬───┬──┘
                       │   │   │
            ┌──────────┘   │   └──────────┐
            ▼              ▼              ▼
    ┌───────────┐  ┌──────────────┐  ┌──────────┐
    │  Scraper  │  │  LLM Gateway │  │  Storage  │
    │ (httpx +  │  │ (TrueFoundry)│  │  Layer    │
    │ playwright│  │              │  │          │
    │  for JS)  │  │ Route:       │  │ Ghost:   │
    │           │  │ - Extract →  │  │  Postgres │
    └─────┬─────┘  │   fast model │  │  (schema, │
          │        │ - Analyze →  │  │  snapshots│
          │        │   strong mdl │  │  diffs)   │
          │        │ - Normalize →│  │          │
          │        │   fast model │  │ Aerospike:│
          │        └──────┬───────┘  │  Vector   │
          │               │          │  embeddings│
          └───────────────┼──────────│  + cache   │
                          │          └─────┬──────┘
                          │                │
                    ┌─────▼────────────────▼─┐
                    │     Overmind SDK        │
                    │  (wraps all LLM calls,  │
                    │   cost/latency tracking)│
                    └─────────────────────────┘
```

**Data flow:**
1. Orchestrator receives trigger, reads competitor list from Ghost DB
2. Scraper fetches each URL (httpx for static, Playwright for JS-rendered)
3. Raw HTML sent to TrueFoundry gateway for extraction (routed to fast model)
4. Extracted structured data stored in Ghost (new snapshot row)
5. Previous snapshot fetched from Ghost, diff computed
6. If changes detected, analysis prompt sent to TrueFoundry (routed to strong model)
7. Embeddings of pricing structure stored in Aerospike for similarity search
8. All LLM calls instrumented by Overmind for cost tracking
9. Results written back to Ghost with timestamp and change metadata

## Buildability Risk Assessment

**Overall buildability: 3/5 (tight but possible with disciplined scoping)**

**What could go wrong in 8 hours:**

1. **Web scraping fragility (HIGH RISK):** Pricing pages are diverse. Some are static HTML, some are React SPAs, some load pricing dynamically. Playwright helps but adds complexity. **Mitigation:** Pre-select 5-6 competitor pages that are known to work. Have fallback HTML snapshots ready. For the demo, the "changed" competitor page is one you control, so it definitely works.

2. **LLM extraction reliability (MEDIUM RISK):** Getting an LLM to consistently extract structured pricing data from varied HTML layouts is non-trivial. Different pages have wildly different structures. **Mitigation:** Build a solid extraction prompt with few-shot examples. Test against the demo competitors during development. Accept ~80% accuracy for the demo.

3. **Normalization complexity (MEDIUM RISK):** Comparing per-seat to per-token pricing is genuinely hard. There's no clean formula. **Mitigation:** For the demo, don't try to solve the general case. Use "estimated monthly cost for a typical 50-person team" as the common unit. Let the LLM reason about the conversion with stated assumptions.

4. **Sponsor integration time (MEDIUM RISK):** Four sponsor tools means four SDKs to learn and integrate. **Mitigation:** Ghost and TrueFoundry first (they're essential). Aerospike next. Overmind last (it's a drop-in SDK, so it's fast to add).

**The hardest part:** Reliable, structured extraction from diverse pricing page HTML. This is the make-or-break technical challenge.

**What to cut if time runs out:**
- Cut Aerospike vector search (fall back to simple keyword matching in Ghost)
- Cut the strategic analysis step (just show the raw diff, skip the "strategic assessment")
- Cut historical trend visualization (just show current snapshot + most recent change)
- Cut Overmind (add it in the last 30 minutes if there's time, or skip it)

**Hour-by-hour plan:**
- Hours 1-2: Ghost DB schema, scraper for 5-6 known pages, basic extraction prompt
- Hours 3-4: TrueFoundry integration, normalization logic, change detection
- Hours 5-6: Aerospike integration, demo flow end-to-end, controlled "change" page
- Hours 7-8: Overmind integration, polish, rehearse demo, build backup paths

## Honest Weaknesses

**1. Scraping is legally and ethically gray.**
A skeptical judge might ask: "Is this just a scraper? Aren't you violating ToS?" Response: we only scrape publicly available pricing pages (the same data any human visitor sees). We don't log in, bypass paywalls, or access anything non-public. But the optics matter, and a judge who works at a company that hates scrapers might ding this.

**2. The normalization problem is not fully solvable in 8 hours.**
Comparing per-seat to per-token pricing requires assumptions (usage per seat, tokens per request, etc.). The agent will state its assumptions, but the output is an estimate, not a precise comparison. A sharp judge will probe this: "Your comparison says Competitor A costs $X/month, but that depends on usage patterns you're guessing at."

**3. The "autonomous loop" is somewhat thin.**
The agent scrapes, extracts, compares, and alerts. That's a valuable pipeline, but the actual decision-making is concentrated in the extraction and analysis steps. The scraping and storage steps are mechanical. A judge focused on "does the agent make interesting decisions?" might say the autonomy is narrower than it appears.

**4. Sponsor tool fit for Aerospike is a stretch.**
Ghost handles the primary storage. Aerospike adds vector search, which is genuinely useful but not strictly necessary. If a judge probes "why do you need two databases?", the answer is honest (vector search vs relational queries) but could feel like sponsor-stacking.

**5. Demo dependency on live web scraping.**
If the wifi is bad or a pricing page changes its structure between rehearsal and demo, things break. Mitigation: pre-cache fallback data and use a controlled staging page for the change-detection moment.

## Final Score Recommendation

Using the evaluation framework from EVALUATION.md:

| Dimension | Score | Justification |
|-----------|-------|---------------|
| **Autonomy** (3x) | **5** | The agent runs end-to-end: trigger, scrape, extract, normalize, compare, detect, analyze, alert, store. No human in the loop. It makes real decisions about how to normalize and what constitutes a meaningful change. This is a genuine autonomous loop, not a human-assisted workflow. |
| **Idea** (3x) | **5** | Carmen manually analyzed 125 pricing pages. That is the user. That is the pain. Every pricing decision at every SaaS company starts with manual competitive research. Judges will immediately recognize the problem because they've done this research themselves. |
| **Technical Implementation** (2x) | **4** | Multi-step agent with real reasoning (extraction, normalization, change detection, strategic analysis). LLM routing via TrueFoundry shows architectural sophistication. Not a 5 because the scraping layer is fundamentally fragile and the normalization problem is hard to solve robustly in 8 hours. |
| **Tool Use** (2x) | **4** | Ghost (essential, load-bearing), TrueFoundry (essential, load-bearing), Aerospike (genuinely useful but not strictly essential), Overmind (thematic but cosmetic). Three strong integrations, one nice-to-have. Not a 5 because Aerospike and Overmind don't fully clear the "breaks without it" bar. |
| **Presentation/Demo** (2x) | **5** | The live price change detection is a visible, dramatic "whoa" moment. The normalized comparison table is immediately legible. The "four cents per scan" cost reveal is memorable. The story setup ("raise your hand if you've been surprised by a competitor price change") is universally relatable. |
| **Buildability** (1x) | **3** | Scraping fragility is real. Four sponsor integrations in 8 hours is tight. The extraction prompt needs iteration to be reliable. Achievable if you scope aggressively (5-6 known pages, pre-tested), but there's no room for debugging rabbit holes. |

**Total: (5x3) + (5x3) + (4x2) + (4x2) + (5x2) + (3x1) = 15 + 15 + 8 + 8 + 10 + 3 = 59/65**

**Verdict: BUILD THIS.** Highest-scoring idea in the triage (57 initial, 59 on deep evaluation). The combination of a universally relatable problem, strong autonomy loop, and a dramatic live demo moment is hard to beat. The main risks are technical (scraping reliability, normalization accuracy) rather than conceptual, and they can be mitigated with disciplined scoping.

**Key risk to monitor:** If scraping proves too fragile during build, have a fallback plan where the agent works from pre-fetched HTML snapshots rather than live URLs. The demo impact barely changes (you still show the controlled pricing change page live), but the "fully autonomous" story weakens slightly.
