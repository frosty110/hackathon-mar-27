# Idea #11: Autonomous Breach Notification Agent

## One-Line Pitch

An AI agent that autonomously handles GDPR/CCPA breach notification end-to-end — from breach report ingestion to personalized customer calls — beating the 72-hour regulatory clock without a war room.

## The Problem

When a data breach occurs, regulations like GDPR (72 hours), CCPA (expedient timeframe), and PIPEDA (as soon as feasible) mandate that affected individuals be notified promptly. In practice, this means:

- **Engineering** scrambles to scope the breach — which systems, which data, which users
- **Legal** determines jurisdiction-specific obligations for each affected user
- **Comms** drafts notification templates that satisfy legal requirements without causing unnecessary panic
- **Customer Success** personally contacts high-value enterprise accounts
- **Compliance** documents everything to prove timely notification to regulators

This war room operates around the clock with spreadsheets, ad-hoc scripts, and frantic Slack threads. It routinely costs $50K-$200K in labor alone, and the risk of missing a regulatory deadline is existential (GDPR fines up to 4% of global revenue). Every CISO at RSAC has lived this nightmare. It is one of the most painful, high-stakes operational workflows in security.

## The Autonomous Agent Loop

The agent operates as a closed loop triggered by a breach report:

```
BREACH REPORT (input)
    |
    v
[1] INGEST & SCOPE — Parse breach report, identify affected systems and data types
    |
    v
[2] PULL AFFECTED USERS — Query CRM/billing for users tied to affected systems
    |
    v
[3] CROSS-REFERENCE — Match users against breach scope (which specific users had which data exposed)
    |
    v
[4] SEGMENT — Classify each user by:
      - Jurisdiction (GDPR, CCPA, PIPEDA, etc.) based on location
      - Customer tier (enterprise, mid-market, self-serve)
      - Data sensitivity (PII, financial, health, credentials)
    |
    v
[5] DRAFT NOTIFICATIONS — Generate legally-compliant, personalized notifications per segment
      - Different templates per regulation
      - Personalized with specific data types exposed
      - Appropriate tone per tier
    |
    v
[6] EXECUTE OUTREACH
      - Email notifications to all affected users
      - Phone calls to enterprise/high-value customers
      - Track delivery status per user
    |
    v
[7] COMPLIANCE REPORTING — Generate audit trail proving:
      - Timeline of actions taken
      - Who was notified, when, via what channel
      - Regulatory requirements met per jurisdiction
    |
    v
COMPLIANCE REPORT (output)
```

**Autonomy depth:** Steps 1-7 execute without human intervention. The agent makes real decisions: which regulation applies, what to say, who to call vs. email, how to personalize. This is not a pipeline — the agent reasons about edge cases (dual-jurisdiction users, data types with special handling requirements, VIP customers who need different messaging).

## Sponsor Stack (3+ required)

### 1. Ghost (Postgres MCP) — ESSENTIAL

**What it does:** Agent autonomously creates a per-breach Postgres database to track affected users, notification status, and compliance timeline.

**Why it's load-bearing:** Every breach is unique — different scope, different affected users, different jurisdictions. The agent needs to create a structured data store on the fly, design its own schema based on breach characteristics, and query it throughout the notification process. This is not a pre-provisioned database — the agent literally spins up infrastructure autonomously.

**Verdict: ESSENTIAL.** This is genuinely agentic database use. The agent decides the schema, creates the DB, and uses it as working memory throughout the entire loop. Removing Ghost would mean pre-building the database, which defeats the autonomy story.

### 2. Bland AI (Voice Calls) — ESSENTIAL

**What it does:** Makes live phone calls to enterprise customers, delivering personalized breach notifications with context about what was exposed and what they should do.

**Why it's load-bearing:** Enterprise customers expect personal outreach during a breach — an email alone signals you don't care about them. The agent calls them, explains the situation, answers basic questions, and logs the call outcome. This is where the demo becomes visceral.

