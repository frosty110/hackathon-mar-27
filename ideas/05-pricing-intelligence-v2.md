# Idea #5: Competitive Pricing Intelligence Agent -- Deep Exploration (v2)

## Original Score: 59/65

---

## 1. Landscape Research

### What Already Exists

The competitive pricing intelligence space breaks into three distinct tiers:

**Tier 1: E-commerce Price Scrapers (Mature, Crowded)**
- **Prisync** ($99-$399/mo): SKU-level price tracking for e-commerce. Monitors product prices and stock levels, suggests repricing. Focused on physical goods and marketplaces.
- **Priceva**: Automated price scraping software with dashboards and alerts. E-commerce focused.
- **Omnia Retail**: Marketplace-specific (Buy Box dynamics, channel pricing). Enterprise retail.
- **Price2Spy** ($67.95/mo+): Cost-effective for small businesses, hospitality.
- **Intelligence Node**: AI-powered, tracks 1B+ products daily. Pure retail/e-commerce.
- **Competera**: Enterprise-grade price scraping for retailers.

These are all focused on **product/SKU pricing** (what does this TV cost on Amazon vs Best Buy). They do NOT handle SaaS pricing structures, tiered plans, or heterogeneous pricing models.

**Tier 2: Website Change Monitors (Generic, Not Purpose-Built)**
- **Visualping** (free tier, paid from ~$14/mo): Monitors any webpage for visual changes. Can alert on pricing page changes but provides no extraction, normalization, or analysis. Just "this page changed" with a screenshot diff and AI summary.
- **Distill.io**: Element-level monitoring, browser-based.
- **Fluxguard** ($99-$499/mo): Developer-focused HTML diffs. Side-by-side comparisons. No semantic understanding.
- **ChangeTower**: Compliance/audit-focused change tracking.

These detect that something changed but have zero understanding of what pricing data means. No normalization, no cross-competitor comparison, no strategic analysis.

**Tier 3: SaaS Pricing Intelligence (Emerging, Direct Competitor)**
- **Tierly** (tierly.app): This is the closest direct competitor. AI-powered competitive pricing tool specifically for SaaS. Scores pricing tiers across 5 dimensions, benchmarks against competitors, provides recommendations. Uses 10+ AI models per analysis. Has a free first analysis.

Tierly is the one to watch. However, from the research, Tierly appears to be a **point-in-time analysis tool** (you submit competitors, get a report) rather than a **continuous monitoring agent**. It does not appear to do autonomous change detection or ongoing surveillance. It is more of a "consultant in a box" than a "pricing radar."

**Tier 4: LLM Pricing Trackers (Adjacent)**
- **PricePerToken.com**: Tracks LLM API pricing across 300+ models. Real-time comparison.
- **Artificial Analysis**: Compares AI models on intelligence, performance, and price.

These are vertical-specific (AI/LLM only) and manually curated, not autonomous agents.

### Where the Conventional Wisdom is Wrong

**Conventional wisdom says:** "Price monitoring is a solved problem -- use Prisync or Visualping."

**Where it's wrong:** Those tools solve a fundamentally different problem. E-commerce price monitoring tracks homogeneous units (same product, different stores, compare dollar amounts). SaaS pricing monitoring requires understanding **heterogeneous pricing models** -- how do you compare $25/seat/month vs $0.006/token vs 100 credits/month vs "contact sales"? This is a reasoning problem, not a scraping problem. The existing tools can scrape but cannot think.

**Conventional wisdom says:** "Just use a website change monitor and set up alerts."

**Where it's wrong:** Knowing a page changed is ~5% of the value. The other 95% is: what specifically changed, what does it mean when normalized against your pricing, and what should you do about it. Visualping gives you "page changed, here's a screenshot." A pricing analyst needs "they dropped their Pro tier 20%, removed the free tier, and this suggests they're moving upmarket -- here's how that affects your competitive position."

**Conventional wisdom says:** "Tierly already does this."

**Where it's wrong (maybe):** Tierly does point-in-time competitive analysis, not continuous monitoring with change detection. The gap is the autonomous loop -- the agent that runs every 6 hours, remembers what it saw last time, and alerts you proactively. Tierly is a telescope; this idea is a radar.

