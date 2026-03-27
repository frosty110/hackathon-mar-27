# Idea #3: AI Spend Forecasting Agent

## One-Line Pitch

An autonomous agent that ingests your LLM usage data, forecasts next month's bill per team, and enforces token budgets in real time before you get a surprise invoice.

## The Problem

**Who has it:** Engineering leaders and FinOps teams at companies spending $50K-$500K+/month on LLM inference. Henry Norris has $500K engineers burning $250K in tokens with zero team-level attribution. Trey Harnden is at 15M tokens/day projecting 50M next year with no forecasting capability. Matias Coca can't even find AI costs because they're buried under generic cloud service names.

**Current workaround:** Manually pulling billing CSVs at month-end, building one-off spreadsheets, setting Slack alerts when the total bill crosses a threshold. Some teams just... don't track it. They find out when the invoice arrives.

**Why it sucks:** It's entirely reactive. By the time you see the bill, the money is spent. There's no per-team attribution, no forecasting, no automated enforcement. It's like running a company with no departmental budgets -- just one big credit card and a prayer.

## The Autonomous Agent Loop

The agent runs a continuous loop with three phases:

**Phase 1 -- Ingest & Attribute (triggered on schedule or webhook)**
- Pulls token usage logs from LLM providers via Airbyte connectors (OpenAI, Anthropic, Azure OpenAI billing APIs)
- Pulls team/project metadata from GitHub (who owns which repos, which repos call which models)
- Joins usage data to teams/projects using API key mappings and repo ownership
- Stores all raw and attributed usage data in Aerospike for sub-ms retrieval

**Phase 2 -- Forecast & Alert (triggered after each ingestion cycle)**
- Runs time-series forecasting on per-team token consumption (exponential smoothing or simple regression -- keep it buildable in 8 hours)
- Compares forecasted spend against per-team budgets stored in Ghost Postgres
- If a team is projected to exceed their monthly budget: autonomously decides severity (warning at 80%, critical at 95%, hard-stop at 100%)
- Routes all LLM calls through TrueFoundry gateway, which the agent can configure to throttle or block teams exceeding budget

**Phase 3 -- Notify & Enforce (triggered by threshold breach)**
- Posts Slack-style alerts to a dashboard (or calls the team lead via Bland AI for critical overages)
- Adjusts TrueFoundry gateway routing to downgrade expensive models to cheaper ones (e.g., GPT-4 -> GPT-3.5) for teams approaching limits
- Logs all autonomous decisions with reasoning to Ghost Postgres for audit trail

**What makes it autonomous:** The agent doesn't just report -- it acts. It decides when to warn vs. throttle vs. block. It decides which model to downgrade to. It adjusts enforcement in real time without a human approving each action.

## Sponsor Stack (3+ required)

### 1. TrueFoundry -- Multi-LLM Gateway (ESSENTIAL)
- **What it does:** Acts as the central proxy through which all LLM calls flow. Provides the unified usage log (tokens, costs, latency per request). The agent configures routing rules on the gateway to enforce budgets -- downgrading models or blocking calls when teams exceed limits.
- **Why it's load-bearing:** Without TrueFoundry, you have no single point of control for enforcement. You'd need to instrument every individual LLM provider separately, and you'd have no way to dynamically reroute calls to cheaper models. The entire "enforcement" half of the product disappears.
- **Essential or cosmetic:** Essential. This is the product's control plane.

### 2. Airbyte -- Data Connectors (ESSENTIAL)
- **What it does:** Pulls billing/usage data from LLM providers and pulls team/project metadata from GitHub. The Python agent connector packages let the agent programmatically trigger syncs and pull data without writing custom API integrations for each provider.
- **Why it's load-bearing:** Without Airbyte, you're writing bespoke API clients for OpenAI, Anthropic, Azure, GitHub, etc. during an 8-hour hackathon. You'd spend all your time on plumbing instead of the forecasting and enforcement logic.
- **Essential or cosmetic:** Essential. It's the data ingestion layer.

