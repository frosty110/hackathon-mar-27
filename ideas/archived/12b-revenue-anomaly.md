# Idea #12b: Autonomous Revenue Anomaly Agent

## One-Line Pitch

An agent that autonomously investigates revenue drops by cross-correlating billing data, code changes, and sales call signals -- connects the dots across silos in seconds instead of days -- then calls the revenue owner with a root cause diagnosis and recommended fix.

## The Problem

**Who has this problem:** CFOs, VP Sales, VP Engineering, revenue operations -- every leadership team at every SaaS company that has ever scrambled to explain a revenue dip.

**Current workaround:** Revenue drops 15% on Wednesday. Finance notices Thursday morning. They ping engineering: "Did anything break?" Engineering checks deploys. Product checks usage metrics. Support checks ticket volume. Sales checks pipeline. Each team looks at their own silo. Slack threads proliferate. A war room gets scheduled for Friday afternoon. By Monday, someone finally connects the dots: a Tuesday deploy broke the checkout flow for mobile users in the EU. Failed payments spiked 340% for that segment. Four days, $200K lost, and the answer was sitting across three systems the entire time.

**Why it sucks:** The signals exist in real time -- failed payment spikes in Stripe, a PR that changed payment flow in GitHub, prospects mentioning checkout problems in Gong calls. But no human monitors all three simultaneously. The cross-correlation is the hard part. Each team sees their slice and waits for someone else to connect the dots. The investigation is mechanical pattern-matching across data sources -- exactly what an autonomous agent should do.

## The Autonomous Agent Loop

**Trigger:** Revenue anomaly detected (a threshold alert fires when MRR, failed payments, or refund rate deviates beyond a configurable threshold -- simulated for demo but architecturally real).

```
Revenue anomaly detected (threshold breach)
    |
    v
[1] MULTI-SOURCE INVESTIGATION (parallel, via Airbyte)
    |-- Stripe: Failed payments, refunds, downgrades, geographic/device breakdown
    |-- GitHub: Recent PRs merged, deploys, what code changed this week
    |-- Gong: Sales call transcripts -- are prospects/customers mentioning issues?
    |
    v
[2] CROSS-CORRELATION (the hard part -- this is where the agent reasons)
    |-- Temporal alignment: What changed when?
    |-- "Failed EU mobile payments up 340% since Tuesday"
    |-- "PR #847 merged Tuesday, changed payment flow"
    |-- "3 Gong calls this week: prospects mentioned checkout issues"
    |-- Agent builds a causal hypothesis with evidence chain
    |
    v
[3] HISTORICAL PATTERN SEARCH (Aerospike vector search)
    |-- Embed the anomaly profile (symptoms, timing, affected segments)
    |-- Search past anomalies: "Last time EU payments dropped, it was a
    |   Stripe webhook config issue after a deploy"
    |-- Augment hypothesis with historical precedent
    |
    v
[4] PERSIST & STRUCTURE (Aerospike)
    |-- Store investigation as structured record + vector embedding
    |-- Future anomalies will find this investigation via similarity search
    |-- The agent builds institutional memory with every incident
    |
    v
[5] GENERATE ROOT CAUSE REPORT
    |-- Evidence chain: timeline, data from each source, correlation logic
    |-- Confidence level and alternative hypotheses
    |-- Recommended fix with specificity ("revert PR #847" not "check the code")
    |
    v
[6] PHONE CALL BRIEFING (Bland AI)
    |-- Agent calls the revenue owner (VP Sales, CFO)
    |-- Delivers: what happened, why, evidence, what to do
    |-- No human wrote the script -- the agent composed it from its investigation
```

**Key autonomous decisions the agent makes:**
- Which data to pull and how to slice it (geographic? device type? plan tier?)
- How to align signals temporally (what correlates with the timing of the drop?)
- Whether the evidence supports a causal hypothesis or just correlation
- Which past anomalies are relevant and what they imply
- What confidence level to assign the root cause
- How to compose the phone briefing -- prioritizing the most actionable information

