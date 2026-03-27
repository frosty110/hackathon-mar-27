# Idea Comparison — Deep Agents Hackathon

All 13 ideas ranked by score. Formula: (Autonomy x3) + (Idea x3) + (Tech x2) + (Tools x2) + (Demo x2) + (Build x1). Max 65.

---

## Final Rankings

| Rank | # | Idea | Auto | Idea | Tech | Tools | Demo | Build | Total | Tier |
|------|---|------|------|------|------|-------|------|-------|-------|------|
| 1 | 5 | Pricing Intelligence | 5 | 5 | 4 | 4 | 5 | 3 | **59** | BUILD THIS |
| 2 | 13 | Deal Intelligence | 5 | 4 | 4 | 5 | 5 | 3 | **58** | BUILD THIS |
| 3 | 11b | Churn Rescue | 5 | 4 | 4 | 5 | 5 | 3 | **58** | BUILD THIS |
| 4 | 12b | Revenue Anomaly | 4 | 4 | 4 | 5 | 5 | 3 | **55** | BUILD THIS |
| 5 | 7 | Pricing Change Comms | 4 | 4 | 3.5 | 4.5 | 4.5 | 3 | **52** | Contender |
| 8 | 8 | Margin-Aware Flags | 3.5 | 4.5 | 3.5 | 4 | 4.5 | 3 | **50** | Contender |
| 9 | 3 | Spend Forecasting | 3.5 | 3.5 | 3 | 4 | 4 | 3.5 | **46** | Contender |
| 10 | 4 | Model Router | 5 | 3 | 3 | 3 | 4 | 4 | **47** | Contender |
| 11 | 9 | War-Game Simulator | 3 | 4 | 4 | 3 | 4 | 2 | **45** | Risky |
| 12 | 1 | Agent Commerce Rails | 4 | 4 | 3 | 3 | 2 | 2 | **40** | Risky |
| 13 | 2 | Pricing Experimentation | 3 | 4 | 3 | 3 | 3 | 2 | **39** | Risky |
| — | 6 | Contract Builder | 2 | 4 | 3 | 2 | 2 | 3 | **34** | Pass |
| — | 10 | Cost Calculator | 2 | 3 | 2 | 2 | 3 | 5 | **34** | Pass |

---

## Top 6 — Detailed Breakdown

---

### #1: Idea #5 — Competitive Pricing Intelligence Agent (59/65)

**Pitch:** An agent that autonomously scrapes competitor pricing pages, normalizes across pricing models, detects changes, and delivers strategic analysis.

**Pros:**
- Universally relatable problem. Carmen manually analyzed 125 pricing pages. Everyone does this in spreadsheets.
- Strongest autonomy loop: scrape, extract, normalize, compare, detect changes, analyze, alert. Zero human intervention.
- Demo moment: change a competitor's price on a staging page, agent detects it in 30 seconds with strategic analysis. Visible, dramatic.

**Cons:**
- Web scraping fragility. Pricing pages vary wildly. Could break during live demo.
- Normalization across pricing models (per-seat vs per-token vs per-credit) is genuinely hard and only approximate in 8 hours.
- Sponsor tools are good but not exceptional. Ghost and TrueFoundry are load-bearing. Aerospike is strong but not strictly essential. Overmind is cosmetic.

**Biggest risk:** Scraping breaks during demo. Mitigate with pre-selected known-good pages + cached HTML fallback + staged change-detection page.

---

### #2 (tie): Idea #13 — Deal Intelligence Agent (58/65)

**Pitch:** New deal enters pipeline. Agent pulls from Gong, Stripe, and GitHub, finds similar won/lost deals, and calls the rep with specific coaching when a deal is at risk.

**Pros:**
- Best Airbyte story across all ideas. All 3 available connectors (Gong, Stripe, GitHub) are naturally load-bearing in one product. This is Airbyte's pitch deck come to life.
- Rich autonomy: trigger, multi-source pull, vector search, risk analysis, strategy generation, phone call. More decision branches than most ideas.
- Phone call with data-derived coaching ("your champion hasn't attended since January") is arguably the single best demo moment.

**Cons:**
- Sales/CRM domain may not resonate with the RSAC-adjacent audience.
- Gong requires enterprise account. Data will almost certainly be synthetic, which weakens the autonomy narrative.
- Six sponsor tools risks looking like sponsor-stacking. Ghost and Overmind are defensible but not essential.

**Biggest risk:** Gong connector access. If unavailable, the most important data source is mocked. Test in hour 1.

---

### #2 (tie): Idea #11b — Customer Churn Rescue Agent (58/65)

**Pitch:** Agent monitors customer health across Stripe, Gong, and GitHub, identifies churn risk with LLM reasoning, and calls high-value at-risk customers with a personalized retention offer.

