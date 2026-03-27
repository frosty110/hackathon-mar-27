# Macroscope — Pricing & Competitive Strategy

**Owner:** Dana Chen, Head of Product
**Last updated:** 2026-03-20
**Status:** Living document — reviewed monthly at Pricing Council
**Distribution:** Internal only — exec team + GTM leads

---

## 1. Company Positioning

Macroscope is a **premium AI-powered codebase analysis and PR review platform**. We are not a linter. We are not a static analysis tool with an LLM bolted on. We are the only product that combines deep AST-level understanding of code structure with large language model reasoning to find real bugs, explain complex changes, and answer natural language questions about any codebase.

Our buyer is a VP of Engineering or Staff+ IC at a company with 20–500 engineers who is tired of shallow AI code review tools that flag style nits while missing actual logic errors. They have tried CodeRabbit or turned on Copilot's review features and found the signal-to-noise ratio unacceptable.

**One-liner:** "AI code review that actually understands your codebase — not just the diff."

We win on **depth of analysis** (AST + LLM means we catch bugs that surface-level pattern matching never will), **integration breadth** (GitHub, Slack, Linear, Jira — we meet teams where they already work), and **actionable output** (we do not just flag problems, we explain why they are problems and suggest fixes with full context of the surrounding code).

The AST layer is the moat. Every competitor either does pure LLM-on-diff (fast but shallow) or pure static analysis (deep but rigid). We do both, and the combination is what lets us catch the class of bugs that matter most: logic errors, race conditions, subtle API misuse, and security issues that only surface when you understand the call graph.

---

## 2. Current Pricing Model

**Structure:** Per-developer seat, simple and predictable.

| Plan | Monthly per dev | Minimum | Target customer |
|------|----------------|---------|-----------------|
| **Team** | $30/dev | 5 seats ($150/mo floor) | Engineering teams of 5–50. Want PR review, codebase Q&A, and Slack/Linear integration. |
| **Business** | $30/dev (volume discounts at 100+) | 25 seats | Mid-market, 50–200 devs. Need SSO, audit logs, custom review policies, Jira integration. |
| **Enterprise** | Custom (typically $25–$30/dev) | 100 seats | 200+ devs. On-prem/VPC deployment, custom model tuning, dedicated support, SLAs. |

**No free tier.** We offer a 14-day free trial with full functionality. Prospects can connect a real repo and see real results before they pay anything.

**Annual discount:** 15% on Team and Business. Enterprise is always annual.

**Current blended ARPU:** ~$29/dev/month (slight drag from early Enterprise deals at lower rates — tightening this).

**Seat mix:** 22% Team, 48% Business, 30% Enterprise (by revenue).

---

## 3. Competitive Positioning