This is not a script. The agent is doing genuine cross-source reasoning. Give it a different anomaly (churn spike instead of payment failures) and it pulls different data, asks different questions, reaches a different conclusion.

## Sponsor Stack (3+ required)

### 1. Airbyte (SaaS Data Connectors) -- ESSENTIAL x3

**What it does:** Pulls structured data from three sources simultaneously via Airbyte's Python agent connector packages:
- **Stripe:** Failed payment events, refund records, subscription changes, geographic and plan-level breakdown. The "what is happening" data.
- **GitHub:** Recent merged PRs, commit history, deploy metadata. The "what changed" data.
- **Gong:** Sales call transcripts from the past week. The "what are people saying" data.

**Why it's load-bearing:** This is the entire investigative backbone. The agent's cross-correlation requires data from all three sources -- remove any one and the investigation has a blind spot. Stripe without GitHub means you see the problem but can't identify the code change that caused it. GitHub without Stripe means you see deploys but don't know which one correlates with the revenue drop. Gong without either means you hear complaints but can't tie them to a root cause. The three-source investigation IS the product.

**Essential or cosmetic:** **ESSENTIAL -- all three connectors are load-bearing.** This is the single best Airbyte use case possible. Three distinct connectors, each providing a different investigative signal, all necessary for the cross-correlation that is the core value proposition. This is Airbyte's sales pitch made real: "connect all your SaaS data and let an agent reason across it."

### 2. Aerospike (Vector Search + Anomaly Memory) -- ESSENTIAL

**What it does:** Two functions:
1. **Vector similarity search:** When a new anomaly is detected, embed its profile (symptoms, timing, affected segments) and search for similar past anomalies. "The last time EU mobile payments dropped after a deploy, the root cause was a Stripe webhook configuration change."
2. **Persistent investigation memory:** Every completed investigation is stored as a vector embedding + structured metadata. The agent accumulates institutional knowledge.

**Why it's load-bearing:** The pattern matching against past anomalies is what separates this from a one-shot analyzer. A human revenue leader's value comes from experience -- "I've seen this pattern before." Aerospike gives the agent that experience. The sub-millisecond latency means the historical search doesn't slow down the investigation. The persistence means every investigation makes the next one smarter.

**Essential or cosmetic:** **ESSENTIAL.** Remove Aerospike and the agent can only analyze the current anomaly in isolation. The historical pattern matching ("last time this happened, the cause was X") is a primary input to the root cause hypothesis and dramatically increases confidence in the diagnosis.

### 3. Bland AI (Voice Phone Call) -- ESSENTIAL

**What it does:** When the investigation is complete, the agent calls the revenue owner's phone and delivers a structured voice briefing: what the anomaly is, what caused it, what the evidence is, and what to do about it. The script is dynamically composed from the investigation findings -- no human wrote it.

**Why it's load-bearing:** The phone call is the autonomous ACTION. Most anomaly detection tools send a Slack notification that gets buried. This agent calls you. It interrupts your day because it has done enough investigation to warrant interruption. The phone call also proves the highest level of autonomy: the agent didn't just analyze data, it decided the findings were urgent enough to call a human and composed a briefing.

**Essential or cosmetic:** **ESSENTIAL for the value prop, ESSENTIAL for the demo.** Without the call, this is a dashboard that generates a report. With the call, it's an autonomous investigator that briefs leadership. The phone ringing in the demo room is also the single most memorable moment.

### 4. TrueFoundry (Multi-LLM Gateway) -- STRONG VALUE-ADD

**What it does:** Routes different investigation subtasks to appropriate models:
- Data extraction and summarization (parsing Stripe events, GitHub PRs, Gong transcripts) goes to a fast, cheap model
- Cross-correlation reasoning and root cause synthesis goes to a strong reasoning model
- Phone script composition goes to a model optimized for natural language generation