**Pros:**
- Same best-in-class Airbyte triple-connector story as #13. All 3 connectors are naturally essential.
- Textbook autonomous agent loop: ingest, reason, act, learn. Genuine LLM reasoning, not threshold rules.
- Phone ringing on the table with a personalized retention offer is an S-tier demo moment.

**Cons:**
- Churn prediction is not a novel concept. The innovation is in autonomous action, not detection.
- Same Gong access risk as #13.
- Seven integrations in 8 hours. Aerospike and Overmind should be cut first if behind schedule.

**Biggest risk:** Same as #13 — Gong connector availability. Also, the outbound call to customers is slightly theatrical. In production, fully autonomous outbound calls are a liability.

---

### #4: Idea #11 — Breach Notification Agent (56/65)

**Pitch:** Breach happens. Agent identifies affected customers, segments by jurisdiction (GDPR/CCPA), drafts notifications, calls enterprise customers, generates compliance report.

**Pros:**
- Best possible RSAC fit. Every CISO in the room has lived through breach notification. The 72-hour clock is real and terrifying.
- Phone call crossing the digital-physical boundary is the single strongest demo moment.
- Ghost (agent creates its own DB), Bland AI, TrueFoundry are all genuinely essential.

**Cons:**
- PII through LLM is a real concern. RSAC judges will question sending breach data to model providers.
- Autonomy is mostly linear — parse, match, segment, draft, send, report. Not deeply branching.
- Six tools, seven pipeline steps, two output modalities in 8 hours is very aggressive.

**Biggest risk:** PII objection from security-minded judges. Also, buildability with 6 tools.

---

### #5: Idea #12b — Revenue Anomaly Agent (55/65)

**Pitch:** Revenue drops 15%. Agent investigates across Stripe, GitHub, and Gong, cross-correlates a bad deploy with payment failures, and calls the revenue owner with the answer.

**Pros:**
- Cross-silo investigation narrative is universally relatable. "Four days to 30 seconds" is a punchy line.
- Same strong Airbyte triple-connector play as #13 and #11b.
- Cross-correlation is more technically demonstrable than risk detection — you can visualize the causal chain.

**Cons:**
- Structurally very similar to #13 (same connectors, same Bland AI call, same Aerospike search) but with a narrower decision space.
- Investigation follows a predictable sequence rather than making branching decisions.
- Revenue investigation is not security-adjacent for RSAC.

**Biggest risk:** This is a refinement of #12, not a leap beyond it. If choosing between this and #13, #13 has a richer autonomy story.

---

### #6: Idea #12 — Incident Response Agent (53/65)

**Pitch:** Alert fires. Agent pulls recent commits, searches past incidents, analyzes the codebase, and calls the on-call engineer with a full briefing.

**Pros:**
- Strongest raw tool use: Airbyte (GitHub), Aerospike (incident memory), Macroscope (codebase analysis), Bland AI (call). Four tools, all essential, each doing something distinct.
- Only idea that uses Macroscope naturally — it's literally what Macroscope is built for.
- Every on-call engineer in the room will connect with this.

**Cons:**
- Not a novel idea. PagerDuty, Rootly, incident.io all exist in this space.
- Demo repo with a planted bug is contrived. Judges will see through it.
- Most "intelligence" comes from sponsor tools, not original engineering.

**Biggest risk:** Judges think "this already exists" and discount the novelty.

---

## The Patterns

Three ideas (#13, #11b, #12b) share nearly identical architecture: Airbyte triple-connector + Aerospike vector search + Bland AI phone call + TrueFoundry routing. They differ only in domain (sales, churn, revenue). Pick one.

All top ideas share the same weakness: **buildability (3/5)**. No way around it — ambitious sponsor integration in 8 hours is the universal constraint. The winner is whoever manages build risk best, not whoever has the best idea on paper.

**Bland AI phone call** appears in 6 of the top 6 ideas. It's the strongest demo differentiator available. Whatever you build, the phone should ring.

---

## Decision Framework

| If you value... | Build... | Why |
|-----------------|----------|-----|
| Highest overall score | #5 Pricing Intelligence | 59/65. Strongest autonomy + idea combo. |
| Best sponsor story | #13 Deal Intelligence or #11b Churn Rescue | 5/5 tool use. Airbyte's dream use case. |
| Best demo moment | #13 or #11b | Phone call with data-derived specifics. |
| RSAC relevance | #11 Breach Notification | Every CISO gets it instantly. |
| Most buildable top idea | All tied at 3/5 | Manage risk, don't avoid it. |
| Safest scraping-free option | #13 or #11b | API-based data, no scraping fragility. |
