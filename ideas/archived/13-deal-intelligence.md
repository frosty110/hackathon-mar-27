# Idea #13: Autonomous Deal Intelligence Agent

## One-Line Pitch

An agent that autonomously researches every new deal entering your pipeline -- pulling call transcripts, revenue history, and tech stack data -- finds pattern-matched similar deals, detects risk signals, generates a strategy brief, and phones the account exec with coaching when a deal is at risk.

## The Problem

**Who has this problem:** Account executives and sales leaders at B2B SaaS companies running enterprise pipelines with 30+ day sales cycles. Every rep at every company with a CRM.

**Current workaround:** A new opportunity hits the pipeline. The rep spends 2-4 hours manually researching: checks Gong for past calls with the prospect, looks up revenue history in Stripe, browses the prospect's GitHub to understand their tech stack, scans Salesforce for similar deals that won or lost. They ask the Slack channel, "anyone sold to a company like this?" They build a strategy in their head based on gut feel and whatever they remember. Meanwhile, risk signals -- competitor mentions in call transcripts, declining email engagement, objection patterns that historically kill deals -- go completely unnoticed until the deal is marked "Closed Lost" and someone asks what happened.

**Why it sucks:** The information exists across Gong, Stripe, GitHub, and the CRM, but nobody synthesizes it. Every deal starts from zero. Tribal knowledge about "what worked on deals like this" lives in people's heads and walks out the door when reps leave. Risk signals are detectable in retrospect ("oh yeah, they mentioned a competitor in the third call") but nobody catches them in real time. Enterprise deals worth $50K-$500K are lost to preventable ignorance.

## The Autonomous Agent Loop

This is the core of the product and the #1 judging criterion.

**Trigger:** A new opportunity is created in the CRM (simulated webhook). No human initiates the research. The agent wakes up on its own.

**Step 1 -- Multi-Source Data Collection (autonomous):**
The agent simultaneously pulls data from three sources via Airbyte connectors:
- **Gong:** Call recordings/transcripts for the prospect account. Past meetings, sentiment, topics discussed, objections raised.
- **Stripe:** Revenue history if this is an existing customer. Payment patterns, plan tier, expansion/contraction signals.
- **GitHub:** The prospect's public repositories. Tech stack, activity level, engineering team size, open-source engagement.

This is not a single API call. The agent decides what to pull based on what's available (new prospect = no Stripe history, adjust accordingly).

**Step 2 -- Pattern Matching via Vector Search (autonomous):**
The agent embeds the deal profile (industry, company size, tech stack, deal size, buying signals) and searches Aerospike for similar historical deals. Returns the top 5 most similar deals with outcomes: "deals like this with annual billing closed 73% of the time" or "3 out of 4 similar deals that stalled at this stage were lost to Competitor X."

The agent decides what dimensions matter for similarity. This is genuine reasoning, not a canned query.

**Step 3 -- Risk Signal Analysis (autonomous):**
The agent analyzes Gong call transcripts for:
- Competitor mentions (frequency and context)
- Objection patterns (pricing pushback, feature gaps, security concerns)
- Sentiment trajectory (enthusiasm declining across calls)
- Engagement patterns (longer gaps between meetings, fewer attendees)

It compares these signals against patterns from the historical deals that were lost. The agent makes a risk assessment: GREEN (on track), YELLOW (watch), RED (at risk, intervention needed).

**Step 4 -- Strategy Brief Generation (autonomous):**
The agent synthesizes everything into a deal strategy brief:
- Deal summary and context
- Similar deal outcomes and what worked
- Identified risk signals with severity
- Recommended pricing strategy (based on what closed similar deals)
- Likely objections and suggested responses
- Recommended next steps

This requires the strong model -- it's reasoning across multiple data sources, historical patterns, and strategic judgment.

**Step 5 -- Proactive Phone Coaching (autonomous, triggered by risk):**
If the deal is RED (at-risk), the agent calls the account executive via Bland AI with specific coaching: "Your Acme Corp deal is showing 3 risk signals. They mentioned DataDog in the last call, engagement has dropped from weekly to biweekly, and the VP of Engineering who was your champion hasn't attended the last two meetings. In similar deals, reps who addressed the competitor concern directly in a 1:1 with the champion recovered 60% of the time. I'm sending you a strategy brief now."

The phone rings. A real voice delivers actionable intelligence. No human wrote the script.

