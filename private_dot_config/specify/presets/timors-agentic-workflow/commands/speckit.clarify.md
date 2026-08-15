---
description: Converge a Feature Definition through one-question human dialogue.
---

## Authority and Inputs

Read canonical `spec.md`, `references/artifact-contracts.md`,
`references/feature-definition.md`, `references/clarification.md`, and
`references/review-lifecycle.md`. `spec.md` is the only authority for product
intent. Do not inspect repository reality or make technical design choices.

## Procedure

1. Start from Draft or Ready. Proportionally challenge the stated solution versus
   underlying need, consequential assumptions, actors, boundaries, constraints,
   failure conditions, and adjacent effects. Do not ask a question already settled
   in `spec.md`.
2. If a Ready artifact reveals a material question, first persist it, set Draft,
   and apply the ordered invalidation before dialogue. A package lock or completion
   refuses a material change.
3. Each assistant turn that needs a human answer asks exactly one direct,
   highest-impact question. Explanatory text cannot hide another question. When
   there are viable choices, give two or three genuine alternatives, concrete
   tradeoffs, and a reasoned recommendation; retain human decision ownership.
4. Immediately after each answer, update all affected Feature Definition sections,
   remove contradictions and resolved questions, record consequential decision
   context, apply invalidation, and rescan the complete Ready predicate. Do not
   batch unanswered questions or silently select an answer.
5. Mark Ready only when the complete predicate holds. A pause or declined material
   choice remains explicitly recorded and Draft. Report either the one next
   question or the Ready handoff to `plan`.