### 3. Ghost -- Agent-Managed Postgres (ESSENTIAL)
- **What it does:** Stores team budgets, budget policies, forecast results, enforcement action logs, and audit trails. The agent creates and queries databases via MCP -- it can spin up new tables for new teams, fork the DB to test policy changes before applying them.
- **Why it's load-bearing:** Without Ghost, you need to manually provision and manage a database. The MCP interface means the agent can autonomously create schema, run migrations, and query data as part of its loop -- no human DBA needed.
- **Essential or cosmetic:** Essential. It's the persistent state store and audit log.

### 4. Aerospike -- Sub-ms Usage Cache (STRONG BUT NOT CRITICAL)
- **What it does:** Caches the most recent usage data and per-team running totals for sub-millisecond lookups. When the TrueFoundry gateway needs to check "is this team over budget?" before routing a call, it hits Aerospike instead of running a Postgres query.
- **Why it's load-bearing:** Real-time budget enforcement requires low-latency lookups on every LLM call. Postgres can't do this at scale. Aerospike makes the enforcement path fast enough to be practical.
- **Essential or cosmetic:** Genuinely useful but you could demo without it. At hackathon scale (not millions of requests), Ghost Postgres alone could handle it. Moves from "essential" to "strong nice-to-have" at demo scale. But the architectural argument for it is real and defensible.

### 5. Bland AI -- Voice Alerts for Critical Overages (COSMETIC BUT DEMO GOLD)
- **What it does:** When a team is about to blow through their budget and the agent escalates to "critical," it calls the team lead's phone and delivers a voice briefing: "Your team has consumed 94% of your March token budget with 8 days remaining. The agent has downgraded your routing from GPT-4 to GPT-3.5. Call back to override."
- **Why it's load-bearing:** It isn't, strictly. You could replace this with a Slack message. But for a 3-minute demo, having the judge's phone ring live is an unforgettable moment.
- **Essential or cosmetic:** Cosmetic for the product. Essential for the demo score.

## The "Whoa" Demo Moment

Midway through the demo, the agent detects that "Team Payments" has just crossed 90% of their monthly budget with 10 days left in the month. The dashboard updates in real time showing the forecast line crossing the budget ceiling. The agent autonomously:

1. Downgrades Team Payments' model routing from Claude Opus to Claude Haiku via TrueFoundry (visible in the gateway config updating live)
2. Logs the enforcement action with reasoning to Ghost Postgres (visible in audit log)
3. Calls the presenter's actual phone via Bland AI -- the phone rings audibly in the room, and the voice says: "Alert: Team Payments has hit 90% of their March AI budget. Forecasted overage: $12,400. The agent has downgraded your model tier. Reply to override."

The phone ringing in a quiet room during a demo is visceral. Judges remember it.

## 3-Minute Demo Script

**0:00 - 0:20 | The Hook**
"Last month, three of our interviewees told us the same story: they're spending $250K a month on AI tokens and they have no idea which team is driving the bill. One of them only finds out when the invoice arrives. We built an agent that doesn't just tell you what you spent -- it predicts what you'll spend and enforces budgets before the money is gone."

**0:20 - 0:50 | The Dashboard**
Show the live dashboard: 5 teams, each with a token budget, a usage bar, and a forecast line. Point out that Team Payments is at 87% with 10 days left. The forecast line clearly crosses the budget ceiling. "The agent ingested usage data from our LLM gateway via Airbyte, attributed it to teams using GitHub repo ownership, and ran a forecast. It's already watching."

**0:50 - 1:40 | The Trigger**
Simulate a burst of API calls from Team Payments (pre-scripted load generator hitting TrueFoundry gateway). The dashboard updates live -- usage climbs to 91%. The agent kicks in:
- Dashboard shows: "ENFORCEMENT ACTION: Team Payments downgraded from claude-opus to claude-haiku"
- Show the TrueFoundry gateway config updating in real time
- Show the Ghost Postgres audit log entry: timestamp, team, action, reasoning, projected savings