**Step 6 -- Learning & Memory (autonomous):**
All analysis is stored in Aerospike (vector embeddings for future pattern matching) and Ghost (structured deal records for querying). Every deal outcome that's eventually recorded makes future pattern matching more accurate. The agent gets smarter with every deal.

**Key autonomous decisions the agent makes:**
- What data to pull based on what's available for each prospect
- Which dimensions matter for deal similarity
- Whether risk signals are noise or genuine threats
- What strategy to recommend based on pattern-matched outcomes
- Whether to escalate via phone call vs. just file the brief
- How to prioritize and frame coaching advice for the phone call

## Sponsor Stack (3+ required)

### 1. Airbyte (SaaS Data Connectors) -- ESSENTIAL

**What it does:** Pulls structured data from Gong (call transcripts, meeting metadata, sentiment), Stripe (customer revenue history, payment patterns, plan data), and GitHub (repositories, tech stack, activity metrics) using Airbyte's agent connector packages. Three real SaaS connectors feeding one agent.

**Why it's load-bearing:** This is the backbone of the entire product. The agent's intelligence depends entirely on having multi-source data about the prospect. Without Airbyte, you'd spend 6 of your 8 hours hand-coding three separate API integrations, each with their own auth, pagination, and data model quirks. Airbyte is doing what it was literally designed for: making SaaS data available to downstream consumers. The agent IS the downstream consumer.

**Essential or cosmetic:** Essential. Remove Airbyte and you have no data, no agent, no product. This is the single best Airbyte use case in the hackathon. Three production-grade SaaS connectors (Gong, Stripe, GitHub) feeding an autonomous agent -- that's Airbyte's pitch deck come to life.

### 2. Aerospike (Sub-ms Vector Search + Agent Memory) -- ESSENTIAL

**What it does:** Stores vector embeddings of every historical deal profile (company attributes, deal characteristics, outcome) and provides sub-millisecond similarity search. When a new deal arrives, the agent queries Aerospike: "find the 5 most similar deals and their outcomes." Also stores the agent's analysis history as persistent memory.

**Why it's load-bearing:** The "deals like this closed 73% of the time" insight IS the product's core value. Without vector search over historical deals, the agent can only analyze the current deal in isolation. With it, the agent draws on every previous deal's outcome to make pattern-based recommendations. This is the difference between "here's what I found about Acme Corp" (a research assistant) and "here's what happened in deals like Acme Corp" (an intelligence system).

**Essential or cosmetic:** Essential. Remove Aerospike and you lose pattern matching, similar-deal lookup, and agent memory. The product degrades from "deal intelligence" to "deal research summary." The strategic recommendations have no empirical backing.

### 3. Bland AI (Voice/Phone Calls) -- ESSENTIAL

**What it does:** When the agent determines a deal is at-risk, it places a real phone call to the account executive using Bland AI's voice API. The agent composes the coaching script dynamically based on the specific risk signals detected and the recommended recovery strategy.

**Why it's load-bearing:** The phone call is the autonomous ACTION. Most agent demos end at "here's a report." This agent acts on its findings by calling a human with specific instructions. The phone ringing during the demo is the highest-impact moment. It also proves the agent makes a real decision (is this deal at-risk enough to warrant a phone call?) rather than just generating text.

**Essential or cosmetic:** Essential. Remove Bland AI and the agent generates a report that sits in a database. The phone call transforms it from a passive analysis tool to an active coaching system. It's also the single best demo moment -- the phone ringing in a judge's hand would be unforgettable.

### 4. TrueFoundry (Multi-LLM Gateway) -- ESSENTIAL

**What it does:** Routes the agent's LLM calls to appropriate models based on the task. Data extraction from Gong transcripts and GitHub repos goes to a fast, cheap model. Risk analysis and strategy generation go to a strong reasoning model. The gateway also provides cost and latency observability for each step.

**Why it's load-bearing:** The agent makes 4-6 LLM calls per deal analysis. Sending raw transcript extraction to Claude Opus would be wasteful; sending strategic reasoning to GPT-4o-mini would be unreliable. TrueFoundry's routing is what makes the pipeline both reliable AND economically viable. It also gives the demo a "show the reasoning" layer -- you can display which model handled each step and why.

**Essential or cosmetic:** Essential. You could hard-code a single model, but you'd either overpay (using the expensive model everywhere) or get bad strategy recommendations (using the cheap model everywhere). The routing is architecturally meaningful.