**Verdict: ESSENTIAL.** The phone ringing live during demo is the single strongest moment. It also represents a genuinely different modality — the agent isn't just writing text, it's speaking to humans. Removing Bland AI reduces this to "another email-drafting agent."

### 3. TrueFoundry (Multi-LLM Gateway) — ESSENTIAL

**What it does:** Routes different tasks to different models based on complexity and sensitivity. Boilerplate CCPA notifications go to a fast/cheap model. Personalized enterprise call scripts go to the best available model. Legal analysis of jurisdiction edge cases goes to a reasoning model.

**Why it's load-bearing:** Breach notification involves wildly different cognitive tasks — classification, legal reasoning, empathetic writing, structured data extraction. Using one model for everything is either too expensive or too dumb. TrueFoundry lets the agent intelligently allocate model resources.

**Verdict: ESSENTIAL.** This is defensible tool use — you can explain exactly why different models handle different parts. Judges will see this as sophisticated engineering, not checkbox-checking.

### 4. Airbyte (Data Connectors) — STRONG NICE-TO-HAVE

**What it does:** Pulls customer and billing data from Stripe (or similar) to identify who is affected and what tier they belong to.

**Why it's load-bearing:** The agent needs customer data to do its job. Airbyte provides a standardized connector to pull from real SaaS systems rather than hardcoded mock data.

**Verdict: STRONG NICE-TO-HAVE.** In an 8-hour build, you could fake this with a pre-loaded CSV and lose very little demo impact. But if you can get the Airbyte Stripe connector working, it adds legitimacy — the agent is pulling from a real billing system. Worth attempting but have a fallback.

### 5. Aerospike (Fast DB + Vector Search) — COSMETIC

**What it does:** Could store vector embeddings of past breach reports for "similar breach" lookups, or serve as a fast cache for customer jurisdiction lookups.

**Why it's load-bearing:** It isn't, really. Ghost already handles the core database needs. Aerospike's sub-millisecond lookups don't matter when the bottleneck is LLM inference and phone calls taking minutes.

**Verdict: COSMETIC.** You could shoehorn in a "find similar past breaches to inform notification strategy" feature, but it would feel forced. Only add this if you need a 6th tool for scoring and have time to spare.

### 6. Overmind (LLM Observability) — NICE-TO-HAVE

**What it does:** Monitors the agent's LLM calls for PII leakage, cost, and performance. In a breach notification context, you absolutely cannot have the agent leaking PII in its own logs or intermediate reasoning.

**Why it's load-bearing:** The security angle is real — an agent handling breach data that itself leaks data would be catastrophic. But the demo impact is low (it's a monitoring dashboard, not a visible action).

**Verdict: NICE-TO-HAVE.** Good narrative fit for RSAC ("we even monitor the agent for PII leakage"), easy to integrate if it's a drop-in SDK, but not a demo centerpiece.

**Tool count: 3 essential + 1-3 nice-to-haves = comfortably clears the 3-tool minimum.**

## The "Whoa" Demo Moment

**The phone rings.** You're demoing breach notification, the agent has been processing affected users, segmenting jurisdictions, drafting emails — and then a real phone on the demo table starts ringing. The agent is calling an enterprise customer. A judge picks up. The AI voice says: "This is an urgent notification from [Company]. We're calling to inform you that a security incident has affected your account. Specifically, your billing email and API keys associated with your enterprise plan were exposed. We have already rotated your API keys and..."

This is the strongest possible demo moment for several reasons:
- **It crosses the digital-physical boundary** — the agent's actions manifest in the real world
- **It's unexpected** — judges have seen a million dashboards and chat UIs
- **It's emotionally resonant** — every CISO knows the dread of making those calls
- **It proves autonomy** — the agent decided who to call, what to say, and when

**Secondary whoa:** Show the Ghost-created database being queried live — the agent designed the schema itself, and you can run arbitrary SQL against it to prove it's real.