| Competitor | Category | How we position against them |
|-----------|----------|------------------------------|
| **CodeRabbit** | AI code review (LLM-on-diff) | The most direct comp. Good UX, fast, popular with small teams. But it is pure LLM — no structural analysis. It reads the diff like a human would, which means it misses anything that requires understanding code outside the diff. We win on depth. In competitive deals, run a head-to-head on a real PR with a known bug — Macroscope catches things CodeRabbit cannot. Their free tier pulls in hobbyists and small OSS projects that are not our market anyway. |
| **Codacy** | Automated code quality | Legacy-feeling product. Rule-based at core with AI features layered on top. Slow to adapt to new languages and frameworks. Position us as the modern replacement: "You outgrew Codacy the day you started writing code that rule engines cannot understand." They compete on breadth of checks; we compete on intelligence of checks. |
| **DeepSource** | Static analysis + AI | Strongest technical competitor in the static analysis lane. Good AST work, but their AI layer is thin — it is mostly used for auto-fix suggestions, not deep reasoning about code intent. We have a real LLM reasoning layer; they have a suggestion engine. In bake-offs, focus on complex PRs where DeepSource flags surface issues while Macroscope identifies the actual logic bug. |
| **Sourcery** | AI code review (Python) | Python-only. That is the positioning. Any team with a polyglot codebase (which is nearly everyone) cannot standardize on Sourcery. Respect their Python depth but position Macroscope as the platform choice. Do not waste cycles competing for Python-only shops — let Sourcery have them. |
| **Qodo / CodiumAI** | AI test generation + review | Different wedge — they lead with test generation, review is secondary. We lead with review, codebase understanding is the platform. In practice we see them in deals where the buyer wants "AI for code quality" broadly. Position Macroscope as the review brain and acknowledge Qodo as complementary for test gen. Do not trash them — some of our best customers use both. |
| **SonarQube / SonarCloud** | Code quality platform | The 800-lb gorilla of code quality. Installed everywhere. But SonarQube is a rule engine — thousands of rules, zero reasoning. It catches code smells and known vulnerability patterns. It does not catch "this PR introduces a race condition because the author did not realize these two services share state." We are not replacing Sonar; we are the layer above it. Position as complementary in Enterprise deals. In Team/Business deals, position as the modern alternative. |
| **Snyk Code** | Security-focused analysis | Security-first positioning. Strong brand in AppSec. We overlap on security findings but we are broader — we catch bugs, summarize changes, answer questions, AND find security issues. Snyk is a must-have for security teams; Macroscope is a must-have for engineering teams. Different buyer, different budget. Avoid head-to-head on pure security — we will lose to their vuln database depth. Win on "security + everything else." |
| **Graphite** | PR workflow + AI review | Graphite's core is PR workflow (stacking, merging, dashboards). AI review is a feature, not the product. We are the opposite — AI analysis IS the product. If a prospect uses Graphite for workflow, great, Macroscope integrates alongside it. If they are evaluating Graphite's AI review vs. ours, run the depth comparison. Their review is surface-level because it is a feature, not a focus. |
| **Ellipsis** | AI code review | Small but sharp. Good product instincts, similar positioning to CodeRabbit but less established. Watch them — they iterate fast. Same playbook as CodeRabbit: they do LLM-on-diff, we do AST + LLM. Depth wins. |
| **Bito** | AI code assistant | Broader scope — code assistant, not just review. They do chat, documentation, review. Jack of all trades. We are purpose-built for codebase analysis and PR review. In competitive situations, the "purpose-built vs. Swiss army knife" framing works well. |

**Key differentiators we protect:**

- **AST + LLM fusion.** This is the core technical moat. Every competitor either does one or the other. We do both and that is why we find bugs they miss. Protect this in messaging, demos, and hiring (keep investing in the AST engine team).
- **Codebase Q&A.** "Ask your codebase a question in English" is a feature no pure-review tool offers at our quality level. It is also the feature that gets champions internally — an engineer asks "where is the rate limiting logic?" and gets an accurate answer in 10 seconds. That is how we spread inside orgs.
- **Integration depth.** GitHub App + Slack + Linear + Jira. We are not just a GitHub check that posts a comment. We push findings into the tools teams actually use to track work.

---

## 4. Q2 2026 Growth Targets

| Metric | Q1 2026 actual | Q2 2026 target | Notes |
|--------|---------------|----------------|-------|
| ARR | $3.2M | $4.0M | 25% QoQ — aggressive but achievable with pipeline |
| New logo customers | 38 (Q1) | 50 | Deep Agents Hackathon (RSAC) sponsorship should generate 15+ qualified leads alone |
| Net revenue retention | 112% | 118% | Seat expansion is the lever — land with one team, spread to org |
| Avg seats per new customer | 12 | 15 | Push Business tier, de-emphasize 5-seat minimums in outbound |
| Enterprise pipeline | $1.4M | $2.5M | 2 Enterprise AEs starting April, focus on fintech and healthtech verticals |
| Trial-to-paid conversion | 22% | 28% | New onboarding flow shipping mid-April — guided first-PR-review experience |
| Blended ARPU | $29/dev | $29.50/dev | Hold the line. Do not discount to win deals. |