**Why it's load-bearing:** The investigation involves 4-6 LLM calls. Sending raw Stripe event parsing to an expensive reasoning model wastes money and time. Sending the cross-correlation synthesis to a cheap model produces weak hypotheses. The routing is architecturally meaningful, not decoration.

**Essential or cosmetic:** **STRONG VALUE-ADD.** You could hardcode a single model and the demo works. But the routing demonstrates real architectural thinking about cost, latency, and capability matching. Easy to integrate (drop-in gateway) and adds a legitimate talking point for technical implementation scoring.

### 5. Overmind (LLM Observability) -- COSMETIC

**What it does:** Drop-in SDK wrapping all LLM calls. Shows token usage, latency, and cost per investigation step.

**Why it's load-bearing:** It isn't, functionally. But showing "this investigation cost $0.14 and took 22 seconds across 5 LLM calls" after the demo is a nice closing beat that demonstrates production awareness.

**Essential or cosmetic:** **COSMETIC.** Low build cost (drop-in SDK), marginal demo impact. Include it if time allows, cut it first if time is tight.

**Tool count: 3 essential (Airbyte x3, Aerospike, Bland AI) + 1 strong value-add (TrueFoundry) + 1 cosmetic (Overmind) = comfortably clears the 3-tool minimum with genuine depth on 3-4.**

## The "Whoa" Demo Moment

The agent has just finished its investigation. On screen, you see the evidence chain materialize: a timeline showing the revenue drop, the Stripe data showing EU mobile payment failures spiking 340%, the GitHub PR that merged Tuesday changing the payment flow, and three Gong call snippets where prospects mention checkout problems. The root cause report appears: "Revenue down 15%. Root cause: PR #847 broke mobile checkout for EU users. Evidence confidence: HIGH."

Then the presenter says: "The agent has determined this is urgent. It's calling the CFO now."

A phone rings in the room. On speaker:

> "Hi, this is your revenue anomaly agent. I've completed an investigation into the 15% revenue decline detected this week. Here's what I found: failed payments from EU mobile users increased 340% starting Tuesday. This correlates with PR #847, merged Tuesday morning, which modified the payment checkout flow. Three sales calls this week had prospects mentioning checkout problems on mobile. I've seen a similar pattern before -- last quarter, a deploy caused EU payment failures that were traced to a Stripe webhook configuration issue. My recommendation: immediately revert PR #847 and verify the EU mobile checkout flow. Estimated revenue impact so far is $200,000. A full investigation report has been sent to your dashboard."

The room processes what just happened. An agent investigated a complex, cross-departmental problem in seconds and called leadership with the answer. That's the demo moment.

## 3-Minute Demo Script

**[0:00-0:25] The Problem (25s)**
"Revenue dropped 15% this week. What happens at every company? Finance pings engineering, engineering checks deploys, product checks usage, support checks tickets, sales checks pipeline. Everyone's looking at their own silo. By the time someone connects the dots -- a deploy broke checkout for EU mobile users -- it's been four days and you've lost $200K. We built an agent that does the entire investigation autonomously, across every data source, in seconds."

**[0:25-0:50] The Trigger (25s)**
Show the anomaly alert firing: "MRR declined 15% week-over-week. Threshold breached." The agent activates. Show the live investigation dashboard:
- "Pulling billing data from Stripe via Airbyte..."
- "Pulling recent code changes from GitHub via Airbyte..."
- "Pulling sales call data from Gong via Airbyte..."

Three Airbyte connectors running in parallel. Data populating in real time.

**[0:50-1:25] Cross-Correlation (35s)**
This is the technical core. Show the agent's reasoning materializing:
- Stripe panel: "Failed payments from EU mobile users up 340% since Tuesday" (highlight the spike on a chart)
- GitHub panel: "PR #847 merged Tuesday 9:14am -- modified payment checkout flow" (highlight the PR)
- Gong panel: "3 calls this week, prospects mentioned 'checkout issues' and 'payment problems on mobile'" (show transcript snippets)