## 3-Minute Demo Script

**[0:00-0:30] Setup — The Breach**
"It's 2am. Your security team just discovered a breach — an exposed S3 bucket containing customer billing data. You have 72 hours under GDPR before regulators come knocking. Let's feed the breach report to our agent."

Paste/upload a breach report into the agent. Show the agent beginning to reason: "Parsing breach scope... affected system: billing-db-prod... exposed data types: email, billing address, payment method last-4..."

**[0:30-1:15] Autonomous Processing**
Show the agent working through its loop in a split-screen or live terminal:
- "Pulling affected customers from Stripe via Airbyte..." (show Airbyte connector running)
- "Cross-referencing 2,847 customers against breach scope... 1,203 affected"
- "Creating breach tracking database..." (show Ghost spinning up a Postgres DB)
- "Segmenting: 847 EU (GDPR), 291 California (CCPA), 65 Canada (PIPEDA)..."
- "Routing notification drafting: GDPR template to claude-haiku, enterprise personalization to claude-opus..." (show TrueFoundry routing)

**[1:15-2:00] The Notifications**
Show sample drafted notifications side by side:
- A GDPR-compliant email (citing Article 34, specifying data categories, including DPO contact)
- A CCPA notification (different structure, different required disclosures)
- An enterprise-tier personalized notification (references their specific account, usage, impact)

"The agent drafted 6 notification variants across 3 jurisdictions and 2 tiers. Let's look at our enterprise customers..."

**[2:00-2:40] The Phone Call**
"For enterprise customers, an email isn't enough. The agent has identified 12 enterprise accounts affected. Let's watch it make the first call."

The phone on the table rings. A judge or planted audience member picks up. The Bland AI voice delivers a personalized breach notification — professional, empathetic, specific to their account.

"That call was generated autonomously — the agent decided who to call, synthesized the script from their account data and the breach scope, and made the call via Bland AI."

**[2:40-3:00] Compliance Report**
"It's been 3 minutes, not 72 hours. Let's check our compliance status."

Show the Ghost database: query showing all 1,203 users with notification status. Show the generated compliance report: timeline of actions, proof of notification per jurisdiction, ready for regulators.

"Every notification. Every jurisdiction. Every enterprise call. Fully documented. In minutes, not days."

## Technical Architecture

```
┌─────────────────────────────────────────────────────┐
│                   AGENT ORCHESTRATOR                 │
│            (Python — LangGraph or similar)           │
│                                                      │
│  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
│  │ Breach   │  │ Customer │  │ Notification      │  │
│  │ Parser   │→ │ Matcher  │→ │ Engine            │  │
│  └──────────┘  └──────────┘  └───────────────────┘  │
│       │              │              │         │       │
│       │              │              │         │       │
└───────┼──────────────┼──────────────┼─────────┼──────┘
        │              │              │         │
   ┌────▼────┐   ┌────▼────┐   ┌────▼───┐ ┌───▼──────┐
   │TrueFound│   │ Airbyte │   │ Email  │ │ Bland AI │
   │   ry    │   │ (Stripe)│   │ (SMTP) │ │ (Calls)  │
   │LLM Gate │   │         │   │        │ │          │
   └─────────┘   └─────────┘   └────────┘ └──────────┘
        │              │              │         │
        │         ┌────▼────┐         │         │
        │         │  Ghost  │◄────────┴─────────┘
        │         │Postgres │  (all status tracked)
        │         └─────────┘
        │              │
   ┌────▼────┐   ┌────▼────────┐
   │Overmind │   │ Compliance  │
   │Observ.  │   │ Report Gen  │
   └─────────┘   └─────────────┘
```

