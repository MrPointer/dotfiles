---
description: Ground, design, review, explain, and obtain acceptance for the normative plan.
---

## Authority and Inputs

Require canonical `spec.md` with status Ready. Resolve the canonical feature
directory from `.specify/feature.json` per
`.specify/presets/timors-agentic-workflow/references/artifact-contracts.md` before
reading any artifact. Read
`.specify/presets/timors-agentic-workflow/references/artifact-contracts.md`, `.specify/presets/timors-agentic-workflow/references/planning-grounding.md`,
`.specify/presets/timors-agentic-workflow/references/plan-lifecycle.md`, `.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`,
`.specify/presets/timors-agentic-workflow/references/model-and-worker-selection.md`, and reviewer packets
`.specify/presets/timors-agentic-workflow/reviewers/plan-design.md`, `.specify/presets/timors-agentic-workflow/reviewers/plan-clarity.md`, and
`.specify/presets/timors-agentic-workflow/reviewers/plan-overview-fidelity.md`, and writer packet
`.specify/presets/timors-agentic-workflow/writers/plan-overview.md`. Use
`.specify/presets/timors-agentic-workflow/templates/grounding-notes-template.md`, `.specify/presets/timors-agentic-workflow/templates/research-template.md`,
`.specify/presets/timors-agentic-workflow/templates/plan-template.md`, `.specify/presets/timors-agentic-workflow/templates/review-report-template.md`,
`.specify/presets/timors-agentic-workflow/templates/plan-overview-template.md`, and
`.specify/presets/timors-agentic-workflow/templates/review-dispositions-template.md`. A project
constitution is optional governance input. `plan.md` is the sole normative plan;
`plan-overview.md` is its non-normative human view.

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
3. Write a self-contained plan that restates every durable decision needed by task
   planning. Classify changes before writing, apply Revision, acceptance, and
   review-pointer transitions exactly as `.specify/presets/timors-agentic-workflow/references/plan-lifecycle.md` requires,
   and refuse a material change after the package lock.
4. Bind the writer under
   `.specify/presets/timors-agentic-workflow/references/model-and-worker-selection.md`
   and have it generate `plan-overview.md` for the current Revision. Present it to
   the human to confirm direction before any review. A change of direction is a
   human decision: record it in `plan.md` under step 3 (or return `spec.md` to
   Draft when intent changes), then regenerate the overview. This overview is not
   fidelity-reviewed.
5. Bind, invoke, and record independent design and clarity reviews under
   `.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`. Remediate
   only findings that review lifecycle permits directly or a human accepted. After
   each human answer, append its dispositions to `review-dispositions.md` before
   any rerun. Rerun only applicable scope; retain unaffected roles only through the
   stated retained pointer transition. At review quiescence present one collected
   blocker and proposed-remediation interaction. A started incomplete review uses
   recovery, never self-review.
6. When `plan_design` and `plan_clarity` are current and non-blocking, have the
   writer regenerate `plan-overview.md` for the current Revision, then bind,
   invoke, and record the overview fidelity review. The writer applies fidelity
   fixes that only bring the overview into agreement with `plan.md`, then the
   review reruns. A fidelity finding located in `plan.md` is a proposed
   remediation and returns to step 5 after the material change.
7. Request Design Acceptance only when all three pointers are current,
   attributable, and non-blocking for the current plan Revision, and present
   `plan-overview.md` as the human's reading basis. Record the human decision and
   rationale in `plan.md`. Only current Revision Accepted permits `tasks`.