The agent draws the connection: "Temporal correlation: PR #847 merged Tuesday morning. EU mobile payment failures began Tuesday afternoon. Gong mentions of checkout issues began Wednesday." Show the causal hypothesis forming with an evidence chain.

**[1:25-1:50] Historical Pattern Match (25s)**
"The agent is now checking its memory -- has it seen anything like this before?" Show the Aerospike vector search. Results appear: "Similar anomaly detected Q3 2025: EU payment failures after a deploy. Root cause was Stripe webhook configuration change. Resolution: revert deploy, verify webhook config." The agent incorporates this into its hypothesis.

Root cause report generated: confidence HIGH, recommended fix: revert PR #847.

**[1:50-2:25] The Phone Call (35s)**
"The agent has classified this as urgent. Watch what happens." A phone rings. Play on speaker. Bland AI delivers the full briefing -- the 340% spike, PR #847, the Gong call mentions, the historical precedent, the recommended fix, the $200K estimated impact. The audience hears a complete investigation delivered by voice.

**[2:25-2:50] The Learning Loop (25s)**
"This investigation is now stored. Next time revenue drops, the agent will find this investigation in its memory and triage even faster." Show the Aerospike entry. Show TrueFoundry's model routing: "5 LLM calls, 3 models, $0.14 total cost." Optionally flash Overmind observability.

**[2:50-3:00] Close (10s)**
"Four days and $200K lost, or 30 seconds and a phone call. The agent never sleeps, never works in a silo, and never forgets a pattern."

## Technical Architecture

```
                   ┌───────────────────┐
                   │  Anomaly Trigger   │
                   │  (threshold alert) │
                   └────────┬──────────┘
                            │
                   ┌────────▼──────────┐
                   │   Orchestrator     │
                   │  (Python / FastAPI)│
                   └──┬─────┬─────┬────┘
                      │     │     │
       ┌──────────────┘     │     └──────────────┐
       │                    │                     │
       ▼                    ▼                     ▼
┌─────────────┐    ┌──────────────┐      ┌──────────────┐
│   Airbyte   │    │   Airbyte    │      │   Airbyte    │
│   Stripe    │    │   GitHub     │      │    Gong      │
│             │    │              │      │              │
│ - Failed    │    │ - Merged PRs │      │ - Call       │
│   payments  │    │ - Commits    │      │   transcripts│
│ - Refunds   │    │ - Deploy     │      │ - Sentiment  │
│ - Geo/device│    │   metadata   │      │ - Issue      │
│   breakdown │    │              │      │   mentions   │
└──────┬──────┘    └──────┬───────┘      └──────┬───────┘
       │                  │                      │
       └──────────────────┼──────────────────────┘
                          │
                 ┌────────▼────────┐
                 │  TrueFoundry    │
                 │  LLM Gateway    │
                 │                 │
                 │ Fast model:     │
                 │  data extraction│
                 │ Strong model:   │
                 │  cross-         │
                 │  correlation &  │
                 │  root cause     │
                 └────────┬────────┘
                          │
              ┌───────────┼───────────┐
              │           │           │
              ▼           ▼           ▼
       ┌────────────┐ ┌────────┐ ┌──────────┐
       │ Aerospike  │ │ Report │ │ Bland AI │
       │            │ │ Engine │ │          │
       │ - Vector   │ │        │ │ Phone    │
       │   search   │ │ Root   │ │ call to  │
       │   past     │ │ cause  │ │ revenue  │
       │   anomalies│ │ report │ │ owner    │
       │ - Store    │ │ + fix  │ │ with     │
       │   this     │ │ reco   │ │ briefing │
       │   analysis │ │        │ │          │
       └────────────┘ └────────┘ └──────────┘
              │
       ┌──────▼──────┐
       │  Overmind   │
       │  (wraps all │
       │  LLM calls) │
       └─────────────┘
```