---

## 2. Builder Mode Questions

### What's the COOLEST version of this idea?

The coolest version is a **"War Room" live dashboard** that makes pricing intelligence feel like mission control. Imagine:

- A real-time map of your competitive landscape where each competitor is a node, sized by how similar their pricing is to yours (via embeddings), colored by threat level (red = just changed pricing aggressively, green = stable).
- When you click a competitor, you see their full pricing history as a timeline, with annotated events: "Dropped starter tier 30% -- 2 weeks after your price increase."
- A "What If" simulator: "If we drop our Pro tier to $39, which competitors would we undercut? Which would we still lose to? What's the projected revenue impact?"
- The agent proactively sends you a daily briefing: "Your market is stable. No changes in 72 hours. One competitor's pricing page now returns a 404, which may indicate an upcoming relaunch."
- The meta-kicker: the dashboard shows what the agent itself costs to run, per scan, per competitor, per month. A pricing intelligence product that understands its own economics.

What makes it genuinely delightful (not just useful): the moment where you realize the agent noticed something you never would have. Like: "Competitor Z added a new 'Teams' tier between their Pro and Enterprise plans. This is the third competitor in your space to add a mid-market tier in the last 60 days. Trend signal: the market is segmenting downward."

### Who would you show this to at the hackathon, and what specific moment would make them say "whoa"?

Show it to **a judge who runs a SaaS company or works in product/growth**. They have personally done the competitor pricing spreadsheet exercise.

The "whoa" moment: You open the dashboard showing 6 real competitors with normalized pricing. You say "watch this" and change a price on a staging competitor page in real time. Within 30 seconds, the dashboard updates: a red alert pulses on that competitor's node, the change is annotated with magnitude and strategic analysis, and the agent says something like: "This is the second price increase from Competitor Y in 4 months. Combined with their recent removal of the free tier, this suggests margin pressure and an upmarket repositioning. Consider: your free tier is now the only one in the market. This is either a vulnerability or an opportunity."

The second "whoa": you ask the agent a natural language question: "Who in my space has the cheapest option for a 10-person team using 1M tokens/month?" and it answers with a normalized comparison table, pulling from its memory of every competitor's pricing structure -- including the ones that price per-seat, per-token, and per-credit.

### What's the fastest path to a working demo in 8 hours?

1. **Hours 1-2**: Set up Ghost DB with schema (competitors, snapshots, changes). Build scraper for 5 pre-selected pricing pages (use httpx, skip Playwright -- pick pages that work without JS). Create a staging page you control for the live change moment.

2. **Hours 3-4**: Build the extraction prompt (this is the hardest part -- spend time here). Wire up TrueFoundry for LLM routing. Get extraction working reliably for your 5 pages. Store results in Ghost.

3. **Hours 5-6**: Build change detection (diff current vs previous snapshot in Ghost). Build the analysis prompt for detected changes. Wire up the end-to-end loop: trigger -> scrape -> extract -> store -> diff -> analyze -> alert.

4. **Hours 7-8**: Build a minimal but visually impressive frontend (even a well-styled CLI output or a simple React dashboard). Integrate Aerospike for the "who prices like us?" query. Add Overmind. Rehearse the demo. Pre-seed the Ghost DB with "historical" data for the trend story.

**Key shortcut**: Do NOT build a scheduler. Just have a "run scan" button. The "runs every 6 hours" is a story you tell, not a feature you demo. Focus all time on making the extraction reliable and the change detection dramatic.

**Key shortcut #2**: Pre-seed the database with 2-3 "historical" snapshots (scraped during development) so you have trend data to show. Don't try to demonstrate history building up in real-time.

### What existing product is closest to this, and how is this different?

**Closest: Tierly** (tierly.app)

Tierly does AI-powered SaaS pricing analysis. It extracts tier data, matches comparable tiers across products, scores them, and generates reports with recommendations.

**How this is different:**

