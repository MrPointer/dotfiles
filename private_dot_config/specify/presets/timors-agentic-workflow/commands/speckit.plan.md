---
description: Ground, design, review, and obtain acceptance for the normative RFC.
---

## Authority and Inputs

Require canonical `spec.md` with status Ready. Read
`references/artifact-contracts.md`, `references/planning-grounding.md`,
`references/rfc-planning.md`, `references/review-lifecycle.md`,
`references/model-and-worker-selection.md`, and reviewer packets
`reviewers/rfc-design.md` and `reviewers/rfc-clarity.md`. Use
`templates/grounding-notes-template.md`, `templates/research-template.md`,
`templates/plan-template.md`, and `templates/review-report-template.md`. A project
constitution is optional governance input. `plan.md` is the sole normative RFC.

## Procedure

1. Gather cheap delegated observations, normalize sources, facts, bounded
   inferences, relevance, and unknowns in `grounding-notes.md`, and centrally
   verify material claims. Keep recommendations out of that artifact. Use
   `research.md` for evidence, genuine alternatives, working decisions,
   consequences, and technical questions.
2. Ask one adaptive consequential technical question at a time when human input is
   needed; routine implementation detail remains planner-owned. Surface repository
   evidence that contradicts intent. Correct only non-semantic metadata or an
   unambiguous transcription in `spec.md`; any other correction returns it to
   Draft and `clarify` before this command continues.
3. Write a self-contained RFC that restates every durable decision needed by task
   planning. Classify changes before writing, apply Revision, acceptance, and
   review-pointer transitions exactly as `references/rfc-planning.md` requires,
   and refuse a material change after the package lock.
4. Bind, invoke, and record independent design and clarity reviews under
   `references/review-lifecycle.md`. Remediate accepted findings. Rerun only
   applicable scope; retain unaffected roles only through the stated retained
   pointer transition. At review quiescence present one collected blocker
   interaction. A started incomplete review uses recovery, never self-review.
5. Request Design Acceptance only when both pointers are current, attributable,
   and non-blocking for the current RFC Revision. Record the human decision and
   rationale in `plan.md`. Only current Revision Accepted permits `tasks`.
