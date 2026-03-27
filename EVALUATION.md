# Hackathon Idea Evaluation Framework

Deep Agents Hackathon, RSAC 2026. March 27, 2026.
Constraint: use 3+ sponsor tools. ~8 hours build time. Live demo + pitch (3 min).

---

## Official Judging Criteria

These are the actual categories judges will score on. Our evaluation mirrors them exactly.

---

## Scoring Dimensions (each 1-5)

### 1. Autonomy (weight: 3x)
How well does the agent act on real-time data without manual intervention?
- **5** — Agent runs end-to-end with zero human intervention. Reacts to live data, makes decisions, takes actions autonomously. You press "go" and watch it work.
- **4** — Agent handles the core loop autonomously, occasional human confirmation for edge cases
- **3** — Agent automates several steps but needs human nudges to keep moving
- **2** — Agent assists a human workflow but doesn't drive it
- **1** — Human does the work, agent is a fancy autocomplete

### 2. Idea (weight: 3x)
Does the solution solve a meaningful problem or demonstrate real-world value?
- **5** — You can name the person who has this problem today and describe their painful workaround. Judges immediately see why this matters.
- **4** — Clear pain point, you've experienced it yourself, easy to explain
- **3** — Reasonable problem but somewhat hypothetical
- **2** — "It would be nice if..." territory
- **1** — Solution looking for a problem

### 3. Technical Implementation (weight: 2x)
How well was the solution implemented?
- **5** — Clean architecture, real agent reasoning (not just chained API calls), handles errors gracefully, non-obvious data flow
- **4** — Solid multi-step agent with real decision-making, well-structured code
- **3** — Competent integration work, straightforward agent loop, works reliably
- **2** — Fragile, demo-path-only, breaks on edge cases
- **1** — Barely wired together, copy-paste from docs

### 4. Tool Use (weight: 2x)
Did the solution effectively use at least 3 sponsor tools?
- **5** — 3+ sponsors, each is essential and load-bearing. Remove one and the product breaks. Sponsors used in ways that show you understand their strengths.
- **4** — 3+ sponsors, most are essential, one is nice-to-have
- **3** — 3 sponsors used but one or two feel cosmetic / interchangeable
- **2** — Sponsors are used but don't add real value
- **1** — Sponsor logos in the README, that's about it

### 5. Presentation / Demo (weight: 2x)
3-minute live demonstration of the solution.
- **5** — Judges say "whoa." Visible, dramatic before/after. Story is one sentence. Demo flows perfectly.
- **4** — Clear value, needs a sentence of setup, demo runs clean
- **3** — Interesting but requires explanation to land, minor demo hiccups
- **2** — Technical achievement that's hard to show, or demo breaks
- **1** — "Trust me, it works" ... slides instead of live demo

### 6. Buildability (weight: 1x, internal only)
NOT a judging criterion, but critical for us. Can we actually ship this in 8 hours?
- **5** — Core loop works in 2-3 hours, rest is polish and edge cases
- **4** — Tight but doable if nothing goes sideways
- **3** — Requires everything to go right, no room for debugging
- **2** — Probably need to cut major features to demo anything
- **1** — This is a weekend project minimum

---

## Scoring Formula

```
Total = (Autonomy x 3) + (Idea x 3) + (TechImpl x 2) + (ToolUse x 2) + (Demo x 2) + (Buildability x 1)
```

**Max score: 65**

| Range  | Verdict                                      |
|--------|----------------------------------------------|
| 52-65  | BUILD THIS. Strong across the board.         |
| 39-51  | Contender. Check which dimension drags and see if fixable. |
| 26-38  | Risky. Likely weak on autonomy or idea.      |
| < 26   | Pass. Find a different idea.                 |

---

## What Wins This Hackathon

Based on the criteria weights, the winning formula is:
1. **High autonomy** — the agent does real work without hand-holding. This is weighted equal to the idea itself. Build something that runs on its own.
2. **Real problem** — not a toy demo. Judges want to see real-world value.
3. **Sponsor tools are load-bearing** — not decorative. Each one should be doing something the product can't do without.
4. **Demo tells a story in 3 minutes** — set up the problem (30 sec), show the agent working (2 min), land the impact (30 sec).

---

## Sponsor Quick-Reference (for Tool Use scoring)

| Sponsor       | Best fit when your agent needs to...                         |
|---------------|--------------------------------------------------------------|
| Auth0/Okta    | Act on behalf of a user, access third-party APIs             |
| Bland AI      | Talk to humans on the phone                                  |
| Airbyte       | Pull structured data from SaaS tools (GitHub, Stripe, Gong)  |
| Aerospike     | Fast storage, vector search, agent memory                    |
| TrueFoundry   | Route across multiple LLMs, deploy with observability        |
| Overmind      | LLM optimization, observability, cost reduction              |
| Macroscope    | Understand and analyze codebases                             |
| Ghost         | Spin up / fork / query Postgres databases autonomously       |
| Kiro          | Agentic IDE, code generation                                 |
| Senso         | TBD                                                          |

---

## Ideas — Triage Scores

Formula: (Autonomy x3) + (Idea x3) + (Tech x2) + (Tools x2) + (Demo x2) + (Build x1). Max 65.

| # | Idea | Auto | Idea | Tech | Tools | Demo | Build | Total | Verdict |
|---|------|------|------|------|-------|------|-------|-------|---------|
| 1 | Agent commerce rails | 4 | 4 | 3 | 3 | 2 | 2 | 40 | Protocol problem. Hard to demo. Infrastructure, not a product. |
| 2 | Pricing experimentation | 3 | 4 | 3 | 3 | 3 | 2 | 39 | A/B test infra is complex. 8 hours isn't enough. |
| 3 | AI spend forecasting | 4 | 5 | 3 | 3 | 3 | 3 | 48 | Real problem, but dashboards don't demo well. Moderate sponsor fit. |
| 4 | Model router | 5 | 3 | 3 | 3 | 4 | 4 | 47 | High autonomy, but TrueFoundry already does this. What's new? |
| 5 | Pricing intelligence | 5 | 5 | 4 | 4 | 5 | 3 | 57 | **TOP.** Agent scrapes, normalizes, alerts autonomously. Great demo. Scraping risk. |
| 6 | Contract builder | 2 | 4 | 3 | 2 | 2 | 3 | 34 | Low autonomy. Contracts need humans. Weak sponsor fit. |
| 7 | Pricing change comms | 4 | 5 | 4 | 4 | 5 | 3 | 56 | **TOP.** Bland AI calling at-risk customers live = "whoa" moment. |
| 8 | Margin-aware feature flags | 4 | 5 | 4 | 3 | 5 | 3 | 53 | **TOP.** $18K surprise bill story is visceral. Strong demo. |
| 9 | War-game simulator | 3 | 4 | 4 | 3 | 4 | 2 | 43 | Monte Carlo + real data in 8 hrs is tight. Human-driven. |
| 10 | Cost calculator | 2 | 3 | 2 | 2 | 3 | 5 | 34 | Widget, not an agent. Low autonomy. Easy build but weak everywhere else. |

### Triage Result

**Advance to deep exploration:**
- **#5 Pricing Intelligence** (57) — highest score, strong autonomy + demo
- **#7 Pricing Change Comms** (56) — Bland AI phone call is a killer demo moment
- **#8 Margin-Aware Feature Flags** (53) — visceral problem, clear before/after
- **#3 AI Spend Forecasting** (48) — real problem, borderline. Worth a look for sponsor fit.

**Cut:**
- #4 (47), #9 (45), #1 (40), #2 (39), #6 (34), #10 (34)