**Key implementation details:**
- **Orchestrator:** LangGraph with tool-calling nodes. Each step is a node that can call tools and make decisions.
- **Breach report parsing:** Structured extraction from free-text breach reports. Use a reasoning model via TrueFoundry.
- **Customer data:** Airbyte Stripe connector pulls customer records. Fallback: pre-loaded CSV with realistic data.
- **Jurisdiction classification:** Rule-based (country/state mapping) augmented by LLM for edge cases (EU customers with US billing addresses).
- **Notification drafting:** Template-guided generation with per-regulation legal requirements injected as system prompts. Different models per complexity via TrueFoundry.
- **Phone calls:** Bland AI API with dynamically generated call scripts. Agent passes customer-specific context into the call prompt.
- **Tracking:** Ghost Postgres DB with tables for: affected_users, notifications_sent, phone_calls, compliance_timeline.
- **Compliance report:** SQL queries against Ghost DB formatted into a regulatory-ready document.

## Buildability Risk Assessment

| Component | Difficulty | Time Est. | Risk | Mitigation |
|-----------|-----------|-----------|------|------------|
| Breach report parser | Low | 1 hr | Low | Structured input format as fallback |
| Airbyte Stripe connector | Medium | 1.5 hr | Medium | Pre-loaded CSV fallback |
| Ghost DB creation via MCP | Medium | 1.5 hr | Medium | Pre-provisioned DB fallback (hurts story) |
| Jurisdiction segmentation | Low | 1 hr | Low | Rule-based with LLM polish |
| Notification drafting | Low | 1.5 hr | Low | Core LLM task, well-understood |
| Bland AI phone calls | Medium | 1.5 hr | Medium | API is straightforward but needs testing |
| TrueFoundry routing | Low | 1 hr | Low | Drop-in gateway, well-documented |
| Overmind integration | Low | 0.5 hr | Low | Drop-in SDK |
| Orchestrator + glue | Medium | 2 hr | Medium | Keep it simple, avoid over-engineering |
| Demo polish | — | 1 hr | — | Essential, do not skip |

**Total estimated: ~12.5 hours of work compressed into 8 hours.**

**Critical path risks:**
1. **Ghost MCP integration** — If the MCP interface is flaky or poorly documented, creating DBs autonomously could eat hours. **Mitigation:** Start here. If it doesn't work in 90 minutes, fall back to a pre-provisioned Postgres and just use Ghost for queries.
2. **Bland AI call quality** — The demo lives or dies on the phone call sounding good. **Mitigation:** Test early, hardcode a known-good call script as fallback while dynamic generation is developed.
3. **Scope creep** — This idea has many moving parts. **Mitigation:** Define the demo-critical path (breach in -> segmentation visible -> phone rings -> compliance report out) and build only that. Everything else is polish.

**Build order (recommended):**
1. Hour 1-2: Ghost DB creation + Bland AI call (prove the two hardest integrations)
2. Hour 2-4: Core agent loop (parser -> matcher -> segmenter -> drafter)
3. Hour 4-5: TrueFoundry routing + Airbyte connector
4. Hour 5-6: End-to-end integration
5. Hour 6-7: Compliance report generation + Overmind
6. Hour 7-8: Demo rehearsal and polish

## Honest Weaknesses

1. **Simulated data, not real breach.** You can't demo a real breach. The customer data will be synthetic, the Stripe account will be a test account, the phone call will go to a planted number. Judges know this, but it still reduces the "is this real?" factor. Every team has this problem, though.

2. **Legal accuracy is hand-wavey.** The GDPR/CCPA notification templates will look plausible but a real lawyer would find gaps. The agent doesn't actually know law — it's pattern-matching against training data. For a hackathon this is fine, but an astute judge might probe here. Mitigation: Be upfront — "in production, templates would be lawyer-reviewed and the agent would select from pre-approved variants."

3. **Many moving parts for 8 hours.** Six sponsor tools, seven pipeline steps, two output modalities (email + phone). The risk of something breaking during live demo is real. Mitigation: Have fallbacks for every sponsor integration. The demo-critical path is: breach report -> segmentation display -> phone call -> compliance report. Everything else can be shown via pre-recorded backup.