**Implementation notes:**
- **Orchestrator:** Python FastAPI. Single endpoint for the anomaly trigger. Fans out three Airbyte data pulls in parallel (asyncio). Pipes results through the LLM reasoning chain.
- **Airbyte:** Use PyAirbyte connector packages. Pre-configure each connector with auth. Stripe: pull failed_payment events, refunds, subscription changes filtered to the anomaly window. GitHub: pull merged PRs and commits from the past week. Gong: pull call transcripts from the past week.
- **Cross-correlation engine:** This is the novel technical work. The LLM receives structured data from all three sources and reasons about temporal alignment, causal relationships, and confidence levels. Use a strong reasoning model via TrueFoundry.
- **Aerospike:** Pre-seed with 15-20 past anomaly investigations as vector embeddings. Each has: symptoms, root cause, resolution, affected segments. On new anomaly, embed the symptom profile and do k-NN search. After investigation, store the new analysis.
- **Bland AI:** Compose the phone script dynamically from investigation findings. Key variables: the anomaly description, the root cause hypothesis, the evidence chain (top 3 data points), the recommended fix, the estimated impact. Trigger the call via Bland's API.
- **Demo data:** Pre-seed a realistic scenario: a demo GitHub repo with a PR that modifies payment code, Stripe test mode with elevated failed payment events for EU, synthetic Gong transcripts mentioning checkout issues. The scenario is contrived but the mechanism is real.

## Buildability Risk Assessment

**Overall buildability: 3/5 (tight but doable with disciplined prioritization)**

| Component | Difficulty | Risk | Mitigation |
|-----------|-----------|------|------------|
| Anomaly trigger + orchestrator | Low | Low | Standard FastAPI, webhook endpoint, well-understood |
| Airbyte Stripe connector | Medium | Medium | Stripe test mode provides good sandbox data. PyAirbyte config can be finicky -- pre-test the night before if possible |
| Airbyte GitHub connector | Low-Medium | Low | Public repos, no auth issues. Easiest of the three connectors. Start here |
| Airbyte Gong connector | Medium | HIGH | Gong requires enterprise API access. If you don't have a Gong account with real call data, you're faking the data source. Mitigation: use synthetic transcript data through the connector, or pre-seed realistic JSON |
| Cross-correlation LLM reasoning | Medium | Medium | This is the novel technical work. The prompt engineering for temporal alignment and causal reasoning needs to be solid. Pre-test with your demo scenario data |
| Aerospike vector search | Medium | Low | Well-documented, straightforward setup. Seed data in advance |
| TrueFoundry multi-model routing | Low | Low | Drop-in gateway, minimal configuration |
| Bland AI phone call | Low | Low | Simple API call. Test the phone number extensively. Pre-compose most of the script template |
| End-to-end timing | High | Medium | All data pulls + LLM reasoning + phone call needs to complete in ~20-30 seconds for demo pacing. Pre-warm connections. Have cached data as fallback |
| Live demo reliability | -- | HIGH | 4-5 external APIs in sequence. Build a deterministic replay mode with cached API responses but a real phone call |

**The critical path:** Airbyte connectors. If all three work in the first 2-3 hours, the rest flows. If Gong fights you, immediately fall back to pre-seeded transcript data and focus the Airbyte story on Stripe + GitHub (still two real connectors, which is still strong).

**Recommended build order:**
1. **Hour 1-2:** Bland AI phone call working (get the demo moment locked). Aerospike schema + seed 15 historical anomaly records.
2. **Hour 2-3:** Airbyte GitHub connector working. Airbyte Stripe connector working against test mode. Attempt Gong connector (if it fails, switch to pre-seeded JSON immediately).
3. **Hour 3-5:** Cross-correlation reasoning engine (the LLM prompting that ties signals together). TrueFoundry integration for model routing. Vector search returning relevant past anomalies.
4. **Hour 5-6:** End-to-end flow: trigger anomaly, data pulls, reasoning, report, phone call. Debug timing.
5. **Hour 6-7:** Demo UI (simple dashboard showing investigation in real time). Polish the phone script.
6. **Hour 7-8:** Overmind integration (if time). Replay/fallback mode. Rehearse full demo 5+ times. Test phone call 5+ times.

