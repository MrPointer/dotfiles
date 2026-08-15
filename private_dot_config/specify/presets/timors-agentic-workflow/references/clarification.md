# Intent Convergence

This phase uses only saved Feature Definition content and current human dialogue.
It proportionally challenges consequential assumptions, distinguishes underlying
need from a proposed solution, and explores relevant actors, boundaries,
constraints, failures, and adjacent effects. It stops when further inquiry cannot
materially change intent, scope, requirements, evidence, constraints, or decision
context. It does not inspect repository facts or choose implementation design.

Each response-requiring assistant turn asks exactly one direct highest-impact
question. Explanation or alternatives may precede it but cannot contain a second
question. When a material choice has viable answers, present two or three genuine
alternatives with concrete tradeoffs and a reasoned recommendation; do not use
straw choices, silently choose, or ask an already answered question.

After every answer, immediately persist and reconcile every affected section,
remove contradictions and resolved questions, record consequential decision
context, apply the invalidation transition in `references/artifact-contracts.md`,
and rescan the complete Ready predicate. A pause or declined decision remains a
material open question and Draft. A lock or completed package refuses a material
write rather than leaving changed content authorized by an older gate.