**1:40 - 2:10 | The Phone Call**
The presenter's phone rings. Put it on speaker. Bland AI voice delivers the budget alert. Pause for audience reaction. "The agent decided this was critical enough to escalate beyond a dashboard notification. It called me."

**2:10 - 2:40 | The Architecture**
Quick slide: Airbyte pulls data -> Agent forecasts and decides -> TrueFoundry enforces -> Ghost logs -> Aerospike caches for real-time checks -> Bland AI escalates. "Five sponsor tools, all load-bearing."

**2:40 - 3:00 | The Landing**
"Every company adopting AI is about to have this problem. The agent doesn't wait for the invoice. It sees the future and acts. FinOps for the AI era."

## Technical Architecture

```
                    LLM Consumers (Teams)
                           |
                           v
                  +------------------+
                  | TrueFoundry      |  <-- All LLM calls routed here
                  | Multi-LLM Gateway|  <-- Agent configures routing rules
                  +------------------+
                     |           |
            Usage Logs      Model Routing
                     |           ^
                     v           |
              +-------------+    |
              | Airbyte     |    |  Enforcement
              | Connectors  |    |  Actions
              +-------------+    |
                     |           |
                     v           |
              +-------------+    |
              | AGENT CORE  |----+
              | (Python)    |
              |             |-----> Bland AI (voice alerts)
              | - Ingest    |
              | - Attribute |-----> Ghost Postgres (budgets,
              | - Forecast  |       audit log, policies)
              | - Enforce   |
              | - Notify    |-----> Aerospike (real-time
              +-------------+       usage cache, running totals)
```

**Data Flow:**
1. All LLM calls from all teams route through TrueFoundry gateway
2. Airbyte syncs usage logs from TrueFoundry + team metadata from GitHub into the agent
3. Agent writes raw usage to Aerospike (fast reads) and structured data to Ghost Postgres (durable store)
4. Agent runs forecasting logic per team, compares to budgets in Ghost Postgres
5. If threshold breached: agent updates TrueFoundry routing rules (enforcement) and optionally triggers Bland AI call
6. All decisions logged to Ghost Postgres audit table

**Key Technical Decisions:**
- Forecasting model: Simple exponential moving average or linear regression. Don't over-engineer this for a hackathon. The demo needs to show the forecast line, not win a Kaggle competition.
- Budget policies: Stored as JSON in Ghost Postgres. Agent reads and interprets them autonomously.
- Pre-scripted load generator: A simple script that fires API calls through TrueFoundry to simulate team usage spikes during the demo.

## Buildability Risk Assessment

**Hardest parts (in order of risk):**

1. **TrueFoundry integration depth (HIGH RISK):** Can the agent actually modify routing rules programmatically via API? If TrueFoundry's gateway config is static or requires manual UI changes, the enforcement story collapses. Need to verify API access to routing/model selection rules on day one, hour one.

2. **Airbyte connector availability (MEDIUM RISK):** Do Airbyte connectors exist for LLM provider billing APIs (OpenAI usage, Anthropic usage)? If not, you're writing custom API clients, which eats 2-3 hours. Fallback: mock the data ingestion and focus on the forecasting + enforcement loop.

3. **End-to-end integration (MEDIUM RISK):** Five sponsor tools means five integration points. Each one can have auth issues, API quirks, or documentation gaps. Budget 1 hour per integration = 5 hours just on plumbing, leaving 3 hours for core logic and demo polish.

4. **Demo reliability (MEDIUM RISK):** The live phone call via Bland AI is high-reward but if it fails during the demo (network issue, API timeout), you lose your best moment. Mitigation: have a pre-recorded backup video of the call working.

**What to cut if time runs out:**
- Cut Aerospike first (use Ghost Postgres for everything)
- Cut live Airbyte sync (pre-load data, show Airbyte config as "this is how it connects")
- Cut real-time dashboard updates (show before/after states instead of live animation)
- NEVER cut: TrueFoundry enforcement + Bland AI phone call. Those are the demo.