**What to cut if time runs out:**
- **Cut first:** Overmind (cosmetic)
- **Cut second:** Third Airbyte connector (Gong) -- pre-seed the transcript data, focus on Stripe + GitHub as live connectors
- **Cut third:** TrueFoundry routing -- hardcode a single model
- **NEVER cut:** The phone call. The cross-correlation reasoning. The Stripe + GitHub data pulls. These are the demo.

## Honest Weaknesses

**1. This is a variant of #12 (Incident Response), and the differentiation is subtle.**
Idea #12 investigates production incidents. Idea #12b investigates revenue anomalies. The architecture is nearly identical: multi-source data pull, cross-correlation, historical pattern match, phone call delivery. The revenue framing is arguably more relatable (everyone understands "revenue dropped") but the underlying agent loop is the same. A judge who saw both would notice.

**2. The Gong data is almost certainly synthetic.**
Gong requires an enterprise API account with real sales calls. Unless the team has an active Gong instance with relevant call data, this connector is running against fabricated transcripts. This weakens the "autonomous investigation" narrative because the most qualitative signal (what customers are saying) is curated. Mitigation: be transparent, or make the demo scenario compelling enough that judges don't probe the data provenance.

**3. The cross-correlation reasoning is a sophisticated prompt, not a sophisticated algorithm.**
The "connecting dots across systems" is fundamentally an LLM prompt that receives structured data from three sources and reasons about temporal correlations. This is genuinely useful and the output can be impressive, but a technical judge may note that the "intelligence" is a well-crafted prompt, not a novel correlation algorithm. The agent's reasoning quality depends entirely on the LLM's ability to identify temporal patterns in structured data -- which current models do well, but it's not your innovation.

**4. The demo scenario is contrived by necessity.**
You plant the bug in PR #847. You configure the Stripe test data to show EU mobile payment failures. You write the Gong transcripts that mention checkout issues. The agent will look brilliant because you designed the investigation to succeed. Every hackathon demo is like this, but a skeptical judge will wonder: "What happens when the signals don't align this cleanly?" Real revenue drops have messy, multi-causal explanations.

**5. The "15% revenue drop" trigger is simplistic.**
The agent triggers on a threshold breach. But real revenue anomalies are nuanced -- is it seasonal? Is it a data lag? Is it a billing cycle artifact? The agent doesn't do anomaly detection, it does anomaly investigation. The detection is a hardcoded threshold. This is fine for a hackathon, but the setup narration needs to acknowledge this.

**6. Bland AI latency during the demo is a single point of failure.**
The phone must ring. If Bland's API has a hiccup, or the call takes 15 seconds to connect, or the voice quality is poor, the demo climax dies. Must have a fallback: pre-recorded audio, a second device ready, or a teammate who "receives" the call on cue.

**7. Revenue investigation is not a security use case.**
This is RSAC 2026. Revenue anomaly investigation is a business operations problem, not a cybersecurity problem. The judges may be predisposed to reward security-adjacent ideas (threat detection, vulnerability management, compliance). The agent's sophistication needs to overcome domain mismatch. Possible counter-framing: "revenue anomalies can be caused by security incidents -- a breached payment flow, a fraudulent refund pattern -- the agent investigates regardless of root cause."

## Final Score Recommendation