**Q2 initiatives:**
- **Hackathon-driven awareness.** Deep Agents Hackathon (RSAC 2026) sponsorship. Goal: 200 trial signups from event attendees. Follow-up sequence with case study content targeting security-conscious buyers (RSAC audience).
- **"Bugs Caught" reporting.** Ship a monthly digest showing each customer how many real issues Macroscope caught. This is the single best retention and expansion lever — hard numbers on value delivered.
- **Linear + Jira integration GA.** Currently in beta with 12 customers. Full launch in April. This unlocks the workflow integration story for Business and Enterprise.

---

## 5. Pricing Principles

These are our standing rules. The Pricing Council (Dana, CFO, CRO) can override, but these are the defaults.

1. **Price on the value of bugs caught, not seats.** $30/dev/month is our current model because it is simple and predictable. But internally, we justify the price by the cost of bugs that reach production. One prevented P1 incident is worth more than a year of Macroscope for a 50-person team. Every sales conversation should anchor on this math, not on per-seat comparisons.

2. **Never race to the bottom with freemium tools.** CodeRabbit, Sourcery, and others offer free tiers or aggressive discounts. We do not. Our product is better, our analysis is deeper, and teams that need what we offer know the difference. If a prospect says "CodeRabbit is free," they are not our buyer yet. Nurture them until they hit the complexity wall.

3. **Hold $30/dev/month for Team and Business through 2026.** No price increases this year — we need market share. No price decreases either — we are already well below the pain threshold for any team that ships production code. $150/month minimum for a 5-person team is a rounding error on an engineering budget.

4. **Enterprise discounts floor at $22/dev/month.** Below that we are giving away margin for a logo. The only exception is a design partner deal with contractual case study rights and product feedback commitments.

5. **Compete on depth, not on price.** If a competitor cuts their price, do not respond with a price cut. Respond with a comparison demo. Ship the "head-to-head" feature (queued for Q3) that lets prospects run Macroscope and a competitor on the same PR and see the difference side by side.

6. **Annual contracts get 15%, nothing more.** We saw the mistake other startups made offering 30–40% annual discounts. It trains buyers to expect it and destroys pricing power at renewal. 15% is enough to incentivize annual without creating a renewal cliff.

7. **Seat minimums exist for a reason.** 5-seat minimum on Team stays. It filters out individual developers and hobbyists who generate support load without meaningful revenue. If an individual developer wants Macroscope, point them to the trial and let them bring it to their team.

---

## 6. Red Lines

Things we will NOT do regardless of competitive pressure:

- **No public free tier.** Free tiers attract users who will never pay and create noise in our usage data, support queue, and infrastructure costs. The 14-day trial with full access is the right wedge. Prospects see real value on real code before they commit. If the data ever shows a free tier would materially improve paid conversion, we revisit — but the bar is high.

- **No per-token or per-analysis pricing.** Engineers hate unpredictable bills. "Your AI code review cost $847 this month because you had a lot of PRs" is a nightmare. $30/dev/month is simple, predictable, and budgetable. Keep it that way.

- **No white-labeling.** We have had inbound interest from consulting firms and agencies wanting to resell Macroscope under their brand. The answer is no. Our brand is our distribution. We will build a partner/referral program, but the user always knows they are using Macroscope.

- **No single-repo pricing.** Pricing per-repo incentivizes monorepos and penalizes good architecture. Per-developer is the right unit because it scales with the team, not with how they choose to organize their code.

- **No "security-only" tier at a lower price.** Snyk and Sonar own the "just security" or "just code quality" lanes. If we offer a stripped-down security-only product at $15/dev, we cannibalize our own full-platform deal. We are the full-stack analysis platform. The buyer gets everything.

- **No matching CodeRabbit's free tier.** If CodeRabbit stays free for open source or small teams, fine. That is their growth strategy. Ours is different — we sell to professional engineering teams who budget for developer tools. We are not the same product at a higher price; we are a different, deeper product.

---

*Next review: April 18, 2026 — Pricing Council monthly sync*
*Escalation for urgent competitive moves: #pricing-war-room Slack channel*
*Deep Agents Hackathon follow-up review: April 4, 2026*