4. **Email sending isn't actually impressive.** Half the pipeline is "draft and send emails" which is table-stakes for any LLM demo in 2026. The differentiation comes from the legal compliance angle, the jurisdiction segmentation, and the phone calls — not the email drafting itself.

5. **Autonomy is somewhat linear.** The loop is sequential (parse -> match -> segment -> draft -> send -> report). It doesn't branch, retry, or handle unexpected situations. A truly autonomous agent would handle things like "this customer's phone number is invalid, try their backup number" or "this enterprise account has a custom notification SLA." Adding these decision points would strengthen the autonomy story but may not be feasible in 8 hours.

6. **PII handling is a real concern.** An agent processing breach data with LLM calls is sending PII to model providers. In a real deployment, this would be a non-starter without on-prem models or strict data processing agreements. The Overmind integration helps the narrative, but doesn't solve the fundamental issue. Be prepared for this question from RSAC judges.

## Final Score Recommendation

| Criterion | Weight | Score | Weighted |
|-----------|--------|-------|----------|
| **Autonomy** | x3 | 4 | 12 |
| **Idea** | x3 | 5 | 15 |
| **Technical Implementation** | x2 | 3.5 | 7 |
| **Tool Use** | x2 | 4.5 | 9 |
| **Presentation/Demo** | x2 | 5 | 10 |
| **Buildability** | x1 | 3 | 3 |
| **TOTAL** | | | **56 / 65** |

**Breakdown:**

- **Autonomy (4/5):** The loop is genuinely autonomous — it makes real decisions about jurisdiction, customer tier, notification content, and outreach channel. Docked one point because the pipeline is mostly linear and the decisions, while real, aren't deeply branching. An agent that handles edge cases, retries, and adapts mid-loop would score 5.

- **Idea (5/5):** This is the best possible idea for an RSAC-colocated hackathon. Every CISO in the room has lived through breach notification. The 72-hour GDPR clock is universally understood and feared. The problem is real, expensive, and painful. The solution is immediately credible. This scores maximum on idea alone.

- **Technical Implementation (3.5/5):** Many integrations but each one individually isn't deeply technical. The orchestrator is a fairly standard tool-calling agent loop. The legal reasoning is prompt engineering, not novel architecture. The risk of something being half-baked due to the breadth of integrations pulls this down. A tighter, deeper implementation of fewer components would score higher.

- **Tool Use (4.5/5):** Three genuinely essential tools (Ghost, Bland AI, TrueFoundry) with clear, defensible reasons for each. Ghost is the standout — an agent creating its own database is a strong autonomy signal. Bland AI adds a physical-world modality. TrueFoundry demonstrates intelligent model routing. Airbyte adds a fourth. The only reason this isn't 5/5 is that some tools (Overmind, Aerospike) would feel forced if added.

- **Presentation/Demo (5/5):** The phone ringing on the demo table is the single strongest demo moment across all ideas evaluated. It's visceral, unexpected, and proves the agent took real-world action. The compliance report at the end ties a bow on it. The narrative arc (breach happens -> chaos usually ensues -> agent handles it in 3 minutes) is compelling and tight.

- **Buildability (3/5):** This is the weakest dimension. Six sponsor tools, seven pipeline steps, and two output modalities in 8 hours is aggressive. The critical path has at least two medium-risk integrations (Ghost MCP, Bland AI). The recommended build order mitigates this, but the team needs to be disciplined about cutting scope. A team that tries to build everything will ship nothing.

**Overall assessment: This is a top-tier idea.** The RSAC relevance is unmatched. The demo moment is the best available. The sponsor fit is natural. The main risk is execution — too many moving parts for 8 hours. A disciplined team that nails the critical path (breach in -> segmentation -> phone call -> compliance report) and has fallbacks for everything else will score extremely well.

**Recommendation: BUILD THIS — but ruthlessly cut scope to the demo-critical path.**