| Criterion | Weight | Score | Weighted | Justification |
|-----------|--------|-------|----------|---------------|
| **Autonomy** | x3 | 4 | 12 | Strong autonomous loop: trigger fires, agent pulls from 3 sources, cross-correlates, searches history, generates report, calls a human. The cross-correlation reasoning is genuine -- the agent is connecting dots, not following a script. Loses a point vs. a 5 because the investigation follows a fairly predictable sequence (pull data, correlate, report, call) rather than making deeply branching decisions. The agent doesn't decide to investigate different angles based on what it finds -- it always pulls all three sources and correlates. |
| **Idea** | x3 | 4 | 12 | The problem is universally relatable -- every company has had a "why did revenue drop?" fire drill. Judges will immediately get it. The cross-silo investigation angle is strong: everyone has experienced the pain of waiting for three departments to coordinate. Loses a point because (a) revenue investigation is not security-adjacent for an RSAC audience, and (b) it's structurally similar to #12 incident response, which is a well-explored problem space. |
| **Technical Implementation** | x2 | 4 | 8 | The cross-correlation engine is legitimate technical work -- temporal alignment, causal reasoning, confidence scoring across heterogeneous data sources. Multi-model routing via TrueFoundry adds architectural depth. Vector search for historical pattern matching is a solid second layer. Not a 5 because the core reasoning is prompt engineering over structured data, not a novel algorithm. |
| **Tool Use** | x2 | 5 | 10 | This is a top-tier tool use story. Airbyte with THREE connectors (Stripe, GitHub, Gong) -- each providing a distinct investigative signal -- is the strongest possible Airbyte integration. Aerospike for historical pattern matching is genuinely load-bearing. Bland AI for the phone call is the autonomous action. TrueFoundry for model routing is architecturally meaningful. Every essential tool maps to a distinct, necessary capability. This is not tool-stuffing. |
| **Presentation/Demo** | x2 | 5 | 10 | The demo arc is outstanding: revenue drops, agent investigates across three data sources in real time (visible on screen), evidence chain materializes showing the correlation, historical match found, phone rings. The phone call delivering a complete investigation briefing with specific numbers ("340% spike", "PR #847", "$200K impact") is visceral and memorable. The "four days or 30 seconds" closing line lands. |
| **Buildability** | x1 | 3 | 3 | The build is tight but achievable. Three Airbyte connectors are aggressive (Gong is high-risk). The cross-correlation prompt engineering needs careful tuning. Five external APIs means fragile live demos. The replay/fallback mode is essential but costs build time. A disciplined team that locks the phone call first and works backward can do it. The Gong connector is the risk -- plan for the fallback from the start. |

**Total: 12 + 12 + 8 + 10 + 10 + 3 = 55/65**

**Verdict: Strong idea, but it's a refinement of #12, not a leap beyond it.**

This scores 55, which puts it in the competitive range but below #13 Deal Intelligence (58) and #5 Pricing Intelligence (59). The core strengths are real: the cross-silo investigation narrative is relatable, the Airbyte triple-connector use is the best possible sponsor alignment, and the phone call is a proven demo moment.

**How it compares to related ideas:**
- **vs. #12 Incident Response (53):** 12b is a stronger version of the same architecture. Revenue anomaly investigation is more relatable to a business audience than production incident triage. The Airbyte triple-connector story (Stripe + GitHub + Gong) is stronger than #12's single GitHub connector + Macroscope. But #12 has a stronger RSAC/security angle.
- **vs. #13 Deal Intelligence (58):** Same Airbyte triple-connector play (Gong, Stripe, GitHub), same Bland AI phone call moment. #13 has a richer autonomy story (more decision points: risk assessment, strategy generation, coaching composition). 12b's cross-correlation reasoning is arguably more technically interesting but the decision space is narrower.
- **vs. #5 Pricing Intelligence (59):** Different domain entirely. #5 has a more universally relatable problem and stronger novelty. 12b has better Airbyte integration depth.

**The honest question: should you build this over #13 or #5?** Probably not, unless the team is specifically excited about the revenue investigation narrative. The architecture is nearly identical to #13 (same three Airbyte connectors, same Aerospike pattern matching, same Bland AI call), and #13 has a richer autonomy story with more decision branching. If you build 12b, you're choosing the more constrained version of the same technical foundation.

**If you DO build it, the key advantages are:** (1) the cross-correlation is more technically demonstrable than #13's risk detection -- you can show a timeline visualization that makes the causal chain visible, and (2) the "four days to 30 seconds" narrative is extremely punchy.