### 5. Ghost (Agent-Managed Postgres) -- STRONG VALUE-ADD

**What it does:** The agent uses Ghost to autonomously create and manage a structured deal analysis database. Per-deal records, strategy briefs, risk assessments, outcome tracking. The agent creates its own schema, inserts analysis results, and queries historical outcomes.

**Why it's load-bearing:** While Aerospike handles vector search (unstructured similarity), Ghost handles structured queries: "show me all deals where competitor X was mentioned," "what's the average close rate for deals in the $100K-$200K range," "which reps have the best win rate on at-risk deals." The agent creates its own tables -- a genuine autonomy showcase.

**Essential or cosmetic:** Strong value-add. You could probably do structured queries in Aerospike too, but Ghost's MCP-based database management (agent creates its own schema, forks DBs for what-if analysis) adds both real capability and a strong autonomy story. Honest assessment: the product works without it, but the team-wide learning layer (Step 6) depends on it.

### 6. Overmind (LLM Observability) -- COSMETIC BUT USEFUL

**What it does:** Drop-in SDK wrapping all LLM calls. Shows cost, latency, and token usage per analysis step. Provides optimization suggestions.

**Why it's load-bearing:** It's not load-bearing in the strictest sense. The product works without it. But it provides the observability layer that makes the demo transparent: "this deal analysis cost $0.08 and took 12 seconds." For judges evaluating technical implementation, seeing the x-ray into the agent's reasoning pipeline is valuable.

**Essential or cosmetic:** Cosmetic. Be honest about this. It's the 6th tool and it's instrumentation, not functionality.

## The "Whoa" Demo Moment

The demo builds to a phone call.

You've shown the agent ingest data, pattern-match, and generate a strategy brief. That's impressive but cerebral. Then the agent finds risk signals in the Gong transcripts: competitor mentions, declining engagement, a missing champion.

The agent's risk assessment flips to RED.

And then a phone rings. In the room. The account exec's phone (your phone, your teammate's phone, a judge's phone if you can set it up) rings with a call from the agent. A voice says:

"Hi, this is your deal intelligence agent. Your Acme Corp opportunity is showing three risk signals. The prospect mentioned DataDog in their last two calls, meeting frequency has dropped from weekly to biweekly, and your executive sponsor hasn't attended since January. In similar deals, reps who scheduled a direct 1:1 with their champion to address the competitive concern recovered 60% of the time. A strategy brief has been sent to your dashboard. Would you like me to draft an email to your contact?"

A phone ringing in the demo room. A voice delivering specific, data-backed coaching. That's the "whoa."

## 3-Minute Demo Script

**0:00-0:25 -- The Problem**
"Enterprise reps spend 3 hours researching every new deal. They check Gong for past calls, Stripe for revenue history, GitHub for the prospect's tech stack, and ask around for anyone who's sold to a similar company. Meanwhile, risk signals in call transcripts go unnoticed until the deal is lost. We built an agent that does all of this autonomously -- and calls you when something is wrong."

**0:25-0:50 -- The Trigger**
Show a simulated CRM webhook: "New Opportunity: Acme Corp, $150K ARR, Enterprise Plan." The agent wakes up. Show the Airbyte connectors pulling data simultaneously from Gong (3 call transcripts), Stripe (existing $12K/year customer), and GitHub (47 public repos, Python/Go stack, active contributor).

**0:50-1:20 -- Pattern Matching**
The agent queries Aerospike. Show the vector search returning 5 similar historical deals. A table appears: "4 out of 5 similar deals closed when annual billing was offered upfront. Average deal cycle: 45 days. Primary risk factor in lost deals: competitor displacement."

**1:20-1:50 -- Risk Detection**
The agent analyzes Gong transcripts via TrueFoundry. Risk signals surface: "Competitor DataDog mentioned in 2 of 3 calls. Meeting frequency declining. VP Engineering (champion) absent from last meeting." Risk assessment: RED. Show the strategy brief being generated and stored in Ghost: recommended pricing, likely objections, recovery playbook.

**1:50-2:25 -- The Phone Call**
"The agent has determined this deal is at risk. Watch what happens next." The phone rings. Play it on speaker. The Bland AI voice delivers the coaching script with the three specific risk signals and the recovery recommendation. The room hears it. Pause.