## Honest Weaknesses

**1. Autonomy depth is debatable.** The agent's "decisions" are really threshold-based rules: if usage > 80%, warn; if > 95%, downgrade. A skeptical judge might say "this is just a cron job with if-statements, not an autonomous agent." Counter-argument: the agent decides WHICH model to downgrade to, forecasts WHEN the budget will be hit, and chooses the escalation path. But the counter-counter is: those are still deterministic rules, not LLM reasoning.

**2. The forecasting is simple.** Linear regression on a few weeks of data isn't impressive ML. The value is in the system -- ingestion, attribution, enforcement -- not the forecasting model itself. Judges who care about technical depth in the AI/ML sense may be underwhelmed.

**3. The problem space is operational, not exciting.** "Budget enforcement" doesn't have the visceral appeal of "catches zero-day exploits" or "replaces your SOC analyst." It's important but unsexy. The Bland AI phone call compensates for this, but the underlying product is a cost dashboard with enforcement -- which FinOps vendors already offer for cloud spend (just not AI-specific).

**4. Simulated data risk.** Unless you actually have 5 teams burning tokens through TrueFoundry during the demo, you're showing pre-loaded or simulated data. Judges may notice and discount the "real-time" claim. Mitigation: the load generator script making real API calls through TrueFoundry during the demo makes it genuinely live, even if the "teams" are synthetic.

**5. Five integrations in 8 hours is ambitious.** Each integration is a potential time sink. If TrueFoundry or Airbyte has poor documentation or auth issues, you could spend half the hackathon debugging plumbing instead of building the product.

## Final Score Recommendation

| Dimension | Score | Justification |
|-----------|-------|---------------|
| **Autonomy** | 7/10 | The agent ingests, forecasts, decides, and enforces without human intervention. But the decisions are mostly threshold-based rules, not rich LLM reasoning. A truly autonomous agent would reason about *why* a team's usage spiked and suggest architectural changes, not just throttle them. The loop is real but shallow. |
| **Idea** | 7/10 | The problem is validated by real interviews and genuinely painful. But "FinOps for AI" is an obvious category -- multiple startups are already here (Helicone, LangSmith, etc.). The enforcement angle (not just monitoring) is the differentiator, but it's an incremental step, not a paradigm shift. |
| **Technical Implementation** | 6/10 | Five integrations in 8 hours is risky. The forecasting model is deliberately simple. The core logic (threshold checks, routing updates) isn't technically complex. The architecture is sound but the individual pieces are shallow. High risk of spending all time on integration plumbing. |
| **Tool Use** | 8/10 | Strong sponsor integration. TrueFoundry, Airbyte, and Ghost are genuinely load-bearing. Aerospike is architecturally justified even if not strictly necessary at demo scale. Bland AI is cosmetic but brilliantly deployed for demo impact. Five tools, three essential -- that's solid. |
| **Presentation** | 8/10 | The phone call moment is a winner. The dashboard with live-updating forecast lines is visually clear. The "before the invoice arrives" narrative is crisp. The demo script is well-structured. Risk: if the phone call fails, you lose your best moment. |
| **Overall Viability** | 7/10 | A solid, buildable idea with a clear demo narrative and strong sponsor usage. Not the most technically ambitious or creative entry, but reliable and well-structured. The main risk is integration time and the autonomy depth feeling thin under scrutiny. |

**Composite: ~72/100**

**Bottom line:** This is a "safe B+" idea. It's buildable, demo-able, and uses sponsors well. But it won't win unless the execution is flawless and the phone call moment lands perfectly. The autonomy story needs to be deeper than threshold-based rules -- consider having the agent use an LLM to reason about usage patterns and generate natural-language explanations for its enforcement decisions. That small addition would push the autonomy score from 7 to 8-9 and could be the difference between "nice project" and "winner."