| Dimension | Tierly | This Agent |
|-----------|--------|------------|
| Mode | On-demand analysis | Continuous autonomous monitoring |
| Change detection | No (point-in-time) | Yes (diffing against history) |
| Alerting | No | Yes (proactive alerts on changes) |
| Normalization | Tier-to-tier matching | Cross-model normalization (seat vs token vs credit) |
| Memory | Per-analysis | Persistent, cumulative history |
| Strategic analysis | Scoring & recommendations | Change-driven strategic assessment |
| Interaction model | "Run an analysis" | "Agent watches for you" |

The key differentiator is **autonomy and continuous monitoring**. Tierly is a tool you use. This is an agent that works for you while you sleep.

### What's the 10x version if you had unlimited time?

- **Multi-source intelligence**: Don't just scrape pricing pages. Monitor job postings (hiring for "pricing manager" = pricing change incoming), press releases, G2 reviews mentioning pricing, social media complaints about price increases, SEC filings for public companies.
- **Predictive pricing intelligence**: "Based on Competitor X's pattern (price decrease every Q1, new tier every Q3), we predict they will announce a new enterprise tier in the next 60 days."
- **Automated response playbook**: When a competitor drops prices, the agent doesn't just alert -- it drafts updated battle cards for the sales team, generates talking points for the next pricing committee meeting, and models revenue impact scenarios.
- **Network effects**: Aggregate anonymized pricing intelligence across all customers to create a "pricing weather map" for each vertical. "SaaS infrastructure pricing is deflationary this quarter. Average price per seat dropped 12% across the sector."
- **Reverse intelligence**: "Here's what your competitors probably know about YOUR pricing, based on what's publicly visible on your pricing page. Here are the gaps in your pricing communication."
- **Negotiation intelligence**: For enterprise deals, show how your pricing compares to what competitors are likely offering, based on their public tiers and typical enterprise discount patterns.

---

## 3. Premise Challenges

### Is this the right problem? Could a different framing yield a dramatically better product?

The current framing is "monitor competitor pricing pages." A potentially better framing:

**"Competitive pricing strategy co-pilot."** Instead of just watching competitors, the agent actively helps you make pricing decisions. It asks: "You're launching a new tier next month. Based on the competitive landscape I've been tracking, here's where to price it, here's the risk if you go too high, and here are the 3 competitors most likely to respond." This shifts from surveillance (reactive) to strategy (proactive). The monitoring is just the data collection layer; the real product is the strategic reasoning.

This reframing could be more compelling to judges because it demonstrates deeper agent reasoning. The scraping/monitoring becomes a means to an end, not the product itself.

**However**, for a hackathon demo, the surveillance framing is actually better because it has a clearer "whoa" moment (live change detection). The strategy framing requires more setup to demonstrate. **Recommendation: keep the monitoring framing for the demo, but hint at the strategy framing as the vision.**

### What happens if we do nothing? Is the manual spreadsheet approach actually fine?

Honestly? For most companies, the manual spreadsheet approach works "well enough" -- and that's the real risk for this idea.

Here's why the status quo survives:
- Most SaaS companies do competitive pricing research 1-2 times per year (before a pricing change, during annual planning). The data being stale for months is tolerable because pricing decisions are infrequent.
- The people who do this research (product managers, pricing analysts) are expensive humans doing it manually, but it takes them 1-2 days, not months. The cost is real but not crushing.
- Competitive pricing is only one input to pricing decisions. Customer willingness-to-pay, unit economics, and strategic positioning matter more.

**Where "do nothing" breaks down:**
- In fast-moving markets (AI, where pricing changes weekly). Carmen's research on 125 AI startups is a perfect example -- that data was probably stale before she finished collecting it.
- At companies with many competitors (15+) where manual monitoring is genuinely infeasible.
- When a competitor makes a surprise move and your sales team finds out from a prospect. That moment of embarrassment is the emotional driver.

**Assessment**: This is a "vitamin, not painkiller" for most companies, but a genuine painkiller for companies in fast-moving, crowded markets (AI/ML, DevTools, cloud infrastructure). The hackathon pitch should emphasize the AI vertical specifically, where pricing changes are frequent and dramatic.

### What's the weakest assumption in this idea?

**"Pricing pages contain enough structured information to extract and normalize meaningfully."**