**2:25-2:50 -- The Learning Loop**
"Every analysis is stored. Every outcome feeds back. The more deals the agent sees, the better its pattern matching gets." Show Aerospike's deal memory growing. Show Ghost's database of deal briefs. Show Overmind's cost dashboard: "This analysis cost $0.08 and took 15 seconds."

**2:50-3:00 -- Landing**
"Sales reps spend hours researching deals and miss risk signals until it's too late. This agent researches every deal in seconds, finds patterns in your history, and calls you before you lose the deal. It never forgets to check. It never misses a signal. And it calls you when it matters."

## Technical Architecture

```
                     ┌──────────────┐
                     │  CRM Webhook │
                     │  (simulated) │
                     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  Orchestrator │  (Python agent)
                     │    Agent      │
                     └──┬──┬──┬──┬──┘
                        │  │  │  │
         ┌──────────────┘  │  │  └────────────────┐
         │                 │  │                    │
         ▼                 ▼  ▼                    ▼
  ┌─────────────┐  ┌────────────────┐     ┌──────────────┐
  │   Airbyte   │  │  TrueFoundry   │     │   Bland AI   │
  │  Connectors │  │  LLM Gateway   │     │  Voice API   │
  │             │  │                │     │              │
  │ - Gong      │  │ Route:         │     │ Phone call   │
  │   (calls,   │  │ - Extract →    │     │ to AE with   │
  │   transcripts) │   fast model   │     │ coaching     │
  │ - Stripe    │  │ - Analyze →    │     │ script       │
  │   (revenue) │  │   strong model │     │              │
  │ - GitHub    │  │ - Strategy →   │     │ (triggered   │
  │   (tech     │  │   strong model │     │  only on RED │
  │   stack)    │  │                │     │  risk deals) │
  └──────┬──────┘  └───────┬────────┘     └──────────────┘
         │                 │
         │     ┌───────────┴───────────┐
         │     │                       │
         ▼     ▼                       ▼
  ┌──────────────┐            ┌──────────────┐
  │  Aerospike   │            │    Ghost     │
  │              │            │   Postgres   │
  │ - Deal       │            │              │
  │   embeddings │            │ - Deal       │
  │ - Vector     │            │   records    │
  │   similarity │            │ - Strategy   │
  │   search     │            │   briefs     │
  │ - Agent      │            │ - Risk       │
  │   memory     │            │   history    │
  └──────────────┘            └──────────────┘
         │                          │
         └──────────┬───────────────┘
                    │
             ┌──────▼──────┐
             │  Overmind   │
             │  SDK        │
             │  (wraps all │
             │  LLM calls) │
             └─────────────┘
```

**Data Flow:**
1. CRM webhook triggers the orchestrator with deal metadata
2. Orchestrator queries Airbyte connectors in parallel: Gong transcripts, Stripe revenue, GitHub repos
3. Raw data sent to TrueFoundry for extraction (fast model): structured deal profile emerges
4. Deal profile embedded and queried against Aerospike: top 5 similar deals returned
5. Gong transcripts sent to TrueFoundry for risk analysis (strong model): risk signals identified
6. All context (deal profile, similar deals, risk signals) sent to TrueFoundry for strategy generation (strong model)
7. Strategy brief and risk assessment stored in Ghost; deal embedding stored in Aerospike
8. If risk = RED, coaching script generated and Bland AI places a phone call to the AE
9. Overmind instruments all LLM calls across the pipeline
10. Deal outcome (when eventually recorded) feeds back into Aerospike for future pattern matching

## Buildability Risk Assessment

**Overall buildability: 3/5 (tight but doable with disciplined scoping)**

**What could go wrong in 8 hours:**

1. **Airbyte connector setup for three sources (HIGH RISK):** Gong, Stripe, and GitHub are all supported Airbyte connectors, but configuring three connectors with proper auth, schema mapping, and data extraction in a hackathon is ambitious. Any single connector that doesn't cooperate eats an hour. **Mitigation:** Set up Airbyte connectors the night before if rules allow. Have pre-extracted JSON fallbacks for each source. GitHub is the easiest (public API, no auth needed for public repos) -- start there. Stripe test mode has good sandbox data. Gong is the highest risk -- have sample transcript data ready.

