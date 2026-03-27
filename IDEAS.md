Mapping all 10 ideas to the transaction lifecycle:

  ┌───────────────────────────┬─────────────────────────────────────────┬────────────────────────────────┐
  │           Phase           │                 Problem                 │              Idea              │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ Before the product exists │ What will this feature cost at scale?   │ #8 Margin-aware feature flags  │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ Before the transaction    │ Agents can't buy things                 │ #1 Agent commerce rails        │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ Before the transaction    │ Buyer can't predict their cost          │ #10 Cost calculator for buyers │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ During pricing design     │ Which model for which request?          │ #4 Model router                │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ During pricing design     │ What do competitors charge?             │ #5 Pricing intelligence        │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ During pricing design     │ What if we change models?               │ #9 War-game simulator          │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ During pricing change     │ How to communicate without losing trust │ #7 Change communication        │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ During pricing change     │ Can we safely A/B test this?            │ #2 Pricing experimentation     │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ After the transaction     │ Can I forecast next month's spend?      │ #3 Spend forecasting           │
  ├───────────────────────────┼─────────────────────────────────────────┼────────────────────────────────┤
  │ After the deal            │ Are the contract terms right?           │ #6 Contract builder            │
  └───────────────────────────┴─────────────────────────────────────────┴────────────────────────────────┘


1. Agent-to-agent commerce rails

  Bharat Kumar (Key0, scored 94), Theodora Lau, and the agentic pricing brainstorm all point to the same gap:
  agents can't buy things. There's no machine-readable pricing, no programmatic payment path, no way for an
  agent to discover what an API costs and pay for it. Bear Lumen measures costs after the transaction. This
  problem is about enabling the transaction in the first place. If the future is agents calling agents calling
   APIs, the entire billing chain needs to be machine-readable. Pricing pages need an API. Payment needs to be
   programmatic. Budget enforcement needs to happen at the agent level, not the human level.

  2. Pricing experimentation (demand-side)

  Fynn Glover (95): "A credit cost. A usage limit. An overage threshold. These are weekly operational
  decisions." Wisam/Stripe (98): companies with 100%+ growth change pricing 3x more frequently. Bear Lumen
  tells you your cost-to-serve, which informs what you should charge. But nobody has good tooling to actually
  run the pricing experiment: A/B test pricing pages, roll out new rates to a customer segment, measure
  conversion and retention impact of the change, and roll back if it hurts. The gap is between "I know what to
   charge" and "I can safely change what I charge and measure whether it worked."

  3. AI spend forecasting and budget enforcement

  Henry Norris (91): $500K engineers burning $250K in tokens, no one can track which team drives the bill.
  Trey Harnden (93): burning 15M tokens/day, projecting 50M next year, no way to forecast. Matias Coca (90):
  AI costs buried under generic service names. Bear Lumen is retrospective: what did it cost? This problem is
  prospective: what will it cost next month? How do I set per-team or per-project token budgets and enforce
  them before the invoice arrives? It's FinOps for AI inference, aimed at the internal cost consumer
  (engineering teams), not the external pricing decision (product/finance teams).

  4. AI model routing / cost-quality optimizer

  Carmen Li (90): "OpenAI pricing isn't a number, it's a system." Output tokens cost 3-5x input. Context
  accumulation multiplies costs silently. Companies are picking models by vibes, not data. The missing tool:
  before you send a request, route it to the model that gives the best cost/quality tradeoff for that specific
   prompt type. Simple prompts go to Flash-Lite at $0.25/M tokens. Complex reasoning goes to Opus. The router
  learns which prompts need which models. This is a pre-transaction decision, not post-transaction
  measurement. Hackathon-friendly because you can demo it with a proxy that sits between the app and the LLM
  providers.

  5. Competitive pricing page intelligence

  Carmen Insignares Newell (95) manually analyzed 125 AI startup pricing pages. Koshima Satija (92) found 55%
  of voice AI companies hide pricing behind "contact sales." That manual research could be a product: scrape
  public pricing pages, normalize the data (per-seat vs per-token vs per-minute vs credits), track changes
  over time, alert when a competitor changes pricing. Companies making pricing decisions want to know "what
  does everyone else in my vertical charge, and how are they packaging it?" Nobody aggregates this. Everyone
  does it manually on spreadsheets.

  6. AI services contract builder

  Daljeet Saran (92): AI is rewriting software contract structures. Rahul Bhuman (92): milestone contracts
  fail because nobody agrees on what "done" means before signing. The pain: legal and procurement teams are
  writing AI service contracts using SaaS templates. They don't have clauses for usage variability, outcome
  SLAs, model-change rights, cost pass-through mechanisms, or quality degradation remedies. A contract clause
  library and builder specific to AI agreements. Pulls from real contract patterns, helps both buyers and
  vendors get the terms right before the deal starts. Especially relevant for the services efficiency paradox,
   where firms need contracts that protect their margin when delivery gets faster.

   7. Pricing change communication platform

  Glenn Turner (94) wrote a detailed post about Figma's AI credit rollout backlash. Alex Smith (85): "this is
  a preview of what's to come in AI pricing across most of your favorite AI tools." The problem isn't deciding
   the new price. It's communicating the change without destroying trust. No tool helps you: segment customers
   by how the change impacts them, draft personalized migration messaging per segment, predict which customers
   will churn vs accept vs upgrade, and time the rollout. Figma had a good product and a good pricing model.
  They lost trust because the communication was wrong. This happens every time a company changes pricing, and
  everyone does it manually with spreadsheets and gut feel.

  8. Margin-aware feature flags

  CostReveal (93): team ships AI email feature, expects $516/month, gets an $18K bill. The feature flag tools
  (LaunchDarkly, etc.) gate on user segments and rollout percentage. None of them gate on cost. The missing
  product: before you roll out an AI feature to 100% of users, it estimates cost impact based on beta usage
  data. "This feature will cost $14 per user per month at full rollout. Your plan charges $29. Margin-safe."
  Or: "This feature will cost $47 per user per month. Your plan charges $29. Do not roll out without a pricing
   change." Connects product decisions to financial outcomes before the CFO gets surprised.

  9. AI pricing war-game simulator

  Matt Green (97): GRR improved when moving BACK from consumption to per-seat. Nobody predicted that. The
  missing tool: take your current customer data, current pricing, and current costs, then simulate "what if we
   switch to usage-based?" or "what if our top 20 customers double usage?" or "what if Anthropic raises token
  prices 40%?" Monte Carlo simulations across your customer base showing revenue, margin, and churn impact of
  pricing model changes. Every company considering a pricing change is doing this in spreadsheets with made-up
   assumptions. A simulator with real customer data would turn a 3-week analysis into a 10-minute session.

  10. AI cost calculator for buyers

  Raj Khera (scored in dataset): prospects anchor to $20/month consumer AI pricing when the real value
  equation is $3K+. Alan Jacobson (95): "AI usage must be bounded, visible, and controllable." From the buyer
  side, the problem is: "what will this AI product actually cost me?" Vendors publish pricing pages with
  per-token or per-credit rates, but buyers can't translate that into a monthly budget number without doing
  math nobody wants to do. An embeddable cost calculator that takes the vendor's pricing model and the buyer's
   expected usage and shows the real monthly cost. Vendors embed it on their pricing page to build trust.
  Buyers use it to get budget approval. Both sides benefit from fewer surprises post-sale.