In reality:
- 55% of voice AI companies hide pricing behind "contact sales" (per the original doc's own research). For enterprise SaaS, this number is even higher. The agent simply cannot work for these competitors.
- Many pricing pages show a starting price ("from $X/month") without enough detail to normalize. What features are included? What are the usage limits? This information is often buried in FAQ pages, comparison tables on separate URLs, or documentation.
- Pricing pages are marketing collateral, not data feeds. They're designed to sell, not to be machine-readable. The "real" pricing (with enterprise discounts, volume commitments, etc.) is never on the page.

**Impact**: The agent works well for the subset of competitors with transparent, detailed pricing pages. For the rest, it produces incomplete data or a "contact sales" flag. This limits the product's usefulness to markets with high pricing transparency (developer tools, SMB SaaS). Enterprise-focused markets are largely opaque.

**Mitigation for demo**: Choose competitors with good pricing pages. Acknowledge the limitation honestly -- "the agent flags opaque competitors as a signal in itself."

---

## 4. Alternative Approaches

### Approach A: Minimal Viable (Fewest Moving Parts, Ships Fastest)

**Core insight**: Skip the scraper entirely. Use pre-fetched HTML.

- **Input**: User pastes 5-6 competitor pricing page URLs into a form. The backend fetches the HTML once via a simple HTTP GET (no Playwright, no JS rendering).
- **Processing**: Send raw HTML to Claude via TrueFoundry with a structured extraction prompt. Get back JSON with tiers, prices, features, and unit of measurement.
- **Storage**: Store extracted data in Ghost Postgres. One table: `snapshots(competitor_id, timestamp, extracted_json, raw_html)`.
- **Change detection**: On re-scan, compare new JSON to most recent snapshot for each competitor. Simple JSON diff.
- **Output**: A clean comparison table and, if changes detected, an alert with analysis.
- **Demo trick**: Pre-seed with a baseline snapshot from the morning. Run a live scan. Show the diff for the staging page you changed.

**Sponsors used**: Ghost (essential), TrueFoundry (essential). That's it. Add Aerospike and Overmind only if time permits.

**Pros**: Can build in 4-5 hours. Very reliable because there are only two moving parts (fetch HTML, call LLM). Easy to debug.
**Cons**: Only 2 sponsors. No vector search. Less impressive architecture. "Autonomous agent" story is thinner.
**Build time**: 5 hours, leaving 3 hours for polish and additional sponsors.

### Approach B: Ideal Architecture (Best Long-Term)

**Core insight**: Build a proper agent framework with memory, planning, and self-improvement.

- **Agent framework**: Use a proper agent loop (e.g., LangGraph or custom state machine). The agent plans its own scraping strategy, decides which pages need Playwright vs simple fetch, and learns from extraction failures.
- **Smart scraping**: Playwright for JS-rendered pages. Cache rendered HTML. Detect anti-bot measures and adapt (switch to cached version, try mobile user-agent, etc.).
- **Two-phase extraction**: Phase 1 (fast model via TrueFoundry): structured extraction of visible pricing. Phase 2 (strong model): normalization, cross-competitor comparison, strategic analysis.
- **Rich storage**: Ghost Postgres for structured data (schemas, snapshots, diffs, competitor metadata). Aerospike for vector embeddings of pricing structures (enabling "who prices like X?" semantic queries) and sub-ms cache for latest snapshots.
- **Historical intelligence**: Agent uses its own history to contextualize new findings. "This is the 3rd change from Competitor X in 6 months" is not a query -- it's the agent reasoning over its memory.
- **Observability**: Overmind wraps every LLM call. The agent can report its own cost per scan cycle.
- **Frontend**: Real-time dashboard with WebSocket updates during scan. Competitor landscape visualization. Historical trend charts.

**Sponsors used**: Ghost, TrueFoundry, Aerospike, Overmind (all 4).

**Pros**: Maximum judge impression. Deep sponsor integration. Genuine autonomous behavior. Strong technical story.
**Cons**: 8 hours is extremely tight. High risk of partial completion. Many integration points that can fail.
**Build time**: 10-12 hours realistically. In 8 hours, you'd need to cut corners.

### Approach C: Creative/Lateral (Unexpected Angle)

**Core insight**: Flip the product. Instead of monitoring competitors, help companies **design** their pricing by understanding the full competitive landscape.

**"Pricing Architect Agent"**: You tell the agent your product, your target market, and your rough cost structure. The agent:

1. Autonomously discovers your competitors (web search + LLM reasoning, not just a pre-defined list).
2. Scrapes and extracts their pricing.
3. Identifies pricing patterns in your market (most common # of tiers, typical feature gating, prevalent pricing model).
4. Generates 3 pricing strategy options for you, each with a competitive positioning rationale:
   - "Undercut" strategy: price below the median, capture volume.
   - "Premium" strategy: price above, justify with feature differentiation.
   - "Disrupt" strategy: use a different pricing model entirely (if everyone charges per-seat, propose per-usage).
5. For each strategy, shows a mock pricing page and a competitive comparison table showing where you'd land.

**Demo moment**: "I told the agent I'm building an AI code review tool. It found 8 competitors, analyzed their pricing, and designed three pricing strategies for me. Here's the one it recommends, and here's why." The agent generates a pricing page mock-up on the fly.

**Sponsors used**: Ghost (store research), TrueFoundry (multi-model routing for search, extraction, analysis, generation), Aerospike (vector search for competitor discovery), Overmind (cost tracking).

**Pros**: More creative, more differentiated from existing tools, stronger "agent reasoning" story. No existing tool does this. The demo output is tangible and visually impressive (a generated pricing page).
**Cons**: Harder to build reliably. Competitor discovery via web search is noisy. Generated pricing strategies might feel generic. Less dramatic "live change detection" moment.
**Build time**: 8-9 hours. Tight but feasible if competitor discovery uses a curated seed list.

---

## 5. Sharpened Recommendation

### The BEST Version for This Hackathon

**Hybrid of Approach A (build speed) + the best elements of B (demo impressiveness) + the reframing from Approach C (strategic, not just surveillance).**

Here is the sharpened version:

**Name: "Pricing Radar"**

**Pitch**: "An autonomous agent that watches your competitors' pricing, understands what changes mean, and tells you what to do about it -- before your sales team hears it from a prospect."

**Key changes from v1:**

1. **Simplify the scraping layer ruthlessly.** Use httpx only (no Playwright). Pre-select 6 competitors with clean, static pricing pages. Have fallback cached HTML for every page. The scraping is NOT the interesting part -- the intelligence is. Do not let scraping fragility eat your build hours.

2. **Double down on the normalization and analysis, because that's the moat.** The extraction prompt is the product. Spend 2 full hours on prompt engineering for structured extraction. Use a specific normalization frame: "monthly cost for a 50-person team using 1M tokens/month" as the universal comparison unit. This is opinionated and imperfect, but it's concrete and demo-able.

3. **Add the "Pricing Architect" twist from Approach C as a bonus feature.** After showing change detection, show the agent generating a strategic recommendation: "Based on the competitive landscape, here are 3 ways you could reposition your pricing." This elevates the demo from "monitoring tool" to "strategic advisor" and demonstrates deeper agent reasoning.

4. **Lean into the Tierly differentiation.** Tierly exists and does point-in-time analysis. Your differentiator is continuous monitoring + change detection + strategic recommendations. Make this explicit in the pitch: "Tierly does pricing analysis. We do pricing surveillance. They're a report. We're a radar."

5. **Pre-seed aggressively.** Before the demo, run 3-4 extraction cycles over the previous days to build up genuine historical data in Ghost. Show real trend data, not simulated. This makes the "memory" story authentic.

6. **Nail the meta-moment.** Use Overmind to show: "This scan of 6 competitors cost $0.04 and took 38 seconds. A human analyst doing the same work takes 2 days and costs $800." This is the killer closing line.

### Revised Hour-by-Hour Plan

| Hour | Focus | Deliverable |
|------|-------|-------------|
| 1 | Ghost DB schema + scraper for 6 pages + staging page | Working fetch + store pipeline |
| 2 | Extraction prompt engineering (THE critical hour) | Reliable structured extraction for all 6 pages |
| 3 | TrueFoundry integration + model routing | Fast model for extraction, strong model for analysis |
| 4 | Change detection + analysis generation | End-to-end loop working |
| 5 | Aerospike vector search ("who prices like us?") | Similarity query working |
| 6 | Frontend/output layer (clean CLI or minimal web UI) | Demo-ready output |
| 7 | Overmind + "Pricing Architect" recommendation feature | Cost tracking + strategic recommendations |
| 8 | Pre-seed historical data, rehearse demo, build fallbacks | Bulletproof demo |

### Revised Score

| Dimension | v1 | v2 | Change | Justification |
|-----------|-----|-----|--------|---------------|
| **Autonomy** (3x) | 5 | 5 | -- | Unchanged. The loop is strong: trigger, scrape, extract, normalize, diff, analyze, recommend, store. |
| **Idea** (3x) | 5 | 5 | -- | Unchanged. Problem is real, relatable, and validated by research (Carmen's 125 pages, Tierly's existence proves market). |
| **Technical Implementation** (2x) | 4 | 4 | -- | Simplified scraping (httpx only) trades robustness for reliability. Better prompt engineering for extraction. Net wash. |
| **Tool Use** (2x) | 4 | 5 | +1 | Clearer story for all 4 sponsors. Ghost is the memory. TrueFoundry is the brain router. Aerospike answers "who prices like us?" Overmind is the cost x-ray. The meta-moment (pricing agent that knows its own cost) ties Overmind in thematically. Stronger narrative = better judge impression. |
| **Presentation/Demo** (2x) | 5 | 5 | -- | Unchanged. Live change detection remains the anchor. Adding the "Pricing Architect" recommendation is bonus, not required. |
| **Buildability** (1x) | 3 | 4 | +1 | Dropping Playwright eliminates the biggest fragility risk. Pre-selecting known-good pages. Pre-seeding historical data. Clearer hour-by-hour plan with explicit fallback points. Still tight but higher confidence. |

**New Total: (5x3) + (5x3) + (4x2) + (5x2) + (5x2) + (4x1) = 15 + 15 + 8 + 10 + 10 + 4 = 62/65**

### What Changed: +3 Points

- **Tool Use 4 -> 5**: The Overmind meta-moment ("a pricing intelligence agent that reports its own economics") is now a first-class demo beat, not an afterthought. Aerospike's role is crisper (semantic similarity for "who prices like us?" is a clear, demo-able feature). The narrative connecting all 4 sponsors is tighter.
- **Buildability 3 -> 4**: Dropping Playwright, pre-selecting pages, and pre-seeding data removes the three biggest risk factors. The hour-by-hour plan has clearer milestones and cut points.

### LinkedIn Demand Validation (from primary research)

**The manual work is confirmed real:**
- Carmen Insignares Newell (scored 95): Manually analyzed 125 AI startup pricing pages. Says "DM me for the detailed analysis sheet" — the output itself has demand.
- Koshima Satija (scored 92): Analyzed 200+ Voice AI companies over 3 months. Found 12 distinct pricing models across 4 stack layers. Only 45% have public pricing. The normalization problem is exactly what Pricing Radar solves.
- Brennan Plaetzer (scored 75): "Reviewed 20+ AI startups in the last 30 days."

**Pricing velocity is accelerating:**
- Wisam Hirzalla/Stripe (scored 98): Companies with 100%+ growth change pricing 3x more frequently. Monitoring window is shrinking from quarterly to weekly.
- Fynn Glover (scored 95): "A credit cost. A usage limit. An overage threshold. These are weekly operational decisions."
- Glenn Turner (scored 94): Figma's pricing change wiped $4B+ in market value. A single competitor pricing change reshapes the landscape.

**Three validated risks to address:**
1. **55% hide pricing behind "contact sales"** (Koshima's data). Coverage skews toward self-serve/SMB. The most interesting competitors are invisible. **MITIGATION: Flip the gap into a feature.** "Competitor A: public, last changed March 15" vs "Competitor B: contact sales, no public pricing." The absence of data is data. A transparency score IS competitive intelligence.
2. **Normalization requires usage assumptions.** 12 distinct pricing models in voice AI alone. Per-seat vs per-token vs per-minute can't be compared without a reference workload. **MITIGATION: Make assumptions explicit.** "Normalized to: 50-person team, 1M tokens/month." Let users change the assumptions.
3. **Public pricing is the least interesting signal.** Biggest shifts happen in sales negotiations. Competitor X keeps $49/seat on website but offers $25/seat in every enterprise deal. **MITIGATION: Scope honestly.** This monitors the public signal. The strategic value is speed + normalization + change detection, not omniscience.

### Key Risks Remaining

1. **Extraction prompt reliability** (MEDIUM): This is now the #1 technical risk. If the LLM can't reliably extract structured data from diverse HTML, the whole product falls apart. Mitigation: spend a full hour on prompt engineering with few-shot examples from each target page.
2. **Normalization quality** (MEDIUM): "Monthly cost for a 50-person team using 1M tokens/month" is a useful frame but requires assumptions the LLM must state explicitly. A sharp judge will probe these assumptions.
3. **Tierly awareness** (LOW): If a judge knows about Tierly, you need the "radar vs telescope" differentiation ready. Tierly does analysis; you do surveillance. This is a strength, not a weakness -- Tierly's existence validates the market.
4. **Coverage ceiling** (MEDIUM): 55% behind "contact sales" means you're monitoring the minority. Frame this as "transparency radar" — knowing who's hiding is itself intelligence.

### Final Verdict

**BUILD THIS. Score: 62/65.** The sharpened version eliminates the biggest buildability risks (Playwright, scraping fragility) while strengthening the demo narrative (4-sponsor story, meta-cost moment, "Pricing Architect" bonus feature). The core insight holds: every SaaS company does competitive pricing research manually, nobody does it continuously, and the normalization problem (per-seat vs per-token vs per-credit) is a genuine reasoning challenge that showcases agent intelligence.

The idea occupies a sweet spot in the competitive landscape: e-commerce price scrapers can't handle SaaS pricing complexity, website change monitors can't reason about what changes mean, and Tierly does analysis but not continuous monitoring. This agent fills a real gap.

---

## Sources

- [Tierly - Competitive Pricing Intelligence for SaaS](https://tierly.app/)
- [How Tierly builds AI-powered pricing intelligence using Trigger.dev](https://trigger.dev/customers/tierly-customer-story)
- [Top Competitor Price Tracking Tools (2026) - Visualping](https://visualping.io/blog/top-tools-competitor-price-tracking)
- [Competitive Pricing Tools: Complete Guide for SaaS (2026) - Tierly](https://tierly.app/blog/competitive-pricing-tools)
- [10 Best Competitive Pricing Tools - Orb](https://www.withorb.com/blog/competitive-pricing-tools)
- [Prisync - Competitor Price Tracking](https://prisync.com/)
- [Omnia Price Monitoring Software](https://www.omniaretail.com/price-monitoring-software)
- [Pricefx - AI Pricing Software](https://www.pricefx.com/)
- [Price Scraping 2026 Guide - Tendem](https://tendem.ai/blog/price-scraping-competitor-monitoring-guide)
- [How to Build Automated Competitor Price Monitoring - Firecrawl](https://www.firecrawl.dev/blog/automated-competitor-price-scraping)
- [LLM API Pricing 2026 - Compare 300+ AI Model Costs](https://pricepertoken.com/)
- [AI Pricing Models Explained - Data-Mania](https://www.data-mania.com/blog/ai-pricing-models-explained-usage-seats-credits-outcome-based-options/)
- [Selling Intelligence: The 2026 Playbook for Pricing AI Agents - Chargebee](https://www.chargebee.com/blog/pricing-ai-agents-playbook/)
- [The AI Pricing and Monetization Playbook - Bessemer Venture Partners](https://www.bvp.com/atlas/the-ai-pricing-and-monetization-playbook)
- [Best Competitive Pricing Analysis Software for SaaS (2026) - Tierly](https://tierly.app/blog/competitive-pricing-analysis-software)
- [Monitoring the Pricing Model of Your SaaS Competitors - Stillio](https://www.stillio.com/blog/monitoring-saas-competitors-pricing)
- [Competitor Pricing Monitoring - Visualping](https://visualping.io/blog/competitor-pricing-change-alerts)
- [15 Best AI Pricing Software Reviewed in 2026 - CRO Club](https://croclub.com/tools/best-ai-pricing-software/)