2. **Gong API access (HIGH RISK):** Gong's API requires an enterprise account and API key. If you don't have a Gong account with call data, you're faking the most important data source. **Mitigation:** Use synthetic but realistic Gong transcript data. Be transparent in the demo: "this is structured like real Gong data." Alternatively, use Airbyte's Gong connector against a trial account with seeded calls.

3. **Bland AI voice quality and latency (MEDIUM RISK):** The phone call is the demo climax. If the voice sounds robotic, the script is garbled, or there's a 30-second delay before the call connects, the moment dies. **Mitigation:** Test Bland AI extensively. Pre-test the exact phone number. Have the coaching script mostly pre-composed with variable slots for the specific risk signals. Bland offers 100 free calls/day -- use 20 of them testing.

4. **Vector search quality (MEDIUM RISK):** The "similar deals" pattern matching is only as good as the embeddings and the historical data. With seeded data, the similarity results might feel contrived. **Mitigation:** Seed 20-30 diverse deal profiles in Aerospike. Make sure the demo deal has obvious similarity to 4-5 of them. The results will be controlled, but the mechanism is real.

5. **Six sponsor integrations (MEDIUM-HIGH RISK):** Six tools is a lot of integration surface. Each SDK has setup time, auth, and debugging. **Mitigation:** Prioritize ruthlessly: Airbyte + Aerospike + Bland AI first (they're the core). TrueFoundry next (it's the LLM layer). Ghost and Overmind last (they're nice-to-haves).

**The hardest part:** Getting three Airbyte connectors working reliably, with real or realistic data, in the first 3 hours. This is the make-or-break.

**What to cut if time runs out:**
- **Cut Overmind** (cosmetic, add in last 30 min or skip)
- **Cut Ghost** (use Aerospike for structured storage too, or just print to console)
- **Reduce Airbyte from 3 connectors to 1** (GitHub only, pre-seed Gong and Stripe data)
- **Simplify the strategy brief** (just show risk signals + similar deals, skip the full playbook)
- **NEVER cut:** The phone call. That IS the demo. Build backward from the phone ringing.

**Hour-by-hour plan:**
- Hours 1-2: Airbyte connectors (GitHub first, then Stripe sandbox, then Gong or mock). Aerospike schema + seed historical deals.
- Hours 3-4: Agent orchestrator. TrueFoundry integration for extraction and analysis. Vector search working end-to-end.
- Hours 5-6: Bland AI phone call integration. Risk detection logic. Coaching script generation.
- Hours 6-7: Ghost integration. End-to-end demo flow. Strategy brief output.
- Hours 7-8: Overmind integration. Polish. Rehearse demo 5+ times. Test phone call 5+ times.

## Honest Weaknesses

**1. The data is almost certainly synthetic.**
Gong requires an enterprise account with real sales calls. Stripe needs a real customer base. GitHub is the only genuinely "live" data source. A skeptical judge will ask: "Is this Gong data real?" And the honest answer is probably no. The architecture is real, the connectors are real, but the data is seeded. This weakens the "autonomous" narrative because the agent is working with curated inputs.

**2. The Airbyte use case, while perfect on paper, may be fragile in practice.**
Airbyte's agent connectors for Gong, Stripe, and GitHub are listed as available, but "available" and "works flawlessly in a hackathon with zero debugging" are different things. If the Airbyte connectors require extensive configuration or have rate-limiting issues, you could burn hours on plumbing. The irony: the sponsor tool that fits best on paper could be the one that causes the most build pain.

**3. Six sponsor tools risks looking like sponsor-stacking.**
Using 6 tools when you need 3 might impress or might look like resume-padding. Ghost and Overmind are genuinely useful but not essential. If a judge probes "why do you need both Aerospike and Ghost?", the answer is defensible (vector vs. relational) but the question itself creates doubt.

**4. The autonomous loop has a simulation problem.**
The trigger is a simulated CRM webhook. The Gong data is likely synthetic. The "similar deals" are seeded. The phone call goes to a known number. Everything works, but nothing is truly production-scale autonomous. A judge focused on "does this work in the real world?" will see the scaffolding.

**5. The sales/CRM domain may not resonate with security-focused judges.**
This is RSAC -- a security conference. The judges may be more impressed by security-adjacent problems (threat detection, vulnerability triage) than sales automation. "Deal intelligence" is a B2B SaaS problem, not a cybersecurity problem. The agent's sophistication has to overcome the domain mismatch.

**6. Bland AI call timing is a single point of demo failure.**
If the phone doesn't ring during the demo, you've lost your "whoa" moment. Network issues, Bland API latency, or a misconfigured phone number turns your climax into an awkward pause. Must have a bulletproof fallback (pre-recorded call audio, or a teammate pretending to receive the call on a second device).

## Final Score Recommendation

| Dimension | Score | Weight | Weighted | Justification |
|-----------|-------|--------|----------|---------------|
| **Autonomy** | 5 | 3x | 15 | This is a textbook autonomous agent loop. Trigger fires, agent pulls data from 3 sources, runs vector search, analyzes risk, generates strategy, and CALLS A HUMAN ON THE PHONE without any human initiation. The phone call is the strongest autonomy proof in the entire idea set -- the agent doesn't just write a report, it takes real-world action. |
| **Idea** | 4 | 3x | 12 | The problem is real -- every sales rep does this research manually. But it's less viscerally relatable than the pricing ideas because not every judge runs a sales pipeline. The "deals like this close 73% of the time" insight is compelling, but the pain isn't as universal as "your competitor just changed their pricing" or "$18K surprise bill." Docked from 5 because the RSAC audience may not connect with sales ops pain. |
| **Technical Implementation** | 4 | 2x | 8 | Multi-source data ingestion, vector similarity search, LLM-driven risk analysis, dynamic voice script generation, multi-model routing. This is a genuinely sophisticated pipeline. Not a 5 because the data is likely synthetic and the "real" agent reasoning (deciding what's risky, what's similar) depends heavily on the quality of seeded data. |
| **Tool Use** | 5 | 2x | 10 | This is the strongest tool-use story in the idea set. Airbyte (3 connectors, exactly what it's built for), Aerospike (vector search, the core intelligence), Bland AI (autonomous phone action), TrueFoundry (multi-model routing), Ghost (structured deal records), Overmind (observability). Three are genuinely essential, two are strong value-adds, one is cosmetic. The Airbyte fit alone is a 5 -- this is their dream use case. |
| **Presentation/Demo** | 5 | 2x | 10 | The phone ringing in the room is the best "whoa" moment of any idea evaluated so far. It's visceral, surprising, and proves autonomy in a way no dashboard can. The story arc builds naturally: data pull, pattern match, risk detection, PHONE RINGS. Even judges who don't care about sales will react to a phone call with specific, data-derived coaching. |
| **Buildability** | 3 | 1x | 3 | Six sponsor tools. Three Airbyte connectors. Gong data access. Bland AI testing. This is buildable but leaves zero margin for error. The Airbyte triple-connector setup is the critical path -- if all three work, the rest flows. If Gong's connector fights you, the whole schedule slides. The aggressive cut list (down to 1 Airbyte connector + pre-seeded data) preserves the demo but weakens the tool-use story. |

**Total: 15 + 12 + 8 + 10 + 10 + 3 = 58/65**

**Verdict: BUILD THIS.** This scores 58, just behind Pricing Intelligence (59) and ahead of Pricing Change Comms (56) and Margin-Aware Feature Flags (53). The phone call is arguably the single best demo moment across all ideas evaluated. The Airbyte fit is the strongest sponsor alignment in the set. The weakness is buildability -- six tools and synthetic data are real risks.

**How this compares to the current top ideas:**
- vs. #5 Pricing Intelligence (59): Slightly lower on Idea (sales ops vs. universally relatable pricing), equal on Demo, stronger on Tool Use. Very close. If you can't get Gong data, #5 wins. If you nail the phone call, #13 wins.
- vs. #7 Pricing Change Comms (56): Both use Bland AI for the "phone ringing" moment. #13 has a richer data pipeline (3 Airbyte sources vs. simulated churn data) but more build risk. #13 is the more ambitious version of the same "agent calls you" thesis.
- vs. #8 Margin-Aware Feature Flags (53): #13 has a better demo moment (phone call vs. dashboard going red) and stronger tool use. #8 has a more visceral problem statement ($18K bill).

**Key risk to monitor:** The Airbyte triple-connector setup. If you can get Gong + Stripe + GitHub flowing in the first 2 hours, this idea is a hackathon winner. If Gong's connector doesn't cooperate, immediately fall back to pre-seeded transcript data and focus the Airbyte story on Stripe + GitHub. The phone call must work -- test it 20 times before the demo.
